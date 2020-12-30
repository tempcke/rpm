package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	t.Run("is unique", func(t *testing.T) {
		a := NewID()
		b := NewID()
		assert.NotEqual(t, a, b)
	})
}
