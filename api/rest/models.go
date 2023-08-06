package rest

import (
	"time"

	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/internal/lib/log"
)

// ErrorResponse response model
type ErrorResponse struct {
	Error string `json:"error"`
}

// PropertyList response model
type PropertyList struct {
	Items []PropertyModel `json:"items"`
}

func (l PropertyList) ToProperties() []entity.Property {
	properties := make([]entity.Property, len(l.Items))
	for i, p := range l.Items {
		createdAt, _ := time.Parse(time.RFC3339, p.CreatedAt)
		properties[i] = entity.Property{
			ID:        p.ID,
			Street:    p.Street,
			City:      p.City,
			StateCode: p.State,
			Zip:       p.Zip,
			CreatedAt: createdAt,
		}
	}
	return properties
}

// NewPropertyListModel constructs a PropertyList model from a set of properties
func NewPropertyListModel(propList ...entity.Property) PropertyList {
	pl := PropertyList{
		Items: make([]PropertyModel, len(propList)),
	}
	for i, p := range propList {
		pl.Items[i] = NewPropertyModel(p)
	}
	return pl
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

func (m PropertyModel) ToProperty() *entity.Property {
	return &entity.Property{
		ID:        m.ID,
		Street:    m.Street,
		City:      m.City,
		StateCode: m.State,
		Zip:       m.Zip,
		CreatedAt: m.createdAtTime(),
	}
}
func (m PropertyModel) createdAtTime() time.Time {
	if m.CreatedAt == "" {
		return time.Time{}
	}
	createdAt, err := time.Parse(time.RFC3339, m.CreatedAt)
	if err != nil {
		log.WithError(err).Errorf("PropertyModel.ToProperty time.Parse failed to parse: [%s]", m.CreatedAt)
		return time.Time{}
	}
	return createdAt
}
