package redis_test

import (
	"context"
	"testing"
	"time"

	"certitrack/internal/cache/redis"
	"certitrack/internal/config"
	"certitrack/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestTokenRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("No se pudo iniciar Redis: %v", err)
	}
	defer redisC.Terminate(ctx)

	endpoint, err := redisC.Host(ctx)
	if err != nil {
		t.Fatalf("No se pudo obtener el host de Redis: %v", err)
	}

	port, err := redisC.MappedPort(ctx, "6379/tcp")
	if err != nil {
		t.Fatalf("No se pudo obtener el puerto de Redis: %v", err)
	}

	client, err := redis.NewClient(&config.RedisConfig{
		URL: "redis://" + endpoint + ":" + port.Port(),
	})
	if err != nil {
		t.Fatalf("No se pudo conectar a Redis: %v", err)
	}
	defer client.Close()

	repo := repositories.NewTokenRepository(client)

	t.Run("debe revocar y verificar token", func(t *testing.T) {
		token := "test-jwt-token-" + time.Now().String()
		expiration := time.Hour

		err := repo.RevokeToken(token, expiration)
		assert.NoError(t, err)

		revoked, err := repo.IsTokenRevoked(token)
		assert.NoError(t, err)
		assert.True(t, revoked)

		// Verificar que un token no revocado devuelve falso
		notRevoked, err := repo.IsTokenRevoked("non-existent-token")
		assert.NoError(t, err)
		assert.False(t, notRevoked)
	})
}
