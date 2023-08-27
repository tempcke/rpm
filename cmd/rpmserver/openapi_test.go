//go:build withDocker
// +build withDocker

package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tempcke/rpm/api/rest"
	"github.com/tempcke/rpm/api/rest/openapi"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/internal/config"
	"github.com/tempcke/rpm/internal/test"
	"github.com/tempcke/rpm/specifications"
)

var _ specifications.Driver = (*oapiDriver)(nil)
var conf = test.GetConfig()
var httpClient = &http.Client{
	Timeout: 1 * time.Second,
}

// TestAcceptanceOpenAPI runs against the binary running in docker, so you can't debug
// through to the server code. on local make sure you `make dockerRestartApp`
// before running this else it won't test the latest copy of the code
// use rest/server_test to debug.
// the main difference is that this tests the server as built by main()
func TestAcceptanceOpenAPI_Property(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	var tests = map[string]struct {
		specTest func(*testing.T, specifications.PropertyDriver)
	}{
		"StoreProperty":  {specifications.AddRental},
		"GetProperty":    {specifications.GetProperty},
		"ListProperties": {specifications.ListProperties},
		"RemoveProperty": {specifications.RemoveProperty},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.specTest(t, oapiClient())
		})
	}
}
func TestAcceptanceOpenAPI_Tenant(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	tests := map[string]struct {
		specTest func(*testing.T, specifications.TenantDriver)
	}{
		"StoreTenant": {specifications.AddTenant},
		"GetTenant":   {specifications.GetTenant},
		"ListTenants": {specifications.ListTenants},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.specTest(t, oapiClient())
		})
	}
}
func oapiClient() *oapiDriver {
	return &oapiDriver{
		BaseURL: "http://localhost:" + conf.GetString(config.AppPort),
		Client:  httpClient,
	}
}

type ID = entity.ID
type oapiDriver struct {
	BaseURL string
	Client  *http.Client
	Logger  logrus.FieldLogger
}

func (d oapiDriver) StoreProperty(ctx context.Context, p entity.Property) (string, error) {
	body := map[string]string{
		"street": p.Street,
		"city":   p.City,
		"state":  p.StateCode,
		"zip":    p.Zip,
	}
	url := d.BaseURL + "/property"
	req := postReq(url, body, d.headers())
	if p.ID != "" {
		url = d.BaseURL + "/property/" + p.GetID()
		req = putReq(url, body, d.headers())
	}
	res, err := d.Client.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	var created openapi.Property
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		return "", err
	}
	return created.GetID(), nil
}
func (d oapiDriver) GetProperty(ctx context.Context, id ID) (*entity.Property, error) {
	url := d.BaseURL + "/property/" + id
	req := getReq(url, d.headers()).WithContext(ctx)
	res, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if code := res.StatusCode; code >= 400 {
		return nil, fmt.Errorf("expected 200 response, got %d", code)
	}
	var p openapi.Property
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		return nil, err
	}
	property := p.ToProperty()
	return &property, nil
}
func (d oapiDriver) ListProperties(ctx context.Context) ([]entity.Property, error) {
	url := d.BaseURL + "/property"
	req := getReq(url, d.headers()).WithContext(ctx)
	res, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	var list openapi.PropertyList
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		return nil, err
	}
	return list.ToProperties(), nil
}
func (d oapiDriver) RemoveProperty(ctx context.Context, id ID) error {
	url := d.BaseURL + "/property/" + id
	req := delReq(url, d.headers()).WithContext(ctx)
	res, err := d.Client.Do(req)
	if err != nil {
		return err
	}
	if code := res.StatusCode; code >= 400 {
		return fmt.Errorf("expected 204 response, got %d", code)
	}
	return nil
}

