package main

import (
	"fmt"
	"net/http"

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

	myLogger.Info("Server starting on %s", cfg.HTTPServer.Addr)
	err := server.ListenAndServe()
	if err != nil {
		myLogger.Fatal("Failed to start server: %v", err)
	}
}
