package main

import (
	"log"
	"net/http"

	"elotus.com/hackathon/auth"
	"elotus.com/hackathon/server"
)

func main() {
	// TODO: load from config file
	serverCfg := &server.ServerConfig{
		AuthSecretKey: "SECRET-KEY",
	}
	server := server.NewServer(serverCfg)
	loginHandler := auth.NewLoginHandler(server)
	http.Handle("/api/auth/login", loginHandler)
	registerHandler := auth.NewRegisterHandler(server)
	http.Handle("/api/auth/register", registerHandler)

	port := ":8080"
	log.Printf("Starting server on %s...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
