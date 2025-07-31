package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"certitrack/internal/config"
	"certitrack/internal/handlers"
	"certitrack/internal/middleware"
	"certitrack/internal/repositories"
	"certitrack/internal/services"
	"certitrack/internal/validators"
	"certitrack/testutils/testcontainer"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type testUser struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type testRouter struct {
	Router      *gin.Engine
	DB          *testcontainer.PostgresContainer
	Middleware  *middleware.Middleware
	AuthHandler *handlers.AuthHandler
}

func setupTestRouter(t *testing.T) *testRouter {
	ctx := context.Background()

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			Name:     "certitrack_test",
			User:     "testuser",
			Password: "testpassword",
			SSLMode:  "disable",
		},
		JWT: config.JWTConfig{
			Secret:             "test_secret_key_must_be_at_least_32_chars_long_123",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 24 * time.Hour,
			Issuer:             "certitrack-test",
			Audience:           "certitrack-test-client",
		},
	}

	pgContainer, err := testcontainer.SetupPostgres(ctx, cfg)
	require.NoError(t, err, "Failed to setup postgres container")

	userRepo := repositories.NewUserRepositoryImpl(pgContainer.DB)
	authService := services.NewAuthService(cfg, userRepo)
	middlewareInstance := middleware.NewMiddleware(authService)
	authHandler := handlers.NewAuthHandler(authService)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	if err := validators.RegisterAll(); err != nil {
		log.Fatal("Failed to register validators:", err)
	}

	authGroup := r.Group("/api/auth")
	{
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

	return &testRouter{
		Router:      r,
		DB:          pgContainer,
		Middleware:  middlewareInstance,
		AuthHandler: authHandler,
	}
}

func registerTestUser(t *testing.T, router *testRouter, user testUser) *httptest.ResponseRecorder {
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.Router.ServeHTTP(w, req)
	t.Cleanup(func() {
		db, _ := router.DB.DB.DB()
		db.Exec("TRUNCATE TABLE users CASCADE")
	})
	return w
}

func loginTestUser(_ *testing.T, router *testRouter, email, password string) *httptest.ResponseRecorder {
	loginData := map[string]string{
		"email":    email,
		"password": password,
	}

	body, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.Router.ServeHTTP(w, req)
	return w
}

func getResponseData(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to parse response body")

	if data, ok := response["data"]; ok {
		return data.(map[string]interface{})
	}
	return response
}

func getAccessToken(t *testing.T, router *testRouter) string {
	user := testUser{
		Email:     "test@example.com",
		Password:  "Password123!",
		FirstName: "Test",
		LastName:  "User",
	}

	registerTestUser(t, router, user)

	w := loginTestUser(t, router, user.Email, user.Password)
	require.Equal(t, http.StatusOK, w.Code, "Login should succeed")

	responseData := getResponseData(t, w)
	accessToken, ok := responseData["accessToken"].(string)
	require.True(t, ok, "Access token should be present in response")

	return accessToken
}

func getRefreshToken(t *testing.T, router *testRouter) string {
	user := testUser{
		Email:     "refresh_test@example.com",
		Password:  "Password123!",
		FirstName: "Refresh",
		LastName:  "Test",
	}

	registerTestUser(t, router, user)

	w := loginTestUser(t, router, user.Email, user.Password)
	require.Equal(t, http.StatusOK, w.Code, "Login should succeed")

	responseData := getResponseData(t, w)
	refreshToken, ok := responseData["refreshToken"].(string)
	require.True(t, ok, "Refresh token should be present in response")

	return refreshToken
}
