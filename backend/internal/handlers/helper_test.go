package handlers_test

import (
	"log"
	"os"
	"testing"

	"certitrack/internal/handlers"
	"certitrack/internal/middleware"
	"certitrack/internal/validators"
	"certitrack/testutils/mocks"

	"github.com/gin-gonic/gin"
)

var (
	testRouter     *gin.Engine
	mockAuthSvc    *mocks.MockAuthService
	testMiddleware *middleware.Middleware
)

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}

func setup() {
	mockAuthSvc = new(mocks.MockAuthService)

	testMiddleware = middleware.NewMiddleware(mockAuthSvc)

	testRouter = setupTestRoute(handlers.NewAuthHandler(mockAuthSvc))
}

func teardown() {
}

func setupTest(t *testing.T) {
	teardown()
	setup()
	t.Cleanup(teardown)
}

func setupTestRoute(handler *handlers.AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	if err := validators.RegisterAll(); err != nil {
		log.Fatal("Failed to register validators:", err)
	}

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
		}

		protected := v1.Group("")
		protected.Use(testMiddleware.AuthMiddleware())
		{
			auth := protected.Group("/auth")
			{
				auth.POST("/refresh", handler.RefreshToken)
				auth.POST("/logout", handler.Logout)
			}
		}
	}

	return r
}
