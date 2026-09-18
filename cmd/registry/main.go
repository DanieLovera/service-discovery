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

	cfg, err := config.LoadRegistry()
	if err != nil {
		return fmt.Errorf("load registry configuration: %w", err)
	}

	logger, err := logging.New(cfg.Log)
	if err != nil {
		return fmt.Errorf("initialize logger: %w", err)
	}

	processMetrics := metrics.New("registry", cfg.NodeID)

	observabilityServer := observability.New(cfg.ObservabilityAddress, processMetrics.Handler(), logger)
	go func() {
		if err := observabilityServer.Start(); err != nil {
			logger.Error("Observability server stopped", "error", err)
		}
	}()
	observabilityServer.SetReady(true)

	logger.Info(
		"Registry process started",
		"node_id", cfg.NodeID,
		"backend", cfg.Backend,
		"grpc_address", cfg.GRPCAddress,
	)

	<-signalCtx.Done()

	logger.Info("Registry process stopping")

	observabilityServer.SetReady(false)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := observabilityServer.Shutdown(shutdownCtx); err != nil {
		return err
	}

	logger.Info("Registry process stopped")
	return nil
}
