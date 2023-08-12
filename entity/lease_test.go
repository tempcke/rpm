package entity_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/entity/fake"
	"github.com/tempcke/schedule"
)

func TestLease(t *testing.T) {
	var (
		property  = fake.Property()
		tenant1   = fake.Tenant()
		tenant2   = fake.Tenant()
		tenant3   = fake.Tenant()
		rent      = rand.Intn(1000) + 1000
		deposit   = rent*3 + rand.Intn(1000)
		nextMonth = schedule.Today().AddDate(0, 1, 0)
		start     = schedule.NewDate(nextMonth.Year(), nextMonth.Month(), 1)
		end       = start.AddDate(0, 12, -1)
	)
	lease := entity.NewLease(property)
	require.NotEmpty(t, lease.ID)
	require.True(t, property.Equal(lease.Property))

	lease = lease.WithTenant(tenant1, tenant2).WithTenant(tenant3)
	require.Equal(t, 3, len(lease.Tenants))
	require.True(t, lease.HasTenant(tenant1.ID))
	require.True(t, lease.HasTenant(tenant2.ID))
	require.True(t, lease.HasTenant(tenant3.ID))

	// don't add duplicate tenants
	lease = lease.WithTenant(tenant2)
	require.Equal(t, 3, len(lease.Tenants))

	lease = lease.WithRent(rent).WithDeposit(deposit)
	assert.Equal(t, rent, lease.RentAmount)
	assert.Equal(t, deposit, lease.Deposit)

	lease = lease.WithTerm(start, end)
	assert.Equal(t, start, lease.StartDate)
	assert.Equal(t, end, lease.EndDate)
}
