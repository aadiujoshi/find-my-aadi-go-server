package main

import (
	"log"
	"net/http"
	"server/config"
	"server/db"
	"server/ws"
)

func main() {
	// Load config (env vars)
	cfg := config.LoadConfig()
	db.InitDb()
	// Initialize router
	router := SetupRoutes(cfg)
	ws.InitHub()
	
	// Start server
	addr := ":" + cfg.Port
	log.Printf("Starting server on %s...", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
	
	db.SaveToDisk()
}
