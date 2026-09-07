package provisioning

import (
	"errors"
	"sync"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

var (
	ErrProvisionerNotFound = errors.New("service provisioner driver not found")
)

// ProvisionerRegistry manages registered hosting and service provisioner drivers
type ProvisionerRegistry struct {
	mu      sync.RWMutex
	drivers map[string]domain.ServiceProvisioner
}

// NewProvisionerRegistry initializes a thread-safe provisioner driver registry
func NewProvisionerRegistry() *ProvisionerRegistry {
	return &ProvisionerRegistry{
		drivers: make(map[string]domain.ServiceProvisioner),
	}
}

// Register binds a provisioner driver to a specific ID (e.g. "cpanel", "plesk")
func (r *ProvisionerRegistry) Register(id string, p domain.ServiceProvisioner) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.drivers[id] = p
}

// Get retrieves a provisioner driver by ID
func (r *ProvisionerRegistry) Get(id string) (domain.ServiceProvisioner, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	prov, exists := r.drivers[id]
	if !exists {
		return nil, ErrProvisionerNotFound
	}
	return prov, nil
}

// List returns all registered service provisioners
func (r *ProvisionerRegistry) List() map[string]domain.ServiceProvisioner {
	r.mu.RLock()
	defer r.mu.RUnlock()
	copy := make(map[string]domain.ServiceProvisioner)
	for k, v := range r.drivers {
		copy[k] = v
	}
	return copy
}
