package test

import (
	"math/rand"
)

// RandString produces a random string of lower case letters [a-z]
func RandString(n int) string {
	var runes = make([]rune, n)
	for i := 0; i < n; i++ {
		runes[i] = 'a' + rand.Int31n(26)
	}
	return string(runes)
}
