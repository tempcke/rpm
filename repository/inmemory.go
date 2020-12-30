package repository

import (
	"errors"

	"github.com/tempcke/rpm/entity"
)

// InMemory repository should NOT be used in production
type InMemory struct {
	properties map[string]entity.Property
}

// NewInMemoryRepo constructs an InMemory repository
func NewInMemoryRepo() InMemory {
	return InMemory{
		properties: make(map[string]entity.Property),
	}
}

func (r InMemory) StoreProperty(property entity.Property) error {
	r.properties[property.GetID()] = property
	return nil
}

func (r InMemory) NewProperty(street, city, state, zip string) entity.Property {
	return entity.NewProperty(street, city, state, zip)
}

func (r InMemory) RetrieveProperty(id string) (entity.Property, error) {
	p, ok := r.properties[id]
	if !ok {
		return p, errors.New("property not found")
	}
	return p, nil
}

func (r InMemory) DeleteProperty(id string) error {
	delete(r.properties, id)
	return nil
}
