package models

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"
	//avoid import error with generic database/sql
)

var pool *sql.DB

func ConnectToDB() (ctx context.Context) {
	pool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(3)

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	appSignal := make(chan os.Signal, 3)
	signal.Notify(appSignal, os.Interrupt)

	go func() {
		<-appSignal
		stop()
	}()

	PingDatabase(ctx)

	return ctx
}

// verify that database credentials are valid
func PingDatabase(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err := pool.PingContext(ctx)
	if err != nil {
		log.Fatalf("Unable to connect to database %v", err)
	}

}

// func ConnectToDB() {
// 	var err error
// 	pool, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))

// 	if err != nil {
// 		log.Fatal("Unable to retrieve database credentials.")
// 	}
// 	fmt.Println("Retrieving Database Credentials...")
// 	defer pool.Close()

// 	pool.SetConnMaxLifetime(0) //how many times a connection can be used
// 	pool.SetMaxIdleConns(3)    //max number of connections in idle pool
// 	pool.SetMaxOpenConns(3)    // max number of open connections in a database

// 	ctx, stop := context.WithCancel(context.Background())
// 	defer stop()

// 	appSignal := make(chan os.Signal, 3)
// 	signal.Notify(appSignal, os.Interrupt)

// 	go func() {
// 		<-appSignal
// 		stop()
// 	}()

// 	//ping database

// 	PingDatabase(ctx)
// 	fmt.Println("Connecting to Database...")
// 	fmt.Println("Connected to Database!")

// }

// // If the query fails exit the program with an error.
// func Query(ctx context.Context, id int64) {
// 	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()

// 	var name string
// 	err := pool.QueryRowContext(ctx, "select p.name from people as p where p.id = :id;", sql.Named("id", id)).Scan(&name)
// 	if err != nil {
// 		log.Fatal("unable to execute search query", err)
// 	}
// 	log.Println("name=", name)
// }
