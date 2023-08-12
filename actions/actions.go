package actions

import (
	"context"

	"github.com/google/uuid"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/specifications"
	"github.com/tempcke/rpm/usecase"
)

var _ specifications.Driver = (*Actions)(nil)

type Actions struct {
	propRepo   usecase.PropertyRepo
	tenantRepo usecase.TenantRepo
}
type Repo interface {
	usecase.PropertyRepo
	usecase.TenantRepo
}

func NewActions() Actions               { return Actions{} }
func NewActionsWithRepo(r Repo) Actions { return Actions{propRepo: r, tenantRepo: r} }
func (a Actions) WithPropertyRepo(r usecase.PropertyRepo) Actions {
	a.propRepo = r
	return a
}
func (a Actions) WithTenantRepo(r usecase.TenantRepo) Actions {
	a.tenantRepo = r
	return a
}

func (a Actions) StoreProperty(ctx context.Context, p entity.Property) (entity.ID, error) {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	uc := usecase.NewStoreProperty(a.propRepo)
	if err := uc.Execute(ctx, p); err != nil {
		return "", err
	}
	return p.ID, nil
}
func (a Actions) RemoveProperty(ctx context.Context, id string) error {
	uc := usecase.NewDeleteProperty(a.propRepo)
	return uc.Execute(ctx, id)
}
func (a Actions) ListProperties(ctx context.Context) ([]entity.Property, error) {
	uc := usecase.NewListProperties(a.propRepo)
	list, err := uc.Execute(ctx)
	if err != nil {
		return nil, err
	}
	return list, nil
}
func (a Actions) GetProperty(ctx context.Context, id string) (*entity.Property, error) {
	uc := usecase.NewGetProperty(a.propRepo)
	p, err := uc.Execute(ctx, id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (a Actions) StoreTenant(ctx context.Context, e entity.Tenant) (*entity.Tenant, error) {
	if e.ID == "" {
		e.ID = uuid.NewString()
	}
	return a.tenantMan().Store(ctx, e)
}
func (a Actions) GetTenant(ctx context.Context, id entity.ID) (*entity.Tenant, error) {
	return a.tenantMan().Get(ctx, id)
}
func (a Actions) ListTenants(ctx context.Context) ([]entity.Tenant, error) {
	return a.tenantMan().List(ctx)
}
func (a Actions) tenantMan() usecase.TenantManager {
	return usecase.NewTenantManager(a.tenantRepo)
}
