package auth

import (
	"net/http"
	"strings"

	"elotus.com/hackathon/server"
)

func AuthMiddleware(s *server.Server, handler http.Handler) http.Handler {
	au := NewAuth(s.Storage, []byte(s.Cfg.AuthSecretKey), s.Cfg.MysqlConnMaxLifetimeInSeconds)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		req := &VerifyRequest{
			AccessToken: token,
		}
		_, err := au.Verify(req)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
