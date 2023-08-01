package entity

import (
	"github.com/tempcke/schedule"
)

type RentalDetails struct {
	PropertyID    string
	AllowSmoking  bool
	AllowPets     bool
	ParkingSpaces int
	ParkingDesc   string
}
type RentalApplication struct {
	PropertyID string
}
type Applicant struct {
	FullName        string
	DateOfBirth     schedule.Date
	SSN             string
	DLNum           string // drivers license number
	DLState         string // drivers license state
	HasPets         bool
	PetsDesc        string
	VehicleCount    int
	VehicleDesc     string
	CrimeConviction bool
	CrimeDesc       string
	BankruptcyFiled bool
	BankruptcyDesc  string
}
