package handlers_test

import (
	"certitrack/internal/services"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const logoutPath = "/api/v1/auth/logout"

func TestAuthHandler_Logout_Success(t *testing.T) {
	setupTest(t)

	expectedToken := "test-token"

	mockAuthSvc.On("GetUserFromToken", "test-token").
		Return(expectedUser, nil)

	mockAuthSvc.On("Logout", &services.LogoutRequest{
		AccessToken: expectedToken,
	}).Return(&services.AuthResponse{}, nil)

	w := performRegisterRequest("", logoutPath, expectedToken)

	response := assertRegisterResponse(t, w, http.StatusOK)

	assert.Equal(t, "Logout successful", response["message"])
	assert.NotNil(t, response["data"])

	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "data should be an object")

	assert.Nil(t, data["user"])
	assert.Empty(t, data["access-token"])
	assert.Empty(t, data["refresh-token"])

	mockAuthSvc.AssertExpectations(t)
}

func TestAuthHandler_Logout_Fail_server(t *testing.T) {
	setupTest(t)

	expectedToken := "test-token"

	mockAuthSvc.On("GetUserFromToken", "test-token").
		Return(expectedUser, nil)

	expectedError := errors.New("error revoking token")
	mockAuthSvc.On("Logout", &services.LogoutRequest{
		AccessToken: expectedToken,
	}).Return(nil, expectedError)

	w := performRegisterRequest("", logoutPath, expectedToken)
	response := assertRegisterResponse(t, w, http.StatusInternalServerError)

	assert.Equal(t, "Logout failed", response["error"])
	mockAuthSvc.AssertExpectations(t)
}

func TestAuthHandler_Logout_with_invalid_token(t *testing.T) {
	setupTest(t)

	expectedToken := "test-token"

	mockAuthSvc.On("GetUserFromToken", "test-token").
		Return(nil, services.ErrInvalidToken)

	w := performRegisterRequest("", logoutPath, expectedToken)

	response := assertRegisterResponse(t, w, http.StatusUnauthorized)

	assert.Equal(t, "Invalid or expired token", response["error"])
	mockAuthSvc.AssertExpectations(t)
}
