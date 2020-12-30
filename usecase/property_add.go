package usecase

import "github.com/tempcke/rpm/entity"

// AddProperty is a use case to add a property
type AddProperty struct {
	propRepo PropertyWriter
}

// NewAddProperty constructs and returns an AddProperty
func NewAddProperty(repo PropertyWriter) AddProperty {
	return AddProperty{propRepo: repo}
}

// Execute the use case
func (c AddProperty) Execute(property entity.Property) error {
	if err := property.Validate(); err != nil {
		return err
	}
	return c.propRepo.StoreProperty(property)
}
