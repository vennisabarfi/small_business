package models

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbpool *pgx.ConnPool

// using a connection pool
func ConnectToDB() {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Successfully created connection pool to database!")
	}
	defer dbpool.Close()

	// Ensure database is connected by pinging
	connection, err := dbpool.Acquire(context.Background())
	if err != nil {
		log.Fatal("Unable to acquire connection to the database %v", err)
		os.Exit(1)
	}
	fmt.Println("Connection to database established!")

	defer connection.Release()

	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'This is a test query'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
}

// look into this and test(causing errors. Fix this or implement differently)
func CreateHttpMiddleware(c *gin.Context) {
	tx, err := dbpool.Begin() //(c.Request.Context)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	fmt.Println("Database Connection initiated with HTTP Request!")
	defer tx.Rollback()

	c.Set("db", tx)
	c.Next()

	if c.Writer.Status() >= http.StatusInternalServerError {
		tx.Rollback()
		return
	}
}
