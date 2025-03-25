package main

import (
	"fmt"
	"log"
	"os"
	"github.com/Sophinaz/go-jwt-project/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println(44)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("env error:", err)
	}

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