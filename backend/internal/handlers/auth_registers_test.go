package handlers_test

import (
	"certitrack/internal/services"
	"certitrack/testutils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const registerPath = "/api/v1/auth/register"

func TestRegister_Success(t *testing.T) {
	setupTest(t)

	reqBuilder := testutils.NewRegisterRequest()

	mockAuthSvc.On("Register", mock.AnythingOfType("*services.RegisterRequest")).Return(expectedResponse, nil)

	w := performRegisterRequest(reqBuilder.ToJSON(), registerPath, "")
	response := assertRegisterResponse(t, w, http.StatusCreated)

	assert.Equal(t, "User registered successfully", response["message"])
	data := response["data"].(map[string]interface{})
	userData := data["user"].(map[string]interface{})
	assert.Equal(t, expectedUser.Email, userData["email"])
	assert.Equal(t, expectedResponse.AccessToken, data["access-token"])

	mockAuthSvc.AssertExpectations(t)
}

func TestRegister_EmailExists(t *testing.T) {
	setupTest(t)

	reqBuilder := testutils.NewRegisterRequest()
	mockAuthSvc.On("Register", mock.AnythingOfType("*services.RegisterRequest")).
		Return(nil, services.ErrUserExists)

	w := performRegisterRequest(reqBuilder.ToJSON(), registerPath, "")
	response := assertRegisterResponse(t, w, http.StatusConflict)

	assert.Equal(t, "User with this email already exists", response["error"])
	mockAuthSvc.AssertExpectations(t)
}

func TestRegister_InvalidInput(t *testing.T) {
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
			setupTest(t)
			w := performRegisterRequest(tc.requestBody, registerPath, "")
			response := assertRegisterResponse(t, w, tc.expectedStatus)
			assert.Contains(t, response["error"], tc.expectedError)
		})
	}
}

func assertRegisterResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) map[string]interface{} {
	assert.Equal(t, expectedStatus, w.Code)
	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	return response
}
