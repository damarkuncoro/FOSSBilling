package provisioning

import (
	"errors"
	"sync"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

var (
	ErrProvisionerNotFound = errors.New("service provisioner not registered for this product type")
)

// ProvisionerRegistry manages registered hosting and service provisioner drivers
type ProvisionerRegistry struct {
	mu          sync.RWMutex
	provisioners map[domain.ProductType]domain.ServiceProvisioner
}

// NewProvisionerRegistry initializes a thread-safe provisioner driver registry
func NewProvisionerRegistry() *ProvisionerRegistry {
	return &ProvisionerRegistry{
		provisioners: make(map[domain.ProductType]domain.ServiceProvisioner),
	}
}

// Register binds a provisioner driver to a specific ProductType
func (r *ProvisionerRegistry) Register(p domain.ServiceProvisioner) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.provisioners[p.Type()] = p
}

// Get retrieves a provisioner driver for a given ProductType
func (r *ProvisionerRegistry) Get(productType domain.ProductType) (domain.ServiceProvisioner, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	prov, exists := r.provisioners[productType]
	if !exists {
		return nil, ErrProvisionerNotFound
	}
	return prov, nil
}

// List returns all registered service provisioners
func (r *ProvisionerRegistry) List() []domain.ServiceProvisioner {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]domain.ServiceProvisioner, 0, len(r.provisioners))
	for _, p := range r.provisioners {
		list = append(list, p)
	}
	return list
}
