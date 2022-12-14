package repository_test

import (
	"testing"

	"github.com/tempcke/rpm/repository"
)

func TestInMemoryRepository(t *testing.T) {
	r := repository.NewInMemoryRepo()

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
