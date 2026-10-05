package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	loadbalancer "tpiii.local/daniel-tpiii/internal/load-balancer"
	"tpiii.local/daniel-tpiii/internal/load-balancer/discovery"
	"tpiii.local/daniel-tpiii/internal/load-balancer/http"
	"tpiii.local/daniel-tpiii/internal/load-balancer/registry"
	"tpiii.local/daniel-tpiii/internal/load-balancer/routing"
	"tpiii.local/daniel-tpiii/internal/metrics"
	"tpiii.local/daniel-tpiii/internal/observability"
)

type Params struct {
	HTTPAddress              string
	ObservabilityAddress     string
	RegistryAddresses        []string
	RegistryRetryInterval    time.Duration
	RegistryKeepaliveTime    time.Duration
	RegistryKeepaliveTimeout time.Duration
	Services                 map[string]string
	Logger                   *slog.Logger
}

type App struct {
	discoveryManager    *discovery.Manager
	httpServer          *http.Server
	observabilityServer *observability.Server
	logger              *slog.Logger
}

func New(params Params) (*App, error) {
	pools := make(map[string]*loadbalancer.ServicePool, len(params.Services))
	upstreams := make(map[string]routing.Upstream, len(params.Services))

	for serviceName, strategyName := range params.Services {
		strategy, err := routing.NewStrategy(strategyName)
		if err != nil {
			return nil, fmt.Errorf("service %q: %w", serviceName, err)
		}

		servicePool := loadbalancer.NewServicePool()
		pools[serviceName] = servicePool
		upstreams[serviceName] = routing.Upstream{Pool: servicePool, Strategy: strategy}
	}

	discoveryManager := discovery.NewManager(
		params.RegistryAddresses,
		pools,
		params.RegistryRetryInterval,
		registry.KeepaliveConfig{
			Time:    params.RegistryKeepaliveTime,
			Timeout: params.RegistryKeepaliveTimeout,
		},
		params.Logger,
	)

	httpServer := http.NewServer(params.HTTPAddress, routing.NewManager(upstreams), params.Logger)

	processMetrics := metrics.New("load-balancer", "load-balancer")
	observabilityServer := observability.NewServer(
		params.ObservabilityAddress,
		processMetrics.Handler(),
		params.Logger,
	)

	return &App{
		discoveryManager:    discoveryManager,
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

	if err := a.discoveryManager.Start(ctx); err != nil {
		return fmt.Errorf("start discovery manager: %w", err)
	}

	a.observabilityServer.SetReady(true)
	a.logger.Info("Load balancer started")

	select {
	case <-ctx.Done():
		return nil

	case err := <-observabilityErrs:
		return fmt.Errorf("observability server failed: %w", err)

	case err := <-httpErrs:
		return fmt.Errorf("load balancer HTTP server failed: %w", err)
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	defer a.logger.Info("Load balancer stopped")

	a.observabilityServer.SetReady(false)

	errs := make([]error, 3)

	var wg sync.WaitGroup

	wg.Go(func() {
		if err := a.observabilityServer.Shutdown(ctx); err != nil {
			errs[0] = fmt.Errorf("shutdown observability server: %w", err)
		}
	})

	wg.Go(func() {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			errs[1] = fmt.Errorf("shutdown load balancer HTTP server: %w", err)
		}
	})

	wg.Go(func() {
		if err := a.discoveryManager.Shutdown(ctx); err != nil {
			errs[2] = fmt.Errorf("shutdown discovery manager: %w", err)
		}
	})

	wg.Wait()

	return errors.Join(errs...)
}
