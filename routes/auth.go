package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samdgea/gin-boilerplate/controllers"
	"github.com/samdgea/gin-boilerplate/middlewares"
)

func SetupAuthRoutes(router *gin.RouterGroup) {
	router.POST("/login", controllers.Login)
	router.POST("/refresh-token", controllers.RefreshToken)

	router.Use(middlewares.IsAuth).GET("/me", controllers.Me)
}
