//go:build withDocker
// +build withDocker

package main_test

import (
	"testing"

	"github.com/tempcke/rpm/internal/config"
	"github.com/tempcke/rpm/specifications"
)

// TestGRPC runs against the binary running in docker, so you can't debug
// through to the server code. on local make sure you `make dockerRestartApp`
// before running this else it won't test the latest copy of the code
// use rpc/server_test to debug.
// the main difference is that this tests the server as built by main()
func TestGRPC_property(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	var driver = grpcDriver{
		Addr: "localhost:" + conf.GetString(config.GrpcPort),
	}
	type (
		specTest func(*testing.T, specifications.PropertyDriver)
	)
	var tests = map[string]struct{ specTest }{
		"StoreProperty":  {specifications.AddRental},
		"GetProperty":    {specifications.GetProperty},
		"ListProperties": {specifications.ListProperties},
		"RemoveProperty": {specifications.RemoveProperty},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.specTest(t, &driver)
		})
	}
}
