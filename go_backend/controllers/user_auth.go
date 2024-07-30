package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq" //avoid import postgres error with sql
	"golang.org/x/crypto/bcrypt"
)

// create user struct
type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// var pool *sql.DB // Database connection pool.

//handler for creating user

func CreateUser(c *gin.Context) {
	// receive user data (validated)
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	// if error with fields

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error Binding JSON Data": err,
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

	ctx := context.Background()

	_, err = pool.ExecContext(ctx, "INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error creating new user": err,
		})
		return

	} else {
		fmt.Println("Inserting user into database")

		// Respond with the user ID
		c.String(http.StatusOK, "User successfully Added")
	}

}

func LoginUser(c *gin.Context) {
	// Parse email and password from body
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	//user data validation
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	// Open database connection
	pool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to connect to database",
		})
		return
	}
	defer pool.Close()

	ctx := context.Background()

	// Get user from database
	var storedEmail, storedHashedPassword string
	row := pool.QueryRowContext(ctx, "SELECT email, password FROM users WHERE email=$1", body.Email)

	err = row.Scan(&storedEmail, &storedHashedPassword)

	// if email and password not found
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Email and/or password is incorrect",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve user data",
		})
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	// jwt authentication (refreshes every 30 days)
	var user User
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create token",
		})

		return
	}

	// set cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	// Successfully authenticated
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})
}
