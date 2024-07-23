package controllers

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// var pool *sql.DB

type Suppliers struct {
}

func ViewSuppliers(c *gin.Context) {
	//open database connection
	pool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error opening database connection")
	}

	defer pool.Close()

	// ctx := context.Background()

}
