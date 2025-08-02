// Package redis provides a Redis-based implementation of the cache repository.
package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"certitrack/internal/config"
)

type Client struct {
	client *redis.Client
}

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

func (c *Client) Close() error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close()
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.client.Set(ctx, key, value, expiration)
}
func (c *Client) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.client.Exists(ctx, keys...)
}
