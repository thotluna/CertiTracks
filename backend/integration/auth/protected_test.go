package auth_test

import (
	"certitrack/testutils"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProtectedRoutes(t *testing.T) {

	t.Run("should access protected route with valid token", func(t *testing.T) {
		rtr := setupTestRouter(t)
		defer rtr.DB.Teardown(context.Background())

		user := testutils.NewRegisterRequest().RegisterRequest
		registerTestUser(t, rtr, user)
		loginTestUser(t, rtr, user.Email, user.Password)
		token := getTokens(t, rtr)

		req, _ := http.NewRequest("GET", "/api/me", nil)
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)

		w := httptest.NewRecorder()
		rtr.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK for valid token")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse response body")
		assert.NotEmpty(t, response["email"], "User email should be present in response")
		assert.NotEmpty(t, response["firstName"], "User first name should be present in response")
		assert.NotEmpty(t, response["lastName"], "User last name should be present in response")
	})

	t.Run("should reject request without token", func(t *testing.T) {
		rtr := setupTestRouter(t)
		defer rtr.DB.Teardown(context.Background())

		req, _ := http.NewRequest("GET", "/api/me", nil)

		w := httptest.NewRecorder()
		rtr.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Should return 401 Unauthorized without token")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse error response")
		assert.Contains(t, response, "error", "Error message should be present in response")
	})

	t.Run("should reject request with invalid token", func(t *testing.T) {
		rtr := setupTestRouter(t)
		defer rtr.DB.Teardown(context.Background())

		req, _ := http.NewRequest("GET", "/api/me", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")

		w := httptest.NewRecorder()
		rtr.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Should return 401 Unauthorized with invalid token")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse error response")
		assert.Contains(t, response, "error", "Error message should be present in response")
	})
}
