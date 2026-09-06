package notification

import (
	"context"
	"errors"
	"sync"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
)

var (
	ErrDriverNotFound = errors.New("notification mail driver not registered")
)

// MailDriver defines the provider contract for email delivery services
type MailDriver interface {
	ID() string
	Name() string
	Send(ctx context.Context, msg mailer.Message) error
}

// DriverRegistry manages registered email and notification transport providers
type DriverRegistry struct {
	mu      sync.RWMutex
	drivers map[string]MailDriver
}

// NewDriverRegistry creates a new empty driver registry
func NewDriverRegistry() *DriverRegistry {
	return &DriverRegistry{
		drivers: make(map[string]MailDriver),
	}
}

// Register adds a new transport driver to the registry
func (r *DriverRegistry) Register(driver MailDriver) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.drivers[driver.ID()] = driver
}

// Get retrieves a driver by its unique ID
func (r *DriverRegistry) Get(id string) (MailDriver, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	drv, exists := r.drivers[id]
	if !exists {
		return nil, ErrDriverNotFound
	}
	return drv, nil
}

// List returns all registered notification drivers
func (r *DriverRegistry) List() []MailDriver {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]MailDriver, 0, len(r.drivers))
	for _, drv := range r.drivers {
		list = append(list, drv)
	}
	return list
}
