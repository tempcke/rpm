package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tempcke/rpm/api/rest"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/internal/test"
	"github.com/tempcke/rpm/specifications"
)

var _ specifications.Driver = (*httpDriver)(nil)

type ID = string
type httpDriver struct {
	BaseURL string
	Client  *http.Client
}

func (d httpDriver) StoreProperty(ctx context.Context, p entity.Property) (string, error) {
	if p.ID != "" {
		return p.ID, d.addRentalWithID(ctx, p)
	}
	url := d.BaseURL + "/property"
	body := map[string]string{
		"street": p.Street,
		"city":   p.City,
		"state":  p.StateCode,
		"zip":    p.Zip,
	}
	req := postReq(url, body, d.headers()).WithContext(ctx)
	res, err := d.Client.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusCreated {
		resData := internal.SPrintData("httpDriver.StoreProperty response body", res.Body)
		return "", fmt.Errorf("expected status 201 but got %v\n%s", res.StatusCode, resData)
	}
	var created rest.PropertyModel
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		return "", err
	}
	return created.ID, nil
}
func (d httpDriver) addRentalWithID(ctx context.Context, p entity.Property) error {
	url := d.BaseURL + "/property/" + p.GetID()
	body := map[string]string{
		"street": p.Street,
		"city":   p.City,
		"state":  p.StateCode,
		"zip":    p.Zip,
	}
	req := putReq(url, body, d.headers()).WithContext(ctx)
	res, err := d.Client.Do(req)
	if err != nil {
		return err
	}
	var created rest.PropertyModel
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		return err
	}
	return nil
}

func (d httpDriver) GetProperty(ctx context.Context, id ID) (*entity.Property, error) {
	url := d.BaseURL + "/property/" + id
	req := getReq(url, d.headers()).WithContext(ctx)
	res, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if code := res.StatusCode; code >= 400 {
		return nil, fmt.Errorf("expected 200 response, got %d", code)
	}
	var p rest.PropertyModel
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		return nil, err
	}
	return p.ToProperty(), nil
}
func (d httpDriver) ListProperties(ctx context.Context) ([]entity.Property, error) {
	url := d.BaseURL + "/property"
	req := getReq(url, d.headers()).WithContext(ctx)
	res, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	var list rest.PropertyList
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		return nil, err
	}
	return list.ToProperties(), nil
}
func (d httpDriver) RemoveProperty(ctx context.Context, id ID) error {
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

func (d httpDriver) headers() map[string]string {
	conf := test.GetConfig()
	headers := map[string]string{
		rest.HeaderAPIKey:    conf.GetString(internal.EnvAPIKey),
		rest.HeaderAPISecret: conf.GetString(internal.EnvAPISecret),
	}
	return headers
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
