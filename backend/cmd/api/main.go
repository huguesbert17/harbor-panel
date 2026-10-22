package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/harbor-panel/harbor-panel/internal/controlplane"
	"github.com/harbor-panel/harbor-panel/internal/health"
	"github.com/harbor-panel/harbor-panel/internal/hostingaccounts"
	"github.com/harbor-panel/harbor-panel/internal/platform/httpserver"
	"github.com/harbor-panel/harbor-panel/internal/servers"
)

const version = "0.1.0-foundation"

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	addr := os.Getenv("HARBOR_API_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		httpserver.JSON(w, http.StatusOK, health.Status{Status: health.StatusOK})
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, _ *http.Request) {
		httpserver.JSON(w, http.StatusOK, health.Status{Status: health.StatusReady})
	})
	mux.HandleFunc("GET /api/v1", func(w http.ResponseWriter, _ *http.Request) {
		httpserver.JSON(w, http.StatusOK, map[string]string{"name": "harbor-panel", "version": version})
	})

	accounts := hostingaccounts.NewMemoryRepository()
	serverStore := servers.NewMemoryRepository()
	apiHandler := controlplane.NewHandler(accounts, serverStore)
	mux.Handle("/api/v1/accounts", apiHandler)
	mux.Handle("/api/v1/servers", apiHandler)

	server := httpserver.New(addr, mux, logger)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("api listening", "addr", addr, "version", version)
		if err := server.HTTP.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("api stopped unexpectedly", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("api shutdown failed", "error", err)
		os.Exit(1)
	}
}
