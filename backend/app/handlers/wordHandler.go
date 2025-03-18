package handlers

/*
TODO: - Move HTTP handler and routing implementation to Gin framework. It will make the API faster and more efficient
	    which will allow it to scale better and decrease potential costs of running the server
	  - Make common error handeling helper functions to reduce code overhead within this file
*/

import (
	"encoding/json"
	"log"
	"net/http"
	"wordBot/dictionary"
	"wordBot/database"
	"wordBot/utils"
)

/* TODOS:
	- Implement input sanitization to prevent SQL injection and ensure valid JSON format.
	- Create common error handling helper functions to reduce code overhead within this file
*/
func WordHandler(w http.ResponseWriter, r *http.Request) {
    req, err := utils.DecodeWordRequest(w, r)
    if err != nil {
        return
    }

	rsp := dictionary.WordResponse{
		Word: req.Word,
	}

    // Check if the word is new
    isNewWord, err := database.IsNewWord(req.Word)
    if err != nil {
        log.Printf("Error checking for word existence: %v", err)
        http.Error(w, "Error checking for word existence", http.StatusInternalServerError)
        return
    }

	// Get the word data from the dictionary API
	wordData, err := dictionary.GetWordData(req.Word)
	if err != nil {
		log.Printf("D'oh: %v", err)
		http.Error(w, "Error getting word definition", http.StatusInternalServerError)
		return
	}

	// Construct the response based on the requested attributes
	// "range" returns both the index and the value of the slice. '_' is used to ignore the index value
	// TODOS: 
	// 		- Implement input sanitization to prevent SQL injection and ensure valid JSON format.
	//      - Add error handling for database connection issues and invalid word attributes.
	//      - Handle cases where multiple definitions are returned for a word/other attributes
	for _, attr := range req.Request {
		switch attr {
		case "definition":
			rsp.Definition = wordData[0].Meanings[0].Definitions[0].Definition // This indexing should really be done in the dictionary api package
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

	// Writing word data to the database
	if isNewWord {
		rsp.ConfidenceRating = 1 // New word so set confidence rating to 1
		if err := database.WriteWordData(rsp); err != nil {
			log.Printf("D'oh: %v", err)
			http.Error(w, "Error writing to database", http.StatusInternalServerError)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(rsp); err != nil {
		log.Printf("D'oh: %v", err)
		http.Error(w, "Error encoding Dictionary API JSON response", http.StatusInternalServerError)
		return
	}
	log.Printf("Dictionary API Response sent successfully: %v", rsp)
}
