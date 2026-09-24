package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"tpiii.local/daniel-tpiii/internal/metrics"
	"tpiii.local/daniel-tpiii/internal/observability"
	"tpiii.local/daniel-tpiii/internal/registry"
	"tpiii.local/daniel-tpiii/internal/registry/grpc"
)

type Params struct {
	NodeID               string
	Backend              string
	GRPCAddress          string
	ObservabilityAddress string
	ClusterMembers       []string
	Logger               *slog.Logger
}

type App struct {
	grpcServer          *grpc.Server
	observabilityServer *observability.Server
	logger              *slog.Logger
}

func New(params Params) (*App, error) {
	backend, err := newBackend(params)
	if err != nil {
		return nil, err
	}

	registryManager := registry.NewRegistryManager(backend)
	grpcHandler := grpc.NewHandler(registryManager)
	grpcServer := grpc.NewServer(params.GRPCAddress, grpcHandler, params.Logger)

	processMetrics := metrics.New("registry", params.NodeID)
	metricsHandler := processMetrics.Handler()
	observabilityServer := observability.NewServer(
		params.ObservabilityAddress,
		metricsHandler,
		params.Logger,
	)

	return &App{
		logger:              params.Logger,
		grpcServer:          grpcServer,
		observabilityServer: observabilityServer,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	if err := a.observabilityServer.Listen(); err != nil {
		return fmt.Errorf("listen observability server: %w", err)
	}
	if err := a.grpcServer.Listen(); err != nil {
		return fmt.Errorf("listen registry gRPC server: %w", err)
	}

	grpcErrs := make(chan error, 1)

	go func() {
		if err := a.observabilityServer.Serve(); err != nil {
			a.logger.Error("Observability server failed", "error", err)
		}
	}()

	go func() {
		if err := a.grpcServer.Serve(); err != nil {
			grpcErrs <- fmt.Errorf("serve registry gRPC server: %w", err)
		}
	}()

	a.observabilityServer.SetReady(true)
	a.logger.Info("Registry started")

	select {
	case <-ctx.Done():
		return nil
	case err := <-grpcErrs:
		return err
	}
}

func (a *App) Stop(ctx context.Context) error {
	defer a.logger.Info("Registry stopped")

	a.observabilityServer.SetReady(false)

	errs := make([]error, 2)

	var wg sync.WaitGroup
	wg.Go(func() {
		if err := a.observabilityServer.Shutdown(ctx); err != nil {
			errs[1] = fmt.Errorf("shutdown observability server: %w", err)
		}
	})

	wg.Go(func() {
		if err := a.grpcServer.Shutdown(ctx); err != nil {
			errs[0] = fmt.Errorf("shutdown registry gRPC server: %w", err)
		}
	})

	wg.Wait()

	return errors.Join(errs...)
}

func newBackend(params Params) (registry.Backend, error) {
	switch params.Backend {
	case "cp":
		// TODO: Implement CP backend
		return nil, nil
	case "ap":
		// TODO: Implement AP backend
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported registry backend %q", params.Backend)
	}
}
