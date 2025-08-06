package handlers_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"certitrack/internal/handlers"
	"certitrack/internal/middleware"
	"certitrack/internal/models"
	"certitrack/internal/services"
	"certitrack/internal/validators"
	"certitrack/testutils"
	"certitrack/testutils/mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	testRouter     *gin.Engine
	mockAuthSvc    *mocks.MockAuthService
	testMiddleware *middleware.Middleware
)

var (
	expectedUser = &models.User{
		ID:        uuid.New(),
		Email:     testutils.NewRegisterRequest().Email,
		FirstName: testutils.NewRegisterRequest().FirstName,
		LastName:  testutils.NewRegisterRequest().LastName,
	}
	expectedResponse = &services.AuthResponse{
		User:         expectedUser,
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour),
	}
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
			auth.POST("/refresh", handler.RefreshToken)
		}

		protected := v1.Group("")
		protected.Use(testMiddleware.AuthMiddleware())
		{
			auth := protected.Group("/auth")
			{

				auth.POST("/logout", handler.Logout)
			}
		}
	}

	return r
}

func performRequest(body string, path string, token string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", path, bytes.NewBufferString(body))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func assertRegisterResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) map[string]interface{} {
	assert.Equal(t, expectedStatus, w.Code)
	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	return response
}
