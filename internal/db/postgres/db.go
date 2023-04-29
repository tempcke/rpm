package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tempcke/rpm/internal/db/postgres/migrate"
	"github.com/tempcke/rpm/internal/lib/log"
)

var _db *sql.DB
var driverName = "postgres"

var (
	ErrConnectionFailed = errors.New("failed to connect to database")
	ErrMigrationsFailed = errors.New("migrate up failed")
)

type Config interface {
	PostgresDSN() string
	PostgresHost() string
	PostgresPort() int
	PostgresUser() string
	PostgresPass() string
	PostgresDB() string
}

func MakeDSN(conf Config) string {
	if dsn := conf.PostgresDSN(); dsn != "" {
		return dsn
	}
	dsn := fmt.Sprintf(
		"host=%s port=%v user=%s password=%s dbname=%s sslmode=disable",
		conf.PostgresHost(),
		conf.PostgresPort(),
		conf.PostgresUser(),
		conf.PostgresPass(),
		conf.PostgresDB())
	return dsn
}

func DB(dsn string) (*sql.DB, error) {
	if _db == nil || _db.Ping() != nil {
		logger := log.WithField("func", "postgres.DB")

		db, err := sql.Open(driverName, dsn)
		if err != nil {
			logger.WithError(err).Error("sql.Open failed")
			return nil, fmt.Errorf("%w: %s", ErrConnectionFailed, err)
		}
		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrConnectionFailed, err)
		}

		if err := migrate.Up(db, logger); err != nil {
			logger.WithError(err).Error("migrate.Up failed")
			return nil, fmt.Errorf("%w: %s", ErrMigrationsFailed, err)
		}

		_db = db
	}
	return _db, _db.Ping()
}
