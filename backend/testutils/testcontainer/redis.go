// Package testcontainer provides utilities for managing Docker containers during testing.
// It simplifies the creation and management of containerized services like Redis and PostgreSQL
// to enable reliable integration testing with real dependencies.
package testcontainer

import (
	"context"
	"fmt"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RedisContainer struct {
	Container testcontainers.Container
}

func SetRedisContainer(ctx context.Context) (*RedisContainer, error) {
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
		return nil, fmt.Errorf("error creating container: %w", err)
	}

	return &RedisContainer{
		Container: redisC,
	}, nil
}

func (rc *RedisContainer) Teardown(ctx context.Context) error {
	if rc.Container != nil {
		if err := rc.Container.Terminate(ctx); err != nil {
			log.Printf("error to stop container: %v", err)
			return err
		}
	}

	return nil
}
