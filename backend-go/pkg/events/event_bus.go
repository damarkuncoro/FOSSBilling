package events

import (
	"context"
	"sync"
)

type EventType string
const (
	EventClientRegistered EventType = "client.registered"; EventInvoiceCreated EventType = "invoice.created"; EventInvoicePaid EventType = "invoice.paid"
	EventOrderActivated EventType = "order.activated"; EventOrderSuspended EventType = "order.suspended"; EventTicketOpened EventType = "ticket.opened"
	EventTicketReplied EventType = "ticket.replied"; EventTicketClosed EventType = "ticket.closed"; EventLowStock EventType = "system.low_stock"
	EventOrderProvisioningFailed EventType = "order.provisioning_failed"
)

type Event struct { Type EventType `json:"type"`; Payload any `json:"payload"` }
type EventHandler func(context.Context, Event) error
type EventBus struct { mu sync.RWMutex; hs map[EventType][]EventHandler }

func NewEventBus() *EventBus { return &EventBus{hs: make(map[EventType][]EventHandler)} }
func (b *EventBus) Subscribe(t EventType, h EventHandler) { b.mu.Lock(); defer b.mu.Unlock(); b.hs[t] = append(b.hs[t], h) }

func (b *EventBus) Publish(ctx context.Context, e Event) error {
	b.mu.RLock(); hs, ok := b.hs[e.Type]; b.mu.RUnlock(); if !ok { return nil }
	for _, h := range hs { if err := h(ctx, e); err != nil { return err } }; return nil
}

func (b *EventBus) PublishAsync(_ context.Context, e Event) {
	b.mu.RLock(); hs, ok := b.hs[e.Type]; b.mu.RUnlock(); if !ok { return }
	for _, h := range hs { go func(fn EventHandler) { _ = fn(context.Background(), e) }(h) }
}
