package usecase

// DeleteProperty Use Case
type DeleteProperty struct {
	propRepo PropertyWriter
}

// NewDeleteProperty constructs a DeleteProperty use case
func NewDeleteProperty(repo PropertyWriter) DeleteProperty {
	return DeleteProperty{propRepo: repo}
}

// Execute the DeleteProperty use case to delete a property by ID
func (uc DeleteProperty) Execute(id string) error {
	return uc.propRepo.DeleteProperty(id)
}
