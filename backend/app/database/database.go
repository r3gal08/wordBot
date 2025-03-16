package database

/* TODOs:
    - Make common error handeling helper functions to reduce code overhead within this file
    - Implement input sanitization to prevent SQL injection and ensure valid JSON format.
    - Add error handling for database connection issues and invalid word attributes.
    - Handle cases where multiple definitions are returned for a word/other attributes
*/

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

// Note: This is currently broken as unmarshalling will not work with the current database schema
// 		 This is because the data is stored as a byte array and not as a JSON object
// 		 This will be fixed in the future
func GetWordData(word string) (dictionary.WordData, error) {    
    // Connect to the database
    conn, err := pgx.Connect(context.Background(), DATABASE_URL)
    if err != nil {
        log.Printf("D'oh: Unable to connect to word database")
        return nil, fmt.Errorf("Unable to connect to word database: %v", err)
    }
    defer conn.Close(context.Background())

    // Query the word data from the database
    query := `SELECT data FROM words WHERE word = $1`
    var data []byte
    err = conn.QueryRow(context.Background(), query, word).Scan(&data)
    if err != nil {
        if err == pgx.ErrNoRows {
            log.Printf("Word not found in DB: %v", err)
            return nil, nil // Return nil, nil to indicate the word was simply not found
        }
        log.Printf("Query failed: %v", err)
        return nil, fmt.Errorf("Query failed: %v", err)
    }

    // Just keep in mind that marshalling and unmarshalling JSON is expensive in terms of memory and CPU. 
    // A better method may be to leave the handeling of this marshalling on the front end (this way we offload that work to the client)
    // Unmarshal the JSON data into a WordData struct
    var wd dictionary.WordData
    err = json.Unmarshal(data, &wd)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal wordData: %v", err)
    }

    log.Printf("Retrieved word data from internal DB: %v", wd)
    return wd, nil
}

func IsNewWord(word string) (bool, error) {
    // Connect to the database
    conn, err := pgx.Connect(context.Background(), DATABASE_URL)
    if err != nil {
        log.Printf("D'oh: Unable to connect to word database")
        return false, fmt.Errorf("Unable to connect to word database: %v", err)
    }
    defer conn.Close(context.Background())

    // Query the word data from the database
    query := `SELECT 1 FROM words WHERE word = $1`
    var exists int
    err = conn.QueryRow(context.Background(), query, word).Scan(&exists)
    if err != nil {
        if err == pgx.ErrNoRows {
            // Word not found in DB
            return true, nil
        }
        log.Printf("Query failed: %v", err)
        return false, fmt.Errorf("Query failed: %v", err)
    }

    // Word found in DB
    return false, nil
}