package testutils

import (
	"context"
	"os"
	"testing"

	"certitrack/internal/config"
	"certitrack/internal/database"
	"certitrack/testutils/testcontainer"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func SetupTestEnvironment(t *testing.T) (*config.Config, *gorm.DB, func()) {
	t.Helper()

	os.Setenv("APP_ENV", "test")

	ctx := context.Background()

	testCfg := &config.Config{
		Database: config.DatabaseConfig{
			Name:     "test_db",
			User:     "testuser",
			Password: "testpassword",
			SSLMode:  "disable",
		},
	}

	pgContainer, err := testcontainer.SetupPostgres(ctx, testCfg)
	require.NoError(t, err, "Failed to set up PostgreSQL container")

	db, err := database.Connect(pgContainer.Config)
	require.NoError(t, err, "Failed to connect to test database")

	err = database.AutoMigrate(db)
	require.NoError(t, err, "Failed to run migrations")
	return pgContainer.Config, db, func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}

		if err := pgContainer.Teardown(ctx); err != nil {
			t.Logf("Error stopping container: %v", err)
		}
	}
}
