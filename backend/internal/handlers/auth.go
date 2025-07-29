package handlers

import (
	"net/http"

	"certitrack/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
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
		case services.ErrUserNotFound, services.ErrInvalidCredentials:
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

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully. Please remove the token from client storage.",
	})
}

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

	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Password change functionality is not yet implemented",
	})
}
