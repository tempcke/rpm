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

type Driver interface {
	PropertyDriver
}
type PropertyDriver interface {
	StoreProperty(context.Context, entity.Property) (entity.ID, error)
	GetProperty(context.Context, entity.ID) (*entity.Property, error)
	ListProperties(context.Context) ([]entity.Property, error)
	RemoveProperty(context.Context, entity.ID) error
}
type TenantDriver interface {
	StoreTenant(context.Context, entity.Tenant) (*entity.Tenant, error)
	GetTenant(context.Context, entity.ID) (*entity.Tenant, error)
	ListTenants(context.Context) ([]entity.Tenant, error)
}

func AddRental(t *testing.T, driver PropertyDriver) {
	t.Run("without ID", func(t *testing.T) {
		var pIn = fake.Property().WithID("")
		id, err := driver.StoreProperty(ctx, pIn)
		require.NoError(t, err)
		require.NotEqual(t, "", id)
		pOut, err := driver.GetProperty(ctx, id)
		require.NoError(t, err)
		assert.NotNil(t, pOut)
		assert.True(t, pIn.Equal(*pOut))
	})
	t.Run("with ID", func(t *testing.T) {
		// we allow client to define ID
		var pIn = fake.Property()
		id, err := driver.StoreProperty(ctx, pIn)
		require.NoError(t, err)
		require.Equal(t, pIn.ID, id)
		pOut, err := driver.GetProperty(ctx, id)
		require.NoError(t, err)
		assert.NotNil(t, pOut)
		assert.True(t, pIn.Equal(*pOut))
	})
}
func GetProperty(t *testing.T, driver PropertyDriver) {
	var pIn = fake.Property()
	_, err := driver.StoreProperty(ctx, pIn)
	require.NoError(t, err)
	pOut, err := driver.GetProperty(ctx, pIn.GetID())
	require.NoError(t, err)
	assert.NotNil(t, pOut)
	assert.True(t, pIn.Equal(*pOut))
}
func ListProperties(t *testing.T, driver PropertyDriver) {
	in1 := fake.Property()
	in2 := fake.Property()
	_, err := driver.StoreProperty(ctx, in1)
	require.NoError(t, err)
	_, err = driver.StoreProperty(ctx, in2)
	require.NoError(t, err)

	list, err := driver.ListProperties(ctx)
	require.NoError(t, err)
	m := entityMap(list...)
	require.Contains(t, m, in1.GetID())
	require.Contains(t, m, in2.GetID())
}
func RemoveProperty(t *testing.T, driver PropertyDriver) {
	in1 := fake.Property()
	in2 := fake.Property()
	_, err := driver.StoreProperty(ctx, in1)
	require.NoError(t, err)
	_, err = driver.StoreProperty(ctx, in2)
	require.NoError(t, err)

	require.NoError(t, driver.RemoveProperty(ctx, in1.ID))
	pOut, err := driver.GetProperty(ctx, in1.ID)
	assert.Error(t, err)
	assert.Nil(t, pOut)

	list, err := driver.ListProperties(ctx)
	require.NoError(t, err)
	m := entityMap(list...)
	require.NotContains(t, m, in1.GetID())
	require.Contains(t, m, in2.GetID())
}

func AddTenant(t *testing.T, driver TenantDriver) {
	t.Run("without ID", func(t *testing.T) {
		var in = fake.Tenant().WithID("")
		require.Equal(t, "", in.ID)
		out, err := driver.StoreTenant(ctx, in)
		require.NoError(t, err)
		require.NotNil(t, out)
		require.NotEmpty(t, out.GetID(), "expected ID to be assigned")
		assert.True(t, in.Equal(*out))
	})
	t.Run("with ID", func(t *testing.T) {
		var in = fake.Tenant()
		require.NotEqual(t, "", in.ID)
		out, err := driver.StoreTenant(ctx, in)
		require.NotNil(t, out)
		require.NoError(t, err)
		require.Equal(t, in.ID, out.ID, "should have used provided ID")
		assert.True(t, in.Equal(*out))
	})
}
func GetTenant(t *testing.T, driver TenantDriver) {
	in := fake.Tenant()
	_, err := driver.StoreTenant(ctx, in)
	require.NoError(t, err)

	out, err := driver.GetTenant(ctx, in.GetID())
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.True(t, out.Equal(in))
}
func ListTenants(t *testing.T, driver TenantDriver) {
	in1 := fake.Tenant()
	in2 := fake.Tenant()
	_, err := driver.StoreTenant(ctx, in1)
	require.NoError(t, err)
	_, err = driver.StoreTenant(ctx, in2)
	require.NoError(t, err)

	list, err := driver.ListTenants(ctx)
	require.NoError(t, err)
	m := entityMap(list...)
	require.Contains(t, m, in1.GetID())
	require.Contains(t, m, in2.GetID())
}

func entityMap[T entity.Entity](entities ...T) map[entity.ID]entity.Entity {
	var m = make(map[entity.ID]entity.Entity)
	for _, e := range entities {
		m[e.GetID()] = e
	}
	return m
}
