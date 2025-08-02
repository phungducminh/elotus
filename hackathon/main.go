package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"elotus.com/hackathon/auth"
)

// Handler for GET /hello
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello, world!")
}

// Handler for POST /echo
func echoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("/auth/register", auth.RegisterHandler)
	http.HandleFunc("/auth/login", auth.LoginHandler)

	port := ":8080"
	log.Printf("Starting server on %s...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
