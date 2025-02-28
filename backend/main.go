package main

// TODO: Reverse proxy should be used for handling to give a user a valid http request. Likely can containerize these things

import (
	"encoding/json"
	"fmt"
	"io"
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
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

type DefinitionResponse []struct {
	Word      string `json:"word"`
	Phonetic  string `json:"phonetic"`
	Phonetics []struct {
		Text  string `json:"text"`
		Audio string `json:"audio,omitempty"`
	} `json:"phonetics"`
	Origin   string `json:"origin"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string `json:"definition"`
			Example    string `json:"example"`
			Synonyms   []any  `json:"synonyms"`
			Antonyms   []any  `json:"antonyms"`
		} `json:"definitions"`
	} `json:"meanings"`
}

// This should either return http errors or a queried word info
// We will want to turn this into an API realistically
// wordHandler route (writer and a reader)
func wordHandler(w http.ResponseWriter, r *http.Request) {

	// Reject anything that is not a POST request
	// Go encourages explict API designs.....
	if r.Method != http.MethodPost {
		log.Printf("D'oh: StatusMethodNotAllowed")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// SecTODO: Will want to sanitize or validate input here to prevent log injection...
	// Decode the JSON request body into req's memory address
	var req WordRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		log.Printf("D'oh: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // Free resources **after** wordHandler function is executed. Fun fact: Goes garbage collector is concurrent

	definition, err := getWordDefinition(req.Word)
	if err != nil {
		log.Printf("D'oh: %v", err)
		http.Error(w, "Error getting word definition", http.StatusInternalServerError)
		return
	}

	// Craft HTTP response...
	rsp := WordResponse{
		Word:       req.Word,
		Definition: definition,
	}

	w.Header().Set("Content-Type", "application/json") // set writer header content type (in this case json)
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(rsp); err != nil {
		log.Printf("D'oh: %v", err)
		http.Error(w, "Error encoding Dictionary API JSON response", http.StatusInternalServerError)
		return
	}
	log.Printf("Dictionary API Response sent successfully: %v", rsp)
}

// TODO: We will likely want some kind of "did you mean?" functionality to account for spelling mistakes.....
// TODO: Improve error handeling and error messages
func getWordDefinition(word string) (string, error) {
	url := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", word)

	response, err := http.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		if err == nil {
			err = fmt.Errorf("API returned status code: %d", response.StatusCode)
		}
		return "", err
	}
	defer response.Body.Close()

	// Dynamically allocate and read data stream. Returns []byte slice and err value
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err // Curious what err actually looks like. Test this by inducing an error sometime
	}

	var definitionResponse DefinitionResponse
	if err := json.Unmarshal(body, &definitionResponse); err != nil {
		return "", err
	}

	// TODO: This should be a seperate method likely that returns specific information
	// Avoid runtime panic by insuring an index-out-of-bound error does not occur and return the first definition
	if len(definitionResponse) > 0 && len(definitionResponse[0].Meanings) > 0 && len(definitionResponse[0].Meanings[0].Definitions) > 0 {
		def := definitionResponse[0].Meanings[0].Definitions[0].Definition
		return def, nil
	}

	return "", fmt.Errorf("Definition not found for word: %s", word)
}

func main() {
	http.HandleFunc("/api/word", wordHandler)
	port := 8080
	log.Printf("Server started on :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
