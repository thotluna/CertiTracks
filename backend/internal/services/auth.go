package services

import (
	"errors"
	"time"

	"certitrack/internal/config"
	"certitrack/internal/models"
	"certitrack/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req *RegisterRequest) (*AuthResponse, error)
	Login(req *LoginRequest) (*AuthResponse, error)
	RefreshToken(req *RefreshRequest) (*AuthResponse, error)
	GetUserFromToken(token string) (*models.User, error)
	ValidateAccessToken(token string) (*JWTClaims, error)
	ValidateRefreshToken(token string) (*JWTClaims, error)
	HashPassword(password string) (string, error)
	CheckPassword(password, hash string) bool
}

type AuthServiceImpl struct {
	repository repositories.UserRepository
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
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"firstName" binding:"required,min=2"`
	LastName  string `json:"lastName" binding:"required,min=2"`
	Phone     string `json:"phone"`
}

type AuthResponse struct {
	User         *models.User `json:"user"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	ExpiresAt    time.Time    `json:"expiresAt"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserExists         = errors.New("user with this email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token has expired")
)

func NewAuthService(config *config.Config, repository repositories.UserRepository) *AuthServiceImpl {
	return &AuthServiceImpl{
		config:     config,
		repository: repository,
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

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.config.JWT.AccessTokenExpiry),
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
	claims := JWTClaims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.JWT.Issuer,
			Audience:  []string{s.config.JWT.Audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *AuthServiceImpl) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
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

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func (s *AuthServiceImpl) ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	return s.ValidateAccessToken(tokenString)
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
