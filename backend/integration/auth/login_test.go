package auth_test

import (
	"certitrack/testutils"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	t.Run("should login with valid credentials", func(t *testing.T) {
		rtr := setupTestRouter(t)
		defer rtr.DB.Teardown(context.Background())

		user := testutils.NewRegisterRequest(
			testutils.WithEmail(GenerateUniqueEmail("login_test")),
			testutils.WithPassword("ValidPass123!"),
		).RegisterRequest
		registerTestUser(t, rtr, user)

		response := loginTestUser(t, rtr, user.Email, "ValidPass123!")

		assert.Equal(t, http.StatusOK, response.Code)
		data := getResponseData(t, response)
		assert.NotEmpty(t, data["access-token"], "Access token should not be empty")
		assert.NotEmpty(t, data["refresh-token"], "Refresh token should not be empty")
	})

	t.Run("should fail with invalid password", func(t *testing.T) {
		rtr := setupTestRouter(t)
		defer rtr.DB.Teardown(context.Background())

		user := testutils.NewRegisterRequest(
			testutils.WithEmail(GenerateUniqueEmail("login_test")),
			testutils.WithPassword("ValidPass123!"),
		).RegisterRequest
		registerTestUser(t, rtr, user)
		response := loginTestUser(t, rtr, user.Email, "wrongpassword")

		assert.Equal(t, http.StatusUnauthorized, response.Code)

		var responseBody map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &responseBody)
		require.NoError(t, err, "Failed to parse error response")
		assert.Contains(t, responseBody, "error", "Error message should be present in response")
	})

	t.Run("should fail with non-existent user", func(t *testing.T) {
		rtr := setupTestRouter(t)
		defer rtr.DB.Teardown(context.Background())

		nonExistentEmail := GenerateUniqueEmail("nonexistent")
		response := loginTestUser(t, rtr, nonExistentEmail, "password123")

		assert.Equal(t, http.StatusUnauthorized, response.Code)

		var responseBody map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &responseBody)
		require.NoError(t, err, "Failed to parse error response")
		assert.Contains(t, responseBody, "error", "Error message should be present in response")
	})
}
