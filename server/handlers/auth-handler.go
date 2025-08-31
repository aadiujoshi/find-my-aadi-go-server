package handlers

import (
	"encoding/json"
	"net/http"

	"server/config"
	"server/auth"
)

type AuthRequest struct {
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func AuthenticateHandler(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// Check password
		if req.Password != cfg.ClientPassword {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Generate JWT
		token, err := auth.GenerateJWT(cfg.JWTSecret)
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(AuthResponse{Token: token})
	}
}
