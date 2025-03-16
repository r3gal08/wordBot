package database

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "wordBot/dictionary"

    "github.com/jackc/pgx/v5"
)

// TODO: use os.Getenv("DATABASE_URL") instead of hardcoding the connection string
const DATABASE_URL = "postgres://postgres:test@localhost:5432"

func WriteWordData(wr dictionary.WordResponse) error {
    log.Printf("word Response test: %v", wr)

    // Connect to the database
    conn, err := pgx.Connect(context.Background(), DATABASE_URL)
    if err != nil {
        log.Printf("D'oh: Unable to connect to database")
        return fmt.Errorf("Unable to connect to database: %v", err)
    }
    defer conn.Close(context.Background())

    // Note: Sorta inefficient that we are re-marshaling this json here
    // 		 but it is a simple solution for now......
    // Convert wordResponse to JSON
    data, err := json.Marshal(wr)
    if err != nil {
        return fmt.Errorf("failed to marshal wordResponse: %v", err)
    }

    // Insert the word response into the database
    // Command inserts word and data into the words table
    // If the word already exists (IE: A conflict exists), it will update the data
    query := `INSERT INTO words (word, data) VALUES ($1, $2) ON CONFLICT (word) DO UPDATE SET data = EXCLUDED.data`
    _, err = conn.Exec(context.Background(), query, wr.Word, data)
    if err != nil {
        log.Printf("D'oh: Insert failed")
        return fmt.Errorf("Insert failed: %v", err)
    }

    log.Println("Insert successful")
    return nil
}