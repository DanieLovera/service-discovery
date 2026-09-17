package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"tpiii.local/daniel-tpiii/internal/config"
	"tpiii.local/daniel-tpiii/internal/logging"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadMockService()
	if err != nil {
		return fmt.Errorf("Load mock-service configuration: %w", err)
	}

	logger, err := logging.New(cfg.Log)
	if err != nil {
		return fmt.Errorf("Initialize logger: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info(
		"Mock Service process started",
		"service_name", cfg.ServiceName,
		"instance_id", cfg.InstanceID,
		"http_address", cfg.HTTPAddress,
	)
	<-ctx.Done()
	logger.Info("Mock Service process stopping")

	return nil
}
