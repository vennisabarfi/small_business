package main

import (
	"fmt"
	"log"
	"small_business/controllers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file.", err)
	}
	fmt.Println(".env file loaded successfully!")

	r := gin.Default()

	//home page
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to the business API")
	})

	r.POST("/register", controllers.CreateUser)
	r.POST("/login", controllers.LoginUser)

	r.Run() //running on port in env due to fresh
}
