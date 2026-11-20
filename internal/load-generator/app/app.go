package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"tpiii.local/daniel-tpiii/internal/load-generator/http"
	"tpiii.local/daniel-tpiii/internal/load-generator/registry"
	"tpiii.local/daniel-tpiii/internal/load-generator/results"
)

type Params struct {
	RegistryAddresses           []string
	LoadBalancerURL             string
	LoadBalancerServices        []string
	Duration                    time.Duration
	RequestTimeout              time.Duration
	RegistryCycleInterval       time.Duration
	VisibilityProbeInterval     time.Duration
	LoadBalancerRequestInterval time.Duration
	OutputDir                   string
	Logger                      *slog.Logger
}

type workload interface {
	Run(ctx context.Context)
}

type App struct {
	nodes     []*registry.Node
	writers   []*results.Writer
	workloads []workload
	duration  time.Duration
	outputDir string
	logger    *slog.Logger
}

func New(params Params) (*App, error) {
	if err := os.MkdirAll(params.OutputDir, 0o755); err != nil {
		return nil, fmt.Errorf("create output directory: %w", err)
	}

	a := &App{
		duration:  params.Duration,
		outputDir: params.OutputDir,
		logger:    params.Logger,
	}

	operationsWriter, err := a.createWriter("operations.csv", registry.OperationsHeader)
	if err != nil {
		return nil, errors.Join(err, a.Close())
	}

	probesWriter, err := a.createWriter("probes.csv", registry.ProbesHeader)
	if err != nil {
		return nil, errors.Join(err, a.Close())
	}

	sightingsWriter, err := a.createWriter("sightings.csv", registry.SightingsHeader)
	if err != nil {
		return nil, errors.Join(err, a.Close())
	}

	trafficWriter, err := a.createWriter("load_balancer.csv", http.TrafficHeader)
	if err != nil {
		return nil, errors.Join(err, a.Close())
	}

	for _, address := range params.RegistryAddresses {
		node, err := registry.Connect(address)
		if err != nil {
			return nil, errors.Join(err, a.Close())
		}

		a.nodes = append(a.nodes, node)
		a.workloads = append(a.workloads, registry.NewOperations(
			node,
			params.RegistryCycleInterval,
			params.RequestTimeout,
			operationsWriter,
		))
	}

	a.workloads = append(a.workloads, registry.NewVisibility(
		a.nodes,
		params.VisibilityProbeInterval,
		params.RequestTimeout,
		probesWriter,
		sightingsWriter,
	))

	for _, service := range params.LoadBalancerServices {
		a.workloads = append(a.workloads, http.NewTraffic(
			params.LoadBalancerURL,
			service,
			params.LoadBalancerRequestInterval,
			params.RequestTimeout,
			trafficWriter,
		))
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, a.duration)
	defer cancel()

	a.logger.Info("Load generator started", "duration", a.duration, "output_dir", a.outputDir)

	var wg sync.WaitGroup
	for _, w := range a.workloads {
		wg.Go(func() {
			w.Run(ctx)
		})
	}

	wg.Wait()

	a.logger.Info("Load generator finished")
}

func (a *App) Close() error {
	errs := make([]error, 0, len(a.nodes)+len(a.writers))

	for _, node := range a.nodes {
		errs = append(errs, node.Close())
	}

	for _, writer := range a.writers {
		errs = append(errs, writer.Close())
	}

	return errors.Join(errs...)
}

func (a *App) createWriter(name string, header []string) (*results.Writer, error) {
	writer, err := results.Create(a.outputDir, name, header...)
	if err != nil {
		return nil, err
	}

	a.writers = append(a.writers, writer)

	return writer, nil
}
