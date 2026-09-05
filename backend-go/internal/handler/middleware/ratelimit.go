package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type clientVisitor struct {
	tokens     int
	lastSeen   time.Time
}

type RateLimiter struct {
	mu         sync.Mutex
	visitors   map[string]*clientVisitor
	limit      int           // max tokens
	refillRate time.Duration // time per token
}

func NewRateLimiter(limit int, refillRate time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors:   make(map[string]*clientVisitor),
		limit:      limit,
		refillRate: refillRate,
	}

	// Periodic cleanup of stale visitors every 5 minutes
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			rl.mu.Lock()
			for ip, v := range rl.visitors {
				if time.Since(v.lastSeen) > 10*time.Minute {
					delete(rl.visitors, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()

	return rl
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &clientVisitor{
			tokens:   rl.limit - 1,
			lastSeen: now,
		}
		return true
	}

	// Refill tokens based on elapsed time
	elapsed := now.Sub(v.lastSeen)
	tokensToAdd := int(elapsed / rl.refillRate)
	if tokensToAdd > 0 {
		v.tokens += tokensToAdd
		if v.tokens > rl.limit {
			v.tokens = rl.limit
		}
		v.lastSeen = now
	}

	if v.tokens > 0 {
		v.tokens--
		return true
	}

	return false
}

// RateLimit middleware limits incoming HTTP requests per remote IP
func (rl *RateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			parts := strings.Split(forwarded, ",")
			ip = strings.TrimSpace(parts[0])
		}

		if !rl.Allow(ip) {
			w.Header().Set("Retry-After", "5")
			response.Error(w, http.StatusTooManyRequests, "TOO_MANY_REQUESTS", "Rate limit exceeded. Please try again in a few seconds.", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}
