package usecase

import (
	"context"

	"github.com/tempcke/rpm/entity"
)

// PropertyReader allows queries regarding properties
type PropertyReader interface {
	GetProperty(ctx context.Context, id string) (entity.Property, error)
	PropertyList(ctx context.Context) ([]entity.Property, error)
}

// PropertyWriter allows for mutations within the property repository
type PropertyWriter interface {
	NewProperty(street, city, state, zip string) entity.Property
	StoreProperty(context.Context, entity.Property) error
	DeleteProperty(ctx context.Context, id string) error
}

// PropertyRepo contains both the PropertyReader and PropertyWriter
type PropertyRepo interface {
	PropertyReader
	PropertyWriter
}
