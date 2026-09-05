package events_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fossbilling/backend-go/pkg/events"
)

func TestEventBus_SyncPublish(t *testing.T) {
	bus := events.NewEventBus()
	ctx := context.Background()

	var callCount int32
	bus.Subscribe(events.EventClientRegistered, func(ctx context.Context, e events.Event) error {
		atomic.AddInt32(&callCount, 1)
		return nil
	})

	bus.Subscribe(events.EventClientRegistered, func(ctx context.Context, e events.Event) error {
		atomic.AddInt32(&callCount, 1)
		return nil
	})

	err := bus.Publish(ctx, events.Event{
		Type:    events.EventClientRegistered,
		Payload: map[string]string{"email": "user@example.com"},
	})
	if err != nil {
		t.Fatalf("unexpected publish error: %v", err)
	}

	if atomic.LoadInt32(&callCount) != 2 {
		t.Errorf("expected callCount 2, got %d", atomic.LoadInt32(&callCount))
	}
}

func TestEventBus_AsyncPublish(t *testing.T) {
	bus := events.NewEventBus()
	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(1)

	bus.Subscribe(events.EventInvoicePaid, func(ctx context.Context, e events.Event) error {
		defer wg.Done()
		return nil
	})

	bus.PublishAsync(ctx, events.Event{
		Type:    events.EventInvoicePaid,
		Payload: map[string]int64{"invoice_id": 123},
	})

	waitChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitChan)
	}()

	select {
	case <-waitChan:
		// success
	case <-time.After(1 * time.Second):
		t.Fatal("async event handler timed out")
	}
}
