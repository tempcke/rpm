package entity

import (
	"github.com/google/uuid"
)

// Entity interface
type Entity interface {
	GetID() string
}

// ID for entities
type ID uuid.UUID

// NewID returns a new ID
func NewID() string {
	return uuid.New().String()
}
