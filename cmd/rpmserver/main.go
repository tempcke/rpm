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
	"github.com/tempcke/rpm/internal/config"
	"github.com/tempcke/rpm/internal/db/postgres"
	"github.com/tempcke/rpm/pkg/log"
	"github.com/tempcke/rpm/repository"
	"google.golang.org/grpc"
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
		c <- restServer(conf, db, log)
	}()

	go func() {
		c <- grpcServer(conf, db, log)
	}()

	return <-c
}

func restServer(conf config.Config, db *sql.DB, log logrus.FieldLogger) error {
	var (
		server = rest.NewServer(repo(db))
		port   = ":" + conf.GetString(config.AppPort)
	)
	if port == ":" {
		return errors.New(config.AppPort + " not configured")
	}

	log.Info("Listening on " + port)
	fmt.Println("Listening on " + port)
	return http.ListenAndServe(port, server)
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
	s := grpc.NewServer()
	rpcServer := rpc.NewServer(actions.NewActions(repo(db)))
	pb.RegisterRPMServer(s, rpcServer)

	log.Info("Listening on " + port)
	fmt.Println("Listening on " + port)
	return s.Serve(lis)
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
