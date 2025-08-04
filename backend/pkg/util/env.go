// Package util provides general-purpose utility functions for the application.
// This file specifically contains environment variable handling utilities that
// simplify working with configuration values across different environments.
package util

import "os"

// GetEnv retrieves the value of an environment variable with a fallback default.
// If the environment variable is not set, it returns the provided default value.
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
