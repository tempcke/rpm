package usecase

import (
	"context"
	"errors"

	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/internal"
)

var noProperty = entity.Property{}

// GetProperty UseCase
type GetProperty struct {
	propRepo PropertyReader
}

// NewGetProperty constructs and returns an GetProperty
func NewGetProperty(repo PropertyReader) GetProperty {
	return GetProperty{propRepo: repo}
}

// Execute GetProperty returns a property by id
func (uc GetProperty) Execute(ctx context.Context, id string) (entity.Property, error) {
	if err := uc.Validate(); err != nil {
		return noProperty, err
	}
	p, err := uc.propRepo.RetrieveProperty(ctx, id)
	if err != nil {
		if errors.Is(err, internal.ErrEntityNotFound) {
			return p, err
		}
		// TODO: make sure the error is logged here or in the repo layer
		return noProperty, internal.NewErrors(internal.ErrInternal, ErrRepo)
	}
	return p, nil
}

// Validate checks the state of this use case, such as if it has a usable repo or not
func (uc GetProperty) Validate() error {
	if uc.propRepo == nil {
		return internal.NewErrors(internal.ErrInternal, ErrRepoNotSet)
	}
	return nil
}
