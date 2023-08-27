package entity

import (
	"github.com/tempcke/schedule"
)

type Lease struct {
	ID           ID
	PropertyID   ID
	TenantIDs    []ID
	StartDate    schedule.Date
	EndDate      schedule.Date
	Deposit      int    // dollars
	RentAmount   int    // dollars
	Currency     string // empty will be considered USD
	RentInterval Interval
}
type Interval = string

const (
	CurrencyUSD     = "USD"
	IntervalDaily   = "daily"
	IntervalWeekly  = "weekly"
	IntervalMonthly = "monthly"
)

func NewLease(propertyID ID) Lease {
	return Lease{
		ID:         NewID(),
		PropertyID: propertyID,
	}
}
func (l Lease) WithTenant(tenantIDs ...ID) Lease {
	for _, id := range tenantIDs {
		if !l.HasTenant(id) {
			l.TenantIDs = append(l.TenantIDs, id)
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
	for i := range l.TenantIDs {
		if l.TenantIDs[i] == id {
			return true
		}
	}
	return false
}
