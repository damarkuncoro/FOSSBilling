package plugins

import (
	"context"
	"sync"
)

// HookHandler is a function that can modify data or perform side effects
type HookHandler func(ctx context.Context, payload interface{}) (interface{}, error)

// HookManager manages dynamic extension points in the system
type HookManager struct {
	mu    sync.RWMutex
	hooks map[string][]HookHandler
}

func NewHookManager() *HookManager {
	return &HookManager{
		hooks: make(map[string][]HookHandler),
	}
}

// Register adds a new handler for a specific hook point (e.g. "filter_invoice_title")
func (m *HookManager) Register(hookName string, handler HookHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hooks[hookName] = append(m.hooks[hookName], handler)
}

// Apply executes all handlers for a hook, passing the result of one to the next (pipelining)
func (m *HookManager) Apply(ctx context.Context, hookName string, data interface{}) (interface{}, error) {
	m.mu.RLock()
	handlers, exists := m.hooks[hookName]
	m.mu.RUnlock()

	if !exists {
		return data, nil
	}

	var err error
	currentData := data
	for _, h := range handlers {
		currentData, err = h(ctx, currentData)
		if err != nil {
			return currentData, err
		}
	}

	return currentData, nil
}
