package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/internal/db/postgres/migrate"
)

var _db *sql.DB
var driverName = "postgres"

var (
	ErrConnectionFailed = errors.New("failed to connect to database")
	ErrMigrationsFailed = errors.New("migrate up failed")
)

func DB(dsn string) (*sql.DB, error) {
	if _db == nil || _db.Ping() != nil {
		logger := slog.Default().With("func", "postgres.DB")

		db, err := sql.Open(driverName, dsn)
		if err != nil {
			logger.With("error", err).Error("sql.Open failed")
			return nil, fmt.Errorf("%w: %s", ErrConnectionFailed, err)
		}
		if err := db.Ping(); err != nil {
			logger.With("error", err).Error("db.Ping failed")
			return nil, fmt.Errorf("%w: %s", ErrConnectionFailed, err)
		}

		if err := migrate.Up(db, logger); err != nil {
			logger.With("error", err).Error("migrate.Up failed")
			return nil, fmt.Errorf("%w: %s", ErrMigrationsFailed, err)
		}

		_db = db
	}
	return _db, _db.Ping()
}

func NewDB(c Config) (*sql.DB, error) {
	if dsn := c.GetString(internal.EnvPostgresDSN); dsn != "" {
		return DB(dsn)
	}
	return DB(fmt.Sprintf(
		"host=%s port=%v user=%s password=%s dbname=%s sslmode=%s",
		c.GetString(internal.EnvPostgresHost),
		c.GetString(internal.EnvPostgresPort),
		c.GetString(internal.EnvPostgresUser),
		c.GetString(internal.EnvPostgresPass),
		c.GetString(internal.EnvPostgresDB),
		c.GetString(internal.EnvPostgresSSLMode)))
}

type Config interface {
	GetString(key string) string
}
