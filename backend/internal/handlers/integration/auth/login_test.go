package auth_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	setupTestUser := func(router *gin.Engine) testUser {
		user := testUser{
			Email:     "login_test@example.com",
			Password:  "password123",
			FirstName: "Login",
			LastName:  "Test",
		}
		registerTestUser(t, router, user)
		return user
	}

	t.Run("should login with valid credentials", func(t *testing.T) {
		router := setupTestRouter()
		user := setupTestUser(router)

		response := loginTestUser(t, router, user.Email, user.Password)

		assert.Equal(t, http.StatusOK, response.Code)
		data := getResponseData(t, response)
		assert.NotEmpty(t, data["accessToken"])
		assert.NotEmpty(t, data["refreshToken"])
	})

	t.Run("should fail with invalid password", func(t *testing.T) {
		router := setupTestRouter()
		user := setupTestUser(router)

		response := loginTestUser(t, router, user.Email, "wrongpassword")

		assert.Equal(t, http.StatusUnauthorized, response.Code)
	})

	t.Run("should fail with non-existent email", func(t *testing.T) {
		router := setupTestRouter()
		// Asegurarse de que el correo no existe
	nonExistentEmail := "nonexistent_" + time.Now().Format("20060102150405") + "@example.com"
		response := loginTestUser(t, router, nonExistentEmail, "password123")

		assert.Equal(t, http.StatusUnauthorized, response.Code)
		
		var responseBody map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Contains(t, responseBody, "error")
	})
}
