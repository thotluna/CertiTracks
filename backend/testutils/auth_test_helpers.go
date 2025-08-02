// Package testutils provides common test helpers and utilities for the application.
// It includes reusable test fixtures, mocks, and helper functions to reduce
// code duplication and improve test consistency across the codebase.
package testutils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"certitrack/internal/handlers"
	"certitrack/internal/middleware"
	"certitrack/internal/repositories"
	"certitrack/internal/services"
	"certitrack/internal/validators"
	"certitrack/testutils/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type TestUser struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func SetupTestRouter(t *testing.T, db *gorm.DB) *gin.Engine {
	t.Helper()
	mockClient := new(mocks.MockRedisClient)
	tokenRepo := repositories.NewTokenRepository(mockClient)

	userRepo := repositories.NewUserRepositoryImpl(db)
	cfg := GetTestConfig()
	authService := services.NewAuthService(cfg, userRepo, tokenRepo)
	middlewareInstance := middleware.NewMiddleware(authService)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	if err := validators.RegisterAll(); err != nil {
		log.Fatal("Failed to register validators:", err)
	}

	authGroup := r.Group("/api/auth")
	{
		authHandler := handlers.NewAuthHandler(authService)
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)
	}

	api := r.Group("/api")
	api.Use(middlewareInstance.AuthMiddleware())
	{
		api.GET("/me", func(c *gin.Context) {
			user, _ := c.Get("user")
			c.JSON(http.StatusOK, user)
		})
	}

	return r
}

func RegisterTestUser(t *testing.T, router *gin.Engine, user services.RegisterRequest) *httptest.ResponseRecorder {
	t.Helper()

	body, err := json.Marshal(user)
	require.NoError(t, err, "Failed to marshal user")

	req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	require.NoError(t, err, "Failed to create request")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

func LoginTestUser(t *testing.T, router *gin.Engine, email, password string) *httptest.ResponseRecorder {
	t.Helper()

	loginData := map[string]string{
		"email":    email,
		"password": password,
	}

	body, err := json.Marshal(loginData)
	require.NoError(t, err, "Failed to marshal login data")

	req, err := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	require.NoError(t, err, "Failed to create login request")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

func GetResponseData(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to decode JSON response")

	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok, "Response does not contain valid data")

	return data
}
