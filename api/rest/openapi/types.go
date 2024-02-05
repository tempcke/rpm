package openapi

import (
	"encoding/json"
	"fmt"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/usecase"
	"github.com/tempcke/schedule"
)

type (
	StorePropertyRes = GetPropertyRes
	StoreTenantRes   = GetTenantRes
)

var (
	NewStorePropertyRes = NewGetPropertyRes
)

func (x Error) Error() string {
	var label = "openapi error"
	if x.Code == 0 {
		return fmt.Sprintf("%s: %s", label, x.Message)
	}
	return fmt.Sprintf("%s: (%d) %s", label, x.Code, x.Message)
}

func NewStorePropertyReq(p entity.Property) *StorePropertyReq {
	return &StorePropertyReq{
		Property: Address{
			Street: p.Street,
			City:   p.City,
			State:  p.StateCode,
			Zip:    p.Zip,
		},
	}
}
func (x *StorePropertyReq) ToProperty() entity.Property {
	return x.Property.ToProperty()
}
func (x *Property) GetID() string { return x.Id }
func (x *Address) ToProperty() entity.Property {
	p := Property{
		Street: x.Street,
		City:   x.City,
		State:  x.State,
		Zip:    x.Zip,
	}
	return p.ToProperty()
}
func (x *Property) ToProperty() entity.Property {
	return entity.Property{
		ID:        x.Id,
		Street:    x.Street,
		City:      x.City,
		StateCode: x.State,
		Zip:       x.Zip,
	}
}
func ToProperty(e entity.Property) *Property {
	return &Property{
		Id:     e.GetID(),
		Street: e.Street,
		City:   e.City,
		State:  e.StateCode,
		Zip:    e.Zip,
	}
}
func NewGetPropertyRes(in entity.Property) GetPropertyRes {
	return GetPropertyRes{
		Property: *ToProperty(in),
	}
}
func NewListPropertiesRes(in ...entity.Property) ListPropertiesRes {
	var list = make([]Property, len(in))
	for i, e := range in {
		list[i] = *ToProperty(e)
	}
	return ListPropertiesRes{
		Properties: list,
	}
}
func (x *ListPropertiesRes) ToProperties() []entity.Property {
	var list = make([]entity.Property, len(x.Properties))
	for i, e := range x.Properties {
		list[i] = e.ToProperty()
	}
	return list
}
func (x *ListPropertiesParams) ToFilter() usecase.PropertyFilter {
	return usecase.PropertyFilter{
		Search: removePointer(x.Search),
	}
}

type Date = types.Date

func ToDate(in schedule.Date) Date {
	return Date{Time: in.ToTime()}
}
func (x *MinTenant) ToTenant() entity.Tenant {
	return entity.Tenant{
		FullName:    x.FullName,
		DateOfBirth: schedule.NewDateFromTime(x.Dob.Time),
		Phones:      FromPhones(x.Phones...),
		DLNum:       x.DlNum,
		DLState:     x.DlState,
	}
}
func (x *Tenant) GetID() string { return x.Id }
func (x *Tenant) ToTenant() *entity.Tenant {
	return &entity.Tenant{
		ID:          x.GetID(),
		FullName:    x.FullName,
		DateOfBirth: schedule.NewDateFromTime(x.Dob.Time),
		Phones:      FromPhones(x.Phones...),
		DLNum:       x.DlNum,
		DLState:     x.DlState,
	}
}
func (x *Tenant) JSON() []byte {
	if x != nil {
		if bytes, err := json.Marshal(*x); err == nil {
			return bytes
		}
	}
	return nil
}
func (x *StoreTenantReq) JSON() []byte {
	if x != nil {
		if bytes, err := json.Marshal(*x); err == nil {
			return bytes
		}
	}
	return nil
}
func (x TenantList) ToTenantMap() map[entity.ID]entity.Tenant {
	var list = make(map[entity.ID]entity.Tenant)
	for _, t := range x.Tenants {
		list[t.GetID()] = *t.ToTenant()
	}
	return list
}
func (x TenantList) ToTenants() []entity.Tenant {
	if len(x.Tenants) == 0 {
		return nil
	}
	var list = make([]entity.Tenant, len(x.Tenants))
	for i, t := range x.Tenants {
		list[i] = *t.ToTenant()
	}
	return list
}
func ToTenantList(in ...entity.Tenant) TenantList {
	var list = make([]Tenant, len(in))
	for i, e := range in {
		list[i] = *ToTenant(e)
	}
	return TenantList{
		Tenants: list,
	}
}
func NewStoreTenantReq(in entity.Tenant) *StoreTenantReq {
	return &StoreTenantReq{
		Tenant: MinTenant{
			DlNum:    in.DLNum,
			DlState:  in.DLState,
			Dob:      ToDate(in.DateOfBirth),
			FullName: in.FullName,
			Phones:   ToPhones(in.Phones...),
		},
	}
}
func NewGetTenantRes(in entity.Tenant) GetTenantRes {
	return GetTenantRes{Tenant: *ToTenant(in)}
}
func ToTenant(in entity.Tenant) *Tenant {
	return &Tenant{
		Id:       in.GetID(),
		DlNum:    in.DLNum,
		DlState:  in.DLState,
		Dob:      ToDate(in.DateOfBirth),
		FullName: in.FullName,
		Phones:   ToPhones(in.Phones...),
	}
}
func ToPhones(in ...entity.Phone) []Phone {
	if len(in) == 0 {
		return nil
	}
	var list = make([]Phone, len(in))
	for i, e := range in {
		list[i] = Phone{
			Desc:   e.Note,
			Number: e.Number,
		}
	}
	return list
}
func FromPhones(in ...Phone) []entity.Phone {
	if len(in) == 0 {
		return nil
	}
	var list = make([]entity.Phone, len(in))
	for i, e := range in {
		list[i] = entity.Phone{
			Number: e.Number,
			Note:   e.Desc,
		}
	}
	return list
}

func removePointer[T any](in *T) T {
	var out T
	if in != nil {
		out = *in
	}
	return out
}
