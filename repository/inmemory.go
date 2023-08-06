package repository

import (
	"context"
	"sync"
	"time"

	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/internal"
)

var rwMutex sync.RWMutex

// InMemory repository should NOT be used in production
type InMemory struct {
	properties  map[string]entity.Property
	propertyErr map[string]error
}

// NewInMemoryRepo constructs an InMemory repository
func NewInMemoryRepo() InMemory {
	return InMemory{
		properties:  make(map[string]entity.Property),
		propertyErr: make(map[string]error),
	}
}
func (r InMemory) WithPropertyErr(id string, err error) InMemory {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	r.propertyErr[id] = err
	return r
}

func (r InMemory) StoreProperty(_ context.Context, property entity.Property) error {
	if property.CreatedAt.IsZero() {
		property.CreatedAt = time.Now()
	}
	rwMutex.Lock()
	defer rwMutex.Unlock()
	if err := r.propertyErr[property.GetID()]; err != nil {
		return err
	}
	r.properties[property.GetID()] = property
	return nil
}
func (r InMemory) NewProperty(street, city, state, zip string) entity.Property {
	return entity.NewProperty(street, city, state, zip)
}
func (r InMemory) RetrieveProperty(_ context.Context, id string) (entity.Property, error) {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	p, ok := r.properties[id]
	if !ok {
		return p, internal.MakeErr(internal.ErrEntityNotFound, "property["+id+"]")
	}
	if err := r.propertyErr[id]; err != nil {
		return p, err
	}
	return p, nil
}
func (r InMemory) PropertyList(_ context.Context) ([]entity.Property, error) {
	pl := make([]entity.Property, len(r.properties))
	i := 0
	for _, p := range r.properties {
		pl[i] = p
		i++
	}
	return pl, nil
}
func (r InMemory) DeleteProperty(_ context.Context, id string) error {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	if err := r.propertyErr[id]; err != nil {
		return err
	}
	delete(r.properties, id)
	return nil
}
