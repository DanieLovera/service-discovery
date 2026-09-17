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
	cfg, err := config.LoadRegistry()
	if err != nil {
		return fmt.Errorf("Load registry configuration: %w", err)
	}

	logger, err := logging.New(cfg.Log)
	if err != nil {
		return fmt.Errorf("Initialize logger: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info(
		"Registry process started",
		"node_id", cfg.NodeID,
		"backend", cfg.Backend,
		"grpc_address", cfg.GRPCAddress,
	)
	<-ctx.Done()
	logger.Info("Registry process stopping")

	return nil
}
