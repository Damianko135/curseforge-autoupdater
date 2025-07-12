package helper

import (
	"os"
	"strconv"
)

func GetEnvVar(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

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
