package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"server/config"
	"server/db"
	"server/util"
	"server/ws"
)

func GetRangeHandler(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startStr := r.URL.Query().Get("start")
		endStr := r.URL.Query().Get("end")

		start, err1 := strconv.ParseInt(startStr, 10, 64)
		end, err2 := strconv.ParseInt(endStr, 10, 64)
		if err1 != nil || err2 != nil {
			http.Error(w, "invalid timestamp params", http.StatusBadRequest)
			return
		}

		entries, err := db.GetEntryRange(start, end)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(entries)
	}
}

type NewLocationRequest struct {
	Timestamp int64   `json:"timestamp"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewLocationHandler(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req NewLocationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		entry := db.LocationEntry{
			Timestamp: req.Timestamp,
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
		}

		util.DebugPrint("Added New location Entry");

		// Save to DB
		if err := db.AddEntry(entry.Timestamp, entry.Latitude, entry.Longitude); err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		util.DebugPrint("Added New location Entry to Db");
		util.DebugPrint("Notifieying clients...");

		// Notify websocket clients
		ws.HubBroadcast(entry)

		w.WriteHeader(http.StatusCreated)
	}
}
