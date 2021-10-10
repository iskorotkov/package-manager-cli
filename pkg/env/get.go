package env

import (
	"os"
	"strconv"
)

func Get(name string, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}

	return value
}

func GetBool(name string, fallback bool) bool {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return b
}

func GetInt(name string, fallback int) int {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return i
}
