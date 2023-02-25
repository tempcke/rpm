package usecase

import (
	"context"

	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/internal"
)

// ListProperties UseCase
type ListProperties struct {
	propRepo PropertyReader
}

// NewListProperties constructs and returns an ListProperties
func NewListProperties(repo PropertyReader) ListProperties {
	return ListProperties{propRepo: repo}
}

// Execute ListProperties returns slice of properties
func (uc ListProperties) Execute(ctx context.Context) ([]entity.Property, error) {
	if err := uc.Validate(); err != nil {
		return nil, err
	}
	return uc.propRepo.PropertyList(ctx)
}

// Validate checks the state of this use case, such as if it has a usable repo or not
func (uc ListProperties) Validate() error {
	if uc.propRepo == nil {
		return internal.NewErrors(internal.ErrInternal, ErrRepoNotSet)
	}
	return nil
}
