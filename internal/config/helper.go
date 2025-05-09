package config

import (
	"os"
	"strconv"
)

func fetchEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	panic("Mandatory env var missing: " + key)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}

		return fallback
	}

	return fallback
}

func getEnvAsInt64(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}

		return fallback
	}

	return fallback
}
