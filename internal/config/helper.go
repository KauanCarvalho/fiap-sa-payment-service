package config

import (
	"os"
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
