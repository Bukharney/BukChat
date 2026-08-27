package ws

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/bukharney/bukchat/middlewares"
	"github.com/bukharney/bukchat/modules/entities"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client struct for websocket connection and message sending
type Client struct {
	Chat      entities.ChatRepository
	User      *entities.UsersClaims
	ID        string
	Conn      *websocket.Conn
	send      chan Message
	hub       *Hub
	closeOnce sync.Once
}

// NewClient creates a new client
func NewClient(id string, conn *websocket.Conn, hub *Hub, user *entities.UsersClaims, chat entities.ChatRepository) *Client {
	return &Client{
		ID:   id,
		Conn: conn,
		send: make(chan Message, 256),
		hub:  hub,
		User: user,
		Chat: chat,
	}
}

// Client goroutine to read messages from client
func (c *Client) Read() {
	defer func() {
		c.hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("WS Read Unexpected Close Error", "error", err, "username", c.User.Username)
			}
			break
		}
		msg.Sender = c.User.Id
		if msg.ID == "" {
			msg.ID = c.ID
		}
		RoomId, _ := strconv.Atoi(msg.ID)

		if msg.Type == "typing" {
			c.hub.broadcast <- msg
			continue
		}

		ctx := context.Background()
		err = c.Chat.SendMessage(ctx, &entities.ChatMessage{
			RoomId:  RoomId,
			Sender:  c.User.Id,
			Message: msg.Content,
		})
		if err != nil {
			slog.Error("Failed to persist WS message", "error", err, "roomId", RoomId)
		}

		c.hub.broadcast <- msg

		c.hub.broadcast <- Message{
			Type:    "new_message",
			Sender:  c.User.Id,
			Content: msg.Content,
			ID:      "notifications",
			Payload: map[string]interface{}{
				"room_id": RoomId,
				"sender":  c.User.Id,
				"content": msg.Content,
			},
		}
	}
}

// Client goroutine to write messages to client
func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.Conn.WriteJSON(message)
			if err != nil {
				slog.Error("WS Write JSON Error", "error", err, "username", c.User.Username)
				return
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Client closing channel to unregister client
func (c *Client) Close() {
	c.closeOnce.Do(func() {
		close(c.send)
	})
}

// Function to handle websocket connection and register client to hub and start goroutines
func ServeWS(c *gin.Context, hub *Hub, chatRepo entities.ChatRepository) {
	roomId := c.Param("roomId")
	if roomId == "" {
		c.JSON(400, gin.H{"error": "roomId is required"})
		return
	}

	tk := c.Query("token")
	if tk == "" {
		c.JSON(400, gin.H{"error": "token is required"})
		return
	}

	user, err := middlewares.GetUserToken(tk)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("WebSocket Upgrade Error", "error", err)
		return
	}

	slog.Info("WebSocket connected", "username", user.Username, "roomId", roomId)

	client := NewClient(roomId, ws, hub, user, chatRepo)

	hub.register <- client

	go client.Write()
	go client.Read()
}
