package controllers

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func createSupplier(db *sql.DB, name string, contact_email string, phone string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := "INSERT INTO supplier (name, contact_email, phone) VALUES ($1, $2, $3)"
	_, err = tx.Exec(query, name, contact_email, phone)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func TestCreateSupplier(t *testing.T) {
	t.Parallel()

	// Create a new mock database connection and sqlmock instance
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Define the query we expect to be executed
	query := "INSERT INTO supplier \\(name, contact_email, phone\\) VALUES \\(\\$1, \\$2, \\$3\\)"

	// Begin a new transaction
	mock.ExpectBegin()

	// mock query
	mock.ExpectExec(query).WithArgs("testsupplier", "testsupplier@gmail.com", "0123456789").WillReturnResult(sqlmock.NewResult(1, 1)) //new row inserted with id 1

	// Commit the transaction
	mock.ExpectCommit()

	// Call the function to create a product
	err = createSupplier(db, "testsupplier", "testsupplier@gmail.com", "0123456789")
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func getSupplierByID(db *sql.DB, id int) (*Supplier, error) {
	query := "SELECT id, name, contact_email, phone FROM supplier WHERE id = $1"
	row := db.QueryRow(query, id)

	var supplier Supplier
	err := row.Scan(&supplier.ID, &supplier.Name, &supplier.ContactEmail, &supplier.Phone)
	if err != nil {
		return nil, err
	}

	return &supplier, nil
}

func TestGetSupplierByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	//mock query
	query := "SELECT id, name, contact_email, phone FROM supplier WHERE id = \\$1"
	rows := sqlmock.NewRows([]string{"id", "name", "contact_email", "phone"}).
		AddRow(1, "testsupplier", "supplier@gmail.com", "012345678")

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	// assert that not nil and equal user email
	supplier, err := getSupplierByID(db, 1)
	assert.NoError(t, err)
	assert.NotNil(t, supplier)
	assert.EqualValues(t, 1, supplier.ID)

	// Ensure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func deleteSupplierByID(db *sql.DB, id int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := "DELETE FROM supplier WHERE id = $1"
	_, err = tx.Exec(query, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func TestDeleteSupplierByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	// mock query
	query := "DELETE FROM supplier WHERE id = \\$1"
	mock.ExpectExec(query).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1)) // return no new row with index 1

	mock.ExpectCommit()

	err = deleteSupplierByID(db, 1)
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func updateSupplierEmail(db *sql.DB, contact_email string, id int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := "UPDATE supplier SET contact_email = $1 WHERE id = $2"
	_, err = tx.Exec(query, contact_email, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func TestUpdateSupplierEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	//mock query

	query := "UPDATE supplier SET contact_email = \\$1 WHERE id = \\$2"
	mock.ExpectExec(query).WithArgs("supplier2@gmail.com", 1).WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	err = updateSupplierEmail(db, "supplier2@gmail.com", 1)
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func updateSupplierPhone(db *sql.DB, phone string, id int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := "UPDATE supplier SET phone = $1 WHERE id = $2"
	_, err = tx.Exec(query, phone, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func TestUpdateSupplierPhone(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	//mock query

	query := "UPDATE supplier SET phone = \\$1 WHERE id = \\$2"
	mock.ExpectExec(query).WithArgs("0123456789", 1).WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	err = updateSupplierPhone(db, "0123456789", 1)
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
