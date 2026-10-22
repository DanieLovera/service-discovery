package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	googlegrpc "google.golang.org/grpc"

	"tpiii.local/daniel-tpiii/internal/metrics"
	"tpiii.local/daniel-tpiii/internal/observability"
	"tpiii.local/daniel-tpiii/internal/registry"
	"tpiii.local/daniel-tpiii/internal/registry/ap"
	"tpiii.local/daniel-tpiii/internal/registry/cp"
	"tpiii.local/daniel-tpiii/internal/registry/grpc"
)

type Params struct {
	NodeID                   string
	Backend                  string
	GRPCBindAddress          string
	ObservabilityBindAddress string
	ClusterMembers           []registry.ClusterMember
	TTLMultiplier            int
	RaftBindAddress          string
	RaftDataDir              string
	ForwardTimeout           time.Duration
	ForwardMaxAttempts       int
	Logger                   *slog.Logger
}

type grpcRegistrar interface {
	RegisterGRPC(registrar googlegrpc.ServiceRegistrar)
}

type App struct {
	backend             registry.Backend
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
	var registrations []func(googlegrpc.ServiceRegistrar)
	if registrar, ok := backend.(grpcRegistrar); ok {
		registrations = append(registrations, registrar.RegisterGRPC)
	}

	grpcServer := grpc.NewServer(params.GRPCBindAddress, grpcHandler, params.Logger, registrations...)

	processMetrics := metrics.New("registry", params.NodeID)
	metricsHandler := processMetrics.Handler()
	observabilityServer := observability.NewServer(
		params.ObservabilityBindAddress,
		metricsHandler,
		params.Logger,
	)

	return &App{
		logger:              params.Logger,
		backend:             backend,
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

	if err := a.backend.Start(ctx); err != nil {
		return fmt.Errorf("start registry backend: %w", err)
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

	errs := make([]error, 3)

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

	if err := a.backend.Shutdown(ctx); err != nil {
		errs[2] = fmt.Errorf("shutdown registry backend: %w", err)
	}

	return errors.Join(errs...)
}

func newBackend(params Params) (registry.Backend, error) {
	switch params.Backend {
	case "cp":
		return cp.NewBackend(cp.Config{
			NodeID:             params.NodeID,
			BindAddress:        params.RaftBindAddress,
			DataDir:            params.RaftDataDir,
			Members:            params.ClusterMembers,
			ForwardTimeout:     params.ForwardTimeout,
			ForwardMaxAttempts: params.ForwardMaxAttempts,
		}), nil
	case "ap":
		return ap.NewBackend(), nil
	default:
		return nil, fmt.Errorf("unsupported registry backend %q", params.Backend)
	}
}
