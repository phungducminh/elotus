package main

import (
	"log"
	"net/http"

	"elotus.com/hackathon/auth"
	"elotus.com/hackathon/server"
)

func main() {
	// TODO: load from config file
	serverCfg := &server.Config{
		AuthSecretKey:                 "SECRET-KEY",
		MysqlDatasource:               "root:elotus@tcp(localhost:3306)/elotus",
		MysqlConnMaxLifetimeInSeconds: 0, // conn are not closed due to a connection's age
		MysqlMaxOpenConns:             1000,
		MysqlMaxIdleConns:             200,
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
