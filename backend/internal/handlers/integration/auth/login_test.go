package auth_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	setupTestUser := func(router *testRouter) testUser {
		user := testUser{
			Email:     "login_test@example.com",
			Password:  "Password123!",
			FirstName: "Login",
			LastName:  "Test",
		}
		w := registerTestUser(t, router, user)
		fmt.Println("Register response:", w.Body.String())
		return user
	}

	t.Run("should login with valid credentials", func(t *testing.T) {
		router := setupTestRouter(t)
		defer router.DB.Teardown(context.Background())

		user := setupTestUser(router)
		response := loginTestUser(t, router, user.Email, user.Password)
		fmt.Println(response.Body.String())

		assert.Equal(t, http.StatusOK, response.Code)
		data := getResponseData(t, response)
		assert.NotEmpty(t, data["accessToken"], "Access token should not be empty")
		assert.NotEmpty(t, data["refreshToken"], "Refresh token should not be empty")
	})

	t.Run("should fail with invalid password", func(t *testing.T) {

		router := setupTestRouter(t)
		defer router.DB.Teardown(context.Background())

		user := setupTestUser(router)
		response := loginTestUser(t, router, user.Email, "wrongpassword")

		assert.Equal(t, http.StatusUnauthorized, response.Code)

		var responseBody map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &responseBody)
		require.NoError(t, err, "Failed to parse error response")
		assert.Contains(t, responseBody, "error", "Error message should be present in response")
	})

	t.Run("should fail with non-existent email", func(t *testing.T) {
		router := setupTestRouter(t)
		defer router.DB.Teardown(context.Background())

		nonExistentEmail := "nonexistent_" + time.Now().Format("20060102150405") + "@example.com"
		response := loginTestUser(t, router, nonExistentEmail, "password123")

		assert.Equal(t, http.StatusUnauthorized, response.Code)

		var responseBody map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &responseBody)
		require.NoError(t, err, "Failed to parse error response")
		assert.Contains(t, responseBody, "error", "Error message should be present in response")
	})
}
