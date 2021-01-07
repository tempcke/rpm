package repository_test

import (
	"database/sql"
	"testing"

	_ "github.com/mlhoyt/ramsql/driver"
	//	_ "github.com/proullon/ramsql/driver"

	"github.com/tempcke/rpm/repository"
)

func TestPropertyRepository(t *testing.T) {
	t.Skip("ramsql fails to Scan a TIMEZONE field into a time.Time, if that is ever fixed we can have faster tests here perhaps...")

	db, err := sql.Open("ramsql", "TestPropertyRepository")
	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}
	defer db.Close()

	if err = loadMigrations(db); err != nil {
		t.Error(err)
	}

	repo := repository.NewPostgresRepo(db)

	runTestsWithRepo(t, repo)
}
