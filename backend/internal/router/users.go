package router

import (
	"certitrack/internal/config"
	"certitrack/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupUserRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	users := rg.Group("/users")
	users.Use(middleware.AdminMiddleware())
	{
		users.GET("", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "Users"})
		})
		users.POST("", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "Create User"})
		})

		userRoutes := users.Group("/:id")
		{
			userRoutes.GET("", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"message": "Get User"})
			})
			userRoutes.PUT("", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"message": "Update User"})
			})
			userRoutes.DELETE("", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"message": "Delete User"})
			})
		}
	}
}
