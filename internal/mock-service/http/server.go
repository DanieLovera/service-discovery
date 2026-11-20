package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
	service    string
	instanceID string
}

type response struct {
	Service    string `json:"service"`
	InstanceID string `json:"instance_id"`
}

func NewServer(address, service, instanceID string, logger *slog.Logger) *Server {
	s := &Server{
		logger:     logger,
		service:    service,
		instanceID: instanceID,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handle)

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
		s.logger.Info("HTTP server started", "address", s.httpServer.Addr)

		if err := s.httpServer.Serve(listener); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			errs <- fmt.Errorf("serve: %w", err)
		}
	}()

	return errs, nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	defer s.logger.Info("HTTP server stopped")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}

	return nil
}

func (s *Server) handle(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response{
		Service:    s.service,
		InstanceID: s.instanceID,
	}); err != nil {
		s.logger.Error("Failed to encode response", "error", err)
	}
}
