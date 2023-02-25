package actions_test

import (
	"testing"

	"github.com/tempcke/rpm/actions"
	"github.com/tempcke/rpm/repository"
	"github.com/tempcke/rpm/specifications"
)

func TestActions(t *testing.T) {
	var (
		repo   = repository.NewInMemoryRepo()
		driver = actions.NewActions(repo)
	)
	specifications.AddRental(t, driver)
	specifications.ListProperties(t, driver)
	specifications.RemoveProperty(t, driver)
}
