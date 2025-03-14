package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/Sophinaz/go-jwt-project/controllers"
	middleware "github.com/Sophinaz/go-jwt-project/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/user", controllers.Getusers())
	incomingRoutes.GET("/user:id", controllers.GetUser())
}