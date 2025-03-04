package main

import (
	"fmt"
	"log"
	"net/http"
	"wordBot/api"
)

func main() {
	http.HandleFunc("/api/word", api.WordHandler)
	port := 8080
	log.Printf("Server started on :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
