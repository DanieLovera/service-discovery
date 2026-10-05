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
	"tpiii.local/daniel-tpiii/internal/mock-service/app"
)

const shutdownTimeout = 5 * time.Second

func main() {
	if err := start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func start() (err error) {
	cfg, err := config.LoadMockService()
	if err != nil {
		return fmt.Errorf("load mock service configuration: %w", err)
	}

	logger, err := logging.New(cfg.Log)
	if err != nil {
		return fmt.Errorf("initialize logger: %w", err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	application, err := app.New(app.Params{
		ServiceName:          cfg.ServiceName,
		InstanceID:           cfg.InstanceID,
		HTTPAddress:          cfg.HTTPAddress,
		AdvertiseAddress:     cfg.AdvertiseAddress,
		ObservabilityAddress: cfg.ObservabilityAddress,
		RegistryAddresses:    cfg.RegistryAddresses,
		Weight:               cfg.Weight,
		HeartbeatInterval:    cfg.HeartbeatInterval,
		Logger:               logger,
	})
	if err != nil {
		return fmt.Errorf("initialize mock service application: %w", err)
	}

	defer func() {
		err = errors.Join(err, shutdown(application))
	}()

	if err := application.Start(ctx); err != nil {
		return fmt.Errorf("start mock service application: %w", err)
	}

	return nil
}

func shutdown(application *app.App) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout,
	)
	defer cancel()

	return application.Shutdown(ctx)
}
