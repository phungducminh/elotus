package auth

import (
	"net/http"
	"strings"

	"elotus.com/hackathon/server"
)

func AuthMiddleware(s *server.Server, handler http.Handler) http.Handler {
	au := NewAuth(s.Logger, s.Storage, []byte(s.Cfg.AuthSecretKey), s.Cfg.MysqlConnMaxLifetimeInSeconds)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.Logger.Warn("Authorization header missing")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			s.Logger.Warn("Invalid authorization header format")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		req := &VerifyRequest{
			AccessToken: token,
		}
		_, err := au.Verify(req)
		if err != nil {
			s.Logger.Warn("Invalid or expired token")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
