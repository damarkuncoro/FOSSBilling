package plugins

import (
	"context"
	"sync"
)

type HookHandler func(context.Context, any) (any, error)
type HookManager struct { mu sync.RWMutex; hs map[string][]HookHandler }

func NewHookManager() *HookManager { return &HookManager{hs: make(map[string][]HookHandler)} }

func (m *HookManager) Register(n string, h HookHandler) { m.mu.Lock(); defer m.mu.Unlock(); m.hs[n] = append(m.hs[n], h) }

func (m *HookManager) Apply(ctx context.Context, n string, d any) (any, error) {
	m.mu.RLock(); hs, ok := m.hs[n]; m.mu.RUnlock(); if !ok { return d, nil }
	var err error; cur := d
	for _, h := range hs { cur, err = h(ctx, cur); if err != nil { return cur, err } }
	return cur, nil
}
