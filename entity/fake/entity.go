package fake

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/internal/test"
)

func Property() entity.Property {
	var (
		scope  = test.RandString(5)
		number = strconv.Itoa(rand.Intn(80000) + 10000)
	)
	return entity.Property{
		ID:        entity.NewID(),
		Street:    fmt.Sprintf("%s N %s st.", number[0:3], ucFirst(scope)),
		City:      ucFirst(scope[1:5]) + " City",
		StateCode: strings.ToUpper(scope[2:4]),
		Zip:       number[0:5],
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
func ucFirst(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
}
