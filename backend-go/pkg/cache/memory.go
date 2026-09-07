package cache

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

var (
	ErrCacheMiss = errors.New("cache: key not found")
)

type item struct {
	Value      []byte
	Expiration int64
}

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]item
}

func NewMemoryCache() *MemoryCache {
	c := &MemoryCache{
		items: make(map[string]item),
	}
	go c.cleanupLoop()
	return c
}

func (c *MemoryCache) Get(ctx context.Context, key string, dest interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	it, ok := c.items[key]
	if !ok {
		return ErrCacheMiss
	}

	if it.Expiration > 0 && time.Now().UnixNano() > it.Expiration {
		return ErrCacheMiss
	}

	return json.Unmarshal(it.Value, dest)
}

func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	var exp int64
	if expiration > 0 {
		exp = time.Now().Add(expiration).UnixNano()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item{
		Value:      data,
		Expiration: exp,
	}
	return nil
}

func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}

func (c *MemoryCache) Flush(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]item)
	return nil
}

func (c *MemoryCache) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		c.mu.Lock()
		now := time.Now().UnixNano()
		for k, it := range c.items {
			if it.Expiration > 0 && now > it.Expiration {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}
