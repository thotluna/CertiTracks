package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"certitrack/internal/handlers"
	"certitrack/internal/models"
	"certitrack/internal/services"
	"certitrack/internal/validators"
	"certitrack/testutils"
	"certitrack/testutils/mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func TestAuthHandler_Register_Success(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	handler := handlers.NewAuthHandler(mockService)

	reqBuilder := testutils.NewRegisterRequest()

	expectedUser := &models.User{
		ID:        uuid.New(),
		Email:     reqBuilder.Email,
		FirstName: reqBuilder.FirstName,
		LastName:  reqBuilder.LastName,
	}

	expectedResponse := &services.AuthResponse{
		User:         expectedUser,
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	mockService.On("Register", mock.AnythingOfType("*services.RegisterRequest")).Return(expectedResponse, nil)

	r := setupTestRouter(handler)
	requestBody := reqBuilder.ToJSON()

	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "User registered successfully", response["message"])
	data := response["data"].(map[string]interface{})
	userData := data["user"].(map[string]interface{})
	assert.Equal(t, expectedUser.Email, userData["email"])
	assert.Equal(t, expectedResponse.AccessToken, data["access-token"])

	mockService.AssertExpectations(t)
}

func TestAuthHandler_Register_EmailExists(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	handler := handlers.NewAuthHandler(mockService)
	reqBuilder := testutils.NewRegisterRequest()

	mockService.On("Register", mock.AnythingOfType("*services.RegisterRequest")).
		Return(nil, services.ErrUserExists)

	r := setupTestRouter(handler)
	requestBody := reqBuilder.ToJSON()

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
	mockService := new(mocks.MockAuthService)
	handler := handlers.NewAuthHandler(mockService)
	reqBuilder := testutils.NewRegisterRequest()

	expectedUser := &models.User{
		ID:        uuid.New(),
		Email:     reqBuilder.Email,
		FirstName: reqBuilder.FirstName,
		LastName:  reqBuilder.LastName,
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
	requestBody := reqBuilder.ToJSON()

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

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	handler := handlers.NewAuthHandler(mockService)

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
	mockService := new(mocks.MockAuthService)
	handler := handlers.NewAuthHandler(mockService)

	mockService.On("Login", mock.AnythingOfType("*services.LoginRequest")).
		Return(nil, services.ErrUserNotFound)
	reqBuilder := testutils.NewRegisterRequest()

	r := setupTestRouter(handler)
	requestBody := reqBuilder.ToJSON()

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

func TestAuthHandler_Logout_Success(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	handler := handlers.NewAuthHandler(mockService)

	expectedToken := "test-token"

	mockService.On("Logout", &services.LogoutRequest{
		AccessToken: expectedToken,
	}).Return(&services.AuthResponse{}, nil)

	route := setupTestRouter(handler)

	req, _ := http.NewRequest("POST", "/api/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+expectedToken)

	res := httptest.NewRecorder()
	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	assert.NotEmpty(t, res.Body.String())

	var response map[string]interface{}
	if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
		t.Fatalf("Could not parse response: %v. Body: %s", err, res.Body.String())
	}

	assert.Equal(t, "Logout successful", response["message"])
	assert.NotNil(t, response["data"])

	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "data should be an object")

	assert.Nil(t, data["user"])
	assert.Empty(t, data["access-token"])
	assert.Empty(t, data["refresh-token"])

	mockService.AssertExpectations(t)
}

func TestAuthHandler_Logout_Fail_server(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	handler := handlers.NewAuthHandler(mockService)

	testToken := "test-token"
	expectedError := errors.New("error revoking token")

	mockService.On("Logout", &services.LogoutRequest{
		AccessToken: testToken,
	}).Return(nil, expectedError)

	route := setupTestRouter(handler)

	req, _ := http.NewRequest("POST", "/api/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)

	res := httptest.NewRecorder()
	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)

	var response map[string]interface{}
	if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
		t.Fatalf("Could not parse response: %v", err)
	}

	assert.Equal(t, "Logout failed", response["error"])
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Register_InvalidEmail(t *testing.T) {
	handler := handlers.NewAuthHandler(new(mocks.MockAuthService))
	r := setupTestRouter(handler)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "invalid email format",
			requestBody:    testutils.NewRegisterRequest(testutils.WithEmail("invalid-email")).ToJSON(),
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request data",
		},
		{
			name:           "missing email",
			requestBody:    testutils.NewRegisterRequest(testutils.WithEmail("")).ToJSON(),
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request data",
		},
		{
			name:           "password too short",
			requestBody:    testutils.NewRegisterRequest(testutils.WithPassword("short")).ToJSON(),
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
