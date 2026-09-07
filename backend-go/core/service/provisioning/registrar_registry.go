package provisioning

import (
	"errors"
	"sync"
)

var (
	ErrRegistrarNotFound = errors.New("registrar driver not found")
)

// RegistrarRegistry manages multiple domain registrar drivers
type RegistrarRegistry struct {
	mu       sync.RWMutex
	drivers  map[string]RegistrarDriver
}

func NewRegistrarRegistry() *RegistrarRegistry {
	return &RegistrarRegistry{
		drivers: make(map[string]RegistrarDriver),
	}
}

func (r *RegistrarRegistry) Register(id string, driver RegistrarDriver) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.drivers[id] = driver
}

func (r *RegistrarRegistry) Get(id string) (RegistrarDriver, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	driver, ok := r.drivers[id]
	if !ok {
		return nil, ErrRegistrarNotFound
	}
	return driver, nil
}

func (r *RegistrarRegistry) List() map[string]RegistrarDriver {
	r.mu.RLock()
	defer r.mu.RUnlock()
	copy := make(map[string]RegistrarDriver)
	for k, v := range r.drivers {
		copy[k] = v
	}
	return copy
}
