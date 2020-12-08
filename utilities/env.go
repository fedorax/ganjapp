package utilities

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv attempts to populate environment variables
func LoadEnv() bool {
	// Try and load environment variables from .env file...
	if FileExists(".env") {
		err := godotenv.Overload()
		if err != nil {
			log.Fatal("Error loading .env file")
			return false
		}
	} else {
		fmt.Println("No .env file to load, using only OS environment variables!")
		return false
	}
	return true
}

// GetEnv attempts to fetch a value from the environment variables, returning it if found, or the passed fallback string if not
func GetEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	fmt.Println("Warning using default / generated value for " + key + "! (" + fallback + ")")
	return fallback
}
