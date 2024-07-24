package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// var pool *sql.DB

type Supplier struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	ContactEmail string `json:"contact_email"`
	Phone        string `json:"phone"`
	Location     string `json:"location"`
}

var pool *sql.DB

func InsertSupplier(c *gin.Context) {
	//open database connection
	pool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error opening database connection")
	}

	defer pool.Close()

	ctx := context.Background()
	var supplier Supplier

	query := "INSERT INTO suppliers (id, name, contact_email, phone, location) VALUES($1, $2, $3, $4, $5)"

	_, err = pool.ExecContext(ctx, query, supplier.ID, supplier.Name, supplier.ContactEmail, supplier.Phone, supplier.Location)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error inserting new supplier": err,
		})
		return

	} else {
		fmt.Println("Inserting supplier into database")

		// Respond with supplier information
		c.IndentedJSON(http.StatusOK, gin.H{
			"message":             "Supplier Successfully Added",
			"Product Information": supplier,
		})
		c.String(http.StatusOK, "Supplier successfully Added")
	}
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
