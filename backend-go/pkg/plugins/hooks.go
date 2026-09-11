package plugins

import (
	"context"
	"sync"
)

type HookHandler func(context.Context, any) (any, error)
type HookManager struct {
	mu    sync.RWMutex
	hs    map[string][]HookHandler
	lua   *LuaEngine
}

func NewHookManager(lua ...*LuaEngine) *HookManager {
	var engine *LuaEngine
	if len(lua) > 0 { engine = lua[0] }
	return &HookManager{
		hs:  make(map[string][]HookHandler),
		lua: engine,
	}
}

func (m *HookManager) Register(n string, h HookHandler) { m.mu.Lock(); defer m.mu.Unlock(); m.hs[n] = append(m.hs[n], h) }

func (m *HookManager) Apply(ctx context.Context, n string, d any) (any, error) {
	// 1. Run Internal Go Hooks
	m.mu.RLock()
	hs, ok := m.hs[n]
	m.mu.RUnlock()

	cur := d
	var err error

	if ok {
		for _, h := range hs {
			cur, err = h(ctx, cur)
			if err != nil { return cur, err }
		}
	}

	// 2. Run Dynamic Lua Hooks
	if m.lua != nil {
		cur, err = m.lua.ExecuteHook(ctx, n, cur)
	}

	return cur, nil
}
