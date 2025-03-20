package handlers

/*
TODOs: 
    - Currently no error handeling or data validation and sanitization
    - TODO: Ollama end point environment variables will need to be set (It currently appears to use the default http://localhost:11434)
*/

import (
	"context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
	"wordBot/utils"

	"github.com/ollama/ollama/api"
)

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

    // Honestly this prompt is a great start. Should export this as a var, and allow for multiple types of system responses
	stream := false
    llmReq := &api.GenerateRequest {
        Model:  "llama3.1",
        System: `You are an AI that generates a structured JSON response for word-based multiple-choice questions.
    Your response must always follow this exact JSON format, with no additional text or explanations:
    
    {
        "answers": [
            "answer_1",
            "answer_2",
            "answer_3",
            "answer_4"
        ],
        "correct_answer": correct_answer_id
    }
    
    Rules:
    - The "answers" array must contain exactly four unique options, one of which is correct.
    - The "correct_answer" field must match the correct answer from the "answers" array exactly.
    - correct_answer_id must be one of the four answer options in numerical format indexed at 0
    - Return only valid JSON without markdown formatting, explanations, or extra characters.`,
    
        Stream: &stream,
        Prompt: req.Word,
    }
    
	// Create an empty context so that the request can be cancelled if needed and processed asynchronously
	ctx := context.Background()
    /*	If you wanted to set a timeout:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
	*/

    var llmRsp []byte
    respFunc := func(resp api.GenerateResponse) error {
        if resp.Response == "" {
            return fmt.Errorf("LLM response is empty")
        }
    
        // Creating a key-value map of type string-interface
        // Try parsing the response to make sure it's valid JSON. Maybe a more effiecient way to do this? marshling is expensive...
        var parsedResponse map[string]interface{}
        if err := json.Unmarshal([]byte(resp.Response), &parsedResponse); err != nil {
            return fmt.Errorf("LLM response is not valid JSON: %v", err)
        }
        
        llmRsp = []byte(resp.Response)
        log.Printf("Valid Learn request response:\n %v\n\n", resp.Response)
        return nil
    }

	err = client.Generate(ctx, llmReq, respFunc)
	if err != nil {
        http.Error(w, "Failed to generate response from LLM", http.StatusInternalServerError)
        log.Printf("Error generating response from LLM: %v", err)
        return
	}

    // Send the valid JSON response directly
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(llmRsp)
}
