package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tpiii.local/daniel-tpiii/internal/config"
	"tpiii.local/daniel-tpiii/internal/logging"
	"tpiii.local/daniel-tpiii/internal/registry/app"
)

const shutdownTimeout = 5 * time.Second

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() (err error) {
	cfg, err := config.LoadRegistry()
	if err != nil {
		return fmt.Errorf("load registry configuration: %w", err)
	}

	logger, err := logging.New(cfg.Log)
	if err != nil {
		return fmt.Errorf("initialize logger: %w", err)
	}

	ctx, stopCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopCtx()

	application, err := app.New(app.Params{
		NodeID:               cfg.NodeID,
		Backend:              cfg.Backend,
		GRPCAddress:          cfg.GRPCAddress,
		ObservabilityAddress: cfg.ObservabilityAddress,
		ClusterMembers:       cfg.ClusterMembers,
		Logger:               logger,
	})
	if err != nil {
		return fmt.Errorf("initialize registry application: %w", err)
	}

	defer func() {
		err = errors.Join(err, stop(application))
	}()

	if err := application.Run(ctx); err != nil {
		return fmt.Errorf("run registry application: %w", err)
	}

	return nil
}

func stop(application *app.App) error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return application.Stop(ctx)
}
