//go:build withDocker
// +build withDocker

package repository_test

import (
	"testing"

	"github.com/tempcke/rpm/internal/test"
	"github.com/tempcke/rpm/repository"
	"github.com/tempcke/rpm/usecase"
)

var _ usecase.PropertyRepository = (*repository.Postgres)(nil)

func TestPropertyRepository(t *testing.T) {
	r := repository.NewPostgresRepo(test.DB(t))

	tests := map[string]struct {
		fn func(*testing.T, Repo)
	}{
		"store":  {testStoreAndRetrieveProperty},
		"update": {testUpdateProperty},
		"list":   {testListProperties},
		"remove": {testRemoveProperty},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.fn(t, r)
		})
	}
}
