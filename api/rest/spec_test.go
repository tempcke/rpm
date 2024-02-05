package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tempcke/rpm/api/rest"
	"github.com/tempcke/rpm/specifications"
)

func TestOpenAPI_specifications(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	driver := restDriver(t) // oapiClient()
	specifications.RunAllTests(t, driver, driver)
}
func restDriver(t testing.TB) rest.Driver {
	var (
		server = newServer(t)
		client = newClient(t, server)
	)

	return rest.Driver{
		// no requests are actually sent so the BaseURL is irrelevant
		BaseURL: "http://example.localhost",
		Client:  client,
	}
}

type restClient struct {
	t      testing.TB
	server *rest.Server
}

func newClient(t testing.TB, server *rest.Server) restClient {
	return restClient{
		t:      t,
		server: server,
	}
}
func (c restClient) Do(req *http.Request) (*http.Response, error) {
	c.t.Helper()
	rr := httptest.NewRecorder()
	c.server.Handler().ServeHTTP(rr, req)
	return rr.Result(), nil
}
