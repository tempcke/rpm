package repository

import (
	"context"

	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/internal"
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

// StoreProperty persists a property
func (r InMemory) StoreProperty(_ context.Context, property entity.Property) error {
	r.properties[property.GetID()] = property
	return nil
}

// NewProperty constructs a property
func (r InMemory) NewProperty(street, city, state, zip string) entity.Property {
	return entity.NewProperty(street, city, state, zip)
}

// RetrieveProperty by id
func (r InMemory) RetrieveProperty(_ context.Context, id string) (entity.Property, error) {
	p, ok := r.properties[id]
	if !ok {
		return p, internal.MakeErr(internal.ErrEntityNotFound, "property["+id+"]")
	}
	return p, nil
}

// PropertyList is used to list properties
func (r InMemory) PropertyList(_ context.Context) ([]entity.Property, error) {
	pl := make([]entity.Property, len(r.properties))
	i := 0
	for _, p := range r.properties {
		pl[i] = p
		i++
	}
	return pl, nil
}

// DeleteProperty by id
func (r InMemory) DeleteProperty(_ context.Context, id string) error {
	delete(r.properties, id)
	return nil
}
