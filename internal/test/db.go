package test

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq" // db driver
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/internal/db/postgres"
)

func DB(t testing.TB) *sql.DB {
	t.Helper()
	db, err := postgres.NewDB(Config())
	require.NoError(t, err)
	require.NoError(t, db.Ping())
	return db
}
