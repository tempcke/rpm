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
	return entity.Tenant{
		ID:          entity.NewID(),
		FullName:    FullName(),
		DateOfBirth: DateOfBirth(),
	}
}
