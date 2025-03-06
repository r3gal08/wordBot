package api

/*
TODO: - Move HTTP handler and routing implementation to Gin framework. It will make the API faster and more efficient
	    which will allow it to scale better and decrease potential costs of running the server
	  - Make common error handeling helper functions to reduce code overhead within this file
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"wordBot/dictionary"

	"github.com/jackc/pgx/v5"
)

// Struct tags such as json:"word" specify what a field’s name should be when the struct’s
// contents are serialized into JSON. Without them, the JSON would use the struct’s
// capitalized field names – a style not as common in JSON.
type wordRequest struct {
	Word    string   `json:"word"`
	Request []string `json:"request"`
}

// Can add in additional fields to the struct as needed
type wordResponse struct {
	Word         string `json:"word,omitempty"`
	Definition   string `json:"definition,omitempty"`
	PartOfSpeech string `json:"partofspeech,omitempty"`
}

// TODO: use os.Getenv("DATABASE_URL") instead of hardcoding the connection string
const DATABASE_URL = "postgres://postgres:test@localhost:5432"

func writeWordData(wr wordResponse) error {
	log.Printf("word Response test: %v", wr)

	// Connect to the database
	conn, err := pgx.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		log.Printf("D'oh: Unable to connect to database")
		return fmt.Errorf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Note: Sorta innefficent that we are re-marshaling this json here
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

func WordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("D'oh: StatusMethodNotAllowed")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req wordRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		log.Printf("D'oh: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	wordData, err := dictionary.GetWordData(req.Word)
	if err != nil {
		log.Printf("D'oh: %v", err)
		http.Error(w, "Error getting word definition", http.StatusInternalServerError)
		return
	}

	rsp := wordResponse{
		Word: wordData[0].Word,
	}

	// Construct the response based on the requested attributes
	// "range" returns both the index and the value of the slice. '_' is used to ignore the index value
	// TODO: Implement input sanitization to prevent SQL injection and ensure valid JSON format.
	//       Add error handling for database connection issues and invalid word attributes.
	// Can add in additional cases here as needed
	for _, attr := range req.Request {
		switch attr {
		case "definition":
			rsp.Definition = wordData[0].Meanings[0].Definitions[0].Definition
		case "partofspeech":
			rsp.PartOfSpeech = wordData[0].Meanings[0].PartOfSpeech
		default:
			log.Printf("Unknown attribute requested: %s", attr)
			http.Error(w, "Error getting word data", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// TODO: This should happen somewhere else.... but we are testing b'ye
	// Writing word data to the database
	if err := writeWordData(rsp); err != nil {
		log.Printf("D'oh: %v", err)
		http.Error(w, "Error writing to database", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(rsp); err != nil {
		log.Printf("D'oh: %v", err)
		http.Error(w, "Error encoding Dictionary API JSON response", http.StatusInternalServerError)
		return
	}
	log.Printf("Dictionary API Response sent successfully: %v", rsp)
}
