package specifications

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/entity/fake"
)

var ctx = context.Background()

type ID = string

type Driver interface {
	AddRental(context.Context, entity.Property) (ID, error)
	GetProperty(context.Context, ID) (*entity.Property, error)
	ListProperties(context.Context) ([]entity.Property, error)
	RemoveProperty(context.Context, ID) error
}

func AddRental(t *testing.T, driver Driver) {
	t.Run("without ID", func(t *testing.T) {
		var pIn = fake.Property().WithID("")
		id, err := driver.AddRental(ctx, pIn)
		require.NoError(t, err)
		require.NotEqual(t, "", id)
		pOut, err := driver.GetProperty(ctx, id)
		require.NoError(t, err)
		assert.NotEmpty(t, pOut)
	})
	t.Run("with ID", func(t *testing.T) {
		// we allow client to define ID
		var pIn = fake.Property()
		id, err := driver.AddRental(ctx, pIn)
		require.NoError(t, err)
		require.Equal(t, pIn.ID, id)
		pOut, err := driver.GetProperty(ctx, id)
		require.NoError(t, err)
		assert.NotEmpty(t, pOut)
	})
	t.Run("with ID", func(t *testing.T) {
		// we allow client to define ID
		var pIn = fake.Property()
		id, err := driver.AddRental(ctx, pIn)
		require.NoError(t, err)
		require.Equal(t, pIn.ID, id)
		pOut, err := driver.GetProperty(ctx, id)
		require.NoError(t, err)
		assert.NotEmpty(t, pOut)
	})
}
func ListProperties(t testing.TB, driver Driver) {
	var (
		p1 = fake.Property()
		p2 = fake.Property()
	)
	addProperty(t, driver, p1)
	addProperty(t, driver, p2)
	properties, err := driver.ListProperties(ctx)
	require.NoError(t, err)
	pMap := propertyMap(properties...)
	require.Contains(t, pMap, p1.ID)
	require.Contains(t, pMap, p2.ID)
	assert.Equal(t, p1, pMap[p1.ID])
	assert.Equal(t, p2, pMap[p2.ID])
}
func RemoveProperty(t testing.TB, driver Driver) {
	var (
		p1 = fake.Property()
		p2 = fake.Property()
	)

	addProperty(t, driver, p1)
	addProperty(t, driver, p2)
	require.NoError(t, driver.RemoveProperty(ctx, p1.ID))
	pOut, err := driver.GetProperty(ctx, p1.ID)
	assert.Error(t, err)
	assert.Nil(t, pOut)

	properties, err := driver.ListProperties(ctx)
	require.NoError(t, err)
	pMap := propertyMap(properties...)
	require.NotContains(t, pMap, p1.ID)
	require.Contains(t, pMap, p2.ID)
}

func addProperty(t testing.TB, driver Driver, p entity.Property) {
	t.Helper()
	id, err := driver.AddRental(ctx, p)
	require.NoError(t, err)
	if p.ID != "" {
		assert.Equal(t, p.ID, id)
	}
}
func propertyMap(properties ...entity.Property) map[ID]entity.Property {
	var pMap = make(map[ID]entity.Property)
	for _, p := range properties {
		pMap[p.ID] = p
	}
	return pMap
}