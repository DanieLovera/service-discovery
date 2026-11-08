package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"tpiii.local/daniel-tpiii/internal/config"
	"tpiii.local/daniel-tpiii/internal/load-generator/app"
	"tpiii.local/daniel-tpiii/internal/logging"
)

func main() {
	if err := start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func start() (err error) {
	cfg, err := config.LoadLoadGenerator()
	if err != nil {
		return fmt.Errorf("load load generator configuration: %w", err)
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
		RegistryAddresses:           cfg.RegistryAddresses,
		LoadBalancerURL:             cfg.LoadBalancerURL,
		LoadBalancerServices:        cfg.LoadBalancerServices,
		Duration:                    cfg.Duration,
		RequestTimeout:              cfg.RequestTimeout,
		RegistryCycleInterval:       cfg.RegistryCycleInterval,
		VisibilityProbeInterval:     cfg.VisibilityProbeInterval,
		LoadBalancerRequestInterval: cfg.LoadBalancerRequestInterval,
		OutputDir:                   cfg.OutputDir,
		Logger:                      logger,
	})
	if err != nil {
		return fmt.Errorf("initialize load generator application: %w", err)
	}

	defer func() {
		err = errors.Join(err, application.Close())
	}()

	application.Run(ctx)

	return nil
}
