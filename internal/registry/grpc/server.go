package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	registrypb "tpiii.local/daniel-tpiii/gen/registry"
)

const keepaliveMinTime = 5 * time.Second

type Server struct {
	grpcServer *googlegrpc.Server
	address    string
	logger     *slog.Logger
}

func NewServer(address string, handler registrypb.RegistryServer, logger *slog.Logger) *Server {
	grpcServer := googlegrpc.NewServer(
		googlegrpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             keepaliveMinTime,
			PermitWithoutStream: true,
		}),
	)
	registrypb.RegisterRegistryServer(grpcServer, handler)

	return &Server{
		grpcServer: grpcServer,
		address:    address,
		logger:     logger,
	}
}

func (s *Server) Serve() (<-chan error, error) {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return nil, fmt.Errorf("listen on %s: %w", s.address, err)
	}

	errs := make(chan error, 1)

	go func() {
		s.logger.Info("Registry gRPC server started", "address", s.address)

		if err := s.grpcServer.Serve(listener); err != nil && !errors.Is(err, googlegrpc.ErrServerStopped) {
			errs <- fmt.Errorf("serve: %w", err)
		}
	}()

	return errs, nil
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
