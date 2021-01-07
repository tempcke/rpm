package repository_test

import (
	"testing"

	"github.com/tempcke/rpm/repository"
)

func TestInMemoryRepository(t *testing.T) {

	repo := repository.NewInMemoryRepo()

	runTestsWithRepo(t, repo)
}
