package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/Sophinaz/go-jwt-project/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/signup", controllers.Signup())
	incomingRoutes.POST("/login", controllers.login())
}