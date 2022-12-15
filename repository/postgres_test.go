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

	t.Run("store and retrieve property", func(t *testing.T) {
		testStoreAndRetrieveProperty(t, r)
	})
	t.Run("list properties", func(t *testing.T) {
		testListProperties(t, r)
	})
	t.Run("remove property", func(t *testing.T) {
		testRemoveProperty(t, r)
	})
}
