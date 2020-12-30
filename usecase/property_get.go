package usecase

import "github.com/tempcke/rpm/entity"

type GetProperty struct {
	propRepo PropertyReader
}

// NewGetProperty constructs and returns an GetProperty
func NewGetProperty(repo PropertyReader) GetProperty {
	return GetProperty{propRepo: repo}
}

func (uc GetProperty) Execute(id string) (entity.Property, error) {
	return uc.propRepo.RetrieveProperty(id)
}
