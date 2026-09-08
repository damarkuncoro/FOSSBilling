package provisioning

import (
	"errors"
	"sync"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type DNSProviderRegistry struct {
	mu        sync.RWMutex
	providers map[string]domain.DNSProvider
}

func NewDNSProviderRegistry() *DNSProviderRegistry {
	return &DNSProviderRegistry{
		providers: make(map[string]domain.DNSProvider),
	}
}

func (r *DNSProviderRegistry) Register(id string, p domain.DNSProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[id] = p
}

func (r *DNSProviderRegistry) Get(id string) (domain.DNSProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[id]
	if !ok {
		return nil, errors.New("dns provider not found: " + id)
	}
	return p, nil
}

func (r *DNSProviderRegistry) List() map[string]domain.DNSProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Return a copy
	res := make(map[string]domain.DNSProvider)
	for k, v := range r.providers {
		res[k] = v
	}
	return res
}
