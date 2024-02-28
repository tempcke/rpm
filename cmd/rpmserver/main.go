package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"sync"

	_ "github.com/lib/pq" // db driver
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
	if err := run(os.Getenv, os.Args[1:]...); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func run(
	envFunc func(string) string,
	args ...string,
) error {
	var (
		errChan = make(chan error)
		conf    = buildConfig(envFunc, args...)
		logger  = initLogger(conf)
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

func initLogger(conf configs.Config) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"appEnv", conf.GetString(internal.EnvAppEnv))
	slog.SetDefault(logger)
	return slog.Default()
}

func openapiServer(conf Config, db *sql.DB, log log.SLogger) error {
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
	return http.ListenAndServe(port, server.Handler())
}
func grpcServer(conf Config, db *sql.DB, log *slog.Logger) error {
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
	options, err := grpcOptions(conf)
	if err != nil {
		return err
	}
	s := grpc.NewServer(options...)
	r := repo(db)
	rpcServer := rpc.NewServer(actions.NewActions().
		WithPropertyRepo(r).WithTenantRepo(r))
	pb.RegisterRPMServer(s, rpcServer)

	log.Info("Listening on " + port)
	fmt.Println("Listening on " + port)
	return s.Serve(lis)
}
func grpcOptions(conf Config) ([]grpc.ServerOption, error) {
	var (
		certFile = conf.GetString(internal.EnvServiceCertFile)
		keyFile  = conf.GetString(internal.EnvServiceKeyFile)
	)
	if certFile == "" || keyFile == "" {
		return nil, errors.New("grpcOptions: failed to setup TLS: cert and key file not configured")
	}
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("grpcOptions: credentials.NewServerTLSFromFile failed: %w", err)
	}
	return []grpc.ServerOption{grpc.Creds(creds)}, nil
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

func buildConfig(envFunc func(string) string, args ...string) configs.Config {
	return configs.New(
		configs.WithFlagSet(getFlagSet()),
		configs.WithEnvFunc(envFunc),
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
