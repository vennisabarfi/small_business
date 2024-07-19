package main

import (
	"fmt"
	"log"
	"small_business/controllers"
	"small_business/models"

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

	// connect to db
	models.ConnectToDB()

	// r.Use(func(c *gin.Context) {
	// 	c.Set("db", db)
	// 	c.Next()
	// })

	// Middleware to pass database connection to context -- work on this
	// r.Use(models.CreateHttpMiddleware)

	//home page
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to the business API")
	})

	r.POST("/register", controllers.CreateUser)

	r.Run() //running on port in env due to fresh
}
