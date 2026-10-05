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
	if err := start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func start() (err error) {
	cfg, err := config.LoadRegistry()
	if err != nil {
		return fmt.Errorf("load registry configuration: %w", err)
	}

	logger, err := logging.New(cfg.Log)
	if err != nil {
		return fmt.Errorf("initialize logger: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application, err := app.New(app.Params{
		NodeID:               cfg.NodeID,
		Backend:              cfg.Backend,
		GRPCAddress:          cfg.GRPCAddress,
		ObservabilityAddress: cfg.ObservabilityAddress,
		ClusterMembers:       cfg.ClusterMembers,
		TTLMultiplier:        cfg.TTLMultiplier,
		Logger:               logger,
	})
	if err != nil {
		return fmt.Errorf("initialize registry application: %w", err)
	}

	defer func() {
		err = errors.Join(err, shutdown(application))
	}()

	if err := application.Start(ctx); err != nil {
		return fmt.Errorf("start registry application: %w", err)
	}
	return nil
}

func shutdown(application *app.App) error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return application.Shutdown(ctx)
}
