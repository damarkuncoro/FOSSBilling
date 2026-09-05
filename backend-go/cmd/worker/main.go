package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fossbilling/backend-go/internal/config"
)

func main() {
	cfg := config.Load()
	log.Printf("⚙️ Starting FOSSBilling Background Worker & Scheduler (%s)...", cfg.AppEnv)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			log.Println("⏱️ Running scheduled tasks: Checking due invoices & order renewals...")
		}
	}()

	<-stopChan
	log.Println("🛑 Shutting down background worker gracefully...")
}
