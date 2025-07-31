package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupUserRoutes(rg *gin.RouterGroup, deps *RouterDeps) {
	users := rg.Group("/users")
	users.Use(deps.Middleware.AdminMiddleware())
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
