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
	if err := uc.Validate(); err != nil {
		return err
	}
	if err := property.Validate(); err != nil {
		return err
	}
	return uc.propRepo.StoreProperty(ctx, property)
}

// Validate checks the state of this use case, such as if it has a usable repo or not
func (uc StoreProperty) Validate() error {
	if uc.propRepo == nil {
		return internal.NewErrors(internal.ErrInternal, ErrRepoNotSet)
	}
	return nil
}
