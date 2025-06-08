package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rohit154041/students-api/internal/config"
	"github.com/rohit154041/students-api/logger"
)

func main() {
	// Load config
	cfg := config.MustLoad()

	// Custom logger setup
	myLogger, loggerErr := logger.NewLogger(logger.INFO, cfg.LogPath)
	if loggerErr != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", loggerErr))
	}
	defer myLogger.Close()

	myLogger.Info("Starting the students API server")

	// Setup router
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			myLogger.Warn("Method %s not allowed on /", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		myLogger.Info("Served welcome message to client")
		w.Write([]byte("welcome to students api"))
	})

	// Setup server
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	// Graceful shutdown handling
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		myLogger.Info("Server listening on %s", cfg.HTTPServer.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			myLogger.Fatal("ListenAndServe error: %v", err)
		}
		myLogger.Info("Server stopped listening")
	}()

	// Wait for shutdown signal
	<-done
	myLogger.Info("shutting down the server")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		myLogger.Error("Error during shutdown: %v", err)
	} else {
		myLogger.Info("Server shut down gracefully")
	}

}
