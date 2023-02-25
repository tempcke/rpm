//go:build withDocker
// +build withDocker

package main_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/tempcke/rpm/internal/config"
	"github.com/tempcke/rpm/internal/test"
	"github.com/tempcke/rpm/specifications"
)

var conf = test.GetConfig()
var httpClient = &http.Client{
	Timeout: 1 * time.Second,
}

func TestHTTP(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	var driver = httpDriver{
		BaseURL: "http://localhost:" + conf.GetString(config.AppPort),
		Client:  httpClient,
	}

	specifications.AddRental(t, driver)
}
