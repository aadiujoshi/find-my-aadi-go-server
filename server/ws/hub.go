package ws

import (
	"log"
	"server/db"
	"server/util"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan db.LocationEntry
	register   chan *Client
	unregister chan *Client
}

var hub *Hub

// Initialize hub (called once in main.go, before starting server)
func InitHub() {
	hub = &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan db.LocationEntry),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	// util.DebugPrint("Send dummy location update243432432")
	// go sendDummyUpdates()
	go hub.run()
}

// func sendDummyUpdates()  {
// 	// util.DebugPrint("Send dummy location update")
// 	for {
// 		time.Sleep(5 * time.Second)
// 		HubBroadcast(db.LocationEntry{
// 			Timestamp: time.Now().Unix(),
// 			Latitude:  37.7749,
// 			Longitude: -122.4194,
// 		})
// 	}
// }

func (h *Hub) run() {
	util.DebugPrint("Starting Hub...")
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client connected, total: %d", len(h.clients))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client disconnected, total: %d", len(h.clients))
			}
		case entry := <-h.broadcast:
			util.DebugPrint("SENDING ENTRY TO CLIENTS, NEW BROADCAST RECEIVED")
			for client := range h.clients {
				select {
				case client.send <- entry:
				default:
					// Drop client if send fails
					delete(h.clients, client)
					close(client.send)
				}
			}
		}
	}
}

// Called from NewLocationHandler to push update
func HubBroadcast(entry db.LocationEntry) {
	util.DebugPrint("--------SEDING NEW HUB BROADCAST")
	if hub != nil {
		hub.broadcast <- entry
	}
}

// Called from LiveUpdatesHandler to attach a client
func HubRegister(c *Client) {
	if hub != nil {
		hub.register <- c
	}
}

// Called from client disconnect
func HubUnregister(c *Client) {
	if hub != nil {
		hub.unregister <- c
	}
}
