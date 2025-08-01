package auth_test

import (
	"bytes"
	"certitrack/testutils"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshToken(t *testing.T) {

	t.Run("should refresh token successfully", func(t *testing.T) {
		rtr := setupTestRouter(t)
		defer rtr.DB.Teardown(context.Background())

		user := testutils.NewRegisterRequest().RegisterRequest
		registerResponse := registerTestUser(t, rtr, user)

		assert.Equal(t, http.StatusCreated, registerResponse.Code, "User registration should succeed")

		loginResponse := loginTestUser(t, rtr, user.Email, user.Password)
		if loginResponse.Code != http.StatusOK {
			t.Logf("Login failed with status: %d, body: %s", loginResponse.Code, loginResponse.Body.String())
		}
		require.Equal(t, http.StatusOK, loginResponse.Code, "Login should succeed")

		var loginResp map[string]interface{}
		err := json.Unmarshal(loginResponse.Body.Bytes(), &loginResp)

		require.NoError(t, err, "Failed to parse login response")

		loginDataResp, ok := loginResp["data"].(map[string]interface{})
		require.True(t, ok, "Login response should contain data object")

		refreshToken, ok := loginDataResp["refresh-token"].(string)
		require.True(t, ok, "Refresh token should be present in login response")
		require.NotEmpty(t, refreshToken, "Refresh token should not be empty")

		refreshData := map[string]string{
			"refresh_token": refreshToken,
		}

		body, err := json.Marshal(refreshData)
		require.NoError(t, err, "Failed to marshal refresh data")

		req, err := http.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(body))
		require.NoError(t, err, "Failed to create request")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		rtr.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK for valid refresh token")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse response body")

		require.Contains(t, response, "data", "Response should contain 'data' field")
		data, ok := response["data"].(map[string]interface{})
		require.True(t, ok, "Response data should be an object")
		assert.NotEmpty(t, data["access-token"], "New access token should be present in response")
		assert.NotEmpty(t, data["refresh-token"], "New refresh token should be present in response")
		assert.NotEqual(t, refreshToken, data["refresh-token"], "New refresh token should be different from the old one")
	})

	t.Run("should fail with invalid refresh token", func(t *testing.T) {
		rtr := setupTestRouter(t)
		defer rtr.DB.Teardown(context.Background())

		refreshData := map[string]string{
			"refresh_token": "invalid.refresh.token",
		}

		body, _ := json.Marshal(refreshData)
		req, _ := http.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		rtr.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Should return 401 Unauthorized for invalid refresh token")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to parse error response")
		assert.Contains(t, response, "error", "Error message should be present in response")
	})
}
