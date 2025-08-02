// Package auth contains integration tests for authentication-related functionality.
// It tests the complete authentication flow including user registration,
// login, token generation, and protected endpoints.
package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"certitrack/integration"
	"certitrack/internal/cache/redis"
	"certitrack/internal/config"
	"certitrack/internal/handlers"
	"certitrack/internal/middleware"
	"certitrack/internal/repositories"
	"certitrack/internal/services"
	"certitrack/internal/validators"
	"certitrack/testutils"
	"certitrack/testutils/testcontainer"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

const (
	testDBHost     = "localhost"
	testDBPort     = "5432"
	testDBName     = "certitrack_test"
	testDBUser     = "testuser"
	testDBPassword = "testpassword"
	testDBSSLMode  = "disable"

	testJWTSecret             = "test_secret_key_must_be_at_least_32_chars_long_123"
	testJWTAccessTokenExpiry  = 15 * time.Minute
	testJWTRefreshTokenExpiry = 24 * time.Hour
	testJWTIssuer             = "certitrack-test"
	testJWTAudience           = "certitrack-test-client"
)

func GetTestConfig() *config.Config {
	return &config.Config{
		Database: config.DatabaseConfig{
			Host:     testDBHost,
			Port:     testDBPort,
			Name:     testDBName,
			User:     testDBUser,
			Password: testDBPassword,
			SSLMode:  testDBSSLMode,
		},
		JWT: config.JWTConfig{
			Secret:             testJWTSecret,
			AccessTokenExpiry:  testJWTAccessTokenExpiry,
			RefreshTokenExpiry: testJWTRefreshTokenExpiry,
			Issuer:             testJWTIssuer,
			Audience:           testJWTAudience,
		},
	}
}

type testRouter struct {
	Router      *gin.Engine
	DB          *testcontainer.PostgresContainer
	Middleware  *middleware.Middleware
	AuthHandler *handlers.AuthHandler
}

func GenerateUniqueEmail(prefix string) string {
	return fmt.Sprintf("%s_%d@example.com", prefix, time.Now().UnixNano())
}

func SetupTestDB(t *testing.T) (*testcontainer.PostgresContainer, func()) {
	t.Helper()
	ctx := context.Background()
	pgContainer, err := testcontainer.SetupPostgres(ctx, GetTestConfig())
	require.NoError(t, err, "Failed to setup postgres container")

	cleanup := func() {
		pgContainer.Teardown(ctx)
	}

	return pgContainer, cleanup
}

func setupTestRouter(t *testing.T) *testRouter {
	dbContainer, cleanup := SetupTestDB(t)
	redisContainer, cleanupRedis := integration.SetupTestDB(t)

	t.Cleanup(func() {
		cleanup()
		cleanupRedis()
		db, err := dbContainer.DB.DB()
		if err == nil {
			db.Exec("TRUNCATE TABLE users CASCADE")
		}
	})

	cfg := GetTestConfig()

	ctx := context.Background()
	endpoint, err := redisContainer.Container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get Redis host: %v", err)
	}

	port, err := redisContainer.Container.MappedPort(ctx, "6379/tcp")
	if err != nil {
		t.Fatalf("Failed to get Redis port: %v", err)
	}

	client, err := redis.NewClient(&config.RedisConfig{
		URL: "redis://" + endpoint + ":" + port.Port(),
	})

	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer client.Close()

	tokenRepo := repositories.NewTokenRepository(client)
	userRepo := repositories.NewUserRepositoryImpl(dbContainer.DB)
	authService := services.NewAuthService(cfg, userRepo, tokenRepo)
	middlewareSvc := middleware.NewMiddleware(authService)
	authHandler := handlers.NewAuthHandler(authService)

	router := gin.Default()
	if err := validators.RegisterAll(); err != nil {
		log.Fatal("Failed to register validators:", err)
	}

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		api.GET("/me", middlewareSvc.AuthMiddleware(), func(c *gin.Context) {
			user, _ := c.Get("user")
			c.JSON(http.StatusOK, user)
		})
	}

	return &testRouter{
		Router:      router,
		DB:          dbContainer,
		Middleware:  middlewareSvc,
		AuthHandler: authHandler,
	}
}

func newTestRequest(t *testing.T, method, path string, body interface{}) (*http.Request, *httptest.ResponseRecorder) {
	var buf bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		require.NoError(t, err, "Failed to encode request body")
	}

	req, err := http.NewRequest(method, path, &buf)
	require.NoError(t, err, "Failed to create request")

	req.Header.Set("Content-Type", "application/json")
	return req, httptest.NewRecorder()
}

func executeRequest(_ *testing.T, router *testRouter, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	router.Router.ServeHTTP(w, req)
	return w
}

func makeAuthRequest(t *testing.T, router *testRouter, path string, data interface{}) *httptest.ResponseRecorder {
	req, _ := newTestRequest(t, "POST", path, data)
	return executeRequest(t, router, req)
}

func registerTestUser(t *testing.T, router *testRouter, user services.RegisterRequest) *httptest.ResponseRecorder {
	return makeAuthRequest(t, router, "/api/auth/register", user)
}

func loginTestUser(t *testing.T, router *testRouter, email, password string) *httptest.ResponseRecorder {
	return makeAuthRequest(t, router, "/api/auth/login", map[string]string{
		"email":    email,
		"password": password,
	})
}

func getResponseData(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Logf("Failed to parse response body: %v, body: %s", err, w.Body.String())
	}
	require.NoError(t, err, "Failed to parse response body")

	if data, ok := response["data"]; ok {
		if dataMap, ok := data.(map[string]interface{}); ok {
			return dataMap
		}
	}

	return response
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func getTokens(t *testing.T, router *testRouter) Tokens {
	user := testutils.NewRegisterRequest().RegisterRequest

	registerResp := registerTestUser(t, router, user)
	require.Equal(t, http.StatusCreated, registerResp.Code, "User registration should succeed")

	w := loginTestUser(t, router, user.Email, user.Password)
	if w.Code != http.StatusOK {
		t.Logf("Login failed with status: %d, body: %s", w.Code, w.Body.String())
	}
	require.Equal(t, http.StatusOK, w.Code, "Login should succeed")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Logf("Failed to parse login response: %v, body: %s", err, w.Body.String())
	}
	require.NoError(t, err, "Failed to parse login response")

	require.Contains(t, response, "data", "Response should contain 'data' field")
	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok, "Response data should be an object")

	accessToken, ok1 := data["access-token"].(string)
	if !ok1 {
		t.Logf("Access token not found in response: %+v", data)
	}
	require.True(t, ok1, "Access token should be present in response")

	refreshToken, ok2 := data["refresh-token"].(string)
	if !ok2 {
		t.Logf("Refresh token not found in response: %+v", data)
	}
	require.True(t, ok2, "Refresh token should be present in response")

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
