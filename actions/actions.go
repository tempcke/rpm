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
	Repo usecase.PropertyRepository
}

func NewActions(r usecase.PropertyRepository) Actions {
	return Actions{Repo: r}
}

func (a Actions) StoreProperty(ctx context.Context, p entity.Property) (specifications.ID, error) {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	uc := usecase.NewStoreProperty(a.Repo)
	if err := uc.Execute(ctx, p); err != nil {
		return "", err
	}
	return p.ID, nil
}
func (a Actions) RemoveProperty(ctx context.Context, id string) error {
	uc := usecase.NewDeleteProperty(a.Repo)
	return uc.Execute(ctx, id)
}

func (a Actions) ListProperties(ctx context.Context) ([]entity.Property, error) {
	uc := usecase.NewListProperties(a.Repo)
	list, err := uc.Execute(ctx)
	if err != nil {
		return nil, err
	}
	return list, nil
}
func (a Actions) GetProperty(ctx context.Context, id string) (*entity.Property, error) {
	uc := usecase.NewGetProperty(a.Repo)
	p, err := uc.Execute(ctx, id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
