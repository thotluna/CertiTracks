package handlers_test

import (
	"certitrack/internal/services"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const refreshTokenPath = "/api/v1/auth/refresh"

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	setupTest(t)

	expectedResponse := &services.AuthResponse{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	mockAuthSvc.On("RefreshToken", mock.AnythingOfType("*services.RefreshRequest")).
		Return(expectedResponse, nil)

	requestBody := `{"refresh_token":"old-refresh-token"}`

	w := performRequest(requestBody, refreshTokenPath, "")
	response := assertRegisterResponse(t, w, http.StatusOK)

	assert.Equal(t, "Token refreshed successfully", response["message"])
	mockAuthSvc.AssertExpectations(t)
}

func TestAuthHandler_RefreshToken_Expired(t *testing.T) {
	setupTest(t)

	mockAuthSvc.On("RefreshToken", mock.AnythingOfType("*services.RefreshRequest")).
		Return(nil, services.ErrTokenExpired)

	requestBody := `{"refresh_token": "expire-token"}`

	w := performRequest(requestBody, refreshTokenPath, "")
	response := assertRegisterResponse(t, w, http.StatusUnauthorized)

	assert.Equal(t, "Invalid or expired refresh token", response["error"])
	mockAuthSvc.AssertExpectations(t)
}

func TestAuthHandler_RefreshToken_Invalid(t *testing.T) {
	setupTest(t)

	mockAuthSvc.On("RefreshToken", mock.AnythingOfType("*services.RefreshRequest")).
		Return(nil, services.ErrInvalidToken)

	requestBody := `{"refresh_token":"invalid-token"}`

	w := performRequest(requestBody, refreshTokenPath, "")
	response := assertRegisterResponse(t, w, http.StatusUnauthorized)

	assert.Equal(t, "Invalid or expired refresh token", response["error"])
	mockAuthSvc.AssertExpectations(t)
}
