package main

import (
	"os"
	"github.com/Sophinaz/go-jwt-project/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8082"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)


	router.Run("localhost:" + PORT)
}