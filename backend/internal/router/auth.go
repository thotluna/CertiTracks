package router

import (
	"github.com/gin-gonic/gin"
)

func setupAuthRoutes(rg *gin.RouterGroup, deps *RouterDeps) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", deps.AuthHandler.Login)
		auth.POST("/register", deps.AuthHandler.Register)
		auth.POST("/refresh", deps.AuthHandler.RefreshToken)
	}
}

func setupProtectedAuthRoutes(rg *gin.RouterGroup, deps *RouterDeps) {
	auth := rg.Group("/auth")
	{
		auth.GET("/profile", deps.AuthHandler.GetProfile)
		auth.POST("/change-password", deps.AuthHandler.ChangePassword)
		auth.POST("/logout", deps.AuthHandler.Logout)
	}
}
