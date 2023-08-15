package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// var r *rand.Rand

// func init() {
// 	r = rand.New(rand.NewSource(time.Now().UnixNano()))
// }

// RandomInt generates a random number between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomStr generates a random string of n length
func RandomStr(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// Generates random owner name
func RandomOwner() string {
	return RandomStr(6)
}

// Generates random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// Generates random currency
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "KZT"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// Generates radom email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomStr(6))
}
