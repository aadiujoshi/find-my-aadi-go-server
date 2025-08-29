package findmyaadigoserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func handlePostLocation(w http.ResponseWriter, r *http.Request) {
	// Read the entire body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Hydrate into struct
	var loc LocationEntry
	err = json.Unmarshal(bodyBytes, &loc)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received location: %+v\n", loc)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Location received"))

}

func updateLocationHistoryDb(entry LocationEntry) {
	err := addEntry(entry.timestamp, entry.latitude, entry.longitude);
	if err != nil {
		fmt.Printf("")
	}
}

func initOriginPort() {
	http.HandleFunc("/https", handlePostLocation)
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
