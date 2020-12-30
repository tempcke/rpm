package repository

import (
	"errors"

	"github.com/tempcke/rpm/entity"
)

type InMemory struct {
	properties map[string]entity.Property
}

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
		return p, errors.New("")
	}
	return p, nil
}
