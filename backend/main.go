package main

import (
	"backend/utils"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	hub := NewHub()
	go hub.Run() // Run the hub in a background goroutine
	logger := utils.NewDefaultLogger(true)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Upgrade error:", err)
			return
		}
		client := &Client{Hub: hub, Conn: conn, ID: r.RemoteAddr}
		hub.Register <- client

		for {
			var msg Message
			if err := conn.ReadJSON(&msg); err != nil {
				hub.Unregister <- client
				break
			}
			msg.SenderID = client.ID
			hub.Broadcast <- msg
		}
	})
	logger.Info("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
