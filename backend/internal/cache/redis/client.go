// Package redis provides a Redis-based implementation of the cache repository.
package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"certitrack/internal/config"
)

// Client wraps a Redis client with our custom configuration.
type Client struct {
	client *redis.Client
}

// NewClient creates a new Redis client with the provided configuration.
func NewClient(cfg *config.RedisConfig) (*Client, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	opt, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, err
	}

	if cfg.Password != "" {
		opt.Password = cfg.Password
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

// Close closes the Redis client connection.
// It is safe to call Close on a nil Client.
func (c *Client) Close() error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close()
}

// Client returns the underlying Redis client.
// This method is primarily for testing purposes.
func (c *Client) Client() *redis.Client {
	return c.client
}