func (d oapiDriver) StoreTenant(ctx context.Context, tenant entity.Tenant) (*entity.Tenant, error) {
	body := openapi.ToStoreTenantReq(tenant)
	route := "/tenant"
	req := postReq(d.url(route), body, d.headers())
	if tenant.ID != "" {
		route = "/tenant/" + tenant.ID
		req = putReq(d.url(route), body, d.headers())
	}
	res, err := d.Client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	return d.getTenantRes(res)
}
func (d oapiDriver) GetTenant(ctx context.Context, id entity.ID) (*entity.Tenant, error) {
	var (
		route = "/tenant/" + id
		req   = getReq(d.url(route), d.headers())
	)
	res, err := d.Client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	return d.getTenantRes(res)
}
func (d oapiDriver) ListTenants(ctx context.Context) ([]entity.Tenant, error) {
	var (
		route = "/tenant"
		req   = getReq(d.url(route), d.headers())
		list  openapi.TenantList
	)
	res, err := d.Client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		return nil, err
	}
	return list.ToTenants(), nil
}
func (d oapiDriver) getTenantRes(r *http.Response) (*entity.Tenant, error) {
	var res openapi.GetTenantRes
	if err := d.decodeResponse(r, &res); err != nil {
		return nil, err
	}
	return res.Tenant.ToTenant(), nil
}

func (d oapiDriver) headers() map[string]string {
	conf := test.GetConfig()
	headers := map[string]string{
		rest.HeaderAPIKey:    conf.GetString(internal.EnvAPIKey),
		rest.HeaderAPISecret: conf.GetString(internal.EnvAPISecret),
	}
	return headers
}
func (d oapiDriver) url(route string) string {
	return d.BaseURL + route
}
func (d oapiDriver) decodeResponse(res *http.Response, v interface{}) error {
	var (
		buf    bytes.Buffer
		tr     = io.TeeReader(res.Body, &buf)
		errRes openapi.ErrorResponse
	)

	if res.StatusCode >= 400 {
		if err := json.NewDecoder(tr).Decode(&errRes); err != nil {
			errRes = openapi.ErrorResponse{
				Error: openapi.Error{
					Code:    int32(res.StatusCode),
					Message: buf.String(),
				},
			}
		}
		return errRes.Error
	}

	if err := json.NewDecoder(tr).Decode(&v); err != nil {
		bodyString := buf.String()
		d.logErr(err, "jsonDecode response failed", logrus.Fields{
			"func":    "oapiDriver.decodeResponse",
			"rawBody": bodyString,
		})
		if len(bodyString) == 0 {
			return errors.New("could not decode empty response body")
		}
		return err
	}

	return nil
}
func (d oapiDriver) logErr(err error, msg string, fields logrus.Fields) {
	d.log().WithFields(fields).WithError(err).Error(msg)
}
func (d oapiDriver) log() logrus.FieldLogger {
	var logger = d.Logger
	if logger == nil {
		logger = logrus.StandardLogger()
	}
	return logger.WithField("object", "oapiDriver")
}
func getReq(route string, headers map[string]string) *http.Request {
	return httpReq(http.MethodGet, route, nil, headers)
}
func delReq(route string, headers map[string]string) *http.Request {
	return httpReq(http.MethodDelete, route, nil, headers)
}
func postReq(route string, body any, headers map[string]string) *http.Request {
	return httpReq(http.MethodPost, route, body, headers)
}
func putReq(route string, body any, headers map[string]string) *http.Request {
	return httpReq(http.MethodPut, route, body, headers)
}
func httpReq(method string, route string, body interface{}, headers map[string]string) *http.Request {
	req, err := newReqBuilder(method, route).
		WithBody(body).WithHeaders(headers).Build()
	if err != nil {
		panic("reqBuilder.Build() failed: " + err.Error())
	}
	return req
}

type reqBuilder struct {
	method, route string
	body          any
	header        http.Header
}

func newReqBuilder(method, route string) *reqBuilder {
	return &reqBuilder{
		method: method,
		route:  route,
		header: make(http.Header),
	}
}
func (b reqBuilder) WithHeaders(headers map[string]string) reqBuilder {
	for k, v := range headers {
		b.header.Add(k, v)
	}
	return b
}
func (b reqBuilder) WithBody(body any) reqBuilder {
	if body != nil {
		b.body = body
	}
	return b
}
func (b reqBuilder) Build() (*http.Request, error) {
	var reqBody = &bytes.Buffer{}
	if b.body != nil {
		if err := json.NewEncoder(reqBody).Encode(b.body); err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(b.method, b.route, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header = b.header
	if reqBody.Len() > 0 {
		req.Header.Add("Content-Type", "application/json")
	}
	return req, nil
}
