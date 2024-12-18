package config

import "os"

// GetEnv helps to look up if the env exist or a default is provided
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
