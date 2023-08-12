package usecase

import (
	"context"
	"errors"

	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/internal/filters"
)

type TenantManager struct {
	repo TenantRepo
}
type TenantRepo interface {
	StoreTenant(context.Context, entity.Tenant) error
	GetTenant(context.Context, entity.ID) (*entity.Tenant, error)
	ListTenants(context.Context, ...filters.TenantFilter) ([]entity.Tenant, error)
}

func NewTenantManager(repo TenantRepo) TenantManager {
	return TenantManager{repo: repo}
}

func (uc TenantManager) Store(ctx context.Context, tenant entity.Tenant) (*entity.Tenant, error) {
	if err := uc.Validate(); err != nil {
		return nil, err
	}
	if err := uc.repo.StoreTenant(ctx, tenant); err != nil {
		// TODO: make sure the error is logged here or in the repo layer
		return nil, internal.NewErrors(internal.ErrInternal, ErrRepo)
	}
	return &tenant, nil
}
func (uc TenantManager) Get(ctx context.Context, id entity.ID) (*entity.Tenant, error) {
	if err := uc.Validate(); err != nil {
		return nil, err
	}
	e, err := uc.repo.GetTenant(ctx, id)
	if err != nil {
		if errors.Is(err, internal.ErrEntityNotFound) {
			return nil, err
		}
		// TODO: make sure the error is logged here or in the repo layer
		return nil, internal.NewErrors(internal.ErrInternal, ErrRepo)
	}
	return e, nil
}
func (uc TenantManager) List(ctx context.Context, filter ...filters.TenantFilter) ([]entity.Tenant, error) {
	if err := uc.Validate(); err != nil {
		return nil, err
	}
	list, err := uc.repo.ListTenants(ctx, filter...)
	if err != nil {
		// TODO: make sure the error is logged here or in the repo layer
		return nil, internal.NewErrors(internal.ErrInternal, ErrRepo)
	}
	return list, nil
}
func (uc TenantManager) Validate() error {
	if uc.repo == nil {
		return internal.NewErrors(internal.ErrInternal, ErrRepoNotSet)
	}
	return nil
}
