package actions_test

import (
	"testing"

	"github.com/tempcke/rpm/actions"
	"github.com/tempcke/rpm/repository"
	"github.com/tempcke/rpm/specifications"
)

func TestActions_Property(t *testing.T) {
	var (
		repo   = repository.NewInMemoryRepo()
		driver = actions.NewActionsWithRepo(repo)
	)

	var tests = map[string]struct {
		specTest func(*testing.T, specifications.PropertyDriver)
	}{
		"StoreProperty":  {specifications.AddRental},
		"GetProperty":    {specifications.GetProperty},
		"ListProperties": {specifications.ListProperties},
		"RemoveProperty": {specifications.RemoveProperty},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.specTest(t, &driver)
		})
	}
}
func TestActions_Tenant(t *testing.T) {
	var (
		repo   = repository.NewInMemoryRepo()
		driver = actions.NewActionsWithRepo(repo)
	)
	var tests = map[string]struct {
		specTest func(*testing.T, specifications.TenantDriver)
	}{
		"StoreTenant": {specifications.AddTenant},
		"GetTenant":   {specifications.GetTenant},
		"ListTenants": {specifications.ListTenants},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.specTest(t, &driver)
		})
	}
}
