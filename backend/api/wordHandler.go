package api

/*
TODO: Move implementation to gin framework. It will make the API faster and more efficient
	  which will allow it to scale better and decrease potential costs of running the server
*/

import (
	"encoding/json"
	"log"
	"net/http"
	"wordBot/dictionary"
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
	// TODO: Input sanitization and error handeling
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

	if err := json.NewEncoder(w).Encode(rsp); err != nil {
		log.Printf("D'oh: %v", err)
		http.Error(w, "Error encoding Dictionary API JSON response", http.StatusInternalServerError)
		return
	}
	log.Printf("Dictionary API Response sent successfully: %v", rsp)
}
