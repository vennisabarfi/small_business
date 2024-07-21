package controllers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" //avoid import postgres error with sql
)

var pool *sql.DB // Database connection pool

type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       []uint8 `json: "price"`
	Stock       int     `json:"stock"`
	SupplierID  int64   `json: "supplier_id"`
}

// convert price from int to decimal in usage

// GET products/:id
func ViewProductsById(c *gin.Context) {

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid product ID",
		})
		return
	}

	//open database connection
	pool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error opening database connection")
	}

	defer pool.Close()

	ctx := context.Background()

	var product Product

	query := "SELECT id, name, description, price, stock, supplier_id FROM products WHERE id = $1"
	// query := "SELECT * FROM products WHERE id = $1"

	row := pool.QueryRowContext(ctx, query, id)

	// map onto database
	err = row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.SupplierID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{
				"message": "No product found",
			})
		} else {
			log.Printf("Error scanning row: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": "Error retrieving product",
			})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"Product Found": product,
	})

}

func ViewProducts() {

}
