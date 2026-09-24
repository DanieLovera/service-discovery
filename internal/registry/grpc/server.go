package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"

	googlegrpc "google.golang.org/grpc"

	registrypb "tpiii.local/daniel-tpiii/gen/registry"
)

type Server struct {
	grpcServer *googlegrpc.Server
	address    string
	listener   net.Listener
	logger     *slog.Logger
}

func NewServer(address string, handler registrypb.RegistryServiceServer, logger *slog.Logger) *Server {
	grpcServer := googlegrpc.NewServer()
	registrypb.RegisterRegistryServiceServer(grpcServer, handler)

	return &Server{
		grpcServer: grpcServer,
		address:    address,
		logger:     logger,
	}
}

func (s *Server) Listen() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", s.address, err)
	}
	s.listener = listener
	return nil
}

func (s *Server) Serve() error {
	s.logger.Info("Registry gRPC server started", "address", s.address)

	if err := s.grpcServer.Serve(s.listener); err != nil && !errors.Is(err, googlegrpc.ErrServerStopped) {
		return fmt.Errorf("serve: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	defer s.logger.Info("Registry gRPC server stopped")

	done := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		s.grpcServer.Stop()
		return ctx.Err()
	}
}
