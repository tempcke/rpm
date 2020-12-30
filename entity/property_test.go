package entity_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tempcke/rpm/entity"
)

func TestNewProperty(t *testing.T) {
	p := entity.NewProperty(
		"1234 N Main st.",
		"Dallas",
		"TX",
		"75401")

	t.Run("NewProperty", func(t *testing.T) {
		assert.Equal(t, "1234 N Main st.", p.Street)
		assert.Equal(t, "Dallas", p.City)
		assert.Equal(t, "TX", p.StateCode)
		assert.Equal(t, "75401", p.Zip)
		assert.NotEmpty(t, p.ID)
		assert.Equal(t, p.ID, p.GetID())
		assert.Equal(t, time.Now().Year(), p.CreatedAt.Year())
		assert.Nil(t, p.Validate())
	})
}

func TestPropertyValidation(t *testing.T) {
	e := ""  // empty value
	v := "v" // non-empty value
	tt := []struct{ street, city, state, zip string }{
		{e, v, v, v},
		{v, e, v, v},
		{v, v, e, v},
		{v, v, v, e},
	}
	for _, tc := range tt {
		p := entity.NewProperty(tc.street, tc.city, tc.state, tc.zip)
		assert.Error(t, p.Validate())
	}
}
