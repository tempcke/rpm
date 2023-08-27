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
func (t Tenant) WithName(name string) Tenant {
	t.FullName = name
	return t
}
func (t Tenant) WithDOB(date schedule.Date) Tenant {
	t.DateOfBirth = date
	return t
}

func (t Tenant) GetID() ID { return t.ID }
func (t Tenant) Equal(t2 Tenant) bool {
	return idEqualOrEmpty(t.ID, t2.ID) &&
		t.FullName == t2.FullName &&
		t.DateOfBirth.Equal(t2.DateOfBirth) &&
		phoneListEqual(t.Phones, t2.Phones)
}

func (t Tenant) Ptr() *Tenant { return &t }

func phoneListEqual(a, b []Phone) bool {
	if len(a) != len(b) {
		return false
	}
	if len(a) == 0 {
		return true
	}
	for i := range a {
		found := false
		for j := range b {
			if a[i].Equal(b[j]) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
