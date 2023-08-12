package entity_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/entity/fake"
	"github.com/tempcke/schedule"
)

func TestTenant(t *testing.T) {
	var (
		name = fake.FullName()
		dob  = fake.DateOfBirth()
	)
	tnt := entity.NewTenant(name, dob)
	require.Equal(t, name, tnt.FullName)
	require.Equal(t, dob, tnt.DateOfBirth)
	require.True(t, tnt.Equal(tnt))
	require.False(t, tnt.Equal(entity.NewTenant("bob", dob)))
	require.False(t, tnt.Equal(entity.NewTenant("bob", schedule.Today().AddDate(-20, 0, 0))))
	require.False(t, tnt.Equal(tnt.WithID(entity.NewID())))
	require.True(t, tnt.Equal(tnt.WithID("")))
	require.Equal(t, tnt.ID, tnt.GetID())
}
