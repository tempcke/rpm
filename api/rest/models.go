package rest

import (
	"time"

	"github.com/tempcke/rpm/entity"
)

// ErrorResponse response model
type ErrorResponse struct {
	Error string `json:"error"`
}

// PropertyModel is a response model for a property
type PropertyModel struct {
	ID        string `json:"id"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	Zip       string `json:"zip"`
	CreatedAt string `json:"createdAt"`
}

// NewPropertyModel is the PropertyModel constructor
func NewPropertyModel(property entity.Property) PropertyModel {
	return PropertyModel{
		ID:        property.ID,
		Street:    property.Street,
		City:      property.City,
		State:     property.StateCode,
		Zip:       property.Zip,
		CreatedAt: property.CreatedAt.Format(time.RFC3339),
	}
}
