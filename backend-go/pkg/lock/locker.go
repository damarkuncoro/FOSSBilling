package lock

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrLockHeld = errors.New("lock already held by another process")

type Locker interface {
	Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error)
	Release(ctx context.Context, key string) error
}

type RedisLocker struct {
	client *redis.Client
}

func NewRedisLocker(client *redis.Client) *RedisLocker {
	return &RedisLocker{client: client}
}

func (l *RedisLocker) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	if l.client == nil {
		return true, nil // Fail open if no redis, but locally it might overlap
	}

	// NX: Set if Not eXists
	ok, err := l.client.SetNX(ctx, "lock:"+key, "1", ttl).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (l *RedisLocker) Release(ctx context.Context, key string) error {
	if l.client == nil {
		return nil
	}
	return l.client.Del(ctx, "lock:"+key).Err()
}

type LocalLocker struct{} // Fallback for single-node development
func (l *LocalLocker) Acquire(_ context.Context, _ string, _ time.Duration) (bool, error) { return true, nil }
func (l *LocalLocker) Release(_ context.Context, _ string) error { return nil }
