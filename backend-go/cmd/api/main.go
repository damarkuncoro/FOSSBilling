package main

import (
	"context"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/config"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/postgres"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/cache"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/plugins"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Database Connection Pool
	pgPool, err := postgres.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err == nil {
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

	rateLimiter := middleware.NewRateLimiter(60, time.Second)
	router := setupRoutes(cfg, handlers, rateLimiter)

	// 5. Server Lifecycle & Graceful Shutdown
	serverLifecycle := NewHTTPServerLifecycle(cfg, router)
	serverLifecycle.StartAndListenWithGracefulShutdown()
}
