package findmyaadigoserver

import (
	// "database/sql"
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)


func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection to WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer ws.Close()
	clients[ws] = true

	for {
		var msg map[string]any
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}

		// Broadcast to all clients
		for client := range clients {
			client.WriteJSON(msg)
		}
	}
}

func main() {
}
