package main

import (
	"context"
	"os"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/config"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/postgres"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/cache"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/logger"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/plugins"
)

func main() {
	cfg := config.Load(); logger.Init(cfg.AppEnv); logger.Info("API Start", "port", cfg.Port)
	os.Setenv("BOOT_TIME", time.Now().UTC().Format(time.RFC3339))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second); defer cancel()

	pool, _ := postgres.NewPostgresPool(ctx, cfg.DatabaseURL)
	if pool != nil {
		defer pool.Close()
		_ = postgres.RunMigrations(ctx, pool, "migrations")
	}

	var appCache cache.Cache
	if cfg.RedisURL != "" {
		rc, err := cache.NewRedisCache(cfg.RedisURL)
		if err == nil {
			appCache = rc
			logger.Info("Connected to Redis cache", "url", cfg.RedisURL)
		}
	}
	if appCache == nil {
		appCache = cache.NewMemoryCache()
		logger.Warn("Redis not available, using in-memory cache")
	}

	eb, hm := events.NewEventBus(), plugins.NewHookManager()
	rs := InitRepositories(ctx, cfg, pool)
	ss := InitServices(cfg, rs, pool, eb, appCache, hm)
	hs := InitHandlers(ss, rs)

	arl := middleware.NewRateLimiter(5, time.Minute/5)
	rl := middleware.NewRateLimiter(60, time.Minute/60)

	NewHTTPServerLifecycle(cfg, setupRoutes(cfg, hs, rl, arl)).StartAndListenWithGracefulShutdown()
}
