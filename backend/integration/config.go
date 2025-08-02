// Package integration provides end-to-end integration tests for the application.
// It includes test helpers and configurations for setting up and managing
// external dependencies like databases and caches during testing.
package integration

import (
	"certitrack/testutils/testcontainer"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func SetupTestDB(t *testing.T) (*testcontainer.RedisContainer, func()) {
	t.Helper()
	ctx := context.Background()
	rContainer, err := testcontainer.SetRedisContainer(ctx)
	require.NoError(t, err, "Failed to setup postgres container")

	cleanup := func() {
		rContainer.Teardown(ctx)
	}

	return rContainer, cleanup
}
