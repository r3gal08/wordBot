package utils

import (
    "encoding/json"
    "log"
	"fmt"
    "net/http"
)

// Struct tags such as json:"word" specify what a field’s name should be when the struct’s
// contents are serialized into JSON. Without them, the JSON would use the struct’s
// capitalized field names – a style not as common in JSON.
type wordRequest struct {
	Word    string   `json:"word"`
	Request []string `json:"request"`
}

// DecodeRequest decodes the JSON request body into a WordRequest struct
func DecodeWordRequest(w http.ResponseWriter, r *http.Request) (*wordRequest, error) {
    if r.Method != http.MethodPost {
        log.Printf("D'oh: StatusMethodNotAllowed")
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return nil, fmt.Errorf("invalid request method")
    }

    var req wordRequest
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&req); err != nil {
        log.Printf("D'oh: %v", err)
        http.Error(w, "Bad request", http.StatusBadRequest)
        return nil, err
    }
    defer r.Body.Close()

    return &req, nil
}