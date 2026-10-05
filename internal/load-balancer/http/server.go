package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"tpiii.local/daniel-tpiii/internal/load-balancer/routing"
)

type Server struct {
	httpServer *http.Server
	router     *routing.Manager
	logger     *slog.Logger
}

func NewServer(address string, router *routing.Manager, logger *slog.Logger) *Server {
	s := &Server{
		router: router,
		logger: logger,
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

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	serviceName, path, _ := strings.Cut(strings.TrimPrefix(r.URL.Path, "/"), "/")

	instance, done, err := s.router.Pick(serviceName)
	if err != nil {
		switch {
		case errors.Is(err, routing.ErrUnknownService):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		}
		return
	}
	defer done()

	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(&url.URL{Scheme: "http", Host: instance.Address})
			pr.Out.URL.Path = "/" + path
			pr.Out.URL.RawPath = ""
			pr.SetXForwarded()
		},
		ErrorHandler: func(w http.ResponseWriter, _ *http.Request, err error) {
			s.logger.Warn(
				"Failed to proxy request",
				"service", serviceName,
				"instance_id", instance.ID,
				"error", err,
			)
			w.WriteHeader(http.StatusBadGateway)
		},
	}

	proxy.ServeHTTP(w, r)
}
