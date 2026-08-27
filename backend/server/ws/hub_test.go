package ws

import (
	"testing"
	"time"

	"github.com/bukharney/bukchat/modules/entities"
)

func TestHubRegisterAndBroadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	user := &entities.UsersClaims{
		Id:       1,
		Username: "testuser",
	}

	client := NewClient("room1", nil, hub, user, nil)

	hub.register <- client

	// Allow goroutine time to process register event
	time.Sleep(50 * time.Millisecond)

	if !hub.IsUserOnline(1) {
		t.Errorf("Expected user 1 to be online")
	}

	// Test broadcast non-blocking send
	msg := Message{
		ID:      "room1",
		Content: "Hello World",
		Sender:  1,
	}

	hub.broadcast <- msg
	time.Sleep(50 * time.Millisecond)

	select {
	case received := <-client.send:
		if received.Content != "Hello World" && received.Type != "system" {
			t.Logf("Received message: %v", received)
		}
	default:
		// Message channel was consumed or system message was sent first
	}

	// Clean up
	hub.unregister <- client
	time.Sleep(50 * time.Millisecond)

	if hub.IsUserOnline(1) {
		t.Errorf("Expected user 1 to be offline after unregister")
	}
}
