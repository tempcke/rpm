package fake

import (
	"math/rand"
	"strconv"

	"github.com/tempcke/rpm/entity"
)

func Property() entity.Property {
	streetNum := strconv.Itoa(rand.Intn(8000) + 1000)
	return entity.Property{
		ID:        entity.NewID(),
		Street:    streetNum + " Main st.",
		City:      "Dallas",
		StateCode: "TX",
		Zip:       "75001",
	}
}
func Tenant() entity.Tenant {
	var phones []entity.Phone
	n := rand.Intn(3)
	for i := 0; i < n; i++ {
		phones = append(phones, Phone())
	}
	return entity.Tenant{
		ID:          entity.NewID(),
		FullName:    FullName(),
		DateOfBirth: DateOfBirth(),
		Phones:      phones,
		DLNum:       LowerString(8),
		DLState:     "TX",
	}
}
func Phone() entity.Phone {
	n := rand.Intn(8000) + 1000
	return entity.Phone{
		Number: "555-555-" + strconv.Itoa(n),
		Note:   LowerString(5),
	}
}
