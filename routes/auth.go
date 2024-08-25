package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samdgea/gin-boilerplate/controllers"
)

func SetupAuthRoutes(router *gin.RouterGroup) {
	router.POST("/login", controllers.Login)
}
