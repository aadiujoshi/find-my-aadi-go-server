package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userContextKey = contextKey("user")

// Middleware: JWT auth for client routes
func JWTAuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// Inject user info into context (optional)
			ctx := context.WithValue(r.Context(), userContextKey, token.Claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Middleware: admin password check (header-based)
func AdminAuthMiddleware(adminPassword string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pass := r.Header.Get("X-Admin-Password")
			if pass == "" || pass != adminPassword {
				http.Error(w, "unauthorized (admin required)", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
