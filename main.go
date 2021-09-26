package main

import (
	"fmt"
	foodController "go-with-mongodb/controllers"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	fmt.Println("Running on", port)
	router := gin.New()
	router.Use(gin.Logger())
	router.POST("/foods-create", foodController.PostFood)
	router.Run(":" + port)
}
