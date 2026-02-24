package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"product-listing/config"
	"product-listing/internal/delivery/router"
	"product-listing/pkg/logger"
	"syscall"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("api")

func main() {
	if err := run(); err != nil {
		log.Errorf("Application stopped with error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	// Configure logger
	logger.ConfigureLogger()

	// Load config
	cfg := config.Load()

	// Connect to database
	db, err := config.NewDatabase(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Initialize schema
	if err := db.InitSchema("sql/schema/schema.sql"); err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}

	// Setup router
	r := router.SetupRouter(db)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	// Initializing the server in a goroutine
	go func() {
		log.Infof("Server starting on %s", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// 5-second timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Info("Server exiting")
	return nil
}
