package repository

import (
	"context"
	"sync"
	"time"

	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/internal/filters"
)

var rwMutex sync.RWMutex

// InMemory repository should NOT be used in production
type InMemory struct {
	entities   map[string]entity.Entity
	entityErrs map[string]error
}

// NewInMemoryRepo constructs an InMemory repository
func NewInMemoryRepo() InMemory {
	return InMemory{
		entities:   make(map[string]entity.Entity),
		entityErrs: make(map[string]error),
	}
}
func (r InMemory) WithEntityErr(id string, err error) InMemory {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	r.entityErrs[id] = err
	return r
}

func (r InMemory) StoreProperty(_ context.Context, property entity.Property) error {
	if property.CreatedAt.IsZero() {
		property.CreatedAt = time.Now()
	}
	return r.storeEntity(property)
}
func (r InMemory) NewProperty(street, city, state, zip string) entity.Property {
	return entity.NewProperty(street, city, state, zip)
}
func (r InMemory) GetProperty(_ context.Context, id string) (entity.Property, error) {
	e, err := r.getEntity(id)
	if err != nil {
		return entity.Property{}, err
	}
	return e.(entity.Property), nil
}
func (r InMemory) PropertyList(_ context.Context) ([]entity.Property, error) {
	for _, err := range r.entityErrs {
		return nil, err
	}
	list := make([]entity.Property, 0, len(r.entities))
	for _, e := range r.entities {
		if p, ok := e.(entity.Property); ok {
			if _, err := r.getEntity(e.GetID()); err != nil {
				return nil, err
			}
			list = append(list, p)
		}
	}
	return list, nil
}
func (r InMemory) DeleteProperty(_ context.Context, id string) error { return r.delEntity(id) }

func (r InMemory) StoreTenant(_ context.Context, e entity.Tenant) error { return r.storeEntity(e) }
func (r InMemory) GetTenant(_ context.Context, id entity.ID) (*entity.Tenant, error) {
	e, err := r.getEntity(id)
	if err != nil {
		return nil, err
	}
	t := e.(entity.Tenant) // only used in tests, we want it to panic if it is wrong
	return &t, nil
}
func (r InMemory) ListTenants(_ context.Context, _ ...filters.TenantFilter) ([]entity.Tenant, error) {
	for _, err := range r.entityErrs {
		return nil, err
	}
	list := make([]entity.Tenant, 0, len(r.entities))
	for _, e := range r.entities {
		if item, ok := e.(entity.Tenant); ok {
			if _, err := r.getEntity(e.GetID()); err != nil {
				return nil, err
			}
			list = append(list, item)
		}
	}
	return list, nil
}

func (r InMemory) storeEntity(e entity.Entity) error {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	if err := r.entityErrs[e.GetID()]; err != nil {
		return err
	}
	r.entities[e.GetID()] = e
	return nil
}
func (r InMemory) getEntity(id entity.ID) (entity.Entity, error) {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	if err := r.entityErrs[id]; err != nil {
		return nil, err
	}
	e, ok := r.entities[id]
	if !ok {
		return nil, internal.MakeErr(internal.ErrEntityNotFound, id)
	}
	return e, nil
}
func (r InMemory) delEntity(id entity.ID) error {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	if err := r.entityErrs[id]; err != nil {
		return err
	}
	delete(r.entities, id)
	return nil
}
