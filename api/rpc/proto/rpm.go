package pb

import (
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/schedule"
)

func (x *Property) ToProperty() entity.Property {
	return entity.Property{
		ID:        x.GetPropertyID(),
		Street:    x.GetStreet(),
		City:      x.GetCity(),
		StateCode: x.GetState(),
		Zip:       x.GetZip(),
	}
}
func ToProperty(e entity.Property) *Property {
	return &Property{
		PropertyID: e.GetID(),
		Street:     e.Street,
		City:       e.City,
		State:      e.StateCode,
		Zip:        e.Zip,
	}
}

func (x *Tenant) ToTenant() entity.Tenant {
	e := entity.Tenant{
		ID:       x.GetTenantID(),
		FullName: x.GetFullName(),
		Phones:   nil,
		DLNum:    x.GetDlNum(),
		DLState:  x.GetDlState(),
		// TODO: phones?
	}

	if dob := schedule.ParseDate(x.GetDob()); dob != nil {
		return e.WithDOB(*dob)
	}
	return e
}
func ToTenant(e entity.Tenant) *Tenant {
	return &Tenant{
		TenantID: e.GetID(),
		FullName: e.FullName,
		DlNum:    e.DLNum,
		DlState:  e.DLState,
		Dob:      e.DateOfBirth.String(),
		// TODO: phones?
	}
}
