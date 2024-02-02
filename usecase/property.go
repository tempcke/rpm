package usecase

import (
	"context"
	"errors"

	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/internal"
)

// PropertyReader allows queries regarding properties
type (
	PropertyManager struct {
		propRepo PropertyRepo
	}
	PropertyReader interface {
		GetProperty(ctx context.Context, id string) (entity.Property, error)
		PropertyList(ctx context.Context) ([]entity.Property, error)
	}
	PropertyWriter interface {
		NewProperty(street, city, state, zip string) entity.Property
		StoreProperty(context.Context, entity.Property) error
		DeleteProperty(ctx context.Context, id string) error
	}
	PropertyRepo interface {
		PropertyReader
		PropertyWriter
	}
)

var (
	noProperty = entity.Property{}
)

func NewPropertyManager(repo PropertyRepo) PropertyManager {
	return PropertyManager{propRepo: repo}
}
func (uc PropertyManager) Store(ctx context.Context, p entity.Property) error {
	if err := uc.Validate(); err != nil {
		return err
	}
	if err := p.Validate(); err != nil {
		return err
	}
	if err := uc.propRepo.StoreProperty(ctx, p); err != nil {
		// TODO: make sure the error is logged here or in the repo layer
		return internal.NewErrors(internal.ErrInternal, ErrRepo)
	}
	return nil
}
func (uc PropertyManager) Get(ctx context.Context, id string) (entity.Property, error) {
	if err := uc.Validate(); err != nil {
		return noProperty, err
	}
	p, err := uc.propRepo.GetProperty(ctx, id)
	if err != nil {
		if errors.Is(err, internal.ErrEntityNotFound) {
			return p, err
		}
		// TODO: make sure the error is logged here or in the repo layer
		return noProperty, internal.NewErrors(internal.ErrInternal, ErrRepo)
	}
	return p, nil
}
func (uc PropertyManager) List(ctx context.Context) ([]entity.Property, error) {
	if err := uc.Validate(); err != nil {
		return nil, err
	}
	return uc.propRepo.PropertyList(ctx)
}
func (uc PropertyManager) Remove(ctx context.Context, id string) error {
	if err := uc.Validate(); err != nil {
		return err
	}
	if err := uc.propRepo.DeleteProperty(ctx, id); err != nil {
		// TODO: make sure the error is logged here or in the repo layer
		return internal.NewErrors(internal.ErrInternal, ErrRepo)
	}
	return nil
}
func (uc PropertyManager) Validate() error {
	if uc.propRepo == nil {
		return internal.NewErrors(internal.ErrInternal, ErrRepoNotSet)
	}
	return nil
}
