package fake

import (
	"math/rand"
	"strconv"

	"github.com/google/uuid"
	"github.com/tempcke/rpm/entity"
)

func Property() entity.Property {
	streetNum := strconv.Itoa(rand.Intn(9000) + 1000)
	return entity.Property{
		ID:        uuid.NewString(),
		Street:    streetNum + " Main st.",
		City:      "Dallas",
		StateCode: "TX",
		Zip:       "75001",
	}
}
