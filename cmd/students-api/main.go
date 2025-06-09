package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/rohit154041/students-api/internal/config"
	"github.com/rohit154041/students-api/internal/router"
	"github.com/rohit154041/students-api/internal/storage"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// Initialize MongoDB connection (just call once)
	if err := storage.Init(cfg.MongoURI); err != nil {
		slog.Error("failed to connect to MongoDB", slog.String("error", err.Error()))
		return
	}
	defer storage.Disconnect() // Clean disconnect on shutdown

	// setup router
	mux := router.Init()

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}

	slog.Info("server started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}