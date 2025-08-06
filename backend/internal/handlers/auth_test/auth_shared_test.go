package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_InvalidJSON(t *testing.T) {

	tests := []struct {
		name        string
		url         string
		method      string
		requestBody string
	}{
		{
			name:        "invalid json register",
			url:         registerPath,
			method:      "POST",
			requestBody: `{"email":"test@example.com",`,
		},
		{
			name:        "invalid json login",
			url:         loginPath,
			method:      "POST",
			requestBody: `{"email":"test@example.com",`,
		},
		{
			name:        "invalid json refresh",
			url:         refreshTokenPath,
			method:      "POST",
			requestBody: `{"refresh_token":"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			setupTest(t)
			w := performRequest(tc.requestBody, tc.url, "")
			response := assertRegisterResponse(t, w, http.StatusBadRequest)
			assert.Contains(t, response["error"], "Invalid request data")
		})
	}
}
