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
	cfg, err := config.LoadLoadBalancer()
	if err != nil {
		return fmt.Errorf("Load load-balancer configuration: %w", err)
	}

	logger, err := logging.New(cfg.Log)
	if err != nil {
		return fmt.Errorf("Initialize logger: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info(
		"Load Balancer process started",
		"http_address", cfg.HTTPAddress,
		"registry_addresses", cfg.RegistryAddresses,
	)
	<-ctx.Done()
	logger.Info("Load Balancer process stopping")

	return nil
}
