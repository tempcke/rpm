package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/internal/config"
)

func noAuthConf(t testing.TB) config.Config {
	t.Helper()
	c, err := config.NewConfBuilder().AutomaticEnv().Build()
	require.NoError(t, err)
	c.Set(internal.EnvAPIKey, "")
	c.Set(internal.EnvAPISecret, "")
	return c
}
func authConf(t testing.TB, key, secret string) config.Config {
	t.Helper()
	c, err := config.NewConfBuilder().AutomaticEnv().Build()
	require.NoError(t, err)
	c.Set(internal.EnvAPIKey, key)
	c.Set(internal.EnvAPISecret, secret)
	return c
}

func getReq(t testing.TB, route string, headers map[string]string) *http.Request {
	return httpReq(t, http.MethodGet, route, nil, headers)
}
func delReq(t testing.TB, route string, headers map[string]string) *http.Request {
	return httpReq(t, http.MethodDelete, route, nil, headers)
}
func postReq(t testing.TB, route string, body any, headers map[string]string) *http.Request {
	return httpReq(t, http.MethodPost, route, body, headers)
}
func putReq(t testing.TB, route string, body any, headers map[string]string) *http.Request {
	return httpReq(t, http.MethodPut, route, body, headers)
}
func httpReq(t testing.TB, method string, route string, body interface{}, headers map[string]string) *http.Request {
	t.Helper()
	req, err := newReqBuilder(method, route).
		WithBody(body).WithHeaders(headers).Build()
	require.NoError(t, err)
	return req
}
func execReq(t testing.TB, req *http.Request) *http.Response {
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return res
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

func assertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf(
			"Values not equal \nWant: %v \t%T\nGot:  %v \t%T\n",
			expected, expected, actual, actual,
		)
	}
}
func assertTimeRecent(t testing.TB, actual time.Time) {
	t.Helper()
	diff := time.Since(actual).Abs().Seconds()
	assert.Less(t, diff, 5.0)
}
func assertApplicationJson(t testing.TB, header http.Header) {
	t.Helper()
	require.Contains(t, header.Get("Content-Type"), "application/json")
}
