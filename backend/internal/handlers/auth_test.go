package handlers_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"certitrack/internal/handlers"
	"certitrack/internal/services"
	"certitrack/internal/validators"
	"certitrack/testutils"
	"certitrack/testutils/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestRouter(handler *handlers.AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	if err := validators.RegisterAll(); err != nil {
		log.Fatal("Failed to register validators:", err)
	}

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.RefreshToken)
		auth.POST("/logout", handler.Logout)
	}
	return r
}

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	handler := handlers.NewAuthHandler(mockService)

	expectedResponse := &services.AuthResponse{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	mockService.On("RefreshToken", mock.AnythingOfType("*services.RefreshRequest")).
		Return(expectedResponse, nil)

	r := setupTestRouter(handler)
	requestBody := `{"refresh_token":"old-refresh-token"}`

	req, _ := http.NewRequest("POST", "/api/auth/refresh", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var refreshResp map[string]interface{}
	errRefresh := json.Unmarshal(w.Body.Bytes(), &refreshResp)
	assert.NoError(t, errRefresh)
	assert.Equal(t, "Token refreshed successfully", refreshResp["message"])
	mockService.AssertExpectations(t)
}

func TestAuthHandler_RefreshToken_Expired(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	handler := handlers.NewAuthHandler(mockService)

	mockService.On("RefreshToken", mock.AnythingOfType("*services.RefreshRequest")).
		Return(nil, services.ErrTokenExpired)

	r := setupTestRouter(handler)
	requestBody := `{"refresh_token":"expired-token"}`

	req, _ := http.NewRequest("POST", "/api/auth/refresh", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid or expired refresh token", response["error"])
	mockService.AssertExpectations(t)
}

func TestAuthHandler_RefreshToken_Invalid(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	handler := handlers.NewAuthHandler(mockService)

	mockService.On("RefreshToken", mock.AnythingOfType("*services.RefreshRequest")).
		Return(nil, services.ErrInvalidToken)

	r := setupTestRouter(handler)
	requestBody := `{"refresh_token":"invalid-token"}`

	req, _ := http.NewRequest("POST", "/api/auth/refresh", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid or expired refresh token", response["error"])
	mockService.AssertExpectations(t)
}

func TestAuthHandler_InternalServerError(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	handler := handlers.NewAuthHandler(mockService)

	mockService.On("Login", mock.AnythingOfType("*services.LoginRequest")).
		Return(nil, assert.AnError)
	reqBuilder := testutils.NewRegisterRequest()

	r := setupTestRouter(handler)
	requestBody := reqBuilder.ToJSON()

	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Login failed", response["error"])
	mockService.AssertExpectations(t)
}

func TestAuthHandler_InvalidJSON(t *testing.T) {
	handler := handlers.NewAuthHandler(new(mocks.MockAuthService))
	r := setupTestRouter(handler)

	tests := []struct {
		name        string
		url         string
		method      string
		requestBody string
	}{
		{
			name:        "invalid json register",
			url:         "/api/auth/register",
			method:      "POST",
			requestBody: `{"email":"test@example.com",`,
		},
		{
			name:        "invalid json login",
			url:         "/api/auth/login",
			method:      "POST",
			requestBody: `{"email":"test@example.com",`,
		},
		{
			name:        "invalid json refresh",
			url:         "/api/auth/refresh",
			method:      "POST",
			requestBody: `{"refresh_token":"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, tc.url, bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Contains(t, response["error"], "Invalid request data")
		})
	}
}
