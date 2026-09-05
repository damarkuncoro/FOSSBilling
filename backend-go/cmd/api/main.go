package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fossbilling/backend-go/internal/config"
	"github.com/fossbilling/backend-go/pkg/response"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]interface{}{
			"status":      "ok",
			"environment": cfg.AppEnv,
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
		}, nil)
	})

	// API v1 root
	mux.HandleFunc("GET /api/v1", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{
			"message": "FOSSBilling Next-Gen API v1",
			"docs":    "/docs",
		}, nil)
	})

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 FOSSBilling API Server starting on port %s (%s)...", cfg.Port, cfg.AppEnv)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
