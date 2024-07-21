package controllers

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// createUser inserts a new user into the database
func createUser(db *sql.DB, email string, password string) error {
	// use transactions
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	//mimic login functionality
	query := "INSERT INTO users (email, password) VALUES ($1, $2)"
	_, err = tx.Exec(query, email, password)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
func TestCreateUser(t *testing.T) {
	t.Parallel()

	// Create a new mock database connection and sqlmock instance
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Define the query we expect to be executed
	query := "INSERT INTO users \\(email, password\\) VALUES \\(\\$1, \\$2\\)"

	// Begin a new transaction
	mock.ExpectBegin()

	// mock query
	mock.ExpectExec(query).WithArgs("test@example.com", "hashedpassword").WillReturnResult(sqlmock.NewResult(1, 1)) //new row inserted with id 1

	// Commit the transaction
	mock.ExpectCommit()

	// Call the function to create a user
	err = createUser(db, "test@example.com", "hashedpassword")
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func deleteUserByEmail(db *sql.DB, email string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := "DELETE FROM users WHERE email = $1"
	_, err = tx.Exec(query, email)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func TestDeleteUserByEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	// mock query
	query := "DELETE FROM users WHERE email = \\$1"
	mock.ExpectExec(query).WithArgs("test@example.com").WillReturnResult(sqlmock.NewResult(0, 1)) // return no new row with index 1

	mock.ExpectCommit()

	err = deleteUserByEmail(db, "test@example.com")
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func updateUserPassword(db *sql.DB, email string, newPassword string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := "UPDATE users SET password = $1 WHERE email = $2"
	_, err = tx.Exec(query, newPassword, email)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func TestUpdateUserPassword(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	//mock query

	query := "UPDATE users SET password = \\$1 WHERE email = \\$2"
	mock.ExpectExec(query).WithArgs("newhashedpassword", "test@example.com").WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	err = updateUserPassword(db, "test@example.com", "newhashedpassword")
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func getUserByEmail(db *sql.DB, email string) (*User, error) {
	query := "SELECT id, email, password FROM users WHERE email = $1"
	row := db.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func TestGetUserByEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	//mock query
	query := "SELECT id, email, password FROM users WHERE email = \\$1"
	rows := sqlmock.NewRows([]string{"id", "email", "password"}).
		AddRow(1, "test@example.com", "hashedpassword")

	mock.ExpectQuery(query).WithArgs("test@example.com").WillReturnRows(rows)

	// assert that not nil and equal user email
	user, err := getUserByEmail(db, "test@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)

	// Ensure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
