package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"tpiii.local/daniel-tpiii/internal/metrics"
	"tpiii.local/daniel-tpiii/internal/mock-service/http"
	"tpiii.local/daniel-tpiii/internal/mock-service/registry"
	"tpiii.local/daniel-tpiii/internal/observability"
)

type Params struct {
	ServiceName          string
	InstanceID           string
	HTTPAddress          string
	AdvertiseAddress     string
	ObservabilityAddress string
	RegistryAddresses    []string
	Weight               int
	HeartbeatInterval    time.Duration
	Logger               *slog.Logger
}

type App struct {
	instance            registry.Instance
	registryClient      *registry.Client
	httpServer          *http.Server
	observabilityServer *observability.Server
	logger              *slog.Logger
}

func New(params Params) (*App, error) {
	if len(params.RegistryAddresses) == 0 {
		return nil, errors.New("at least one registry address is required")
	}

	registryClient, err := registry.NewClient(params.RegistryAddresses)
	if err != nil {
		return nil, fmt.Errorf("initialize registry client: %w", err)
	}

	processMetrics := metrics.New("mock-service", params.InstanceID)
	metricsHandler := processMetrics.Handler()

	httpServer := http.NewServer(
		params.HTTPAddress,
		params.ServiceName,
		params.InstanceID,
		params.Logger,
	)

	observabilityServer := observability.NewServer(
		params.ObservabilityAddress,
		metricsHandler,
		params.Logger,
	)

	return &App{
		instance: registry.Instance{
			ServiceName:       params.ServiceName,
			InstanceID:        params.InstanceID,
			Address:           params.AdvertiseAddress,
			Weight:            params.Weight,
			HeartbeatInterval: params.HeartbeatInterval,
		},
		registryClient:      registryClient,
		httpServer:          httpServer,
		observabilityServer: observabilityServer,
		logger:              params.Logger,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	observabilityErrs, err := a.observabilityServer.Serve()
	if err != nil {
		return fmt.Errorf("serve observability server: %w", err)
	}

	httpErrs, err := a.httpServer.Serve()
	if err != nil {
		return fmt.Errorf("serve HTTP server: %w", err)
	}

	if err := a.registryClient.Register(ctx, a.instance); err != nil {
		return fmt.Errorf("register service instance: %w", err)
	}

	go a.runHeartbeats(ctx)

	a.observabilityServer.SetReady(true)
	a.logger.Info("Mock service started")

	select {
	case <-ctx.Done():
		return nil

	case err := <-observabilityErrs:
		return fmt.Errorf("observability server failed: %w", err)

	case err := <-httpErrs:
		return fmt.Errorf("mock HTTP server failed: %w", err)
	}
}

func (a *App) Stop(ctx context.Context) error {
	defer a.logger.Info("Mock service stopped")

	a.observabilityServer.SetReady(false)

	errs := make([]error, 4)

	if err := a.registryClient.Deregister(ctx, a.instance); err != nil {
		errs[0] = fmt.Errorf("deregister service instance: %w", err)
	}

	var wg sync.WaitGroup

	wg.Go(func() {
		if err := a.observabilityServer.Shutdown(ctx); err != nil {
			errs[1] = fmt.Errorf("shutdown observability server: %w", err)
		}
	})

	wg.Go(func() {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			errs[2] = fmt.Errorf("shutdown mock HTTP server: %w", err)
		}
	})

	wg.Wait()

	if err := a.registryClient.Close(); err != nil {
		errs[3] = fmt.Errorf("close registry connection: %w", err)
	}

	return errors.Join(errs...)
}

func (a *App) runHeartbeats(ctx context.Context) {
	ticker := time.NewTicker(a.instance.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := a.registryClient.Heartbeat(ctx, a.instance); err != nil {
				if ctx.Err() != nil {
					return
				}

				a.logger.Error("Failed to send heartbeat", "error", err)
			}

		case <-ctx.Done():
			return
		}
	}
}
