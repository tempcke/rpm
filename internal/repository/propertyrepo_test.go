package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/entity/fake"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/usecase"
)

func testStoreAndRetrieveProperty(t *testing.T, r propertyRepo) {
	pIn := newPropertyFixture(r)
	require.NoError(t, r.StoreProperty(ctx, pIn))
	pOut, err := r.GetProperty(ctx, pIn.ID)
	require.NoError(t, err)
	assert.Equal(t, pIn.ID, pOut.ID)
	assert.Equal(t, pIn.Street, pOut.Street)
	assert.Equal(t, pIn.City, pOut.City)
	assert.Equal(t, pIn.StateCode, pOut.StateCode)
	assert.Equal(t, pIn.Zip, pOut.Zip)
	assertTimestampMatch(t, pIn.CreatedAt, pOut.CreatedAt)
}
func testUpdateProperty(t *testing.T, r propertyRepo) {
	// insert
	pIn := newPropertyFixture(r)
	require.NoError(t, r.StoreProperty(ctx, pIn))

	// update
	pIn.Street = "1" + pIn.Street
	require.NoError(t, r.StoreProperty(ctx, pIn))

	// select
	pOut, err := r.GetProperty(ctx, pIn.ID)
	require.NoError(t, err)
	assert.Equal(t, pIn.ID, pOut.ID)
	assert.Equal(t, pIn.Street, pOut.Street)
	assert.Equal(t, pIn.City, pOut.City)
	assert.Equal(t, pIn.StateCode, pOut.StateCode)
	assert.Equal(t, pIn.Zip, pOut.Zip)
	assertTimestampMatch(t, pIn.CreatedAt, pOut.CreatedAt)
}
func testListProperties(t *testing.T, r propertyRepo) {
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
	propList, err := r.PropertyList(ctx, usecase.AllProperties)
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
func testRemoveProperty(t *testing.T, r propertyRepo) {
	// create and store property
	p := newPropertyFixture(r)
	require.NoError(t, r.StoreProperty(ctx, p))
	// remove property
	err := r.DeleteProperty(ctx, p.ID)
	assert.NoError(t, err)
	// try to retrieve property
	_, err = r.GetProperty(ctx, p.ID)
	assert.Error(t, err)
}
func testGetProperty(t *testing.T, r propertyRepo) {
	id := "id-does-not-exist"
	pOut, err := r.GetProperty(ctx, id)
	require.Error(t, err)
	assert.ErrorIs(t, err, internal.ErrEntityNotFound)
	assert.Empty(t, pOut)
}

func newPropertyFixture(r propertyRepo) entity.Property {
	p := fake.Property()
	return r.NewProperty(p.Street, p.City, p.StateCode, p.Zip)
}
