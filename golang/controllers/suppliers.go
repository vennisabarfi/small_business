package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// var pool *sql.DB

type Supplier struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	ContactEmail string `json:"contact_email"`
	Phone        string `json:"phone"`
	Createdat    string `json: created_at`
	Deletedat    string `json: deleted_at`
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

	// var supplier Supplier

	query := "SELECT * FROM supplier"

	rows, err := pool.Query(query) //uses ctx internally
	if err != nil {
		print(err)
	}
	defer rows.Close()

	var suppliers []Supplier

	// Loop through rows and map onto databases
	for rows.Next() {
		var supplier Supplier
		if err := rows.Scan(&supplier.ID, &supplier.Name, &supplier.ContactEmail, &supplier.Phone, &supplier.Createdat, &supplier.Deletedat); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"Error retrieving suppliers": err,
			})
			log.Print("Error inserting new supplier", err)
			return
		}
		suppliers = append(suppliers, supplier)
	}
	if err == sql.ErrNoRows {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "No suppliers found",
		})
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"Suppliers Found": suppliers,
	})

}

func ViewSuppliersById(c *gin.Context) {

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid supplier ID",
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

	var supplier Supplier

	query := "SELECT id, name, contact_email, phone FROM supplier WHERE id = $1"

	row := pool.QueryRowContext(ctx, query, id)

	// map onto database
	err = row.Scan(&supplier.ID, &supplier.Name, &supplier.ContactEmail, &supplier.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{
				"message": "No supplier found with this ID",
			})
		} else {
			log.Printf("Error scanning row: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": "Error retrieving supplier",
			})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"Supplier Found": supplier,
	})

}

func DeleteSupplierByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid supplier ID",
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

	query := "DELETE FROM supplier WHERE id = $1 RETURNING id"

	res, err := pool.ExecContext(ctx, query, id)

	if err != nil {
		log.Fatal(err)
	}
	rows, err := res.RowsAffected()
	if err == nil {
		if rows != 1 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"No Supplier with this ID Exists": err,
			})
			return
		}

	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"Message": "Supplier successfully removed!",
	})
}

// update supplier email by id
func UpdateSupplierEmail(c *gin.Context) {
	var body struct {
		ID           int64  `json:"id"`
		ContactEmail string `json:"contact_email"`
	}

	// if error with fields

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error Binding JSON Data": err,
		})
		return
	}

	supplier := Supplier{ID: body.ID, ContactEmail: body.ContactEmail}

	//open database connection
	pool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error opening database connection")
	}

	defer pool.Close()

	ctx := context.Background()

	query := "UPDATE supplier SET contact_email = $1 WHERE id = $2"

	_, err = pool.ExecContext(ctx, query, supplier.ContactEmail, supplier.ID)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error updating supplier email address": err,
		})
		log.Print("Error updating supplier email address", err)
		return

	} else {
		fmt.Println("Updating Supplier Email Address into database...")

		// Respond with product information
		c.IndentedJSON(http.StatusOK, gin.H{
			"message":           "Supplier Email Updated Successfully",
			"New Email Address": supplier.ContactEmail,
		})
	}

}

// update supplier phone number by id
func UpdateSupplierPhone(c *gin.Context) {
	var body struct {
		ID    int64  `json:"id"`
		Phone string `json:"phone"`
	}

	// if error with fields

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error Binding JSON Data": err,
		})
		return
	}

	supplier := Supplier{ID: body.ID, Phone: body.Phone}

	//open database connection
	pool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error opening database connection")
	}

	defer pool.Close()

	ctx := context.Background()

	query := "UPDATE supplier SET phone = $1 WHERE id = $2"

	_, err = pool.ExecContext(ctx, query, supplier.Phone, supplier.ID)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error updating supplier phone number": err,
		})
		log.Print("Error updating supplier phone number", err)
		return

	} else {
		fmt.Println("Updating Supplier Phone Number into database...")

		// Respond with product information
		c.IndentedJSON(http.StatusOK, gin.H{
			"message":          "Supplier Phone Number Updated Successfully",
			"New Phone Number": supplier.Phone,
		})
	}

}
