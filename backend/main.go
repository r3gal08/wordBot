package main

// TODO: Reverse proxy should be used for handling to give a user a valid http request. Likely can containerize these things

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Go doesn’t use classes like in object-oriented languages; instead, you define structs to group related fields.
// back ticks are used for assigning struct tags to attach metadata to a struct field
// Struct tags must be enclosed in backticks (`), not quotes (""), because they are raw string literals that Go's reflection system reads at runtime.
type WordRequest struct {
	Word string `json:"word"`
}

func wordHandler(w http.ResponseWriter, r *http.Request) {

    // Reject anything that is not a POST request
    // Go encourages explict API designs.....
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

    // decode the JSON request body into req's memory address
	var req WordRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()    // Free resources **after** wordHandler function is executed

	log.Printf("Received word: %s", req.Word)

    // Respond to http request
    // TODO: It is at this point we should be returning contents back to the user such as word definition
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Word received successfully"})
}

func main() {
	http.HandleFunc("/api/word", wordHandler)
	port := 8080
	log.Printf("Server started on :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
