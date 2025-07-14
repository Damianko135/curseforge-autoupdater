package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadDotenv loads environment variables from a .env file if present.
func LoadDotenv() error {
	return godotenv.Load()
}

// GetEnvVar returns the value of the environment variable or fallback if not set.
func GetEnvVar(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// GetIntEnvVar returns the int value of the environment variable or fallback if not set or invalid.
func GetIntEnvVar(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return intValue
}

// GetBoolEnvVar returns the bool value of the environment variable or fallback if not set or invalid.
func GetBoolEnvVar(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	if value == "true" || value == "1" {
		return true
	}
	if value == "false" || value == "0" {
		return false
	}
	return fallback
}

// GetBoolVar returns a boolean env var or fallback.
func GetBoolVar(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	if value == "true" || value == "1" {
		return true
	}
	return false
}

// GetIntVar returns an int env var or fallback.
func GetIntVar(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return intValue
}
