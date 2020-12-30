package usecase

import "github.com/tempcke/rpm/entity"

// PropertyReader allows queries regarding properties
type PropertyReader interface {
	RetrieveProperty(id string) (entity.Property, error)
}

// PropertyWriter allows for mutations within the property repository
type PropertyWriter interface {
	NewProperty(street, city, state, zip string) entity.Property
	StoreProperty(entity.Property) error
	DeleteProperty(id string) error
}

// PropertyRepository contains both the PropertyReader and PropertyWriter
type PropertyRepository interface {
	PropertyReader
	PropertyWriter
}
