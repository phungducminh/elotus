package main

import (
	"fmt"
	"log"
	"net/http"

	"elotus.com/hackathon/auth"
	"elotus.com/hackathon/file"
	"elotus.com/hackathon/server"
)

func main() {
	// TODO: load from config file
	serverCfg := &server.Config{
		AuthSecretKey:                 "SECRET-KEY",
		TokenExpiresInSeconds:         60,
		MysqlDatasource:               "root:elotus@tcp(localhost:3306)/elotus",
		MysqlConnMaxLifetimeInSeconds: 0, // conn are not closed due to a connection's age
		MysqlMaxOpenConns:             1000,
		MysqlMaxIdleConns:             200,
		ServerPort:                    8080,
	}
	server := server.NewServer(serverCfg)
	loginHandler := auth.NewLoginHandler(server)
	http.Handle("/api/auth/login", loginHandler)
	registerHandler := auth.NewRegisterHandler(server)
	http.Handle("/api/auth/register", registerHandler)
	fileHandler := file.NewFileHandler(server)
	http.Handle("/api/file/upload", fileHandler)

	log.Printf("Starting server on %d...", serverCfg.ServerPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", serverCfg.ServerPort), nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
