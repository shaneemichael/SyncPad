package main

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Client = user connected to the hub
type Client struct {
	ID   string
	Conn *websocket.Conn
	Hub  *Hub
}

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
}

// Message = structure of messages exchanged in the hub
type Message struct {
	Type     string `json:"type"`
	Data     string `json:"data"`
	SenderID string `json:"sender_id"`
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
			fmt.Println("Client registered:", client.ID)
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				client.Conn.Close()
			}
			h.mu.Unlock()
		case message := <-h.Broadcast:
			h.mu.Lock()
			for client := range h.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					client.Conn.Close()
					delete(h.Clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}
