package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"certitrack/internal/config"
	"certitrack/internal/database"
	"certitrack/internal/handlers"
	"certitrack/internal/middleware"
	"certitrack/internal/repositories"
	"certitrack/internal/services"

	"github.com/gin-gonic/gin"
)

type testUser struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func setupTestRouter() *gin.Engine {
	os.Setenv("APP_ENV", "test")
	cfg, _ := config.Load()
	db, _ := database.Connect(cfg)
	_ = database.AutoMigrate(db)

	userRepo := repositories.NewUserRepositoryImpl(db)
	authService := services.NewAuthService(cfg, userRepo)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", handlers.NewAuthHandler(authService).Register)
		authGroup.POST("/login", handlers.NewAuthHandler(authService).Login)
		authGroup.POST("/refresh", handlers.NewAuthHandler(authService).RefreshToken)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(authService))
	{
		api.GET("/me", func(c *gin.Context) {
			user, _ := c.Get("user")
			c.JSON(http.StatusOK, user)
		})
	}

	return r
}

func registerTestUser(_ *testing.T, router *gin.Engine, user testUser) *httptest.ResponseRecorder {
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func loginTestUser(_ *testing.T, router *gin.Engine, email, password string) *httptest.ResponseRecorder {
	loginData := map[string]string{
		"email":    email,
		"password": password,
	}

	body, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func getResponseData(_ *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	return response["data"].(map[string]interface{})
}
