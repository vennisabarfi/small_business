package controllers

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func createProduct(db *sql.DB, name string, description string, supplierID int, price decimal.Decimal, stock int, minimumStock int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := "INSERT INTO products (name, description, supplier_id, price, stock, minimum_stock) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err = tx.Exec(query, name, description, supplierID, price, stock, minimumStock)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func TestCreateProduct(t *testing.T) {
	t.Parallel()

	// Create a new mock database connection and sqlmock instance
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Define the query we expect to be executed
	query := "INSERT INTO products \\(name, description, supplier_id, price, stock, minimum_stock\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6\\)"

	// Begin a new transaction
	mock.ExpectBegin()

	// mock query
	mock.ExpectExec(query).WithArgs("testproduct", "testdescription", 1, decimal.NewFromFloat(123.99), 2, 3).WillReturnResult(sqlmock.NewResult(1, 1)) //new row inserted with id 1

	// Commit the transaction
	mock.ExpectCommit()

	// Call the function to create a product
	err = createProduct(db, "testproduct", "testdescription", 1, decimal.NewFromFloat(123.99), 2, 3) //avoid decimal errors
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func getProductByID(db *sql.DB, id int) (*Product, error) {
	query := "SELECT id, name, description, supplier_id, price, stock, minimum_stock FROM products WHERE id = $1"
	row := db.QueryRow(query, id)

	var product Product
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.SupplierID, &product.Price, &product.Stock, &product.MinimumStock)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func TestGetProductByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	//mock query
	query := "SELECT id, name, description, supplier_id, price, stock, minimum_stock FROM products WHERE id = \\$1"
	rows := sqlmock.NewRows([]string{"id", "name", "description", "supplier_id", "price", "stock", "minimum_stock"}).
		AddRow(1, "testproduct", "testdescription", 1, 999, 2, 3)

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	// assert that not nil and equal user email
	product, err := getProductByID(db, 1)
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.EqualValues(t, 1, product.ID)

	// Ensure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func deleteProductByID(db *sql.DB, id int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := "DELETE FROM products WHERE id = $1"
	_, err = tx.Exec(query, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func TestDeleteProductByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	// mock query
	query := "DELETE FROM products WHERE id = \\$1"
	mock.ExpectExec(query).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1)) // return no new row with index 1

	mock.ExpectCommit()

	err = deleteProductByID(db, 1)
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func updateProductPrice(db *sql.DB, price decimal.Decimal, id int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := "UPDATE products SET price = $1 WHERE id = $2"
	_, err = tx.Exec(query, price, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func TestUpdateProductPrice(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	//mock query

	query := "UPDATE products SET price = \\$1 WHERE id = \\$2"
	mock.ExpectExec(query).WithArgs(decimal.NewFromFloat(999.99), 1).WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	err = updateProductPrice(db, decimal.NewFromFloat(999.99), 1)
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func updateProductStock(db *sql.DB, stock int, id int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := "UPDATE products SET stock = $1 WHERE id = $2"
	_, err = tx.Exec(query, stock, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func TestUpdateProductStock(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	//mock query

	query := "UPDATE products SET stock = \\$1 WHERE id = \\$2"
	mock.ExpectExec(query).WithArgs(2, 1).WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	err = updateProductStock(db, 2, 1)
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
