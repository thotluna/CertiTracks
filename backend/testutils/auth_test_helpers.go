package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"certitrack/internal/handlers"
	"certitrack/internal/middleware"
	"certitrack/internal/repositories"
	"certitrack/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type TestUser struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func SetupTestRouter(t *testing.T, db *gorm.DB) *gin.Engine {
	t.Helper()

	userRepo := repositories.NewUserRepositoryImpl(db)
	cfg := GetTestConfig()
	authService := services.NewAuthService(cfg, userRepo)
	middlewareInstance := middleware.NewMiddleware(authService)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

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

func RegisterTestUser(t *testing.T, router *gin.Engine, user TestUser) *httptest.ResponseRecorder {
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
