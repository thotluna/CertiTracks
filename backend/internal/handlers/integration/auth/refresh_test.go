package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRefreshToken(t *testing.T) {
	getRefreshToken := func(router *gin.Engine) string {
		user := testUser{
			Email:     "refresh@example.com",
			Password:  "password123",
			FirstName: "Refresh",
			LastName:  "Token",
		}

		registerTestUser(t, router, user)
		response := loginTestUser(t, router, user.Email, user.Password)
		data := getResponseData(t, response)
		return data["refreshToken"].(string)
	}

	t.Run("should refresh token successfully", func(t *testing.T) {
		router := setupTestRouter()
		refreshToken := getRefreshToken(router)

		refreshData := map[string]string{
			"refreshToken": refreshToken,
		}

		body, _ := json.Marshal(refreshData)
		req, _ := http.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})

		assert.NotEmpty(t, data["accessToken"])
		assert.NotEmpty(t, data["refreshToken"])
	})

	t.Run("should fail with invalid refresh token", func(t *testing.T) {
		router := setupTestRouter()

		refreshData := map[string]string{
			"refreshToken": "invalid.refresh.token",
		}

		body, _ := json.Marshal(refreshData)
		req, _ := http.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
