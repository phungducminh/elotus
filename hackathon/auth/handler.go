package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"elotus.com/hackathon/db"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
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

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
	userRecord := db.UserRecord{
		ID:       0,
		Username: req.Username,
		Password: string(hashed),
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

	userRecord := &db.UserRecord{
		ID:       100,
		Username: "elotus",
		Password: "elotus",
	}

	key := []byte("SECRET")
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "elotus",
		Subject:   userRecord.Username,
	})
	s, err := t.SignedString(key)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{ \"message\" : \"%s\" }\n", s)
}
