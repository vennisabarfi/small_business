package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" //avoid import postgres error with sql
	"golang.org/x/crypto/bcrypt"
)

//create user struct

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var pool *sql.DB // Database connection pool.

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

	pool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(3)

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	appSignal := make(chan os.Signal, 3)
	signal.Notify(appSignal, os.Interrupt)

	go func() {
		<-appSignal
		stop()
	}()

	_, err = pool.ExecContext(ctx, "INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)

	if err != nil {
		log.Fatal("Failed to insert user", err)
	}
	fmt.Println("Inserting user into database")

	// Respond with the user ID
	c.JSON(http.StatusOK, gin.H{
		"userID": user.ID,
	})
}

func CreateUserQuery() {

}
