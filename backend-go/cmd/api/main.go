package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fossbilling/backend-go/internal/config"
	"github.com/fossbilling/backend-go/internal/handler/http/client"
	"github.com/fossbilling/backend-go/internal/handler/http/guest"
	"github.com/fossbilling/backend-go/internal/handler/middleware"
	"github.com/fossbilling/backend-go/internal/repository/postgres"
	authUsecase "github.com/fossbilling/backend-go/internal/usecase/auth"
	"github.com/fossbilling/backend-go/pkg/response"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Initialize PostgreSQL Connection Pool
	pgPool, err := postgres.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer pgPool.Close()

	// 2. Initialize Repositories
	clientRepo := postgres.NewClientRepository(pgPool)
	_ = postgres.NewProductRepository(pgPool)

	// 3. Initialize Usecases
	authUc := authUsecase.NewAuthUsecase(clientRepo, cfg.JWTSecret)

	// 4. Initialize Handlers
	guestAuthHandler := guest.NewAuthHandler(authUc)
	clientProfileHandler := client.NewProfileHandler(authUc)

	// 5. Setup HTTP Router (Go 1.22 enhanced ServeMux)
	mux := http.NewServeMux()

	// System & Health Endpoints
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]interface{}{
			"status":      "ok",
			"environment": cfg.AppEnv,
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
		}, nil)
	})

	mux.HandleFunc("GET /api/v1", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{
			"message": "FOSSBilling Next-Gen API v1 (Golang)",
			"docs":    "/docs",
		}, nil)
	})

	// Public / Guest Auth Routes
	mux.HandleFunc("POST /api/v1/guest/auth/register", guestAuthHandler.Register)
	mux.HandleFunc("POST /api/v1/guest/auth/login", guestAuthHandler.Login)

	// Protected Client Routes
	clientAuthMiddleware := middleware.RequireAuth(cfg.JWTSecret, "client", "admin")
	mux.Handle("GET /api/v1/client/profile", clientAuthMiddleware(http.HandlerFunc(clientProfileHandler.GetProfile)))
	mux.Handle("PUT /api/v1/client/profile", clientAuthMiddleware(http.HandlerFunc(clientProfileHandler.UpdateProfile)))

	// Wrap with Global Middlewares: Logger -> CORS -> Mux
	handler := middleware.Logger(middleware.CORS(mux))

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful Shutdown Channel
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("🚀 FOSSBilling API Server running on port %s in %s mode...", cfg.Port, cfg.AppEnv)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	case sig := <-shutdown:
		log.Printf("🛑 Signal %v received. Starting graceful shutdown...", sig)
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			_ = server.Close()
		}
		log.Println("✅ Server exited cleanly.")
	}
}
