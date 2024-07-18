package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-playground/assert"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

// load .env file
func loadENV() {
	err := godotenv.Load("models.env")
	if err != nil {
		log.Fatal("Error loading .env file.", err)
	}
	fmt.Println(".env file loaded successfully!")
}

// use pgxpoolmock

// test connection
func TestConnect(t *testing.T) {

	loadENV()

	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	connString := os.Getenv("PGX_TEST_DATABASE")
	pool, err := pgxpool.New(ctx, connString)
	require.NoError(t, err, "Failed to create connection pool")

	defer pool.Close()

	assert.Equal(t, connString, pool.Config().ConnString())
	pool.Close()
}
