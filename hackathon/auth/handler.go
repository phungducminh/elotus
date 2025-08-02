package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"elotus.com/hackathon/db"
)

type RegisterHandler struct {
	dbClient db.Client
}

func NewRegisterHandler() *RegisterHandler {
	dbConfig := &db.Config{}
	dbClient := db.NewClient(dbConfig)

	return &RegisterHandler{
		dbClient: dbClient,
	}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
	var req RegisterRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// TODO: handle password hashing
	userRecord := db.UserRecord{
		ID:       0,
		Username: req.Username,
		Password: req.Password,
	}
	id, err := h.dbClient.InsertUser(userRecord)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := &RegisterResponse{
		UserId: strconv.FormatInt(id, 10),
	}
	json.NewEncoder(w).Encode(res)
}

type LoginHandler struct {
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{ "message" : "Login" }`)
}
