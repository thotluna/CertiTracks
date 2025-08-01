package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"certitrack/internal/middleware"
	"certitrack/internal/models"
	"certitrack/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Register(req *services.RegisterRequest) (*services.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.AuthResponse), args.Error(1)
}

func (m *mockAuthService) Login(req *services.LoginRequest) (*services.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.AuthResponse), args.Error(1)
}

func (m *mockAuthService) RefreshToken(req *services.RefreshRequest) (*services.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.AuthResponse), args.Error(1)
}

func (m *mockAuthService) GetUserFromToken(token string) (*models.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockAuthService) ValidateAccessToken(token string) (*services.JWTClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.JWTClaims), args.Error(1)
}

func (m *mockAuthService) ValidateRefreshToken(token string) (*services.JWTClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.JWTClaims), args.Error(1)
}

func (m *mockAuthService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *mockAuthService) CheckPassword(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}

func setupTestRouter(authService services.AuthService, mwFunc gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	mw := middleware.NewMiddleware(authService)

	authGroup := r.Group("")
	authGroup.Use(mwFunc)
	{
		authGroup.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		admin := authGroup.Group("/admin")
		admin.Use(mw.AdminMiddleware())
		{
			admin.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})
		}

		authGroup.GET("/optional", mw.OptionalAuthMiddleware(), func(c *gin.Context) {
			_, exists := c.Get("userID")
			if exists {
				c.Status(http.StatusOK)
			} else {
				c.Status(http.StatusNoContent)
			}
		})
	}

	return r
}

func TestAuthMiddleware(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name           string
		setupAuth      func() string
		setupMock      func(*mockAuthService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful authentication",
			setupAuth: func() string {
				return "valid.token.here"
			},
			setupMock: func(m *mockAuthService) {
				m.On("GetUserFromToken", "valid.token.here").
					Return(&models.User{
						ID:        userID,
						Email:     "test@example.com",
						Role:      "user",
						FirstName: "Test",
						LastName:  "User",
						IsActive:  true,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing authorization header",
			setupAuth:      func() string { return "" },
			setupMock:      func(*mockAuthService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Authorization header is required",
		},
		{
			name: "invalid token format",
			setupAuth: func() string {
				return "invalidtoken"
			},
			setupMock:      func(*mockAuthService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Authorization header must start with 'Bearer '",
		},
		{
			name: "empty token",
			setupAuth: func() string {
				return ""
			},
			setupMock:      func(*mockAuthService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Authorization header is required",
		},
		{
			name: "invalid token",
			setupAuth: func() string {
				return "invalid.token.here"
			},
			setupMock: func(m *mockAuthService) {
				m.On("GetUserFromToken", "invalid.token.here").
					Return(nil, services.ErrInvalidToken)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid or expired token",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			authService := new(mockAuthService)
			tc.setupMock(authService)

			mw := middleware.NewMiddleware(authService)
			r := setupTestRouter(authService, mw.AuthMiddleware())

			req := httptest.NewRequest("GET", "/test", nil)
			token := tc.setupAuth()
			if token != "" {
				if tc.name == "invalid token format" {
					req.Header.Set("Authorization", token)
				} else if tc.name == "empty token" {
					req.Header.Set("Authorization", "Bearer "+token)
				} else {
					req.Header.Set("Authorization", "Bearer "+token)
				}
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tc.expectedStatus, w.Code)
			if tc.expectedError != "" {
				require.Contains(t, w.Body.String(), tc.expectedError)
			}

			authService.AssertExpectations(t)
		})
	}
}

func TestAdminMiddleware(t *testing.T) {
	adminID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name           string
		authToken      string
		setupMock      func(*mockAuthService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "admin access",
			authToken: "admin.token",
			setupMock: func(m *mockAuthService) {
				m.On("GetUserFromToken", "admin.token").
					Return(&models.User{
						ID:        adminID,
						Email:     "admin@example.com",
						Role:      "admin",
						FirstName: "Admin",
						LastName:  "User",
						IsActive:  true,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "user access denied",
			authToken: "user.token",
			setupMock: func(m *mockAuthService) {
				m.On("GetUserFromToken", "user.token").
					Return(&models.User{
						ID:        userID,
						Email:     "user@example.com",
						Role:      "user",
						FirstName: "Regular",
						LastName:  "User",
						IsActive:  true,
					}, nil)
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  "Admin access required",
		},
		{
			name:           "missing token",
			authToken:      "",
			setupMock:      func(*mockAuthService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Authorization header is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			authService := new(mockAuthService)
			tc.setupMock(authService)

			mw := middleware.NewMiddleware(authService)
			r := setupTestRouter(authService, mw.AuthMiddleware())

			req := httptest.NewRequest("GET", "/admin/test", nil)
			if tc.authToken != "" {
				req.Header.Set("Authorization", "Bearer "+tc.authToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tc.expectedStatus, w.Code)
			if tc.expectedError != "" {
				require.Contains(t, w.Body.String(), tc.expectedError)
			}
			authService.AssertExpectations(t)
		})
	}
}

func TestOptionalAuthMiddleware(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name           string
		setupAuth      func() string
		setupMock      func(*mockAuthService)
		expectedStatus int
		shouldSetUser  bool
		shouldCallAuth bool
	}{
		{
			name: "with valid token",
			setupAuth: func() string {
				return "valid.token.here"
			},
			setupMock: func(m *mockAuthService) {
				m.On("GetUserFromToken", "valid.token.here").
					Return(&models.User{
						ID:        userID,
						Email:     "test@example.com",
						Role:      "user",
						FirstName: "Test",
						LastName:  "User",
						IsActive:  true,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			shouldSetUser:  true,
			shouldCallAuth: true,
		},
		{
			name: "without token",
			setupAuth: func() string {
				return ""
			},
			setupMock:      func(*mockAuthService) {},
			expectedStatus: http.StatusNoContent,
			shouldSetUser:  false,
			shouldCallAuth: false,
		},
		{
			name: "with invalid token",
			setupAuth: func() string {
				return "invalid.token"
			},
			setupMock: func(m *mockAuthService) {
				m.On("GetUserFromToken", "invalid.token").
					Return(nil, errors.New("invalid token"))
			},
			expectedStatus: http.StatusNoContent,
			shouldSetUser:  false,
			shouldCallAuth: true,
		},
		{
			name: "with empty token after Bearer",
			setupAuth: func() string {
				return ""
			},
			setupMock:      func(*mockAuthService) {},
			expectedStatus: http.StatusNoContent,
			shouldSetUser:  false,
			shouldCallAuth: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			authService := new(mockAuthService)
			tc.setupMock(authService)

			r := gin.New()
			mw := middleware.NewMiddleware(authService)

			r.GET("/optional", mw.OptionalAuthMiddleware(), func(c *gin.Context) {
				_, exists := c.Get("userID")
				if exists {
					c.Status(http.StatusOK)
				} else {
					c.Status(http.StatusNoContent)
				}
			})

			req := httptest.NewRequest("GET", "/optional", nil)
			token := tc.setupAuth()
			if token != "" {
				req.Header.Set("Authorization", "Bearer "+token)
			} else if tc.name == "with empty token after Bearer" {
				req.Header.Set("Authorization", "Bearer ")
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tc.expectedStatus, w.Code)
			
			
			if tc.shouldCallAuth {
				authService.AssertExpectations(t)
			}
		})
	}
}
