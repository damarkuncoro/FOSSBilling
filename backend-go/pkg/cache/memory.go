package cache

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

var ErrMiss = errors.New("cache miss")

type itm struct { Val []byte; Exp int64 }
type MemoryCache struct { mu sync.RWMutex; its map[string]itm }

func NewMemoryCache() *MemoryCache {
	c := &MemoryCache{its: make(map[string]itm)}
	go func() {
		for range time.Tick(5 * time.Minute) {
			c.mu.Lock(); now := time.Now().UnixNano()
			for k, v := range c.its { if v.Exp > 0 && now > v.Exp { delete(c.its, k) } }
			c.mu.Unlock()
		}
	}()
	return c
}

func (c *MemoryCache) Get(_ context.Context, k string, d any) error {
	c.mu.RLock(); defer c.mu.RUnlock()
	v, ok := c.its[k]; if !ok || (v.Exp > 0 && time.Now().UnixNano() > v.Exp) { return ErrMiss }
	return json.Unmarshal(v.Val, d)
}

func (c *MemoryCache) Set(_ context.Context, k string, v any, e time.Duration) error {
	b, err := json.Marshal(v); if err != nil { return err }
	exp := int64(0); if e > 0 { exp = time.Now().Add(e).UnixNano() }
	c.mu.Lock(); defer c.mu.Unlock(); c.its[k] = itm{b, exp}; return nil
}

func (c *MemoryCache) Delete(_ context.Context, k string) error { c.mu.Lock(); defer c.mu.Unlock(); delete(c.its, k); return nil }
func (c *MemoryCache) Flush(_ context.Context) error { c.mu.Lock(); defer c.mu.Unlock(); c.its = make(map[string]itm); return nil }
