//go:build withDocker
// +build withDocker

package main_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/tempcke/rpm/api/rest"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/internal/test"
	"github.com/tempcke/rpm/specifications"
)

var _ specifications.Driver = rest.Driver{}
var conf = test.Config()

// TestAcceptanceOpenAPI runs against the binary running in docker, so you can't debug
// through to the server code. on local make sure you `make dockerRestartApp`
// before running this else it won't test the latest copy of the code
// use rest/server_test to debug.
// the main difference is that this tests the server as built by main()
func TestAcceptanceOpenAPI_specifications(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	driver := restDriver() // oapiClient()
	specifications.RunAllTests(t, driver, driver)
}
func restDriver() rest.Driver {
	return rest.Driver{
		BaseURL: "http://localhost:" + conf.GetString(internal.EnvAppPort),
		Client: &http.Client{
			Timeout: 1 * time.Second,
		},
	}
}
