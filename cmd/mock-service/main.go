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

	cfg, err := config.LoadMockService()
	if err != nil {
		return fmt.Errorf("load mock-service configuration: %w", err)
	}

	logger, err := logging.New(cfg.Log)
	if err != nil {
		return fmt.Errorf("initialize logger: %w", err)
	}

	processMetrics := metrics.New("mock-service", cfg.InstanceID)

	observabilityServer := observability.New(cfg.ObservabilityAddress, processMetrics.Handler(), logger)
	go func() {
		if err := observabilityServer.Start(); err != nil {
			logger.Error("Observability server stopped", "error", err)
		}
	}()
	observabilityServer.SetReady(true)

	logger.Info(
		"Mock Service process started",
		"service_name", cfg.ServiceName,
		"instance_id", cfg.InstanceID,
		"http_address", cfg.HTTPAddress,
	)

	<-signalCtx.Done()

	logger.Info("Mock Service process stopping")

	observabilityServer.SetReady(false)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := observabilityServer.Shutdown(shutdownCtx); err != nil {
		return err
	}

	logger.Info("Mock Service process stopped")
	return nil
}
