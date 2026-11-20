package observability

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
	ready      atomic.Bool
}

func NewServer(address string, metricsHandler http.Handler, logger *slog.Logger) *Server {
	s := &Server{
		logger: logger,
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", metricsHandler)
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/readyz", s.handleReady)

	s.httpServer = &http.Server{
		Addr:    address,
		Handler: mux,
	}

	return s
}

func (s *Server) Serve() (<-chan error, error) {
	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return nil, fmt.Errorf("listen on %s: %w", s.httpServer.Addr, err)
	}

	errs := make(chan error, 1)

	go func() {
		s.logger.Info("Observability server started", "address", s.httpServer.Addr)

		if err := s.httpServer.Serve(listener); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			errs <- fmt.Errorf("serve: %w", err)
		}
	}()

	return errs, nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	defer s.logger.Info("Observability server stopped")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown observability HTTP: %w", err)
	}

	return nil
}

func (s *Server) SetReady(ready bool) {
	s.ready.Store(ready)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleReady(w http.ResponseWriter, _ *http.Request) {
	if !s.ready.Load() {
		http.Error(w, "Not ready", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}
