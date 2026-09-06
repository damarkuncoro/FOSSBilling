package drivers

import (
	"context"
	"sync"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
)

type MockDriver struct {
	mu           sync.RWMutex
	SentMessages []mailer.Message
}

func NewMockDriver() *MockDriver {
	return &MockDriver{
		SentMessages: make([]mailer.Message, 0),
	}
}

func (d *MockDriver) ID() string   { return "mock" }
func (d *MockDriver) Name() string { return "In-Memory Test Mock Driver" }

func (d *MockDriver) Send(ctx context.Context, msg mailer.Message) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.SentMessages = append(d.SentMessages, msg)
	return nil
}

func (d *MockDriver) Count() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.SentMessages)
}
