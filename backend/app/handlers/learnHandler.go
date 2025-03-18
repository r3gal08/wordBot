package handlers

import (
	"context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
	"wordBot/utils"

	"github.com/ollama/ollama/api"
)

// TODO: Not currently used
// OllamaRequest represents the request payload for the LLM
type ollamaRequest struct {
    Model  string `json:"model"`
    System string `json:"system"`
    Prompt string `json:"prompt"`
}

// TODO: Not currently used
// OllamaResponse represents the response from the LLM
type ollamaResponse struct {
    Response string `json:"response"`
}

// TODO: Ollama end point environment variables will need to be set (It currently appears to use the default http://localhost:11434)
// Good starting point for next time im on. Basically want to adjust this to allow my request to come through
func LearnHandler(w http.ResponseWriter, r *http.Request) {
    req, err := utils.DecodeWordRequest(w, r)
    if err != nil {
        return
    }

	client, err := api.ClientFromEnvironment()
	if err != nil {
		http.Error(w, "Failed to create ollama API client", http.StatusInternalServerError)
		log.Printf("Error creating ollama API client: %v", err)
		return
	}

	stream := false
	llmReq := &api.GenerateRequest {
        Model:  "llama3.1",
        System: "Provide a word definition with four possible answers (one correct and three incorrect). The correct answer should be clearly marked. Format the response as follows:\n\nWord: [word]\nDefinition: [definition]\n\nAnswers:\n1. [answer 1]\n2. [answer 2]\n3. [answer 3]\n4. [answer 4]\n\nCorrect Answer: [answer number]",
		Stream: &stream,
		Prompt: req.Word,
	}

	// Create an empty context so that the request can be cancelled if needed and processed asynchronously
	/*	If you wanted to set a timeout:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
	*/
	// TODO: Currently no error handeling or data validation and sanitization
	ctx := context.Background()

	// Slice to accumulate responses
    var responses []api.GenerateResponse	
    respFunc := func(resp api.GenerateResponse) error {
        // Accumulate the responses
        responses = append(responses, resp)
        return nil
    }
	
	err = client.Generate(ctx, llmReq, respFunc)
	if err != nil {
        http.Error(w, "Failed to generate response from ollama API", http.StatusInternalServerError)
        log.Printf("Error generating response from ollama API: %v", err)
        return
	}

    // Marshal the accumulated responses to JSON
    responseJSON, err := json.MarshalIndent(responses, "", "  ")
    if err != nil {
        http.Error(w, "Failed to marshal accumulated responses", http.StatusInternalServerError)
        log.Printf("Error marshaling accumulated responses: %v", err)
        return
    }

    // Print the JSON to the terminal
    fmt.Println(string(responseJSON))

    // Write the JSON to the HTTP response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(responseJSON)
}
