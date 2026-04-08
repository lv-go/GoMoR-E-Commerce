package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"gomor-e-commerce/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}
	logLevelVar := slog.LevelVar{}
	if err := logLevelVar.UnmarshalText([]byte(logLevel)); err != nil {
		slog.Error("Error unmarshalling log level", "error", err)
	}
	slog.SetLogLoggerLevel(logLevelVar.Level())

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
