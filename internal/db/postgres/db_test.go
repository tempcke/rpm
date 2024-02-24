//go:build withDocker
// +build withDocker

package postgres_test

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq" // db driver
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/internal/db/postgres"
	"github.com/tempcke/rpm/internal/test"
)

func TestDB(t *testing.T) {
	tblName := test.RandString(6)
	db, err := postgres.NewDB(test.Config())
	require.NoError(t, err)
	query := fmt.Sprintf(`create table if not exists %s (id text primary key);`, tblName)
	_, err = db.Exec(query)
	require.NoError(t, err)
	defer func() {
		_, err = db.Exec("drop table if exists " + tblName)
		require.NoError(t, err)
	}()
}
