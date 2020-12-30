package usecase

import "github.com/tempcke/rpm/entity"

// PropertyRepository is used by the usecase for property CRUD operations
type PropertyRepository interface {
	NewProperty(street, city, state, zip string) entity.Property
	StoreProperty(entity.Property) error
	RetrieveProperty(id string) (entity.Property, error)
}

type propertyUseCase struct {
	propRepo PropertyRepository
}

// AddPropertyCommand is a use case to add a property
type AddPropertyCommand struct {
	propertyUseCase
}

// NewAddPropertyCommand constructs and returns an AddPropertyCommand
func NewAddPropertyCommand(propRepo PropertyRepository) AddPropertyCommand {
	uc := AddPropertyCommand{}
	uc.propRepo = propRepo
	return uc
}

// Execute the use case
func (c AddPropertyCommand) Execute(property entity.Property) error {
	if err := property.Validate(); err != nil {
		return err
	}
	return c.propRepo.StoreProperty(property)
}

type GetPropertyQuery struct {
	propertyUseCase
}

// NewGetPropertyQuery constructs and returns an GetPropertyQuery
func NewGetPropertyQuery(r PropertyRepository) GetPropertyQuery {
	uc := GetPropertyQuery{}
	uc.propRepo = r
	return uc
}

func (uc GetPropertyQuery) Execute(id string) (entity.Property, error) {
	return uc.propRepo.RetrieveProperty(id)
}
