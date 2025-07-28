package services

import (
	"errors"
	"time"

	"certitrack/internal/config"
	"certitrack/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db     *gorm.DB
	config *config.Config
}

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

func NewAuthService(db *gorm.DB, config *config.Config) *AuthService {
	return &AuthService{
		db:     db,
		config: config,
	}
}

func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := s.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      "user", // Default role
		IsActive:  true,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := s.GenerateAccessToken(&user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.GenerateRefreshToken(&user)
	if err != nil {
		return nil, err
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	s.db.Save(&user)

	return &AuthResponse{
		User:         &user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.config.JWT.AccessTokenExpiry),
	}, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	var user models.User
	if err := s.db.Where("email = ? AND is_active = ?", req.Email, true).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check password
	if !s.CheckPassword(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, err := s.GenerateAccessToken(&user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.GenerateRefreshToken(&user)
	if err != nil {
		return nil, err
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	s.db.Save(&user)

	return &AuthResponse{
		User:         &user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.config.JWT.AccessTokenExpiry),
	}, nil
}

// RefreshToken generates new tokens using a refresh token
func (s *AuthService) RefreshToken(req *RefreshRequest) (*AuthResponse, error) {
	claims, err := s.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := s.db.Where("id = ? AND is_active = ?", claims.UserID, true).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Generate new tokens
	accessToken, err := s.GenerateAccessToken(&user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.GenerateRefreshToken(&user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         &user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.config.JWT.AccessTokenExpiry),
	}, nil
}

// HashPassword hashes a plain text password
func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a plain text password with a hash
func (s *AuthService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateAccessToken creates a new access token
func (s *AuthService) GenerateAccessToken(user *models.User) (string, error) {
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

// GenerateRefreshToken creates a new refresh token
func (s *AuthService) GenerateRefreshToken(user *models.User) (string, error) {
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

// ValidateAccessToken validates and parses an access token
func (s *AuthService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
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

// ValidateRefreshToken validates and parses a refresh token
func (s *AuthService) ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	return s.ValidateAccessToken(tokenString) // Same validation logic
}

// GetUserFromToken extracts user information from a valid token
func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	claims, err := s.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := s.db.Where("id = ? AND is_active = ?", claims.UserID, true).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
