package services_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"certitrack/internal/config"
	"certitrack/internal/models"
	"certitrack/internal/repositories"
	"certitrack/internal/services"
	"certitrack/testutils"
	"certitrack/testutils/mocks"
)

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

type testAuthService struct {
	services.AuthService
	mockClient *mocks.MockRedisClient
	repo       *repositories.MockUserRepository
}

func newAuthService() (*testAuthService, *repositories.MockUserRepository) {
	mockClient := new(mocks.MockRedisClient)
	repo := repositories.NewMockUserRepository()
	tokenRepo := repositories.NewTokenRepository(mockClient)
	svc := services.NewAuthService(testConfig(), repo, tokenRepo)
	return &testAuthService{
		AuthService: svc,
		mockClient:  mockClient,
		repo:        repo,
	}, repo
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

func TestLogout_Success(t *testing.T) {
	testSvc, _ := newAuthService()
	service := testSvc.AuthService
	mockClient := testSvc.mockClient

	mockClient.On("Exists", mock.Anything, mock.Anything).
		Return(redis.NewIntResult(0, nil)).Once()

	mockClient.On("Set", mock.Anything, mock.Anything, "1", mock.Anything).
		Return(redis.NewStatusResult("OK", nil)).Once()

	mockClient.On("Exists", mock.Anything, mock.Anything).
		Return(redis.NewIntResult(1, nil)).Once()

	reqBuilder := testutils.NewRegisterRequest()
	_, _ = service.Register(&reqBuilder.RegisterRequest)
	resp, _ := service.Login(&services.LoginRequest{
		Email:    reqBuilder.Email,
		Password: reqBuilder.Password,
	})

	_, err := service.ValidateAccessToken(resp.AccessToken)
	require.NoError(t, err)

	logoutRes, err := service.Logout(&services.LogoutRequest{AccessToken: resp.AccessToken})
	require.NoError(t, err)
	require.Nil(t, logoutRes.User, "User is not nil")
	require.Empty(t, logoutRes.AccessToken, "Access Token is not empty")

	_, err = service.ValidateAccessToken(resp.AccessToken)
	require.Error(t, err, "token should be invalid after logout")

	mockClient.AssertExpectations(t)
}

func TestRefreshToken_Success(t *testing.T) {
	testSvc, _ := newAuthService()
	svc := testSvc.AuthService
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

	cfg := testConfig()
	expiredToken := generateTestToken(t, cfg, "test-user-id", -time.Hour)

	_, err := svc.RefreshToken(&services.RefreshRequest{RefreshToken: expiredToken})
	require.ErrorIs(t, err, services.ErrTokenExpired)
}

func TestValidateToken_Valid(t *testing.T) {
	svc, _ := newAuthService()
	mockClient := svc.mockClient
	mockClient.On("Exists", mock.Anything, mock.Anything).
		Return(redis.NewIntResult(0, nil)).Once()
	reqBuilder := testutils.NewRegisterRequest()
	regResp, _ := svc.Register(&reqBuilder.RegisterRequest)

	claims, err := svc.ValidateAccessToken(regResp.AccessToken)
	require.NoError(t, err)
	require.Equal(t, regResp.User.ID.String(), claims.UserID)
}

func TestValidateToken_Invalid(t *testing.T) {
	svc, _ := newAuthService()
	mockClient := svc.mockClient
	mockClient.On("Exists", mock.Anything, mock.Anything).
		Return(redis.NewIntResult(1, nil)).Once()

	_, err := svc.ValidateAccessToken("invalid.token.here")
	require.ErrorIs(t, err, services.ErrInvalidToken)
}

func TestGetUserFromToken_Success(t *testing.T) {
	svc, _ := newAuthService()
	mockClient := svc.mockClient
	mockClient.On("Exists", mock.Anything, mock.Anything).
		Return(redis.NewIntResult(0, nil)).Once()
	reqBuilder := testutils.NewRegisterRequest()
	regResp, _ := svc.Register(&reqBuilder.RegisterRequest)

	user, err := svc.GetUserFromToken(regResp.AccessToken)
	require.NoError(t, err)
	require.Equal(t, reqBuilder.Email, user.Email)
}

func TestGetUserFromToken_UserNotFound(t *testing.T) {
	svc, _ := newAuthService()
	mockClient := svc.mockClient
	mockClient.On("Exists", mock.Anything, mock.Anything).
		Return(redis.NewIntResult(0, nil)).Once()

	cfg := testConfig()
	token := generateTestToken(t, cfg, uuid.New().String(), time.Hour)

	_, err := svc.GetUserFromToken(token)
	require.Error(t, err)
	require.Contains(t, err.Error(), "user not found")
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, cfg *config.Config) string
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
			return token[:len(token)-2] + "xx"
		},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testSvc, _ := newAuthService()
			svc := testSvc.AuthService
			mockClient := testSvc.mockClient

			// Configurar el mock para la validación del token
			mockClient.On("Exists", mock.Anything, mock.Anything).
				Return(redis.NewIntResult(0, nil)).Once()

			token := tc.setup(t, testConfig())

			_, err := svc.ValidateAccessToken(token)
			require.Error(t, err)
			require.Equal(t, services.ErrInvalidToken, err, "expected invalid token error for case: %s", tc.name)

			// Verificar que se llamaron todas las expectativas
			mockClient.AssertExpectations(t)
		})
	}
}

func TestValidateToken_MissingClaims(t *testing.T) {
	testSvc, _ := newAuthService()
	svc := testSvc.AuthService
	mockClient := testSvc.mockClient

	// Configurar el mock para la validación del token
	mockClient.On("Exists", mock.Anything, mock.Anything).
		Return(redis.NewIntResult(0, nil)).Once()

	// Crear un token con claims pero sin audiencia
	claims := services.JWTClaims{
		UserID: "test-user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Issuer:    "test-issuer",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	cfg := testConfig()
	tokenString, _ := token.SignedString([]byte(cfg.JWT.Secret))

	_, err := svc.ValidateAccessToken(tokenString)
	require.Error(t, err)
	require.Equal(t, services.ErrInvalidAudience, err, "expected invalid audience error")

	// Verificar que se llamaron todas las expectativas
	mockClient.AssertExpectations(t)
}

func TestValidateToken_WrongAudience(t *testing.T) {
	testSvc, _ := newAuthService()
	svc := testSvc.AuthService
	mockClient := testSvc.mockClient

	// Configurar el mock para la validación del token
	mockClient.On("Exists", mock.Anything, mock.Anything).
		Return(redis.NewIntResult(0, nil)).Once()

	cfg := testConfig()

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

	_, err := svc.ValidateAccessToken(tokenString)
	require.Error(t, err)
	require.Equal(t, services.ErrInvalidAudience, err, "expected invalid audience error")

	// Verificar que se llamaron todas las expectativas
	mockClient.AssertExpectations(t)
}

func TestValidateToken_Expired(t *testing.T) {
	testSvc, _ := newAuthService()
	svc := testSvc.AuthService
	mockClient := testSvc.mockClient

	// Configurar el mock para la validación del token
	mockClient.On("Exists", mock.Anything, mock.Anything).
		Return(redis.NewIntResult(0, nil)).Maybe()

	cfg := testConfig()

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

	// Verificar que se llamaron todas las expectativas
	mockClient.AssertExpectations(t)
}

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
