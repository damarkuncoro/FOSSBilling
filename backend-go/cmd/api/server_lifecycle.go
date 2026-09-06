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

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/config"
)

// HTTPServerLifecycle manages startup, listener binding, and graceful drain lifecycle
type HTTPServerLifecycle struct {
	server *http.Server
	cfg    *config.Config
}

// NewHTTPServerLifecycle instantiates a server lifecycle manager with standard timeouts
func NewHTTPServerLifecycle(cfg *config.Config, handler http.Handler) *HTTPServerLifecycle {
	return &HTTPServerLifecycle{
		cfg: cfg,
		server: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// StartAndListenWithGracefulShutdown launches the server in a goroutine and blocks on OS interrupt signals
func (l *HTTPServerLifecycle) StartAndListenWithGracefulShutdown() {
	go func() {
		log.Printf("🚀 FOSSBilling API Server running on port %s in %s mode...", l.cfg.Port, l.cfg.AppEnv)
		if err := l.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server listener error: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	log.Println("🛑 Shutting down HTTP server gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := l.server.Shutdown(shutdownCtx); err != nil {
		log.Printf("❌ Graceful shutdown error: %v", err)
	} else {
		log.Println("✅ Server exited cleanly.")
	}
}
