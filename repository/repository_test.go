package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/usecase"
)

type Repo usecase.PropertyRepository

var ctx = context.Background()

func testStoreAndRetrieveProperty(t *testing.T, r Repo) {
	pIn := newPropertyFixture(r)
	require.NoError(t, r.StoreProperty(ctx, pIn))
	pOut, err := r.RetrieveProperty(ctx, pIn.ID)
	require.NoError(t, err)
	assert.Equal(t, pIn.ID, pOut.ID)
	assert.Equal(t, pIn.Street, pOut.Street)
	assert.Equal(t, pIn.City, pOut.City)
	assert.Equal(t, pIn.StateCode, pOut.StateCode)
	assert.Equal(t, pIn.Zip, pOut.Zip)
	assertTimestampMatch(t, pIn.CreatedAt, pOut.CreatedAt)
}

func testUpdateProperty(t *testing.T, r Repo) {
	// insert
	pIn := newPropertyFixture(r)
	require.NoError(t, r.StoreProperty(ctx, pIn))

	// update
	pIn.Street = "1" + pIn.Street
	require.NoError(t, r.StoreProperty(ctx, pIn))

	// select
	pOut, err := r.RetrieveProperty(ctx, pIn.ID)
	require.NoError(t, err)
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
		err := r.StoreProperty(ctx, p)
		assert.NoError(t, err)
	}

	// list entities, this is what we want to test!
	propList, err := r.PropertyList(ctx)
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
	require.NoError(t, r.StoreProperty(ctx, p))
	// remove property
	err := r.DeleteProperty(ctx, p.ID)
	assert.NoError(t, err)
	// try to retrieve property
	_, err = r.RetrieveProperty(ctx, p.ID)
	assert.Error(t, err)
}

// helper functions
var streetNum = 100

func newPropertyFixture(r Repo) entity.Property {
	streetNum++
	street := fmt.Sprintf("%v N Main st.", streetNum)
	return r.NewProperty(street, "Dallas", "TX", "75401")
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
