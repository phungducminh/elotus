package auth

import (
	"encoding/json"
	"net/http"

	. "elotus.com/hackathon/pkg/logutil/httputil"
	"elotus.com/hackathon/server"
)

type RegisterHandler struct {
	server *server.Server
	auth   Auth
}

func NewRegisterHandler(s *server.Server) *RegisterHandler {
	return &RegisterHandler{
		auth: NewAuth(s.Logger, s.Storage, []byte(s.Cfg.AuthSecretKey), s.Cfg.TokenExpiresInSeconds),
	}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ResponseMethodNotAllowed(w)
		return
	}

	defer r.Body.Close()
	var req RegisterRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		ResponseBadRequest(w, "INVALID_JSON", "invalid json")
		return
	}

	resp, err := h.auth.Register(&req)
	if err == ErrUsernameNotUnique {
		ResponseBadRequest(w, "USERNAME_NOT_UNIQUE", "username is not unique")
		return
	} else if err != nil {
		ResponseInternalServerError(w)
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
		auth: NewAuth(s.Logger, s.Storage, []byte(s.Cfg.AuthSecretKey), s.Cfg.TokenExpiresInSeconds),
	}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ResponseMethodNotAllowed(w)
		return
	}

	var req LoginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		ResponseBadRequest(w, "INVALID_JSON", "invalid json")
		return
	}

	resp, err := h.auth.Login(&req)
	if err != nil {
		ResponseInternalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
