package entity

import (
	"github.com/google/uuid"
)

// Entity interface
type Entity interface {
	GetID() string
}

// ID for entities
type ID = string

// NewID returns a new ID
func NewID() string               { return uuid.NewString() }
func idEqualOrEmpty(a, b ID) bool { return a == "" || b == "" || a == b }
