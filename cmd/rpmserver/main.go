package main

import (
	"fmt"
	"net/http"

	_ "github.com/lib/pq" // db driver
	"github.com/sirupsen/logrus"
	"github.com/tempcke/rpm/api/rest"
	"github.com/tempcke/rpm/internal/config"
	"github.com/tempcke/rpm/internal/db/postgres"
	"github.com/tempcke/rpm/pkg/log"
	"github.com/tempcke/rpm/repository"
)

func main() {
	if err := run(log.Entry()); err != nil {
		log.Fatal(err)
	}
}

func run(log logrus.FieldLogger) error {
	var conf = config.GetConfig()

	db, err := postgres.DB(postgres.MakeDSN(conf))
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer func() { _ = db.Close() }()

	var (
		propRepo = repository.NewPostgresRepo(db)
		server   = rest.NewServer(propRepo)
		port     = ":" + conf.GetString(config.AppPort)
	)

	log.Info("Listening on " + port)
	fmt.Println("Listening on " + port)
	return http.ListenAndServe(port, server)
}
