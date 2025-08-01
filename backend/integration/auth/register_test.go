package auth_test

import (
	"net/http"
	"testing"

	"certitrack/testutils"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	pgContainer, cleanup := SetupTestDB(t)
	defer cleanup()

	router := testutils.SetupTestRouter(t, pgContainer.DB)

	t.Run("should register a new user successfully", func(t *testing.T) {
		user := testutils.NewRegisterRequest(testutils.WithEmail(GenerateUniqueEmail("test_register"))).RegisterRequest

		response := testutils.RegisterTestUser(t, router, user)

		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("should fail with duplicate email", func(t *testing.T) {
		user := testutils.NewRegisterRequest(testutils.WithEmail(GenerateUniqueEmail("duplicate"))).RegisterRequest

		firstResponse := testutils.RegisterTestUser(t, router, user)
		assert.Equal(t, http.StatusCreated, firstResponse.Code)

		response := testutils.RegisterTestUser(t, router, user)

		assert.Equal(t, http.StatusConflict, response.Code)
	})

	t.Run("should fail with invalid email format", func(t *testing.T) {
		user := testutils.NewRegisterRequest(testutils.WithEmail("invalid-email-format")).RegisterRequest

		response := testutils.RegisterTestUser(t, router, user)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("should fail with short password", func(t *testing.T) {
		user := testutils.NewRegisterRequest(testutils.WithPassword("short")).RegisterRequest

		response := testutils.RegisterTestUser(t, router, user)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}
