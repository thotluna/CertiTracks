package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshToken(t *testing.T) {
	setupTestUser := func(router *testRouter) testUser {
		user := testUser{
			Email:     "refresh@example.com",
			Password:  "password123",
			FirstName: "Refresh",
			LastName:  "Token",
		}
		registerTestUser(t, router, user)
		return user
	}

	t.Run("should refresh token successfully", func(t *testing.T) {
		router := setupTestRouter(t)
		defer router.DB.Teardown(context.Background())

		setupTestUser(router)
		refreshToken := getRefreshToken(t, router)

		refreshData := map[string]string{
			"refreshToken": refreshToken,
		}

		body, _ := json.Marshal(refreshData)
		req, _ := http.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK for valid refresh token")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse response body")

		data, ok := response["data"].(map[string]interface{})
		require.True(t, ok, "Response should contain data object")

		assert.NotEmpty(t, data["accessToken"], "New access token should be present in response")
		assert.NotEmpty(t, data["refreshToken"], "New refresh token should be present in response")
		assert.NotEqual(t, refreshToken, data["refreshToken"], "New refresh token should be different from the old one")
	})

	t.Run("should fail with invalid refresh token", func(t *testing.T) {
		router := setupTestRouter(t)
		defer router.DB.Teardown(context.Background())

		refreshData := map[string]string{
			"refreshToken": "invalid.refresh.token",
		}

		body, _ := json.Marshal(refreshData)
		req, _ := http.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Should return 401 Unauthorized for invalid refresh token")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse error response")
		assert.Contains(t, response, "error", "Error message should be present in response")
	})
}
