package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"server/db"
	"server/util"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO: tighten this if you want domain restrictions
	},
}

type Client struct {
	conn *websocket.Conn
	send chan db.LocationEntry
}

// Handle a new websocket connection
func ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan db.LocationEntry, 256),
	}

	HubRegister(client)

	// Start read & write pumps
	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		HubUnregister(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		// We don’t expect clients to send anything, just keep the connection alive
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second) // ping interval
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case entry, ok := <-c.send:
			util.DebugPrint("-DISPATCHING TOSOME CLINET")
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Hub closed channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := json.Marshal(entry)
			if err != nil {
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}

		case <-ticker.C:
			// Send ping
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
