package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	_ "github.com/lib/pq" // db driver
	"github.com/sirupsen/logrus"
	"github.com/tempcke/rpm/actions"
	"github.com/tempcke/rpm/api/rest"
	"github.com/tempcke/rpm/api/rpc"
	pb "github.com/tempcke/rpm/api/rpc/proto"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/internal/config"
	"github.com/tempcke/rpm/internal/db/postgres"
	"github.com/tempcke/rpm/internal/lib/log"
	"github.com/tempcke/rpm/internal/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	if err := run(log.Entry()); err != nil {
		log.Fatal(err)
	}
}

func run(log logrus.FieldLogger) error {
	var (
		conf = config.GetConfig()
		c    = make(chan error)
	)

	db, err := postgres.DB(postgres.MakeDSN(conf))
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer func() { _ = db.Close() }()

	// TODO: if either of these fail the other is orphaned but we
	//  return to main and the program exits.  we must find a way
	//  to gracefully shut down the other
	go func() {
		c <- openapiServer(conf, db, log)
	}()

	go func() {
		c <- grpcServer(conf, db, log)
	}()

	return <-c
}

func openapiServer(conf config.Config, db *sql.DB, log logrus.FieldLogger) error {
	var (
		r    = repo(db)
		acts = actions.NewActions().
			WithPropertyRepo(r).WithTenantRepo(r)
		port = ":" + conf.GetString(config.AppPort)
	)
	if port == ":" {
		return errors.New(config.AppPort + " not configured")
	}

	server := rest.NewServer(acts).WithConfig(conf)

	log.Info("Listening on " + port)
	fmt.Println("Listening on " + port)
	return http.ListenAndServe(port, server.Handler())
}
func grpcServer(conf config.Config, db *sql.DB, log logrus.FieldLogger) error {
	var (
		port = ":" + conf.GetString(config.GrpcPort)
	)
	if port == ":" {
		return errors.New(config.GrpcPort + " not configured")
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
func grpcOptions(conf config.Config, log logrus.FieldLogger) []grpc.ServerOption {
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
