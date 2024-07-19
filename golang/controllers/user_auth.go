package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

//create user struct

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//handler for creating user

func CreateUser(c *gin.Context) {
	// receive user data (validated)
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	//Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		log.Fatal("Password was not successfully hashed", err)
	}

	user := User{Email: body.Email, Password: string(hashedPassword)}

	// Log the received user data
	log.Printf("Received user data: %+v", user)

	// Retrieve the database connection from Gin's context
	dbpool, exists := c.Get("db")
	if !exists {
		log.Fatal("No database connection in Gin context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection is not initialized",
		})
		return
	}

	// Type assertion to convert db to the correct type
	conn, ok := dbpool.(*pgxpool.Pool)
	if !ok {
		log.Fatal("Invalid database connection")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid database connection",
		})
		return
	}

	// Execute the SQL query to insert the user and retrieve the ID
	err = conn.QueryRow(context.Background(), "INSERT INTO small_business (email, password) VALUES ($1, $2) RETURNING id", user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to insert user",
		})
		return
	}

	// Respond with the user ID
	c.JSON(http.StatusOK, gin.H{"userID": user.ID})
}
