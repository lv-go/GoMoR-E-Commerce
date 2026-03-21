package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"gomor-e-commerce/internal/app"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := app.Setup(ctx)

	// Start server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	slog.Info("Server started on :" + port)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}
