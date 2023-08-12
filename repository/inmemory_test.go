package repository_test

import (
	"testing"

	"github.com/tempcke/rpm/repository"
	"github.com/tempcke/rpm/usecase"
)

var _ interface {
	usecase.PropertyRepo
	usecase.TenantRepo
} = (*repository.InMemory)(nil)

func TestPropertyRepo_InMemory(t *testing.T) {
	var propertyTests = map[string]struct {
		fn func(*testing.T, propertyRepo)
	}{
		"store":  {testStoreAndRetrieveProperty},
		"update": {testUpdateProperty},
		"list":   {testListProperties},
		"remove": {testRemoveProperty},
		"get":    {testGetProperty},
	}

	r := repository.NewInMemoryRepo()

	for name, tc := range propertyTests {
		t.Run(name, func(t *testing.T) {
			tc.fn(t, r)
		})
	}
}
