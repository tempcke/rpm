//go:build withDocker
// +build withDocker

package repository_test

import (
	"testing"

	"github.com/tempcke/rpm/internal/test"
	"github.com/tempcke/rpm/repository"
)

func TestPropertyRepo_Postgres(t *testing.T) {
	var tests = map[string]struct {
		fn func(*testing.T, propertyRepo)
	}{
		"store":  {testStoreAndRetrieveProperty},
		"update": {testUpdateProperty},
		"list":   {testListProperties},
		"remove": {testRemoveProperty},
		"get":    {testGetProperty},
	}

	r := repository.NewPostgresRepo(test.DB(t))
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.fn(t, r)
		})
	}
}

func TestTenantRepo_Postgres(t *testing.T) {
	var tests = map[string]struct{ fn func(*testing.T, tenantRepo) }{
		"store get list": {testTenant},
	}

	r := repository.NewPostgresRepo(test.DB(t))
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.fn(t, r)
		})
	}
}
