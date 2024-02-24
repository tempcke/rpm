package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	_ "github.com/lib/pq" // db driver
	"github.com/sirupsen/logrus"
	"github.com/tempcke/rpm/actions"
	"github.com/tempcke/rpm/api/rest"
	"github.com/tempcke/rpm/api/rpc"
	pb "github.com/tempcke/rpm/api/rpc/proto"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/internal/configs"
	"github.com/tempcke/rpm/internal/db/postgres"
	"github.com/tempcke/rpm/internal/lib/log"
	"github.com/tempcke/rpm/internal/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	var (
		ctx = context.Background()
	)
	if err := run(ctx, os.Getenv, os.Args[1:]...); err != nil {
		log.Fatal(err)
	}
}

func run(
	ctx context.Context,
	getenv func(string) string,
	args ...string,
) error {
	var (
		errChan = make(chan error)
		logger  = log.Entry()
		conf    = buildConfig(getenv, args...)
	)

	db, err := postgres.NewDB(conf)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer func() { _ = db.Close() }()

	// TODO: graceful shut down
	go func() { errChan <- openapiServer(conf, db, logger) }()

	go func() { errChan <- grpcServer(conf, db, logger) }()

	return <-errChan
}

func openapiServer(conf Config, db *sql.DB, log logrus.FieldLogger) error {
	var (
		r    = repo(db)
		acts = actions.NewActions().
			WithPropertyRepo(r).WithTenantRepo(r)
		port      = ":" + conf.GetString(internal.EnvAppPort)
		apiKey    = conf.GetString(internal.EnvAPIKey)
		apiSecret = conf.GetString(internal.EnvAPISecret)
	)
	if port == ":" {
		return errors.New(internal.EnvAppPort + " not configured")
	}

	server := rest.NewServer(acts).WithCredentials(apiKey, apiSecret)

	log.Info("Listening on " + port)
	fmt.Println("Listening on " + port)
	return http.ListenAndServe(port, server.Handler())
}
func grpcServer(conf Config, db *sql.DB, log logrus.FieldLogger) error {
	var (
		port = ":" + conf.GetString(internal.EnvGrpcPort)
	)
	if port == ":" {
		return errors.New(internal.EnvGrpcPort + " not configured")
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	s := grpc.NewServer(grpcOptions(conf, log)...)
	r := repo(db)
	rpcServer := rpc.NewServer(actions.NewActions().
		WithPropertyRepo(r).WithTenantRepo(r))
	pb.RegisterRPMServer(s, rpcServer)

	log.Info("Listening on " + port)
	fmt.Println("Listening on " + port)
	return s.Serve(lis)
}
func grpcOptions(conf Config, log logrus.FieldLogger) []grpc.ServerOption {
	var (
		certFile = conf.GetString(internal.EnvServiceCertFile)
		keyFile  = conf.GetString(internal.EnvServiceKeyFile)
	)
	if certFile == "" || keyFile == "" {
		log.WithField("func", "main.grpcCreds").Fatal("Failed to setup TLS: cert and key file not configured")
	}
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.WithField("func", "main.grpcCreds").Fatalf("Failed to setup TLS: %v", err)
	}
	return []grpc.ServerOption{grpc.Creds(creds)}
}

var (
	repoOnce sync.Once
	_repo    repository.Postgres
)

func repo(db *sql.DB) repository.Postgres {
	repoOnce.Do(func() {
		_repo = repository.NewPostgresRepo(db)
	})
	return _repo
}

type Config interface {
	GetString(string) string
}

func buildConfig(envFn func(string) string, args ...string) configs.Config {
	return configs.New(
		configs.WithFlagSet(getFlagSet()),
		configs.WithEnvFunc(envFn),
		configs.WithArgs(args), // os.Args[1:] from main()
	)
}
func getFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.String(internal.EnvLogLevel, "info", "debug|info|warn|error")
	fs.String(internal.EnvAppEnv, "local", "local|dev|stage|prod")
	fs.String(internal.EnvAppPort, "8080", "http service port")
	fs.String(internal.EnvGrpcPort, "8443", "grpc service port")
	fs.String(internal.EnvAPIKey, "", "api key")
	fs.String(internal.EnvAPISecret, "", "api secret")
	fs.String(internal.EnvServiceCertFile, "", "service cert file")
	fs.String(internal.EnvServiceKeyFile, "", "service key file")
	fs.String(internal.EnvPostgresHost, "localhost", "postgres host")
	fs.String(internal.EnvPostgresPort, "5432", "postgres port")
	fs.String(internal.EnvPostgresUser, "postgres", "postgres user")
	fs.String(internal.EnvPostgresPass, "password", "postgres password")
	fs.String(internal.EnvPostgresDB, "rpm", "postgres database")
	fs.String(internal.EnvPostgresSSLMode, "disable", "postgres sslmode")
	return fs
}
