package specifications

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/entity/fake"
	"github.com/tempcke/rpm/internal/test"
	"github.com/tempcke/rpm/usecase"
)

var ctx = context.Background()

type Driver interface {
	PropertyDriver
}
type PropertyDriver interface {
	StoreProperty(context.Context, entity.Property) (entity.ID, error)
	GetProperty(context.Context, entity.ID) (*entity.Property, error)
	ListProperties(context.Context, usecase.PropertyFilter) ([]entity.Property, error)
	RemoveProperty(context.Context, entity.ID) error
}
type TenantDriver interface {
	StoreTenant(context.Context, entity.Tenant) (*entity.Tenant, error)
	GetTenant(context.Context, entity.ID) (*entity.Tenant, error)
	ListTenants(context.Context) ([]entity.Tenant, error)
}

func RunAllTests(t *testing.T, pDriver PropertyDriver, tDriver TenantDriver) {
	t.Run("property", func(t *testing.T) {
		RunAllPropertyTests(t, pDriver)
	})
	t.Run("tenant", func(t *testing.T) {
		RunAllTenantTests(t, tDriver)
	})
}
func RunAllPropertyTests(t *testing.T, driver PropertyDriver) {
	var PropertyTests = map[string]struct {
		SpecTest func(*testing.T, PropertyDriver)
	}{
		"StoreProperty":    {AddRental},
		"GetProperty":      {GetProperty},
		"ListProperties":   {ListProperties},
		"SearchProperties": {SearchProperties},
		"RemoveProperty":   {RemoveProperty},
	}
	for name, tc := range PropertyTests {
		t.Run(name, func(t *testing.T) {
			tc.SpecTest(t, driver)
		})
	}
}
func RunAllTenantTests(t *testing.T, driver TenantDriver) {
	var TenantTests = map[string]struct {
		SpecTest func(*testing.T, TenantDriver)
	}{
		"StoreTenant": {AddTenant},
		"GetTenant":   {GetTenant},
		"ListTenants": {ListTenants},
	}
	for name, tc := range TenantTests {
		t.Run(name, func(t *testing.T) {
			tc.SpecTest(t, driver)
		})
	}
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

	list, err := driver.ListProperties(ctx, usecase.AllProperties)
	require.NoError(t, err)
	m := entityMap(list...)
	require.Contains(t, m, in1.GetID())
	require.Contains(t, m, in2.GetID())
}

func SearchProperties(t *testing.T, driver PropertyDriver) {
	var (
		scope  = test.RandString(5)
		street = "Main st"
		city1  = "city1" + scope
		city2  = "City2" + scope
		state1 = "OH"
		zip1   = "10001"
		p1     = entity.NewProperty("100 "+street, city1, state1, zip1)
		p2     = entity.NewProperty("101 "+street, city2, state1, zip1)
	)
	for _, p := range []entity.Property{p1, p2} {
		_, err := driver.StoreProperty(ctx, p)
		require.NoError(t, err, "failed to store property: "+p.City)
	}
	f := usecase.NewPropertyFilter().WithSearch(city1)
	list, err := driver.ListProperties(ctx, f)
	require.NoError(t, err)
	assert.Len(t, list, 1, city1)
	assert.Equal(t, p1.ID, list[0].ID)
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

	list, err := driver.ListProperties(ctx, usecase.AllProperties)
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
