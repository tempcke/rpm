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

	for name, tc := range repoTests {
		t.Run(name, func(t *testing.T) {
			tc.fn(t, r)
		})
	}
}
