// Package httpserver contains framework-free HTTP primitives shared by the
// control-plane entry points.
package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	HTTP *http.Server
}

func New(addr string, handler http.Handler, logger *slog.Logger) *Server {
	return &Server{HTTP: &http.Server{
		Addr:              addr,
		Handler:           requestLogger(handler, logger),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}}
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.HTTP.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func JSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func requestLogger(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		logger.Info("http request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(started))
	})
}
