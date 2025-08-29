package findmyaadigoserver

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} 
var clients = make(map[*websocket.Conn]bool)

func handleNewClient(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil);
	if err != nil {
		
	}
	defer ws.Close()
	clients[ws] = true

	for {
		
	}
}

func initClientPort() {
	http.HandleFunc("/ws", handleNewClient)
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
