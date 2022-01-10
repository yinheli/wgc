package env

import (
	"os"
	"strings"
)

func GetOrDefaultString(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func GetOrDefaultBool(key string, defaultValue bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		return strings.EqualFold(value, "true")
	}
	return defaultValue
}
