package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"small_business/controllers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// load env file
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file.", err)
	}
	fmt.Println(".env file loaded successfully!")
}

func securityHeaders(c *gin.Context) {

	port := os.Getenv("PORT")
	expectedHost := "localhost:" + port

	if c.Request.Host != expectedHost {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
		return
	}
	c.Header("X-Frame-Options", "DENY")
	c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	c.Header("Referrer-Policy", "strict-origin")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
	c.Next()
}

func main() {

	LoadEnv()

	r := gin.Default()

	//Security headers
	r.Use(securityHeaders)

	//home page
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to the business API")
	})

	// user handlers
	user := r.Group("/user-auth")
	{
		user.POST("/register", controllers.CreateUser)
		user.POST("/login", controllers.LoginUser)
	}

	//products handlers
	products := r.Group("/products")
	{
		products.GET("/:id", controllers.ViewProductsById)
		products.GET("/", controllers.ViewProducts) // "/products"
		products.POST("insert", controllers.InsertProduct)
		products.PUT("/change-price", controllers.UpdateProductPrice)
		products.DELETE("/remove/:id", controllers.DeleteProductByID)
	}

	//supplier handlers
	suppliers := r.Group("/suppliers")
	{
		suppliers.GET("/", controllers.ViewSuppliers)
		suppliers.GET("/:id", controllers.ViewSuppliersById)
		suppliers.POST("/insert", controllers.InsertSupplier)
		suppliers.PUT("/change-email", controllers.UpdateSupplierEmail)
		suppliers.PUT("/change-phone", controllers.UpdateSupplierPhone)
		suppliers.DELETE("/remove/:id", controllers.DeleteSupplierByID)
	}

	r.Run() //running on port in env due to fresh
}
