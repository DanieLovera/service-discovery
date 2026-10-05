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
	TTLMultiplier        int
	Logger               *slog.Logger
}

type App struct {
	healthManager       *registry.HealthManager
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
	healthManager := registry.NewHealthManager(backend, params.TTLMultiplier, params.Logger)
	watchManager := registry.NewWatchManager(backend)
	grpcHandler := grpc.NewHandler(registryManager, healthManager, watchManager)
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
		healthManager:       healthManager,
		grpcServer:          grpcServer,
		observabilityServer: observabilityServer,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	observabilityErrs, err := a.observabilityServer.Serve()
	if err != nil {
		return fmt.Errorf("serve observability server: %w", err)
	}

	grpcErrs, err := a.grpcServer.Serve()
	if err != nil {
		return fmt.Errorf("serve registry gRPC server: %w", err)
	}

	a.healthManager.Start(ctx)

	a.observabilityServer.SetReady(true)
	a.logger.Info("Registry started")

	select {
	case <-ctx.Done():
		return nil
	case err := <-observabilityErrs:
		return fmt.Errorf("observability server failed: %w", err)

	case err := <-grpcErrs:
		return fmt.Errorf("registry gRPC server failed: %w", err)
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	defer a.logger.Info("Registry stopped")

	a.observabilityServer.SetReady(false)

	errs := make([]error, 2)

	var wg sync.WaitGroup
	wg.Go(func() {
		if err := a.observabilityServer.Shutdown(ctx); err != nil {
			errs[0] = fmt.Errorf("shutdown observability server: %w", err)
		}
	})

	wg.Go(func() {
		if err := a.grpcServer.Shutdown(ctx); err != nil {
			errs[1] = fmt.Errorf("shutdown registry gRPC server: %w", err)
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
