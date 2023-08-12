package entity

import (
	"github.com/tempcke/schedule"
)

type Tenant struct {
	ID          ID
	FullName    string
	DateOfBirth schedule.Date
	Phones      []Phone
	DLNum       string // drivers license number
	DLState     string // drivers license state
}

func NewTenant(name string, dob schedule.Date) Tenant {
	return Tenant{
		ID:          NewID(),
		FullName:    name,
		DateOfBirth: dob,
	}
}
func (t Tenant) WithID(id ID) Tenant {
	t.ID = id
	return t
}
func (t Tenant) WithPhone(p Phone) Tenant {
	t.Phones = append(t.Phones, p)
	return t
}
func (t Tenant) GetID() ID { return t.ID }
func (t Tenant) Equal(t2 Tenant) bool {
	return idEqualOrEmpty(t.ID, t2.ID) &&
		t.FullName == t2.FullName &&
		t.DateOfBirth.Equal(t2.DateOfBirth)
}
