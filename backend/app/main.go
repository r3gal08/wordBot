package main

import (
	"fmt"
	"log"
	"net/http"
	"wordBot/handlers"
)

func registerRoutes() *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/wordHandler", handlers.WordHandler)
    mux.HandleFunc("/learnHandler", handlers.LearnHandler)
    return mux
}

func main() {
    mux := registerRoutes()
    port := 8080
    log.Printf("Server started on :%d", port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
