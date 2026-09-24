package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tpiii.local/daniel-tpiii/internal/config"
	"tpiii.local/daniel-tpiii/internal/logging"
	"tpiii.local/daniel-tpiii/internal/metrics"
	"tpiii.local/daniel-tpiii/internal/observability"
)

const shutdownTimeout = 5 * time.Second

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadLoadBalancer()
	if err != nil {
		return fmt.Errorf("load load-balancer configuration: %w", err)
	}

	logger, err := logging.New(cfg.Log)
	if err != nil {
		return fmt.Errorf("initialize logger: %w", err)
	}

	processMetrics := metrics.New("load-balancer", "load-balancer")
	observabilityServer := observability.NewServer(cfg.ObservabilityAddress, processMetrics.Handler(), logger)

	if err := observabilityServer.Listen(); err != nil {
		return fmt.Errorf("listen observability server: %w", err)
	}

	serveErrs := make(chan error, 1)
	go func() {
		if err := observabilityServer.Serve(); err != nil {
			serveErrs <- fmt.Errorf("serve observability server: %w", err)
		}
	}()

	observabilityServer.SetReady(true)

	logger.Info(
		"Load Balancer started",
		"http_address", cfg.HTTPAddress,
		"registry_addresses", cfg.RegistryAddresses,
	)

	select {
	case <-signalCtx.Done():
	case err := <-serveErrs:
		return err
	}

	logger.Info("Load Balancer process stopping")

	observabilityServer.SetReady(false)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := observabilityServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown observability server: %w", err)
	}

	logger.Info("Load Balancer process stopped")
	return nil
}
