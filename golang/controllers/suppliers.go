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
}

var pool *sql.DB

func InsertSupplier(c *gin.Context) {

	var body struct {
		Name         string `json:"name"`
		ContactEmail string `json:"contact_email"`
		Phone        string `json:"phone"`
	}

	// if error with fields

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error Binding JSON Data": err,
		})
		return
	}

	supplier := Supplier{Name: body.Name, ContactEmail: body.ContactEmail, Phone: body.Phone}

	//open database connection
	pool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error opening database connection")
	}

	defer pool.Close()

	ctx := context.Background()

	query := "INSERT INTO supplier (name, contact_email, phone) VALUES($1, $2, $3)"

	_, err = pool.ExecContext(ctx, query, supplier.Name, supplier.ContactEmail, supplier.Phone)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error inserting new supplier": err,
		})
		log.Print("Error inserting new supplier", err)
		return

	} else {
		fmt.Println("Inserting supplier information into database...")

		// Respond with product information
		c.IndentedJSON(http.StatusOK, gin.H{
			"message":             "Supplier Information Successfully Added",
			"Product Information": supplier,
		})
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
