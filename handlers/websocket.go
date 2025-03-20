package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections for now
	},
}

// Global WebSocket connections list
var clients = make(map[*websocket.Conn]bool)

// Broadcast channel
var broadcast = make(chan string)

// HandleWebSocket upgrades HTTP to WebSocket and listens for messages
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP request to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Register new client
	clients[conn] = true

	// Listen for messages
	for msg := range broadcast {
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println("Error writing to WebSocket:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

// NotifyWebSocketClients sends updates when a bike is locked/unlocked
func NotifyWebSocketClients(event string) {
	broadcast <- event
}
