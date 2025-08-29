package main

import (
	"fmt"
	"net/http"
	"os"

	"elotus.com/hackathon/auth"
	"elotus.com/hackathon/file"
	"elotus.com/hackathon/server"
)

func main() {
	// TODO: load from config file
	server, err := server.NewServer()
	if err != nil {
		fmt.Printf("error when create server: %v", err)
		os.Exit(-1)
	}
	loginHandler := auth.NewLoginHandler(server)
	http.Handle("/api/auth/login", loginHandler)
	registerHandler := auth.NewRegisterHandler(server)
	http.Handle("/api/auth/register", registerHandler)
	fileHandler := file.NewFileHandler(server)
	http.Handle("/api/file/upload", auth.AuthMiddleware(server, fileHandler))
	err = server.Serve()
	if err != nil {
		fmt.Printf("error when start the server: %v", err)
		os.Exit(-1)
	}
}
