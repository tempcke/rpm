package usecase

import "github.com/tempcke/rpm/entity"

type PropertyReader interface {
	RetrieveProperty(id string) (entity.Property, error)
}

type PropertyWriter interface {
	NewProperty(street, city, state, zip string) entity.Property
	StoreProperty(entity.Property) error
	DeleteProperty(id string) error
}

type PropertyRepository interface {
	PropertyReader
	PropertyWriter
}
