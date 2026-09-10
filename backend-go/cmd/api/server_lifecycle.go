package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/config"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/logger"
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
		logger.Info("FOSSBilling API Server is live", "port", l.cfg.Port, "env", l.cfg.AppEnv)
		if err := l.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server listener failed", err)
			os.Exit(1)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	logger.Warn("Shutdown signal received, draining connections...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := l.server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Graceful shutdown failed", err)
	} else {
		logger.Info("Server exited cleanly")
	}
}
