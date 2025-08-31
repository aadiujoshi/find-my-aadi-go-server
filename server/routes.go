package main

import (
	"github.com/gorilla/mux"
	"server/config"
	"server/handlers"
)

func SetupRoutes(cfg config.Config) *mux.Router {
	r := mux.NewRouter()

	// Public
	r.HandleFunc("/api/authenticate-client", handlers.AuthenticateHandler(cfg)).Methods("POST")
	
	// Authenticated (JWT required)
	auth := r.PathPrefix("/api").Subrouter()
	auth.Use(JWTAuthMiddleware(cfg.JWTSecret))
	auth.HandleFunc("/get-range", handlers.GetRangeHandler(cfg)).Methods("GET")
	auth.HandleFunc("/get-live-updates", handlers.LiveUpdatesHandler(cfg)).Methods("GET")

	// Admin only
	admin := r.PathPrefix("/api").Subrouter()
	admin.Use(AdminAuthMiddleware(cfg.AdminPassword))
	admin.HandleFunc("/new-location", handlers.NewLocationHandler(cfg)).Methods("POST")

	return r
}
