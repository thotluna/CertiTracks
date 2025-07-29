package auth_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	generateUniqueEmail := func(prefix string) string {
		return prefix + "_" + time.Now().Format("20060102150405") + "@example.com"
	}

	t.Run("should register a new user successfully", func(t *testing.T) {
		router := setupTestRouter()
		user := testUser{
			Email:     generateUniqueEmail("test_register"),
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
		}

		response := registerTestUser(t, router, user)

		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("should fail with duplicate email", func(t *testing.T) {
		router := setupTestRouter()
		user := testUser{
			Email:     generateUniqueEmail("duplicate"),
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
		}

		// Primer registro debe ser exitoso
		firstResponse := registerTestUser(t, router, user)
		assert.Equal(t, http.StatusCreated, firstResponse.Code)

		// Segundo registro con el mismo email debe fallar
		response := registerTestUser(t, router, user)

		assert.Equal(t, http.StatusConflict, response.Code)
	})

	t.Run("should fail with invalid email format", func(t *testing.T) {
		router := setupTestRouter()
		user := testUser{
			Email:     "invalid-email-format",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
		}

		response := registerTestUser(t, router, user)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("should fail with short password", func(t *testing.T) {
		router := setupTestRouter()
		user := testUser{
			Email:     generateUniqueEmail("shortpass"),
			Password:  "12345",
			FirstName: "Test",
			LastName:  "User",
		}

		response := registerTestUser(t, router, user)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}
