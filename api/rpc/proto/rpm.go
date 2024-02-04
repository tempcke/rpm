package pb

import (
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/usecase"
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
		DLNum:    x.GetDlNum(),
		DLState:  x.GetDlState(),
		Phones:   FromPhones(x.GetPhones()),
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
		Phones:   ToPhones(e.Phones),
	}
}
func FromPhones(phones []*Phone) []entity.Phone {
	var list []entity.Phone
	for _, p := range phones {
		list = append(list, entity.Phone{
			Number: p.Number,
			Note:   p.Note,
		})
	}
	return list
}
func ToPhones(phones []entity.Phone) []*Phone {
	var list []*Phone
	for _, e := range phones {
		list = append(list, ToPhone(e))
	}
	return list
}
func ToPhone(e entity.Phone) *Phone {
	return &Phone{
		Number: e.Number,
		Note:   e.Note,
	}
}

func (x *ListPropertiesReq) ToPropertyFilter() usecase.PropertyFilter {
	return usecase.PropertyFilter{
		Search: x.GetSearch(),
	}
}
func FromPropertyFilter(f usecase.PropertyFilter) *ListPropertiesReq {
	return &ListPropertiesReq{
		Search: f.Search,
	}
}
