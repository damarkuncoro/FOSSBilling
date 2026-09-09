package main

import (
	"context"
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
	cfg := config.Load()
	logger.Init(cfg.AppEnv)

	logger.Info("Starting FOSSBilling API", "env", cfg.AppEnv, "port", cfg.Port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Database Connection Pool
	pgPool, err := postgres.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to PostgreSQL", err)
	} else {
		logger.Info("Connected to PostgreSQL successfully")
		defer pgPool.Close()
	}

	// 2. Event Bus
	eventBus := events.NewEventBus()

	// 3. Cache System
	appCache := cache.NewMemoryCache()

	// 4. Hook/Plugin System
	hookManager := plugins.NewHookManager()

	// 5. Data Access Layer (Repositories)
	repos := InitRepositories(ctx, cfg, pgPool)

	// 6. Domain Business Logic Layer (Services & Use Cases)
	services := InitServices(cfg, repos, eventBus, appCache, hookManager)

	// 7. HTTP Presentation Layer (Handlers & Router)
	handlers := InitHandlers(services, repos)

	// API Rate Limiting: 60 requests per minute
	apiRateLimiter := middleware.NewRateLimiter(60, time.Minute/60)
	// Auth Rate Limiting: 5 attempts per minute (Brute force protection)
	authRateLimiter := middleware.NewRateLimiter(5, time.Minute/5)

	router := setupRoutes(cfg, handlers, apiRateLimiter, authRateLimiter)

	// 5. Server Lifecycle & Graceful Shutdown
	serverLifecycle := NewHTTPServerLifecycle(cfg, router)
	serverLifecycle.StartAndListenWithGracefulShutdown()
}
