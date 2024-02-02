package usecase_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/entity/fake"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/internal/filters"
	"github.com/tempcke/rpm/internal/repository"
	"github.com/tempcke/rpm/usecase"
)

func TestTenantUC(t *testing.T) {
	var (
		in1  = fake.Tenant()
		in2  = fake.Tenant()
		repo = repository.NewInMemoryRepo()
		uc   = usecase.NewTenantManager(&repo)

		// force repo to implement interface
		_ usecase.TenantRepo = (*repository.InMemory)(nil)
	)

	// store tenant
	out, err := uc.Store(ctx, in1)
	require.NoError(t, err)
	require.NotNil(t, out)

	// get tenant
	out, err = uc.Get(ctx, in1.GetID())
	require.NoError(t, err)
	require.NotNil(t, out)
	require.True(t, in1.Equal(*out))

	_, err = uc.Store(ctx, in2)
	require.NoError(t, err)

	filter := filters.NewTenantFilter()
	tenants, err := uc.List(ctx, filter)
	require.NoError(t, err)
	require.Equal(t, 2, len(tenants))
}

func TestTenantUC_fail(t *testing.T) {
	t.Run("uc without a repo", func(t *testing.T) {
		var (
			in   = fake.Tenant()
			repo usecase.TenantRepo
			uc   = usecase.NewTenantManager(repo)
		)

		out, err := uc.Store(ctx, in)
		require.Nil(t, out)
		require.ErrorIs(t, err, internal.ErrInternal)
		require.ErrorIs(t, err, usecase.ErrRepoNotSet)

		out, err = uc.Get(ctx, in.ID)
		require.Nil(t, out)
		require.ErrorIs(t, err, internal.ErrInternal)
		require.ErrorIs(t, err, usecase.ErrRepoNotSet)

		_, err = uc.List(ctx)
		require.ErrorIs(t, err, internal.ErrInternal)
		require.ErrorIs(t, err, usecase.ErrRepoNotSet)
	})
	t.Run("repo error", func(t *testing.T) {
		var (
			in      = fake.Tenant()
			repoErr = errors.New(t.Name() + "_" + uuid.NewString())
			repo    = repository.NewInMemoryRepo().WithEntityErr(in.ID, repoErr)
			uc      = usecase.NewTenantManager(&repo)
		)

		out, err := uc.Store(ctx, in)
		require.Nil(t, out)
		require.ErrorIs(t, err, internal.ErrInternal)
		require.ErrorIs(t, err, usecase.ErrRepo)

		out, err = uc.Get(ctx, in.GetID())
		require.Nil(t, out)
		require.ErrorIs(t, err, internal.ErrInternal)
		require.ErrorIs(t, err, usecase.ErrRepo)

		_, err = uc.List(ctx)
		require.ErrorIs(t, err, internal.ErrInternal)
		require.ErrorIs(t, err, usecase.ErrRepo)
	})
	t.Run("entity not found", func(t *testing.T) {
		var (
			in   = fake.Tenant()
			repo = repository.NewInMemoryRepo()
			uc   = usecase.NewTenantManager(repo)
		)

		out, err := uc.Get(ctx, in.GetID())
		require.Nil(t, out)
		require.ErrorIs(t, err, internal.ErrEntityNotFound)
		require.NotErrorIs(t, err, internal.ErrInternal)
	})
}
