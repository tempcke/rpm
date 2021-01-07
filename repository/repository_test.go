package repository_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/usecase"
)

type Repo usecase.PropertyRepository

// These tests are setup this way so that the same tests
// can be run with different *sql.DB instances
// postgres_integration_test.go uses a real pg instance
// in docker which are slow (2 to 5 seconds)
func runTestsWithRepo(t *testing.T, r Repo) {
	t.Run("store and retrieve property", func(t *testing.T) {
		testStoreAndRetrieveProperty(t, r)
	})
	t.Run("list properties", func(t *testing.T) {
		testListProperties(t, r)
	})
	t.Run("remove property", func(t *testing.T) {
		testRemoveProperty(t, r)
	})
}

func testStoreAndRetrieveProperty(t *testing.T, r Repo) {
	pIn := newPropertyFixture(r)
	r.StoreProperty(pIn)
	pOut, err := r.RetrieveProperty(pIn.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, pIn.ID, pOut.ID)
	assert.Equal(t, pIn.Street, pOut.Street)
	assert.Equal(t, pIn.City, pOut.City)
	assert.Equal(t, pIn.StateCode, pOut.StateCode)
	assert.Equal(t, pIn.Zip, pOut.Zip)
	assertTimestampMatch(t, pIn.CreatedAt, pOut.CreatedAt)
}

func testListProperties(t *testing.T, r Repo) {
	props := make(map[string]entity.Property, 3)
	for i := 0; i < 3; i++ {
		// create entity
		p := newPropertyFixture(r)
		props[p.ID] = p

		// store entity
		err := r.StoreProperty(p)
		assert.NoError(t, err)
	}

	// list entities, this is what we want to test!
	propList, err := r.PropertyList()
	assert.NoError(t, err)

	// iterate over list counting the times each id is seen
	seen := make(map[string]int, 3)
	for _, p := range propList {
		delete(props, p.ID)
		seen[p.ID]++
	}

	// ensure the list does not repeat any properties
	for _, n := range seen {
		assert.Equal(t, 1, n, "property ids are repeated in the list")
	}

	// properties are deleted as they are found so there should be none left
	assert.Len(t, props, 0)
}

func testRemoveProperty(t *testing.T, r Repo) {
	// create and store property
	p := newPropertyFixture(r)
	r.StoreProperty(p)
	// remove property
	err := r.DeleteProperty(p.ID)
	assert.NoError(t, err)
	// try to retrieve property
	_, err = r.RetrieveProperty(p.ID)
	assert.Error(t, err)
}

// helper functions
var streetNum = 100

func newPropertyFixture(r Repo) entity.Property {
	streetNum++
	street := fmt.Sprintf("%v N Main st.", streetNum)
	return r.NewProperty(street, "Dallas", "TX", "75401")
}

func loadMigrations(db *sql.DB) error {
	query := `
	  CREATE TABLE IF NOT EXISTS properties (
			id         VARCHAR(36)  PRIMARY KEY,
			street     VARCHAR(255),
			city       VARCHAR(32),
			state      VARCHAR(32),
			zip        VARCHAR(10),
			created_at TIMESTAMPTZ
		)
	`
	_, err := db.Exec(query)
	return err
}

// custom assertions
func assertTimestampMatch(t *testing.T, t1, t2 time.Time) {
	t.Helper()

	// ensure matches within 1 second
	timeDiff := t1.Sub(t2).Seconds()
	assert.LessOrEqual(t, timeDiff, 1.0)

	t1Zone, _ := t1.Zone()
	t2Zone, _ := t2.Zone()

	// ensure time converted back to local time
	assert.Equal(t, t1Zone, t2Zone)
}
