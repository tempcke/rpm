package entity

import (
	"github.com/tempcke/schedule"
)

type Lease struct {
	ID         ID
	Property   Property
	Tenants    []Tenant
	StartDate  schedule.Date
	EndDate    schedule.Date
	RentAmount int // dollars
	Deposit    int // dollars
}

func NewLease(p Property) Lease {
	return Lease{
		ID:       NewID(),
		Property: p,
	}
}
func (l Lease) WithTenant(tenants ...Tenant) Lease {
	for _, t := range tenants {
		if !l.HasTenant(t.ID) {
			l.Tenants = append(l.Tenants, t)
		}
	}
	return l
}
func (l Lease) WithRent(v int) Lease {
	l.RentAmount = v
	return l
}
func (l Lease) WithDeposit(v int) Lease {
	l.Deposit = v
	return l
}
func (l Lease) WithTerm(start, end schedule.Date) Lease {
	l.StartDate = start
	l.EndDate = end
	return l
}

func (l Lease) HasTenant(id ID) bool {
	for _, t := range l.Tenants {
		if t.ID == id {
			return true
		}
	}
	return false
}
