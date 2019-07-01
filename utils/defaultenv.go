package utils

import "os"

func DefaultEnv(key, def string) string {
	value := os.Getenv(key)
	if value == "" {
		value = def
	}
	return value
}
