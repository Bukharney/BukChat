package ws

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bukharney/giga-chat/middlewares"
	"github.com/bukharney/giga-chat/modules/entities"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	User *entities.UsersClaims
	ID   string
	Conn *websocket.Conn
	send chan Message
	hub  *Hub
}

// NewClient creates a new client
func NewClient(id string, conn *websocket.Conn, hub *Hub, user *entities.UsersClaims) *Client {
	return &Client{
		ID:   id,
		Conn: conn,
		send: make(chan Message, 256),
		hub:  hub,
		User: user,
	}
}

// Client goroutine to read messages from client
func (c *Client) Read() {
	defer func() {
		c.hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg Message
		msg.Sender = c.User.Username
		msg.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error: ", err)
			break
		}
		c.hub.broadcast <- msg
	}
}

// Client goroutine to write messages to client
func (c *Client) Write() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.send {
		_ = c.Conn.WriteJSON(message)
	}
}

// Client closing channel to unregister client
func (c *Client) Close() {
	close(c.send)
}

// Function to handle websocket connection and register client to hub and start goroutines
func ServeWS(c *gin.Context, hub *Hub) {
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
		fmt.Println(err)
		return
	}

	fmt.Println("Client connected")

	client := NewClient(roomId, ws, hub, user)

	hub.register <- client

	go client.Write()
	go client.Read()
}
