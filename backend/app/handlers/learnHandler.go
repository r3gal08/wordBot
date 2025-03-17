package handlers

import (
	"context"
    //"encoding/json"
    "fmt"
    "log"
    "net/http"

	"github.com/ollama/ollama/api"
)

// Potentially make private?
// OllamaRequest represents the request payload for the LLM
type OllamaRequest struct {
    Model  string `json:"model"`
    System string `json:"system"`
    Prompt string `json:"prompt"`
}

// Potentially make private?
// OllamaResponse represents the response from the LLM
type OllamaResponse struct {
    Response string `json:"response"`
}

// Good starting point for next time im on. Basically want to adjust this to allow my request to come through
func LearnHandler(w http.ResponseWriter, r *http.Request) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// Not parsing this quite right
	req := &api.GenerateRequest{
        Model:  "llama3.1",
        System: "Provide a word definition with four possible answers (one correct and three incorrect). The correct answer should be clearly marked. Format the response as follows:\n\nWord: [word]\nDefinition: [definition]\n\nAnswers:\n1. [answer 1]\n2. [answer 2]\n3. [answer 3]\n4. [answer 4]\n\nCorrect Answer: [answer number]",
        Prompt: "hubris", // TODO: Should be an input from the client
	}

	ctx := context.Background()
	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.
		fmt.Println(resp.Response)
		return nil
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}
}

// // LearnHandler queries the locally running LLM and returns the response
// func LearnHandlerTest(w http.ResponseWriter, r *http.Request) {
//     log.Println("LearnHandler called")

//     // Define the prompt
//     requestBody := OllamaRequest{
//         Model:  "llama3.1",
//         System: "Provide a word definition with four possible answers (one correct and three incorrect). The correct answer should be clearly marked. Format the response as follows:\n\nWord: [word]\nDefinition: [definition]\n\nAnswers:\n1. [answer 1]\n2. [answer 2]\n3. [answer 3]\n4. [answer 4]\n\nCorrect Answer: [answer number]",
//         Prompt: "hubris", // TODO: Should be an input from the client
//     }

//     // Convert the request to JSON
//     jsonData, err := json.Marshal(requestBody)
//     if err != nil {
//         http.Error(w, "Failed to encode request", http.StatusInternalServerError)
//         log.Printf("Error encoding JSON: %v", err)
//         return
//     }

//     // Use the ollama package to send the request to the LLM
//     client := ollama.NewClient("http://localhost:11434") // Replace with your LLM's API endpoint
//     resp, err := client.Generate(jsonData)
//     if err != nil {
//         http.Error(w, "Failed to contact LLM", http.StatusInternalServerError)
//         log.Printf("Error contacting LLM: %v", err)
//         return
//     }

//     // Parse the JSON response
//     var llmResponse OllamaResponse
//     if err := json.Unmarshal(resp, &llmResponse); err != nil {
//         http.Error(w, "Failed to parse LLM response", http.StatusInternalServerError)
//         log.Printf("Error decoding LLM response: %v", err)
//         return
//     }

//     // Return LLM response to client
//     w.Header().Set("Content-Type", "application/json")
//     w.WriteHeader(http.StatusOK)
//     fmt.Fprint(w, llmResponse.Response)
// }