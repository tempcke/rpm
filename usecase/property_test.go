package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/repository"
	"github.com/tempcke/rpm/usecase"
)

func TestAddProperty(t *testing.T) {
	repo := repository.NewInMemoryRepo()
	uc := usecase.NewAddProperty(repo)
	t.Run("sunny day", func(t *testing.T) {
		p := repo.NewProperty("1234 N Main st.", "Dallas", "TX", "75401")

		e1 := uc.Execute(p)
		_, e2 := repo.RetrieveProperty(p.ID)

		assert.Nil(t, e1)
		assert.Nil(t, e2)
	})

	t.Run("invalid property can not be saved", func(t *testing.T) {
		p := repo.NewProperty("", "a", "b", "c")
		assert.Error(t, p.Validate())
		e1 := uc.Execute(p)
		_, e2 := repo.RetrieveProperty(p.ID)
		assert.Error(t, e1)
		assert.Error(t, e2)
	})
}

func TestListProperties(t *testing.T) {
	r := repository.NewInMemoryRepo()

	t.Run("empty set", func(t *testing.T) {
		propList, err := usecase.NewListProperties(r).Execute()
		assert.NoError(t, err)
		assert.Len(t, propList, 0)
	})

	t.Run("two properties", func(t *testing.T) {
		p1 := r.NewProperty("100 N Main st.", "Dallas", "TX", "75401")
		p2 := r.NewProperty("100 N Main st.", "Dallas", "TX", "75401")
		r.StoreProperty(p1)
		r.StoreProperty(p2)

		propList, err := usecase.NewListProperties(r).Execute()
		assert.NoError(t, err)
		assert.Len(t, propList, 2)
	})
}

func TestGetProperty(t *testing.T) {
	repo := repository.NewInMemoryRepo()

	t.Run("get newly added property", func(t *testing.T) {
		pIn := newPropertyFixture(repo)
		err1 := usecase.NewAddProperty(repo).Execute(pIn)

		c := usecase.NewGetProperty(repo)
		pOut, err2 := c.Execute(pIn.ID)
		assert.Nil(t, err1)
		assert.Nil(t, err2)
		assert.Equal(t, pIn.ID, pOut.ID)
		assert.Equal(t, pIn.Street, pOut.Street)
		assert.Equal(t, pIn.City, pOut.City)
		assert.Equal(t, pIn.StateCode, pOut.StateCode)
		assert.Equal(t, pIn.Zip, pOut.Zip)
	})

	t.Run("unknown property, expect error", func(t *testing.T) {
		_, err := usecase.NewGetProperty(repo).Execute("doesNotExist")
		assert.Error(t, err)
	})
}

func TestDelProperty(t *testing.T) {
	repo := repository.NewInMemoryRepo()

	t.Run("delete newly added property", func(t *testing.T) {
		// add property
		p := newPropertyFixture(repo)
		err := usecase.NewAddProperty(repo).Execute(p)
		assert.Nil(t, err)

		// delete property
		err = usecase.NewDeleteProperty(repo).Execute(p.ID)
		assert.Nil(t, err)

		// get property should fail
		_, err = usecase.NewGetProperty(repo).Execute(p.ID)
		assert.Error(t, err)
	})

	t.Run("Delete property that does not exist should not error, you want it gone, and it isn't there?", func(t *testing.T) {
		id := "doesNotExist"
		err := usecase.NewDeleteProperty(repo).Execute(id)
		assert.Nil(t, err)
	})
}

func newPropertyFixture(r usecase.PropertyRepository) entity.Property {
	return r.NewProperty("1234 N Main st.", "Dallas", "TX", "75401")
}
