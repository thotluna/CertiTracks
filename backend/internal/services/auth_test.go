package services_test

import (
	"testing"
	"time"

	"certitrack/internal/config"
	"certitrack/internal/models"
	"certitrack/internal/repositories"
	"certitrack/internal/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// helper to build minimal Config for tests
func testConfig() *config.Config {
	return &config.Config{
		JWT: config.JWTConfig{
			Secret:             "abcdefghijklmnopqrstuvwxyz123456", // 32+ chars
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 24 * time.Hour,
			Issuer:             "test-suite",
			Audience:           "test-client",
		},
	}
}

func newAuthService() (services.AuthService, *repositories.MockUserRepository) {
	repo := repositories.NewMockUserRepository()
	svc := services.NewAuthService(testConfig(), repo)
	return svc, repo
}

func TestRegister_Success(t *testing.T) {
	svc, repo := newAuthService()

	req := &services.RegisterRequest{
		Email:     "alice@example.com",
		Password:  "s3cr3tPwd",
		FirstName: "Alice",
		LastName:  "Wonder",
	}

	resp, err := svc.Register(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, req.Email, resp.User.Email)
	require.True(t, repo.EmailExists(req.Email))
	require.NotEmpty(t, resp.AccessToken)
	require.NotEmpty(t, resp.RefreshToken)
}

func TestRegister_EmailExists(t *testing.T) {
	svc, repo := newAuthService()
	// Seed existing user
	_ = repo.CreateUser(&models.User{Email: "taken@mail.com", Password: "hash", IsActive: true})

	_, err := svc.Register(&services.RegisterRequest{
		Email:     "taken@mail.com",
		Password:  "pwd",
		FirstName: "T",
		LastName:  "K",
	})
	require.ErrorIs(t, err, services.ErrUserExists)
}

func TestLogin_Success(t *testing.T) {
	svc, _ := newAuthService()
	// Create user via register to ensure hashed password stored
	_, _ = svc.Register(&services.RegisterRequest{
		Email:     "bob@mail.com",
		Password:  "pw12345",
		FirstName: "Bob",
		LastName:  "Builder",
	})

	resp, err := svc.Login(&services.LoginRequest{Email: "bob@mail.com", Password: "pw12345"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "bob@mail.com", resp.User.Email)
}

func TestLogin_InvalidPassword(t *testing.T) {
	svc, _ := newAuthService()
	_, _ = svc.Register(&services.RegisterRequest{Email: "eve@mail.com", Password: "correct", FirstName: "Eve", LastName: "H"})

	_, err := svc.Login(&services.LoginRequest{Email: "eve@mail.com", Password: "wrong"})
	require.ErrorIs(t, err, services.ErrInvalidCredentials)
}

func TestLogin_UserNotFound(t *testing.T) {
	svc, _ := newAuthService()

	_, err := svc.Login(&services.LoginRequest{Email: "ghost@mail.com", Password: "pwd"})
	require.ErrorIs(t, err, services.ErrUserNotFound)
}

func TestRefreshToken_Success(t *testing.T) {
	svc, _ := newAuthService()
	regResp, _ := svc.Register(&services.RegisterRequest{
		Email:     "refresh@test.com",
		Password:  "pass123",
		FirstName: "Refresh",
		LastName:  "Test",
	})

	refreshReq := &services.RefreshRequest{RefreshToken: regResp.RefreshToken}
	refreshResp, err := svc.RefreshToken(refreshReq)

	require.NoError(t, err)
	require.NotEmpty(t, refreshResp.AccessToken)
	require.NotEmpty(t, refreshResp.RefreshToken)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	svc, _ := newAuthService()

	_, err := svc.RefreshToken(&services.RefreshRequest{RefreshToken: "invalid.token.here"})
	require.ErrorIs(t, err, services.ErrInvalidToken)
}

func TestRefreshToken_ExpiredToken(t *testing.T) {
	svc, _ := newAuthService()

	// Create a token that expired 1 hour ago
	cfg := testConfig()
	expiredToken := generateTestToken(t, cfg, "test-user-id", -time.Hour)

	_, err := svc.RefreshToken(&services.RefreshRequest{RefreshToken: expiredToken})
	require.ErrorIs(t, err, services.ErrTokenExpired)
}

func TestValidateToken_Valid(t *testing.T) {
	svc, _ := newAuthService()
	regResp, _ := svc.Register(&services.RegisterRequest{
		Email:     "validate@test.com",
		Password:  "pass123",
		FirstName: "Validate",
		LastName:  "Test",
	})

	claims, err := svc.ValidateAccessToken(regResp.AccessToken)
	require.NoError(t, err)
	require.Equal(t, regResp.User.ID.String(), claims.UserID)
}

func TestValidateToken_Invalid(t *testing.T) {
	svc, _ := newAuthService()

	_, err := svc.ValidateAccessToken("invalid.token.here")
	require.ErrorIs(t, err, services.ErrInvalidToken)
}

func TestGetUserFromToken_Success(t *testing.T) {
	svc, _ := newAuthService()
	regResp, _ := svc.Register(&services.RegisterRequest{
		Email:     "getuser@test.com",
		Password:  "pass123",
		FirstName: "Get",
		LastName:  "User",
	})

	user, err := svc.GetUserFromToken(regResp.AccessToken)
	require.NoError(t, err)
	require.Equal(t, "getuser@test.com", user.Email)
}

func TestGetUserFromToken_UserNotFound(t *testing.T) {
	svc, _ := newAuthService()

	// Create a valid token for a non-existent user
	cfg := testConfig()
	token := generateTestToken(t, cfg, uuid.New().String(), time.Hour)

	_, err := svc.GetUserFromToken(token)
	require.Error(t, err)
	require.Contains(t, err.Error(), "user not found")
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	svc, _ := newAuthService()

	// Create a token with a different secret
	badCfg := testConfig()
	badCfg.JWT.Secret = "different-secret-key-123"
	token := generateTestToken(t, badCfg, "test-user", time.Hour)

	_, err := svc.ValidateAccessToken(token)
	require.Error(t, err)
	require.Equal(t, services.ErrInvalidToken, err)
}

func TestValidateToken_MissingClaims(t *testing.T) {
	svc, _ := newAuthService()

	// Create a token with a different signing method that will be rejected
	token := jwt.New(jwt.SigningMethodRS256) // Using RS256 instead of HS256
	tokenString, _ := token.SignedString([]byte("some-rsa-key"))

	_, err := svc.ValidateAccessToken(tokenString)
	require.Error(t, err)
	require.Equal(t, services.ErrInvalidToken, err)
}

func TestValidateToken_WrongAudience(t *testing.T) {
	svc, _ := newAuthService()
	cfg := testConfig()

	// Create token with wrong audience
	claims := services.JWTClaims{
		UserID: "test-user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Issuer:    cfg.JWT.Issuer,
			Audience:  []string{"wrong-audience"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(cfg.JWT.Secret))

	// Audience is not validated in the current implementation
	validatedClaims, err := svc.ValidateAccessToken(tokenString)
	require.NoError(t, err)
	require.NotNil(t, validatedClaims)
	require.Equal(t, "test-user", validatedClaims.UserID)
}

func TestValidateToken_Expired(t *testing.T) {
	svc, _ := newAuthService()
	cfg := testConfig()

	// Create expired token
	claims := services.JWTClaims{
		UserID: "test-user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			Issuer:    cfg.JWT.Issuer,
			Audience:  []string{cfg.JWT.Audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(cfg.JWT.Secret))

	_, err := svc.ValidateAccessToken(tokenString)
	require.Error(t, err)
	require.Equal(t, services.ErrTokenExpired, err)
}

// generateTestToken creates a signed JWT for testing purposes
func generateTestToken(t *testing.T, cfg *config.Config, userID string, expiry time.Duration) string {
	t.Helper()
	claims := services.JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			Issuer:    cfg.JWT.Issuer,
			Audience:  []string{cfg.JWT.Audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(cfg.JWT.Secret))
	require.NoError(t, err)
	return signedToken
}
