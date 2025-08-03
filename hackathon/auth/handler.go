package auth

import (
	"encoding/json"
	"net/http"

	"elotus.com/hackathon/server"
)

type RegisterHandler struct {
	server *server.Server
	auth   Auth
}

func NewRegisterHandler(s *server.Server) *RegisterHandler {
	return &RegisterHandler{
		auth: NewAuth(s.Storage, []byte(s.Cfg.AuthSecretKey), s.Cfg.TokenExpiresInSeconds),
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

	resp, err := h.auth.Register(&req)
	if err == ErrUsernameNotUnique {
		w.Header().Set("Content-Type", "application/json")
		resp := &Response{
			Error: ErrorResponse{
				Code:    "USERNAME_NOT_UNIQUE",
				Message: "username is not unique",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	} else if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

type LoginHandler struct {
	auth Auth
}

func NewLoginHandler(s *server.Server) *LoginHandler {
	return &LoginHandler{
		auth: NewAuth(s.Storage, []byte(s.Cfg.AuthSecretKey), s.Cfg.TokenExpiresInSeconds),
	}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	resp, err := h.auth.Login(&req)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

type Response struct {
	Error ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
