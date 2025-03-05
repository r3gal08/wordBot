package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

// Test code for connecting to a PostgreSQL database using pgx
// This code is not used in the application
func main() {
	// Connect to the database
	// TODO: Un-hardcode the connection string
	connStr := "postgres://postgres:test@localhost:5432"
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	// Test the connection
	var greeting string
	err = conn.QueryRow(context.Background(), "SELECT 'Hello, world!'").Scan(&greeting)
	if err != nil {
		log.Fatalf("QueryRow failed: %v\n", err)
	}

	fmt.Println(greeting)
}
