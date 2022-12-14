package entity

import (
	"errors"

	"github.com/google/uuid"
)

var (
	// ErrInvalidEntity is an invalid entity error
	ErrInvalidEntity = errors.New("invalid entity")
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
