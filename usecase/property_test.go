package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/repository"
	"github.com/tempcke/rpm/usecase"
)

func TestAddProperty(t *testing.T) {
	r := repository.NewInMemoryRepo()
	uc := usecase.NewAddPropertyCommand(r)
	t.Run("sunny day", func(t *testing.T) {
		p := r.NewProperty("1234 N Main st.", "Dallas", "TX", "75401")

		e1 := uc.Execute(p)
		_, e2 := r.RetrieveProperty(p.ID)

		assert.Nil(t, e1)
		assert.Nil(t, e2)
	})

	t.Run("invalid property can not be saved", func(t *testing.T) {
		p := r.NewProperty("", "a", "b", "c")
		assert.Error(t, p.Validate())
		e1 := uc.Execute(p)
		_, e2 := r.RetrieveProperty(p.ID)
		assert.Error(t, e1)
		assert.Error(t, e2)
	})
}

func TestGetProperty(t *testing.T) {
	r := repository.NewInMemoryRepo()
	pIn := newPropertyFixture(r)
	err1 := usecase.NewAddPropertyCommand(r).Execute(pIn)

	c := usecase.NewGetPropertyQuery(r)
	pOut, err2 := c.Execute(pIn.ID)
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Equal(t, pIn.ID, pOut.ID)
	assert.Equal(t, pIn.Street, pOut.Street)
	assert.Equal(t, pIn.City, pOut.City)
	assert.Equal(t, pIn.StateCode, pOut.StateCode)
	assert.Equal(t, pIn.Zip, pOut.Zip)
}

func newPropertyFixture(r repository.InMemory) entity.Property {
	return r.NewProperty("1234 N Main st.", "Dallas", "TX", "75401")
}
