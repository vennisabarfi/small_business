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
	"github.com/shopspring/decimal"
)

// var pool *sql.DB // Database connection pool

type Product struct {
	ID           int64           `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	SupplierID   int64           `json: "supplier_id"`
	Price        decimal.Decimal `json: "price"`
	Stock        int             `json:"stock"`
	MinimumStock int             `json:"minimum_stock"`
	CreatedAt    string          `json: created_at`
	DeletedAt    string          `json: deleted_at`
}

func InsertProduct(c *gin.Context) {
	var body struct {
		Name         string          `json:"name"`
		Description  string          `json:"description"`
		SupplierID   int64           `json:"supplier_id"`
		Price        decimal.Decimal `json:"price"`
		Stock        int             `json:"stock"`
		MinimumStock int             `json:"minimum_stock"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error binding JSON data",
			"details": err.Error(),
		})
		return
	}

	product := Product{
		Name:         body.Name,
		Description:  body.Description,
		SupplierID:   body.SupplierID,
		Price:        body.Price,
		Stock:        body.Stock,
		MinimumStock: body.MinimumStock,
	}

	// Open database connection
	pool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error opening database connection")
	}

	defer pool.Close()

	ctx := context.Background()
	query := "INSERT INTO products (name, description, supplier_id, price, stock, minimum_stock) VALUES($1, $2, $3, $4, $5, $6) RETURNING id"
	err = pool.QueryRowContext(ctx, query, product.Name, product.Description, product.SupplierID, product.Price, product.Stock, product.MinimumStock).Scan(&product.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error inserting new product",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":             "Product successfully added",
		"product_information": product,
	})
}

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

	query := "SELECT id, name, description, supplier_id, price, stock, minimum_stock FROM products WHERE id = $1"

	row := pool.QueryRowContext(ctx, query, id)

	// map onto database
	err = row.Scan(&product.ID, &product.Name, &product.Description, &product.SupplierID, &product.Price, &product.Stock, &product.MinimumStock)
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

// View all products

func ViewProducts(c *gin.Context) {

	//open database connection
	pool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error opening database connection")
	}

	defer pool.Close()

	query := "SELECT * FROM products"

	rows, err := pool.Query(query) //uses ctx internally
	if err != nil {
		print(err)
	}
	defer rows.Close()

	var products []Product

	// Loop through rows and map onto databases
	for rows.Next() {
		var product Product

		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.SupplierID, &product.Price, &product.Stock, &product.MinimumStock, &product.CreatedAt, &product.DeletedAt); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"Error retrieving products": err,
			})
			log.Print("Error inserting new product", err)
			return
		}
		products = append(products, product)
	}
	if err == sql.ErrNoRows {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "No products found",
		})
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"Products Found": products,
	})

}

func DeleteProductByID(c *gin.Context) {
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

	// var supplier Supplier

	query := "DELETE FROM products WHERE id = $1 RETURNING id"

	res, err := pool.ExecContext(ctx, query, id)

	if err != nil {
		log.Fatal(err)
	}
	rows, err := res.RowsAffected()
	if err == nil {
		if rows != 1 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"No Product with this ID Exists": err,
			})
			return
		}

	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"Message": "Product successfully removed!",
	})
}
