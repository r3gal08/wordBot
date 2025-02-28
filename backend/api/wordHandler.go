package api

import (
	"encoding/json"
	"log"
	"net/http"
	"wordBot/dictionary"
)

type WordRequest struct {
	Word string `json:"word"`
}

type WordResponse struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

func WordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("D'oh: StatusMethodNotAllowed")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req WordRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		log.Printf("D'oh: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	definition, err := dictionary.GetWordDefinition(req.Word)
	if err != nil {
		log.Printf("D'oh: %v", err)
		http.Error(w, "Error getting word definition", http.StatusInternalServerError)
		return
	}

	rsp := WordResponse{
		Word:       req.Word,
		Definition: definition,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(rsp); err != nil {
		log.Printf("D'oh: %v", err)
		http.Error(w, "Error encoding Dictionary API JSON response", http.StatusInternalServerError)
		return
	}
	log.Printf("Dictionary API Response sent successfully: %v", rsp)
}
