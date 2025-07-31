package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"certitrack/internal/models"
	"certitrack/internal/services"
	"certitrack/internal/validators"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

var _ services.AuthService = (*MockAuthService)(nil)

func (m *MockAuthService) Register(req *services.RegisterRequest) (*services.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Login(req *services.LoginRequest) (*services.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.AuthResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(req *services.RefreshRequest) (*services.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.AuthResponse), args.Error(1)
}

func (m *MockAuthService) GetUserFromToken(token string) (*models.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) ValidateAccessToken(token string) (*services.JWTClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.JWTClaims), args.Error(1)
}

func (m *MockAuthService) ValidateRefreshToken(token string) (*services.JWTClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.JWTClaims), args.Error(1)
}

func (m *MockAuthService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) CheckPassword(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}

func setupTestRouter(handler *AuthHandler) *gin.Engine {
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
	}
	return r
}

func TestAuthHandler_Register_Success(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	expectedUser := &models.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}

	expectedResponse := &services.AuthResponse{
		User:         expectedUser,
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	mockService.On("Register", mock.AnythingOfType("*services.RegisterRequest")).Return(expectedResponse, nil)

	r := setupTestRouter(handler)
	requestBody := `{"email":"test@example.com","password":"Password123!","first_name":"Test","last_name":"User"}`

	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	fmt.Println("Request:", req)
	fmt.Println("Response:", w)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "User registered successfully", response["message"])
	data := response["data"].(map[string]interface{})
	userData := data["user"].(map[string]interface{})
	assert.Equal(t, expectedUser.Email, userData["email"])
	assert.Equal(t, expectedResponse.AccessToken, data["accessToken"])

	mockService.AssertExpectations(t)
}

func TestAuthHandler_Register_EmailExists(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	mockService.On("Register", mock.AnythingOfType("*services.RegisterRequest")).
		Return(nil, services.ErrUserExists)

	r := setupTestRouter(handler)
	requestBody := `{"email":"exists@example.com","password":"Password123!","first_name":"Test","last_name":"User"}`

	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "User with this email already exists", resp["error"])
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	expectedUser := &models.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}
	expectedResponse := &services.AuthResponse{
		User:         expectedUser,
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	mockService.On("Login", mock.AnythingOfType("*services.LoginRequest")).
		Return(expectedResponse, nil)

	r := setupTestRouter(handler)
	requestBody := `{"email":"test@example.com","password":"Password123!"}`

	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var loginResp map[string]interface{}
	errLogin := json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.NoError(t, errLogin)
	assert.Equal(t, "Login successful", loginResp["message"])
	mockService.AssertExpectations(t)
}

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	expectedResponse := &services.AuthResponse{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	mockService.On("RefreshToken", mock.AnythingOfType("*services.RefreshRequest")).
		Return(expectedResponse, nil)

	r := setupTestRouter(handler)
	requestBody := `{"refreshToken":"old-refresh-token"}`

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

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	mockService.On("Login", mock.AnythingOfType("*services.LoginRequest")).
		Return(nil, services.ErrInvalidCredentials)

	r := setupTestRouter(handler)
	requestBody := `{"email":"test@example.com","password":"WrongPass123!"}`

	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid email or password", response["error"])
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Login_UserNotFound(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	mockService.On("Login", mock.AnythingOfType("*services.LoginRequest")).
		Return(nil, services.ErrUserNotFound)

	r := setupTestRouter(handler)
	requestBody := `{"email":"nonexistent@example.com","password":"Password123!"}`

	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid email or password", response["error"])
	mockService.AssertExpectations(t)
}

func TestAuthHandler_RefreshToken_Expired(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	mockService.On("RefreshToken", mock.AnythingOfType("*services.RefreshRequest")).
		Return(nil, services.ErrTokenExpired)

	r := setupTestRouter(handler)
	requestBody := `{"refreshToken":"expired-token"}`

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
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	mockService.On("RefreshToken", mock.AnythingOfType("*services.RefreshRequest")).
		Return(nil, services.ErrInvalidToken)

	r := setupTestRouter(handler)
	requestBody := `{"refreshToken":"invalid-token"}`

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

func TestAuthHandler_Register_InvalidEmail(t *testing.T) {
	handler := NewAuthHandler(new(MockAuthService))
	r := setupTestRouter(handler)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "invalid email format",
			requestBody:    `{"email":"invalid-email","password":"Password123!","firstName":"Test","lastName":"User"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request data",
		},
		{
			name:           "missing email",
			requestBody:    `{"password":"Password123!","firstName":"Test","lastName":"User"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request data",
		},
		{
			name:           "password too short",
			requestBody:    `{"email":"test@example.com","password":"123","firstName":"Test","lastName":"User"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request data",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Contains(t, response["error"], tc.expectedError)
		})
	}
}

func TestAuthHandler_InternalServerError(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	mockService.On("Login", mock.AnythingOfType("*services.LoginRequest")).
		Return(nil, assert.AnError)

	r := setupTestRouter(handler)
	requestBody := `{"email":"test@example.com","password":"Password123!"}`

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
	handler := NewAuthHandler(new(MockAuthService))
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
			requestBody: `{"refreshToken":"`,
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
