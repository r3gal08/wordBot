package main

// TODO: Reverse proxy should be used for handling to give a user a valid http request. Likely can containerize these things

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Create WordRequest struct (similar to classes/objects)
// Another interesting note here is we are enforcing type safety as we ensure the data we are receiving is a json string
// Go doesn’t use classes like in object-oriented languages; instead, you define structs to group related fields.
// back ticks are used for assigning struct tags to attach metadata to a struct field
// Struct tags must be enclosed in backticks (`), not quotes (""), because they are raw string literals that Go's reflection system reads at runtime.
type WordRequest struct {
	Word string `json:"word"`
}

type WordResponse struct {
	Word    string `json:"word"`
	// Definition string `json:"message"`
}

// wordHandler route (writer and a reader)
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
	defer r.Body.Close()                        // Free resources **after** wordHandler function is executed. Fun fact: Goes garbage collector is concurrent
  	log.Printf("Received word: %s", req.Word)   // SecTODO: Will want to sanitize or validate input here to prevent log injection...

    // Craft HTTP response ...
    rsp := WordResponse{Word: req.Word} // TODO: We will get word def here

    w.Header().Set("Content-Type", "application/json")  // set writer header content type (in this case json)
    w.WriteHeader(http.StatusOK)

    if err := json.NewEncoder(w).Encode(rsp); err != nil {
        log.Printf("Error encoding JSON response: %v", err)
        http.Error(w, "D'oh!", http.StatusInternalServerError)
        return
    }
    log.Printf("Response sent successfully: %v", rsp)
}

func main() {
	http.HandleFunc("/api/word", wordHandler)
	port := 8080
	log.Printf("Server started on :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
