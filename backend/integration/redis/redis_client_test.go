package redis_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"certitrack/internal/cache/redis"
	"certitrack/internal/repositories"
	"certitrack/testutils"
)

func TestRedisClient_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Establecer entorno de prueba
	os.Setenv("APP_ENV", "test")
	
	// Obtener configuración de prueba
	cfg := testutils.GetTestConfig()

	client, err := redis.NewClient(&cfg.Redis)
	require.NoError(t, err, "should create Redis client")
	defer client.Close()

	var redisClient repositories.RedisClient = client

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := "test:key:" + time.Now().Format(time.RFC3339Nano)
	value := "test value"
	expiration := 1 * time.Minute

	err = redisClient.Set(ctx, key, value, expiration).Err()
	require.NoError(t, err, "should set value in Redis")

	exists := redisClient.Exists(ctx, key).Val()
	assert.Equal(t, int64(1), exists, "key should exist after setting")
}
