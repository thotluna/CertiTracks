package auth_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"certitrack/internal/config"
	"certitrack/testutils"
	"certitrack/testutils/testcontainer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	ctx := context.Background()
	testCfg := &config.Config{
		Database: config.DatabaseConfig{
			Name:     "test_db",
			User:     "testuser",
			Password: "Password123!",
			SSLMode:  "disable",
		},
	}

	pgContainer, err := testcontainer.SetupPostgres(ctx, testCfg)
	require.NoError(t, err, "Failed to setup postgres container")
	defer pgContainer.Teardown(ctx)

	router := testutils.SetupTestRouter(t, pgContainer.DB)

	generateUniqueEmail := func(prefix string) string {
		return prefix + "_" + time.Now().Format("20060102150405") + "@example.com"
	}

	t.Run("should register a new user successfully", func(t *testing.T) {
		user := testutils.NewRegisterRequest(testutils.WithEmail(generateUniqueEmail("test_register"))).RegisterRequest

		response := testutils.RegisterTestUser(t, router, user)

		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("should fail with duplicate email", func(t *testing.T) {
		user := testutils.NewRegisterRequest(testutils.WithEmail(generateUniqueEmail("duplicate"))).RegisterRequest

		firstResponse := testutils.RegisterTestUser(t, router, user)
		assert.Equal(t, http.StatusCreated, firstResponse.Code)

		response := testutils.RegisterTestUser(t, router, user)

		assert.Equal(t, http.StatusConflict, response.Code)
	})

	t.Run("should fail with invalid email format", func(t *testing.T) {
		user := testutils.NewRegisterRequest(testutils.WithEmail("invalid-email-format")).RegisterRequest

		response := testutils.RegisterTestUser(t, router, user)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("should fail with short password", func(t *testing.T) {
		user := testutils.NewRegisterRequest(testutils.WithPassword("short")).RegisterRequest

		response := testutils.RegisterTestUser(t, router, user)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}
