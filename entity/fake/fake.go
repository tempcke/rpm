package fake

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/tempcke/schedule"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// LowerString returns a lowercase string of length n
func LowerString(n int) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = 'a' + rand.Int31n(26)
	}
	return string(s)
}
func FullName() string {
	var (
		fnameLen = 2 + rand.Intn(5)
		lnameLen = 3 + rand.Intn(5)
		fName    = LowerString(fnameLen)
		lName    = LowerString(lnameLen)
		fLetter  = strings.ToUpper(LowerString(1))
	)
	return fmt.Sprintf("%s%s %s%s", fLetter, fName, fLetter, lName)
}
func DateOfBirth() schedule.Date {
	var (
		age   = rand.Intn(20) + 18
		year  = time.Now().Year() - age
		month = rand.Intn(11) + 1
		day   = rand.Intn(25) + 1
	)
	return schedule.NewDate(year, time.Month(month), day)
}
