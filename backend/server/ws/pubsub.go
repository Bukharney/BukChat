package ws

import (
	"context"
	"sync"
)

// PubSubAdapter abstracts message broadcasting for single-node (in-memory) or multi-node (Redis / NATS) setups.
type PubSubAdapter interface {
	Publish(ctx context.Context, channel string, msg Message) error
	Subscribe(ctx context.Context, channel string, handler func(msg Message)) error
	Unsubscribe(ctx context.Context, channel string) error
}

// InMemoryPubSub is the default single-instance PubSub implementation.
type InMemoryPubSub struct {
	subscribers map[string][]func(msg Message)
	mu          sync.RWMutex
}

func NewInMemoryPubSub() *InMemoryPubSub {
	return &InMemoryPubSub{
		subscribers: make(map[string][]func(msg Message)),
	}
}

func (p *InMemoryPubSub) Publish(ctx context.Context, channel string, msg Message) error {
	p.mu.RLock()
	handlers, ok := p.subscribers[channel]
	p.mu.RUnlock()

	if !ok {
		return nil
	}

	for _, handler := range handlers {
		handler(msg)
	}
	return nil
}

func (p *InMemoryPubSub) Subscribe(ctx context.Context, channel string, handler func(msg Message)) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.subscribers[channel] = append(p.subscribers[channel], handler)
	return nil
}

func (p *InMemoryPubSub) Unsubscribe(ctx context.Context, channel string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.subscribers, channel)
	return nil
}
