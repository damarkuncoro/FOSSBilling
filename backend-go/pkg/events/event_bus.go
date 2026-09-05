package events

import (
	"context"
	"sync"
)

type EventType string

const (
	EventClientRegistered EventType = "client.registered"
	EventInvoiceCreated    EventType = "invoice.created"
	EventInvoicePaid       EventType = "invoice.paid"
	EventOrderActivated    EventType = "order.activated"
	EventOrderSuspended    EventType = "order.suspended"
	EventTicketOpened      EventType = "ticket.opened"
	EventTicketReplied     EventType = "ticket.replied"
)

type Event struct {
	Type    EventType   `json:"type"`
	Payload interface{} `json:"payload"`
}

type EventHandler func(ctx context.Context, event Event) error

type EventBus struct {
	mu       sync.RWMutex
	handlers map[EventType][]EventHandler
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]EventHandler),
	}
}

// Subscribe registers a listener for a specific event type
func (b *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish dispatches an event synchronously to all registered listeners
func (b *EventBus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	handlers, exists := b.handlers[event.Type]
	b.mu.RUnlock()

	if !exists {
		return nil
	}

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// PublishAsync dispatches an event asynchronously in background goroutines
func (b *EventBus) PublishAsync(ctx context.Context, event Event) {
	b.mu.RLock()
	handlers, exists := b.handlers[event.Type]
	b.mu.RUnlock()

	if !exists {
		return
	}

	for _, handler := range handlers {
		h := handler
		go func() {
			_ = h(context.Background(), event)
		}()
	}
}
