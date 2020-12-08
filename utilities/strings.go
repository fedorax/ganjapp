package utilities

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GetRandomString generates and returns a random string of the passed length
func GetRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

// UUID returns a UUID as a string
func UUID() string {
	return uuid.New().String()
}
