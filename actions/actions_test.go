package actions_test

import (
	"testing"

	"github.com/tempcke/rpm/actions"
	"github.com/tempcke/rpm/internal/repository"
	"github.com/tempcke/rpm/specifications"
)

func TestActions_specifications(t *testing.T) {
	var (
		repo   = repository.NewInMemoryRepo()
		driver = actions.NewActionsWithRepo(repo)
	)
	specifications.RunAllTests(t, driver, driver)
}
