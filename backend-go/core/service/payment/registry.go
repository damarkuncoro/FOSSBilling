package payment

import (
	"errors"
	"sync"
)

var (
	ErrGatewayNotFound = errors.New("payment gateway not registered")
)

// GatewayRegistry manages registered payment provider adapters
type GatewayRegistry struct {
	mu       sync.RWMutex
	gateways map[string]PaymentGateway
}

// NewGatewayRegistry instantiates an empty gateway registry
func NewGatewayRegistry() *GatewayRegistry {
	return &GatewayRegistry{
		gateways: make(map[string]PaymentGateway),
	}
}

// Register adds a new payment gateway driver to the registry
func (r *GatewayRegistry) Register(gateway PaymentGateway) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.gateways[gateway.ID()] = gateway
}

// Get retrieves a payment gateway by ID
func (r *GatewayRegistry) Get(id string) (PaymentGateway, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	gw, exists := r.gateways[id]
	if !exists {
		return nil, ErrGatewayNotFound
	}
	return gw, nil
}

// List returns all registered payment gateways
func (r *GatewayRegistry) List() []PaymentGateway {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]PaymentGateway, 0, len(r.gateways))
	for _, gw := range r.gateways {
		list = append(list, gw)
	}
	return list
}
