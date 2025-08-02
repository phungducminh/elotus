package main

import (
	"log"
	"net/http"

	"elotus.com/hackathon/auth"
)



func main() {
	// http.HandleFunc("/api/auth/register", auth.RegisterHandler)
	loginHandler := &auth.LoginHandler {

	}
	http.Handle("/api/auth/login", loginHandler)
	registerHandler := auth.NewRegisterHandler()
	http.Handle("/api/auth/register", registerHandler)

	port := ":8080"
	log.Printf("Starting server on %s...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
