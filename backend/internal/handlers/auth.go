package handlers

import (
	"net/http"

	"certitrack/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		switch err {
		case services.ErrUserExists:
			c.JSON(http.StatusConflict, gin.H{
				"error": "User with this email already exists",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    response,
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email or password",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Login failed",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    response,
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req services.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	response, err := h.authService.RefreshToken(&req)
	if err != nil {
		switch err {
		case services.ErrInvalidToken, services.ErrTokenExpired:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired refresh token",
			})
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Token refresh failed",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"data":    response,
	})
}

// GetProfile returns the current user's profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"data":    user,
	})
}

// Logout handles user logout (client-side token removal)
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a JWT-based system, logout is typically handled client-side
	// by removing the token from storage. However, we can provide
	// an endpoint for consistency and future token blacklisting.

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// ChangePassword handles password changes for authenticated users
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	type ChangePasswordRequest struct {
		CurrentPassword string `json:"currentPassword" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required,min=6"`
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	_, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// This would require additional methods in the auth service
	// For now, we'll return a placeholder response
	c.JSON(http.StatusOK, gin.H{
		"message": "Password change endpoint - implementation pending",
	})
}
