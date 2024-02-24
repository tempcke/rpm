//go:build withDocker
// +build withDocker

package main_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/api/rpc"
	pb "github.com/tempcke/rpm/api/rpc/proto"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/specifications"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// TestGRPC runs against the binary running in docker, so you can't debug
// through to the server code. on local make sure you `make dockerRestartApp`
// before running this else it won't test the latest copy of the code
// use rpc/server_test to debug.
// the main difference is that this tests the server as built by main()
func TestAcceptanceGRPC_specifications(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	driver := rpcDriver(t)
	specifications.RunAllTests(t, driver, driver)
}
func rpcDriver(t testing.TB) rpc.Driver {
	var (
		addr        = "localhost:" + conf.GetString(internal.EnvGrpcPort)
		client, err = newClient(t, addr)
	)
	require.NoError(t, err)
	return rpc.NewDriver(client)
}
func newClient(t testing.TB, addr string) (pb.RPMClient, error) {
	var (
		conn    *grpc.ClientConn
		connErr error
		client  pb.RPMClient

		certFile      = conf.GetString(internal.EnvServiceCertFile)
		credentialOpt = grpc.WithTransportCredentials(insecure.NewCredentials())
	)

	if file := findCertFile(certFile); file != "" {
		creds, err := credentials.NewClientTLSFromFile(file, "")
		if err != nil {
			return nil, err
		}
		credentialOpt = grpc.WithTransportCredentials(creds)
	}

	dialOpts := []grpc.DialOption{credentialOpt}
	conn, connErr = grpc.Dial(addr, dialOpts...)
	t.Cleanup(func() {
		if conn != nil {
			_ = conn.Close()
		}
	})
	client = pb.NewRPMClient(conn)

	if err := connErr; err != nil {
		return nil, err
	}
	return client, nil
}
func findCertFile(relPath string) string {
	if relPath == "" {
		return ""
	}
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(relPath); err == nil {
			_ = godotenv.Load(relPath)
			return relPath
		}
		relPath = fmt.Sprintf("../%s", relPath)
	}
	return ""
}
