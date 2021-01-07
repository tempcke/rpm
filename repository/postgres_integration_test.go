// These tests spawn a postgres docker container
// as a consequence they are more reliably true
// however they are also slow (2 to 5 seconds)
package repository_test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/tempcke/rpm/repository"
)

func TestPropertyRepositoryIntegration(t *testing.T) {
	// Setup docker postgres instance
	var (
		db *sql.DB

		user     = "postgres"
		password = "secret"
		dbname   = "postgres"
		port     = "54322"
		dialect  = "postgres"
		dsn      = fmt.Sprintf(
			"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
			user, password, port, dbname,
		)
	)

	pool, err := dockertest.NewPool("")

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.3",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbname,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err.Error())
	}

	if err = pool.Retry(func() error {
		db, err = sql.Open(dialect, dsn)
		if err != nil {
			return err
		}
		if err = db.Ping(); err != nil {
			return err
		}
		if err = loadMigrations(db); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err.Error())
	}

	defer func() {
		db.Close()
	}()

	// construct repository
	repo := repository.NewPostgresRepo(db)

	// run the tests
	runTestsWithRepo(t, repo)

	// cleanup
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
