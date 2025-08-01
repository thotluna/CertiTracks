package redis_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	rediscache "certitrack/internal/cache/redis"
	"certitrack/internal/config"
)

func TestRedisClient_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	cfg := &config.Config{
		Redis: config.RedisConfig{
			URL:      getRedisURL(),
			Password: getRedisPassword(),
		},
	}

	client, err := rediscache.NewClient(&cfg.Redis)
	require.NoError(t, err, "should create Redis client")
	defer client.Close()
	rdb := client.Client()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := "test:key"
	value := "test value"
	expiration := 1 * time.Minute

	err = rdb.Set(ctx, key, value, expiration).Err()
	assert.NoError(t, err, "should set value in Redis")

	got, err := rdb.Get(ctx, key).Result()
	assert.NoError(t, err, "should get value from Redis")
	assert.Equal(t, value, got, "should get the same value that was set")

	err = rdb.Del(ctx, key).Err()
	assert.NoError(t, err, "should delete key from Redis")

	_, err = rdb.Get(ctx, key).Result()
	assert.Equal(t, redis.Nil, err, "should return redis.Nil when key doesn't exist")
}

func getRedisURL() string {
	if url := os.Getenv("REDIS_URL"); url != "" {
		return url
	}
	return "redis://localhost:6379/0"
}

func getRedisPassword() string {
	if pass := os.Getenv("REDIS_PASSWORD"); pass != "" {
		return pass
	}
	return "dev_redis_password"
}
