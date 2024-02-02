package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/entity/fake"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/internal/repository"
	"github.com/tempcke/rpm/usecase"
)

var ctx = context.Background()

func TestAddProperty(t *testing.T) {
	repo := repository.NewInMemoryRepo()
	uc := usecase.NewPropertyManager(repo)

	t.Run("sunny day", func(t *testing.T) {
		p := repo.NewProperty("1234 N Main st.", "Dallas", "TX", "75401")

		require.NoError(t, uc.Store(ctx, p))
		_, err := repo.GetProperty(ctx, p.ID)

		require.NoError(t, err)
	})
	t.Run("invalid property can not be saved", func(t *testing.T) {
		p := repo.NewProperty("", "a", "b", "c")
		require.Error(t, p.Validate())
		err := uc.Store(ctx, p)
		require.Error(t, err)
		assert.ErrorIs(t, err, internal.ErrEntityInvalid)
		_, err = repo.GetProperty(ctx, p.ID)
		require.Error(t, err)
	})
	t.Run("no repo", func(t *testing.T) {
		uc := usecase.NewPropertyManager(nil)
		p := fake.Property()
		err := uc.Store(ctx, p)
		require.Error(t, err)
		assert.ErrorIs(t, err, internal.ErrInternal)
		assert.ErrorIs(t, err, usecase.ErrRepoNotSet)
		_, err = repo.GetProperty(ctx, p.ID)
		require.Error(t, err)
	})
	t.Run("any repo error should be internal", func(t *testing.T) {
		var (
			p       = entity.NewProperty("1234 N Main st.", "Dallas", "TX", "75401")
			repoErr = fmt.Errorf("any database error %s", uuid.NewString())
			repo    = repository.NewInMemoryRepo().WithEntityErr(p.ID, repoErr)
			uc      = usecase.NewPropertyManager(repo)
			err     error
		)
		_, err = repo.GetProperty(ctx, p.ID)
		require.Error(t, err)

		err = uc.Store(ctx, p)
		require.Error(t, err)
		require.ErrorIs(t, err, internal.ErrInternal)
		require.ErrorIs(t, err, usecase.ErrRepo)
	})
}
func TestListProperties(t *testing.T) {
	r := repository.NewInMemoryRepo()

	t.Run("empty set", func(t *testing.T) {
		propList, err := usecase.NewPropertyManager(r).List(ctx)
		assert.NoError(t, err)
		assert.Len(t, propList, 0)
	})
	t.Run("two properties", func(t *testing.T) {
		p1 := r.NewProperty("100 N Main st.", "Dallas", "TX", "75401")
		p2 := r.NewProperty("100 N Main st.", "Dallas", "TX", "75401")
		require.NoError(t, r.StoreProperty(ctx, p1))
		require.NoError(t, r.StoreProperty(ctx, p2))

		propList, err := usecase.NewPropertyManager(r).List(ctx)
		assert.NoError(t, err)
		assert.Len(t, propList, 2)
	})
	t.Run("no repo", func(t *testing.T) {
		_, err := usecase.NewPropertyManager(nil).List(ctx)
		require.Error(t, err)
		assert.ErrorIs(t, err, internal.ErrInternal)
		assert.ErrorIs(t, err, usecase.ErrRepoNotSet)
	})
}
func TestGetProperty(t *testing.T) {
	repo := repository.NewInMemoryRepo()

	t.Run("get newly added property", func(t *testing.T) {
		pIn := newPropertyFixture(repo)
		err1 := usecase.NewPropertyManager(repo).Store(ctx, pIn)

		c := usecase.NewPropertyManager(repo)
		pOut, err2 := c.Get(ctx, pIn.ID)
		assert.Nil(t, err1)
		assert.Nil(t, err2)
		assert.Equal(t, pIn.ID, pOut.ID)
		assert.Equal(t, pIn.Street, pOut.Street)
		assert.Equal(t, pIn.City, pOut.City)
		assert.Equal(t, pIn.StateCode, pOut.StateCode)
		assert.Equal(t, pIn.Zip, pOut.Zip)
	})
	t.Run("unknown property, expect error", func(t *testing.T) {
		p, err := usecase.NewPropertyManager(repo).Get(ctx, "doesNotExist")
		require.Error(t, err)
		assert.ErrorIs(t, err, internal.ErrEntityNotFound)
		assert.Empty(t, p)
	})
	t.Run("no repo", func(t *testing.T) {
		p, err := usecase.NewPropertyManager(nil).Get(ctx, "anything")
		require.Error(t, err)
		assert.ErrorIs(t, err, internal.ErrInternal)
		assert.ErrorIs(t, err, usecase.ErrRepoNotSet)
		assert.Empty(t, p)
	})
}
func TestDelProperty(t *testing.T) {
	repo := repository.NewInMemoryRepo()

	t.Run("delete newly added property", func(t *testing.T) {
		// add property
		p := newPropertyFixture(repo)
		err := usecase.NewPropertyManager(repo).Store(ctx, p)
		assert.Nil(t, err)

		// delete property
		err = usecase.NewPropertyManager(repo).Remove(ctx, p.ID)
		assert.Nil(t, err)

		// get property should fail
		_, err = usecase.NewPropertyManager(repo).Get(ctx, p.ID)
		assert.Error(t, err)
	})
	t.Run("Delete property that does not exist should not error, you want it gone, and it isn't there?", func(t *testing.T) {
		id := "doesNotExist"
		err := usecase.NewPropertyManager(repo).Remove(ctx, id)
		assert.Nil(t, err)
	})
}

func newPropertyFixture(r usecase.PropertyRepo) entity.Property {
	return r.NewProperty("1234 N Main st.", "Dallas", "TX", "75401")
}
