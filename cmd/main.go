package main

import (
	"EcomTechGo/internal/store"
	"EcomTechGo/internal/transport/rest"
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var persistFile string

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	flag.StringVar(&persistFile, "file", "todos.json", "Path to persistence file")
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	storage := store.New()
	if err := storage.Load(persistFile); err != nil {
		slog.Error("failed to load store", "file", persistFile, "error", err)
	}

	handler := rest.NewHandler(storage)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: rest.LoggingMiddleware(handler),
	}

	go func() {
		slog.Info("Server is starting", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exiting")
}
