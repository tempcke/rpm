package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/entity/fake"
	"github.com/tempcke/rpm/internal"
)

func testTenant(t *testing.T, driver tenantRepo) {
	// store 1,2
	in1 := fake.Tenant()
	in2 := fake.Tenant()
	err := driver.StoreTenant(ctx, in1)
	require.NoError(t, err)
	err = driver.StoreTenant(ctx, in2)
	require.NoError(t, err)

	// get 1
	out1, err := driver.GetTenant(ctx, in1.GetID())
	require.NoError(t, err)
	require.NotNil(t, out1)
	require.True(t, out1.Equal(in1))

	// get 2
	out2, err := driver.GetTenant(ctx, in1.GetID())
	require.NoError(t, err)
	require.NotNil(t, out2)
	require.True(t, out2.Equal(in1))

	// get not exists
	out3, err := driver.GetTenant(ctx, entity.NewID())
	require.ErrorIs(t, err, internal.ErrEntityNotFound)
	require.Nil(t, out3)

	// list
	list, err := driver.ListTenants(ctx)
	require.NoError(t, err)
	assertEntityInSet(t, in1.GetID(), list...)
	assertEntityInSet(t, in2.GetID(), list...)

	// update
	in1b := in1.WithName(fake.FullName())
	assert.Equal(t, in1.GetID(), in1b.GetID())
	assert.NotEqual(t, in1.FullName, in1b.FullName)
	require.NoError(t, driver.StoreTenant(ctx, in1b))
	out1b, err := driver.GetTenant(ctx, in1b.ID)
	require.NoError(t, err)
	require.NotNil(t, out1b)
	require.True(t, out1b.Equal(in1b))
}
