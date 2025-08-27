package util

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomInt(min, max int64) int64 {
	return min + rand.Int64N(max-min+1)
}

func RandomString(length int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < length; i++ {
		sb.WriteByte(alphabet[rand.IntN(k)])
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "GBP", "JPY"}
	return currencies[rand.IntN(len(currencies))]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@%s.com", RandomOwner(), RandomString(4))
}
