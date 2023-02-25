package repository_test

import (
	"testing"

	"github.com/tempcke/rpm/repository"
	"github.com/tempcke/rpm/usecase"
)

var _ usecase.PropertyRepository = (*repository.InMemory)(nil)

func TestInMemoryRepository(t *testing.T) {
	r := repository.NewInMemoryRepo()

	for name, tc := range repoTests {
		t.Run(name, func(t *testing.T) {
			tc.fn(t, r)
		})
	}
}
