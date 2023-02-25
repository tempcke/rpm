package entity

import (
	"time"

	"github.com/tempcke/rpm/internal"
)

// Property entity
type Property struct {
	ID        string
	Street    string // 1234 N Main st
	City      string // Dallas
	StateCode string // TX
	Zip       string // 75401
	CreatedAt time.Time
}

// NewProperty returns a new entity.Property
func NewProperty(street, city, state, zip string) Property {
	return Property{
		ID:        NewID(),
		Street:    street,
		City:      city,
		StateCode: state,
		Zip:       zip,
		CreatedAt: time.Now(),
	}
}

// GetID of entity
// method needed to implement entity.Entity
func (p Property) GetID() string {
	return p.ID
}

// Validate is used to validate the entity
func (p Property) Validate() error {
	if p.ID == "" || p.Street == "" || p.City == "" || p.StateCode == "" || p.Zip == "" {
		return internal.ErrEntityInvalid
	}
	return nil
}

func (p Property) WithID(id string) Property {
	p.ID = id
	return p
}
