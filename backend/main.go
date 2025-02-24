package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type WordRequest struct {
	Word string `json:"word"`
}

func wordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req WordRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Received word: %s", req.Word)

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
