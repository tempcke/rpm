package repository_test

import (
	"testing"

	"github.com/tempcke/rpm/repository"
)

func TestPropertyRepo_InMemory(t *testing.T) {
	var tests = map[string]struct {
		fn func(*testing.T, propertyRepo)
	}{
		"store":  {testStoreAndRetrieveProperty},
		"update": {testUpdateProperty},
		"list":   {testListProperties},
		"remove": {testRemoveProperty},
		"get":    {testGetProperty},
	}

	r := repository.NewInMemoryRepo()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.fn(t, r)
		})
	}
}
func TestTenantRepo_InMemory(t *testing.T) {
	var tests = map[string]struct{ fn func(*testing.T, tenantRepo) }{
		"store get list": {testTenant},
	}

	r := repository.NewInMemoryRepo()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.fn(t, r)
		})
	}
}
