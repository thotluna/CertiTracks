package mocks

import (
	"certitrack/internal/models"
	"certitrack/internal/services"

	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req *services.RegisterRequest) (*services.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Login(req *services.LoginRequest) (*services.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Logout(req *services.LogoutRequest) (*services.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.AuthResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(req *services.RefreshRequest) (*services.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.AuthResponse), args.Error(1)
}

func (m *MockAuthService) GetUserFromToken(token string) (*models.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) IsTokenRevoked(tokenString string) (bool, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return true, args.Error(1)
	}
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthService) ValidateAccessToken(token string) (*services.JWTClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.JWTClaims), args.Error(1)
}

func (m *MockAuthService) ValidateRefreshToken(token string) (*services.JWTClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.JWTClaims), args.Error(1)
}

func (m *MockAuthService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) CheckPassword(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}
