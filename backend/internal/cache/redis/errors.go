package redis

import "errors"

// Common Redis client errors
var (
	// ErrNilConfig is returned when a nil config is provided to NewClient
	ErrNilConfig = errors.New("redis: config cannot be nil")
)
