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
	type (
		Driver   = specifications.Driver
		specTest func(*testing.T, Driver)
	)
	var tests = map[string]struct{ specTest }{
		"StoreProperty":  {specifications.AddRental},
		"GetProperty":    {specifications.GetProperty},
		"ListProperties": {specifications.ListProperties},
		"RemoveProperty": {specifications.RemoveProperty},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.specTest(t, &driver)
		})
	}
}
