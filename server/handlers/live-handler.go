package handlers

import (
	"net/http"

	"server/config"
	"server/ws"
)

func LiveUpdatesHandler(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(w, r)
	}
}
