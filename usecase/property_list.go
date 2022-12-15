package usecase

import (
	"context"

	"github.com/tempcke/rpm/entity"
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
	return uc.propRepo.PropertyList(ctx)
}
