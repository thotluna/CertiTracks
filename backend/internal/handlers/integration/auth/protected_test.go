package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProtectedRoutes(t *testing.T) {
	setupTestUser := func(router *testRouter) testUser {
		user := testUser{
			Email:     "protected@example.com",
			Password:  "password123",
			FirstName: "Protected",
			LastName:  "Route",
		}
		registerTestUser(t, router, user)
		return user
	}

	t.Run("should access protected route with valid token", func(t *testing.T) {
		router := setupTestRouter(t)
		defer router.DB.Teardown(context.Background())

		setupTestUser(router)
		token := getAccessToken(t, router)

		req, _ := http.NewRequest("GET", "/api/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK for valid token")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse response body")
		assert.NotEmpty(t, response["email"], "User email should be present in response")
		assert.NotEmpty(t, response["firstName"], "User first name should be present in response")
		assert.NotEmpty(t, response["lastName"], "User last name should be present in response")
	})

	t.Run("should reject request without token", func(t *testing.T) {
		router := setupTestRouter(t)
		defer router.DB.Teardown(context.Background())

		req, _ := http.NewRequest("GET", "/api/me", nil)

		w := httptest.NewRecorder()
		router.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Should return 401 Unauthorized without token")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse error response")
		assert.Contains(t, response, "error", "Error message should be present in response")
	})

	t.Run("should reject request with invalid token", func(t *testing.T) {
		router := setupTestRouter(t)
		defer router.DB.Teardown(context.Background())

		req, _ := http.NewRequest("GET", "/api/me", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")

		w := httptest.NewRecorder()
		router.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Should return 401 Unauthorized with invalid token")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse error response")
		assert.Contains(t, response, "error", "Error message should be present in response")
	})
}
