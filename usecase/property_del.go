package usecase

type DeleteProperty struct {
	propRepo PropertyWriter
}

func NewDeleteProperty(repo PropertyWriter) DeleteProperty {
	return DeleteProperty{propRepo: repo}
}

func (uc DeleteProperty) Execute(id string) error {
	return uc.propRepo.DeleteProperty(id)
}
