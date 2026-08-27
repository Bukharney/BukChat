package ws

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type Hub struct {
	clients     map[string]map[*Client]bool
	userClients map[int]map[*Client]bool
	unregister  chan *Client
	register    chan *Client
	broadcast   chan Message
	pubSub      PubSubAdapter
	mu          sync.RWMutex
}

type Message struct {
	Type      string      `json:"type,omitempty"`
	Sender    int         `json:"sender"`
	Content   string      `json:"content"`
	ID        string      `json:"id"`
	Timestamp string      `json:"timestamp"`
	IsTyping  bool        `json:"is_typing,omitempty"`
	Payload   interface{} `json:"payload,omitempty"`
}

func NewHub(pubSub ...PubSubAdapter) *Hub {
	var ps PubSubAdapter = NewInMemoryPubSub()
	if len(pubSub) > 0 && pubSub[0] != nil {
		ps = pubSub[0]
	}

	return &Hub{
		clients:     make(map[string]map[*Client]bool),
		userClients: make(map[int]map[*Client]bool),
		unregister:  make(chan *Client),
		register:    make(chan *Client),
		broadcast:   make(chan Message),
		pubSub:      ps,
	}
}

func (h *Hub) Run() {
	slog.Info("WebSocket Hub event loop started")
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

func (h *Hub) IsUserOnline(userID int) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	conns, ok := h.userClients[userID]
	return ok && len(conns) > 0
}

func (h *Hub) RegisterNewClient(client *Client) {
	h.mu.Lock()
	connections := h.clients[client.ID]
	if connections == nil {
		connections = make(map[*Client]bool)
		h.clients[client.ID] = connections
	}
	h.clients[client.ID][client] = true

	var isFirstConn bool
	if client.User != nil && client.User.Id != 0 {
		userConns := h.userClients[client.User.Id]
		if userConns == nil {
			userConns = make(map[*Client]bool)
			h.userClients[client.User.Id] = userConns
		}
		h.userClients[client.User.Id][client] = true
		if len(h.userClients[client.User.Id]) == 1 {
			isFirstConn = true
		}
	}
	h.mu.Unlock()

	slog.Info("WebSocket client registered", "username", client.User.Username, "roomId", client.ID)

	h.HandleMessage(Message{
		Type:      "system",
		Sender:    0,
		Content:   client.User.Username + " has joined the room",
		ID:        client.ID,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	})

	if isFirstConn && client.User != nil {
		h.HandleMessage(Message{
			Type:   "user_status",
			Sender: client.User.Id,
			ID:     "notifications",
			Payload: map[string]interface{}{
				"user_id":   client.User.Id,
				"is_online": true,
			},
		})
	}
}

func (h *Hub) RemoveClient(client *Client) {
	h.mu.Lock()
	if _, ok := h.clients[client.ID]; ok {
		delete(h.clients[client.ID], client)
		if len(h.clients[client.ID]) == 0 {
			delete(h.clients, client.ID)
		}
	}

	var isLastConn bool
	if client.User != nil && client.User.Id != 0 {
		if userConns, ok := h.userClients[client.User.Id]; ok {
			delete(userConns, client)
			if len(userConns) == 0 {
				delete(h.userClients, client.User.Id)
				isLastConn = true
			}
		}
	}
	h.mu.Unlock()

	client.Close()
	slog.Info("WebSocket client unregistered", "username", client.User.Username, "roomId", client.ID)

	if isLastConn && client.User != nil {
		h.HandleMessage(Message{
			Type:   "user_status",
			Sender: client.User.Id,
			ID:     "notifications",
			Payload: map[string]interface{}{
				"user_id":   client.User.Id,
				"is_online": false,
			},
		})
	}
}

func (h *Hub) HandleMessage(message Message) {
	h.mu.RLock()
	clients, ok := h.clients[message.ID]
	if !ok {
		h.mu.RUnlock()
		return
	}

	var stale []*Client
	for client := range clients {
		select {
		case client.send <- message:
		default:
			stale = append(stale, client)
		}
	}
	h.mu.RUnlock()

	// Safely clean up stale clients under a write lock
	if len(stale) > 0 {
		h.mu.Lock()
		for _, client := range stale {
			if _, exists := h.clients[message.ID][client]; exists {
				close(client.send)
				delete(h.clients[message.ID], client)
			}
		}
		h.mu.Unlock()
	}
}

func (h *Hub) SendToUser(userID int, eventType string, payload interface{}) {
	msg := Message{
		Type:      eventType,
		Sender:    0,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Payload:   payload,
	}

	h.mu.RLock()
	clients, ok := h.userClients[userID]
	if !ok {
		h.mu.RUnlock()
		return
	}

	var stale []*Client
	for client := range clients {
		select {
		case client.send <- msg:
		default:
			stale = append(stale, client)
		}
	}
	h.mu.RUnlock()

	if len(stale) > 0 {
		h.mu.Lock()
		for _, client := range stale {
			if userConns, exists := h.userClients[userID]; exists {
				close(client.send)
				delete(userConns, client)
			}
		}
		h.mu.Unlock()
	}
}

func (h *Hub) Publish(ctx context.Context, channel string, msg Message) error {
	return h.pubSub.Publish(ctx, channel, msg)
}
