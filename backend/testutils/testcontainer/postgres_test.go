package testcontainer

import (
	"context"
	"testing"

	"certitrack/internal/config"
	"github.com/stretchr/testify/require"
)

func TestSetupPostgres(t *testing.T) {
	ctx := context.Background()

	t.Run("should initialize postgres container with test configuration", func(t *testing.T) {
		testCfg := &config.Config{
			Database: config.DatabaseConfig{
				Name:     "test_db",
				User:     "testuser",
				Password: "testpassword",
				SSLMode:  "disable",
			},
		}

		pgContainer, err := SetupPostgres(ctx, testCfg)
		require.NoError(t, err, "Failed to setup PostgreSQL container")
		defer pgContainer.Teardown(ctx)

		require.NotEmpty(t, pgContainer.Config.Database.Port, "A dynamic port should be assigned")
		require.NotEqual(t, "0", pgContainer.Config.Database.Port, "Port cannot be 0 after assignment")

		require.NotEmpty(t, pgContainer.Config.Database.User, "Database user should not be empty")
		require.NotEmpty(t, pgContainer.Config.Database.Password, "Database password should not be empty")
		require.NotEmpty(t, pgContainer.Config.Database.Name, "Database name should not be empty")
		require.Equal(t, "disable", pgContainer.Config.Database.SSLMode, "SSL mode should be disabled for tests")

		db := pgContainer.DB
		require.NotNil(t, db, "Database connection should not be nil")

		var result int
		err = db.Raw("SELECT 1").Scan(&result).Error
		require.NoError(t, err, "Failed to execute test query")
		require.Equal(t, 1, result, "Query result should be 1")

		var tableExists bool
		err = db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&tableExists).Error
		require.NoError(t, err, "Failed to verify 'users' table existence")
		require.True(t, tableExists, "'users' table should exist after migrations")
	})

	t.Run("should use custom port when specified", func(t *testing.T) {
		customPort := "5434"
		t.Setenv("POSTGRES_TEST_PORT", customPort)

		testCfg := &config.Config{
			Database: config.DatabaseConfig{
				Name:     "test_db_port",
				User:     "testuser",
				Password: "testpassword",
				SSLMode:  "disable",
			},
		}
		
		pgContainer, err := SetupPostgres(ctx, testCfg)
		require.NoError(t, err, "Failed to setup PostgreSQL container")
		defer pgContainer.Teardown(ctx)

		require.Equal(t, customPort, pgContainer.Config.Database.Port, "Should use port specified in POSTGRES_TEST_PORT")
		
		db := pgContainer.DB
		require.NotNil(t, db, "Database connection should not be nil")

		var result int
		err = db.Raw("SELECT 1").Scan(&result).Error
		require.NoError(t, err, "Failed to execute test query on custom port")
		require.Equal(t, 1, result, "Query result should be 1")
	})
}
