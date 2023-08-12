package usecase

import (
	"context"

	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/internal"
)

// StoreProperty is a use case to add or update a property
type StoreProperty struct {
	propRepo PropertyWriter
}

// NewStoreProperty constructs and returns an StoreProperty
func NewStoreProperty(repo PropertyWriter) StoreProperty {
	return StoreProperty{propRepo: repo}
}

// Execute the use case
func (uc StoreProperty) Execute(ctx context.Context, property entity.Property) error {
	if uc.propRepo == nil {
		return internal.NewErrors(internal.ErrInternal, ErrRepoNotSet)
	}
	if err := property.Validate(); err != nil {
		return err
	}
	if err := uc.propRepo.StoreProperty(ctx, property); err != nil {
		// TODO: make sure the error is logged here or in the repo layer
		return internal.NewErrors(internal.ErrInternal, ErrRepo)
	}
	return nil
}
