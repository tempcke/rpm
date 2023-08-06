package usecase

import (
	"context"

	"github.com/tempcke/rpm/internal"
)

// DeleteProperty Use Case
type DeleteProperty struct {
	propRepo PropertyWriter
}

// NewDeleteProperty constructs a DeleteProperty use case
func NewDeleteProperty(repo PropertyWriter) DeleteProperty {
	return DeleteProperty{propRepo: repo}
}

// Execute the DeleteProperty use case to delete a property by ID
func (uc DeleteProperty) Execute(ctx context.Context, id string) error {
	if uc.propRepo == nil {
		return internal.NewErrors(internal.ErrInternal, ErrRepoNotSet)
	}
	if err := uc.propRepo.DeleteProperty(ctx, id); err != nil {
		// TODO: make sure the error is logged here or in the repo layer
		return internal.NewErrors(internal.ErrInternal, ErrRepo)
	}
	return nil
}
