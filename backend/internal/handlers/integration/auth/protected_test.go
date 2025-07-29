package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestProtectedRoutes(t *testing.T) {
	getAccessToken := func(router *gin.Engine) string {
		user := testUser{
			Email:     "protected@example.com",
			Password:  "password123",
			FirstName: "Protected",
			LastName:  "Route",
		}

		registerTestUser(t, router, user)
		response := loginTestUser(t, router, user.Email, user.Password)
		data := getResponseData(t, response)
		return data["accessToken"].(string)
	}

	t.Run("should access protected route with valid token", func(t *testing.T) {
		router := setupTestRouter()
		token := getAccessToken(router)

		req, _ := http.NewRequest("GET", "/api/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should reject request without token", func(t *testing.T) {
		router := setupTestRouter()

		req, _ := http.NewRequest("GET", "/api/me", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should reject request with invalid token", func(t *testing.T) {
		router := setupTestRouter()

		req, _ := http.NewRequest("GET", "/api/me", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
