package main

import (
	"context"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/config"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/postgres"
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

	// 2. Data Access Layer (Repositories)
	repos := InitRepositories(ctx, cfg, pgPool)

	// 3. Domain Business Logic Layer (Services & Use Cases)
	services := InitServices(cfg, repos)

	// 4. HTTP Presentation Layer (Handlers & Router)
	handlers := InitHandlers(services, repos)

	rateLimiter := middleware.NewRateLimiter(60, time.Second)
	router := setupRoutes(cfg, handlers, rateLimiter)

	// 5. Server Lifecycle & Graceful Shutdown
	serverLifecycle := NewHTTPServerLifecycle(cfg, router)
	serverLifecycle.StartAndListenWithGracefulShutdown()
}
