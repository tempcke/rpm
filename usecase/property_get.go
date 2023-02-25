package usecase

import (
	"context"

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
	return uc.propRepo.RetrieveProperty(ctx, id)
}

// Validate checks the state of this use case, such as if it has a usable repo or not
func (uc GetProperty) Validate() error {
	if uc.propRepo == nil {
		return internal.NewErrors(internal.ErrInternal, ErrRepoNotSet)
	}
	return nil
}
