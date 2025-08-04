// Package mailer_test provides test helpers and utilities for mailer integration tests.
// It includes configuration loading, test setup, and verification functions
// to facilitate testing email sending functionality with MailHog.
package mailer_test

import (
	"certitrack/internal/mailer"
	"certitrack/pkg/util"
	"strconv"
	"testing"
)

// TestConfig loads and returns the mailer configuration for integration tests.
// It reads environment variables and provides default values for testing.
func TestConfig(t *testing.T) mailer.Config {
	port, _ := strconv.Atoi(util.GetEnv("SMTP_PORT", "1025"))
	return mailer.Config{
		Host:     util.GetEnv("SMTP_HOST", "localhost"),
		Port:     port,
		Username: util.GetEnv("SMTP_USER", ""),
		Password: util.GetEnv("SMTP_PASSWORD", ""),
		From:     util.GetEnv("SMTP_FROM", "test@example.com"),
	}
}
