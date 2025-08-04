// Package services contains the business logic of the application
package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"certitrack/internal/config"
	"certitrack/internal/models"
	"certitrack/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req *RegisterRequest) (*AuthResponse, error)
	Login(req *LoginRequest) (*AuthResponse, error)
	Logout(req *LogoutRequest) (*AuthResponse, error)
	RevokeToken(req *LogoutRequest) (*AuthResponse, error)
	RefreshToken(req *RefreshRequest) (*AuthResponse, error)
	GetUserFromToken(token string) (*models.User, error)
	IsTokenRevoked(tokenString string) (bool, error)
	ValidateAccessToken(token string) (*JWTClaims, error)
	ValidateRefreshToken(token string) (*JWTClaims, error)
	HashPassword(password string) (string, error)
	CheckPassword(password, hash string) bool
}

type AuthServiceImpl struct {
	repository repositories.UserRepository
	tokenRepo  repositories.TokenRepository
	config     *config.Config
}

var _ AuthService = (*AuthServiceImpl)(nil)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type LogoutRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=72,strong_password"`
	FirstName string `json:"first_name" binding:"required,min=2,max=50"`
	LastName  string `json:"last_name" binding:"required,min=2,max=50"`
	Phone     string `json:"phone,omitempty" binding:"omitempty,min=8,max=20"`
}

type AuthResponse struct {
	User         *models.User `json:"user"`
	AccessToken  string       `json:"access-token"`
	RefreshToken string       `json:"refresh-token"`
	ExpiresAt    time.Time    `json:"expiresAt"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserExists         = errors.New("user with this email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTokenExpired       = errors.New("token has expired")
	ErrRefreshTokenFailed = errors.New("failed to refresh token")
	ErrInvalidAudience    = errors.New("invalid token audience")
)

func NewAuthService(config *config.Config, repository repositories.UserRepository, tokenRepo repositories.TokenRepository) *AuthServiceImpl {
	return &AuthServiceImpl{
		config:     config,
		repository: repository,
		tokenRepo:  tokenRepo,
	}
}

func (s *AuthServiceImpl) Register(req *RegisterRequest) (*AuthResponse, error) {

	if s.repository.EmailExists(req.Email) {
		return nil, ErrUserExists
	}

	hashedPassword, err := s.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      "user",
		IsActive:  true,
	}

	if err := s.repository.CreateUser(&user); err != nil {
		return nil, err
	}

	accessToken, err := s.GenerateAccessToken(&user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.GenerateRefreshToken(&user)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user.LastLogin = &now
	s.repository.UpdateLastLogin(user.ID.String(), now)

	return &AuthResponse{
		User:         &user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.config.JWT.AccessTokenExpiry),
	}, nil
}

func (s *AuthServiceImpl) Login(req *LoginRequest) (*AuthResponse, error) {
	user, err := s.repository.FindActiveByEmail(req.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !s.CheckPassword(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user.LastLogin = &now
	s.repository.UpdateLastLogin(user.ID.String(), now)

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.config.JWT.AccessTokenExpiry),
	}, nil
}

func (s *AuthServiceImpl) Logout(req *LogoutRequest) (*AuthResponse, error) {
	return s.RevokeToken(req)
}

func (s *AuthServiceImpl) RevokeToken(req *LogoutRequest) (*AuthResponse, error) {

	if req.AccessToken == "" && req.RefreshToken == "" {
		log.Printf("RevokeToken failed: no tokens provided")
		return nil, fmt.Errorf("%w: neither access_token nor refresh_token provided", ErrInvalidToken)
	}

	if req.AccessToken != "" {
		log.Printf("Revoking access token")
		if err := s.tokenRepo.RevokeToken(
			req.AccessToken,
			s.config.JWT.AccessTokenExpiry,
		); err != nil {
			log.Printf("Failed to revoke access token: %v", err)
			return nil, fmt.Errorf("failed to revoke access token: %w", err)
		}
	}

	if req.RefreshToken != "" {
		log.Printf("Revoking refresh token")
		if err := s.tokenRepo.RevokeToken(
			req.RefreshToken,
			s.config.JWT.RefreshTokenExpiry,
		); err != nil {
			log.Printf("Failed to revoke refresh token: %v", err)
			return nil, fmt.Errorf("failed to revoke refresh token: %w", err)
		}
	}

	return &AuthResponse{}, nil
}

func (s *AuthServiceImpl) RefreshToken(req *RefreshRequest) (*AuthResponse, error) {
	claims, err := s.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}
	user, err := s.repository.FindActiveByID(claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	accessToken, err := s.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	if refreshToken == req.RefreshToken {
		refreshToken, err = s.GenerateRefreshToken(user)
		if err != nil {
			return nil, err
		}
	}

	now := time.Now()
	user.LastLogin = &now
	s.repository.UpdateLastLogin(user.ID.String(), now)

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    now.Add(s.config.JWT.AccessTokenExpiry),
	}, nil
}

func (s *AuthServiceImpl) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *AuthServiceImpl) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthServiceImpl) GenerateAccessToken(user *models.User) (string, error) {
	claims := JWTClaims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.JWT.Issuer,
			Audience:  []string{s.config.JWT.Audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *AuthServiceImpl) GenerateRefreshToken(user *models.User) (string, error) {
	jti := uuid.New().String()
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"jti":     jti,
		"exp":     time.Now().Add(s.config.JWT.RefreshTokenExpiry).Unix(),
		"iat":     time.Now().Unix(),
		"nbf":     time.Now().Unix(),
		"iss":     s.config.JWT.Issuer,
		"aud":     s.config.JWT.Audience,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *AuthServiceImpl) IsTokenRevoked(tokenString string) (bool, error) {
	log.Println("entrando a preguntar ===============> ")
	return s.tokenRepo.IsTokenRevoked(tokenString)
}

func (s *AuthServiceImpl) ValidateAccessToken(tokenString string) (*JWTClaims, error) {

	log.Println("entrando a validar ===============> ")
	isRevoked, err := s.IsTokenRevoked(tokenString)

	if err != nil {
		return nil, err
	}

	log.Println("Revoked ===============> ", isRevoked, tokenString)
	if isRevoked {
		return nil, ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Validar la audiencia del token
	audienceValid := false
	expectedAudience := s.config.JWT.Audience
	for _, aud := range claims.Audience {
		if aud == expectedAudience {
			audienceValid = true
			break
		}
	}

	if !audienceValid {
		return nil, ErrInvalidAudience
	}

	return claims, nil
}

func (s *AuthServiceImpl) ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, ErrInvalidToken
		}

		email, ok := claims["email"].(string)
		if !ok {
			return nil, ErrInvalidToken
		}

		role, ok := claims["role"].(string)
		if !ok {
			return nil, ErrInvalidToken
		}

		return &JWTClaims{
			UserID: userID,
			Email:  email,
			Role:   role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Unix(int64(claims["exp"].(float64)), 0)),
				IssuedAt:  jwt.NewNumericDate(time.Unix(int64(claims["iat"].(float64)), 0)),
				NotBefore: jwt.NewNumericDate(time.Unix(int64(claims["nbf"].(float64)), 0)),
				Issuer:    claims["iss"].(string),
				Audience:  []string{claims["aud"].(string)},
			},
		}, nil
	}

	return nil, ErrInvalidToken
}

func (s *AuthServiceImpl) GetUserFromToken(tokenString string) (*models.User, error) {
	claims, err := s.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	user, err := s.repository.FindActiveByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
