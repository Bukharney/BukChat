package ws

import (
	"fmt"
	"time"
)

type Hub struct {
	clients    map[string]map[*Client]bool
	unregister chan *Client
	register   chan *Client
	broadcast  chan Message
}

type Message struct {
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		unregister: make(chan *Client),
		register:   make(chan *Client),
		broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.RegisterNewClient(client)
		case client := <-h.unregister:
			h.RemoveClient(client)
		case message := <-h.broadcast:
			h.HandleMessage(message)
		}
	}
}

func (h *Hub) RegisterNewClient(client *Client) {
	connections := h.clients[client.ID]
	if connections == nil {
		connections = make(map[*Client]bool)
		h.clients[client.ID] = connections
	}
	h.clients[client.ID][client] = true
	h.HandleMessage(Message{
		Sender:    "Server",
		Content:   client.User.Username + " has joined the room",
		ID:        client.ID,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	})
	fmt.Println("Registered new client")
}

func (h *Hub) RemoveClient(client *Client) {
	if _, ok := h.clients[client.ID]; ok {
		delete(h.clients[client.ID], client)
		close(client.send)
		fmt.Println("Unregistered client")
	}
}

func (h *Hub) HandleMessage(message Message) {
	clients := h.clients[message.ID]
	for client := range clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients[message.ID], client)
		}
	}
}
