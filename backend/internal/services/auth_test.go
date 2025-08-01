package services_test

import (
	"testing"
	"time"

	"certitrack/internal/config"
	"certitrack/internal/models"
	"certitrack/internal/repositories"
	"certitrack/internal/services"
	"certitrack/testutils"

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
	reqBuilder := testutils.NewRegisterRequest()

	req := &services.RegisterRequest{
		Email:     reqBuilder.Email,
		Password:  reqBuilder.Password,
		FirstName: reqBuilder.FirstName,
		LastName:  reqBuilder.LastName,
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
	reqBuilder := testutils.NewRegisterRequest()
	// Seed existing user
	_ = repo.CreateUser(&models.User{Email: reqBuilder.Email, Password: reqBuilder.Password, IsActive: true})

	_, err := svc.Register(&reqBuilder.RegisterRequest)
	require.ErrorIs(t, err, services.ErrUserExists)
}

func TestLogin_Success(t *testing.T) {
	svc, _ := newAuthService()
	reqBuilder := testutils.NewRegisterRequest()
	_, _ = svc.Register(&reqBuilder.RegisterRequest)

	resp, err := svc.Login(&services.LoginRequest{Email: reqBuilder.Email, Password: reqBuilder.Password})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, reqBuilder.Email, resp.User.Email)
}

func TestLogin_InvalidPassword(t *testing.T) {
	svc, _ := newAuthService()
	reqBuilder := testutils.NewRegisterRequest()
	_, _ = svc.Register(&reqBuilder.RegisterRequest)

	_, err := svc.Login(&services.LoginRequest{Email: reqBuilder.Email, Password: "wrong-password"})
	require.ErrorIs(t, err, services.ErrInvalidCredentials)
}

func TestLogin_UserNotFound(t *testing.T) {
	svc, _ := newAuthService()
	reqBuilder := testutils.NewRegisterRequest()

	_, err := svc.Login(&services.LoginRequest{Email: "nonexistent-" + reqBuilder.Email, Password: reqBuilder.Password})
	require.ErrorIs(t, err, services.ErrUserNotFound)
}

func TestRefreshToken_Success(t *testing.T) {
	svc, _ := newAuthService()
	reqBuilder := testutils.NewRegisterRequest()
	regResp, _ := svc.Register(&reqBuilder.RegisterRequest)

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
	reqBuilder := testutils.NewRegisterRequest()
	regResp, _ := svc.Register(&reqBuilder.RegisterRequest)

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
	reqBuilder := testutils.NewRegisterRequest()
	regResp, _ := svc.Register(&reqBuilder.RegisterRequest)

	user, err := svc.GetUserFromToken(regResp.AccessToken)
	require.NoError(t, err)
	require.Equal(t, reqBuilder.Email, user.Email)
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
	tests := []struct {
		name    string
		setup   func(t *testing.T, cfg *config.Config) string
	}{{
		name: "different secret",
		setup: func(t *testing.T, cfg *config.Config) string {
			badCfg := *cfg
			badCfg.JWT.Secret = "different-secret-key-123"
			return generateTestToken(t, &badCfg, "test-user", time.Hour)
		},
	}, {
		name: "none algorithm",
		setup: func(t *testing.T, cfg *config.Config) string {
			token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
				"sub": "test-user",
				"exp": time.Now().Add(time.Hour).Unix(),
			})
			tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
			require.NoError(t, err)
			return tokenString
		},
	}, {
		name: "empty token",
		setup: func(t *testing.T, _ *config.Config) string {
			return ""
		},
	}, {
		name: "malformed token",
		setup: func(t *testing.T, _ *config.Config) string {
			return "not.a.valid.jwt"
		},
	}, {
		name: "tampered token",
		setup: func(t *testing.T, cfg *config.Config) string {
			token := generateTestToken(t, cfg, "test-user", time.Hour)
			// Cambiar un carácter en el token para invalidar la firma
			return token[:len(token)-2] + "xx"
		},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, _ := newAuthService()
			token := tc.setup(t, testConfig())

			_, err := svc.ValidateAccessToken(token)
			require.Error(t, err)
			require.Equal(t, services.ErrInvalidToken, err, "expected invalid token error for case: %s", tc.name)
		})
	}
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
