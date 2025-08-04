package auth_test

import (
	"certitrack/testutils"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogout(t *testing.T) {

	t.Run("should logout", func(t *testing.T) {
		testRouter := setupTestRouter(t)
		defer testRouter.DB.Teardown(context.Background())

		authTokens := getValidAuthToken(t, testRouter)

		logoutResponse := performLogout(t, testRouter, authTokens.AccessToken)

		assert.Equal(t, http.StatusOK, logoutResponse.Code, "Should return 200 OK for logout")

		protectedEndpointResponse := verifyTokenIsInvalid(t, testRouter, authTokens.AccessToken)

		assert.Equal(t, http.StatusUnauthorized, protectedEndpointResponse.Code, "Should return 401 OK for invalid token")
	})

	t.Run("should return 401 when no token is provided", func(t *testing.T) {
		testRouter := setupTestRouter(t)
		defer testRouter.DB.Teardown(context.Background())

		req, _ := http.NewRequest("POST", "/api/auth/logout", nil)
		response := httptest.NewRecorder()
		testRouter.Router.ServeHTTP(response, req)

		assert.Equal(t, http.StatusUnauthorized, response.Code)
	})

	t.Run("should return 401 with invalid token", func(t *testing.T) {
		testRouter := setupTestRouter(t)
		defer testRouter.DB.Teardown(context.Background())

		invalidToken := "invalid.token.here"
		response := performLogout(t, testRouter, invalidToken)

		assert.Equal(t, http.StatusUnauthorized, response.Code)
	})

	t.Run("should handle multiple logout requests with same token", func(t *testing.T) {
		testRouter := setupTestRouter(t)
		defer testRouter.DB.Teardown(context.Background())

		authTokens := getValidAuthToken(t, testRouter)

		response1 := performLogout(t, testRouter, authTokens.AccessToken)
		assert.Equal(t, http.StatusOK, response1.Code)

		response2 := performLogout(t, testRouter, authTokens.AccessToken)
		assert.Equal(t, http.StatusUnauthorized, response2.Code)
	})
}

func getValidAuthToken(t *testing.T, rtr *testRouter) Tokens {
	t.Helper()
	user := testutils.NewRegisterRequest().RegisterRequest
	registerTestUser(t, rtr, user)
	return getTokens(t, rtr)
}

func verifyTokenIsInvalid(t *testing.T, rtr *testRouter, token string) *httptest.ResponseRecorder {
	t.Helper()
	req, _ := http.NewRequest("GET", "/api/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	rtr.Router.ServeHTTP(w, req)
	return w
}

func performLogout(t *testing.T, rtr *testRouter, token string) *httptest.ResponseRecorder {
	t.Helper()
	req, _ := http.NewRequest("POST", "/api/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	rtr.Router.ServeHTTP(w, req)
	return w
}
