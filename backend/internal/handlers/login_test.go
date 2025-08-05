package handlers_test

import (
	"certitrack/internal/services"
	"certitrack/testutils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const loginPath = "/api/v1/auth/login"

func TestAuthHandler_Login_Success(t *testing.T) {
	setupTest(t)

	mockAuthSvc.On("Login", mock.AnythingOfType("*services.LoginRequest")).
		Return(expectedResponse, nil)

	reqBuilder := testutils.NewRegisterRequest()
	requestBody := reqBuilder.ToJSON()

	w := performRegisterRequest(requestBody, loginPath, "")
	response := assertRegisterResponse(t, w, http.StatusOK)

	assert.Equal(t, "Login successful", response["message"])
	mockAuthSvc.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	setupTest(t)

	mockAuthSvc.On("Login", mock.AnythingOfType("*services.LoginRequest")).
		Return(nil, services.ErrInvalidCredentials)

	requestBody := `{"email":"test@example.com","password":"WrongPass123!"}`

	w := performRegisterRequest(requestBody, loginPath, "")
	response := assertRegisterResponse(t, w, http.StatusUnauthorized)

	assert.Equal(t, "Invalid email or password", response["error"])
	mockAuthSvc.AssertExpectations(t)
}

func TestAuthHandler_Login_UserNotFound(t *testing.T) {
	setupTest(t)

	mockAuthSvc.On("Login", mock.AnythingOfType("*services.LoginRequest")).
		Return(nil, services.ErrUserNotFound)
	reqBuilder := testutils.NewRegisterRequest()

	requestBody := reqBuilder.ToJSON()
	w := performRegisterRequest(requestBody, loginPath, "")

	response := assertRegisterResponse(t, w, http.StatusUnauthorized)
	assert.Equal(t, "Invalid email or password", response["error"])

	mockAuthSvc.AssertExpectations(t)
}
