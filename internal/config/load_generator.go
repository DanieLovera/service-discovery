package config

import (
	"errors"
	"os"
	"strings"
	"time"

	"tpiii.local/daniel-tpiii/internal/logging"
)

const (
	defaultLoadGeneratorDuration       = time.Minute
	defaultRequestTimeout              = 2 * time.Second
	defaultRegistryCycleInterval       = 200 * time.Millisecond
	defaultVisibilityProbeInterval     = time.Second
	defaultLoadBalancerRequestInterval = 50 * time.Millisecond
	defaultOutputDir                   = "/results"
)

type LoadGenerator struct {
	RegistryAddresses           []string
	LoadBalancerURL             string
	LoadBalancerServices        []string
	Duration                    time.Duration
	RequestTimeout              time.Duration
	RegistryCycleInterval       time.Duration
	VisibilityProbeInterval     time.Duration
	LoadBalancerRequestInterval time.Duration
	OutputDir                   string
	Log                         logging.Config
}

func LoadLoadGenerator() (LoadGenerator, error) {
	duration, err := durationFromEnvOrDefault("DURATION", defaultLoadGeneratorDuration)
	if err != nil {
		return LoadGenerator{}, err
	}

	requestTimeout, err := durationFromEnvOrDefault("REQUEST_TIMEOUT", defaultRequestTimeout)
	if err != nil {
		return LoadGenerator{}, err
	}

	cycleInterval, err := durationFromEnvOrDefault("REGISTRY_CYCLE_INTERVAL", defaultRegistryCycleInterval)
	if err != nil {
		return LoadGenerator{}, err
	}

	probeInterval, err := durationFromEnvOrDefault("VISIBILITY_PROBE_INTERVAL", defaultVisibilityProbeInterval)
	if err != nil {
		return LoadGenerator{}, err
	}

	requestInterval, err := durationFromEnvOrDefault(
		"LOAD_BALANCER_REQUEST_INTERVAL",
		defaultLoadBalancerRequestInterval,
	)
	if err != nil {
		return LoadGenerator{}, err
	}

	cfg := LoadGenerator{
		RegistryAddresses:           splitCSV(os.Getenv("REGISTRY_ADDRESSES")),
		LoadBalancerURL:             strings.TrimSpace(os.Getenv("LOAD_BALANCER_URL")),
		LoadBalancerServices:        splitCSV(os.Getenv("LOAD_BALANCER_SERVICES")),
		Duration:                    duration,
		RequestTimeout:              requestTimeout,
		RegistryCycleInterval:       cycleInterval,
		VisibilityProbeInterval:     probeInterval,
		LoadBalancerRequestInterval: requestInterval,
		OutputDir:                   stringFromEnvOrDefault("OUTPUT_DIR", defaultOutputDir),
		Log: logging.Config{
			Level:  stringFromEnvOrDefault("LOG_LEVEL", "info"),
			Format: stringFromEnvOrDefault("LOG_FORMAT", "json"),
		},
	}

	var errs []error

	if len(cfg.RegistryAddresses) == 0 {
		errs = append(errs, errors.New("REGISTRY_ADDRESSES is required"))
	}

	if cfg.LoadBalancerURL == "" {
		errs = append(errs, errors.New("LOAD_BALANCER_URL is required"))
	}

	if len(cfg.LoadBalancerServices) == 0 {
		errs = append(errs, errors.New("LOAD_BALANCER_SERVICES is required"))
	}

	if cfg.Duration <= 0 {
		errs = append(errs, errors.New("DURATION must be greater than zero"))
	}

	if cfg.RequestTimeout <= 0 {
		errs = append(errs, errors.New("REQUEST_TIMEOUT must be greater than zero"))
	}

	if cfg.RegistryCycleInterval <= 0 {
		errs = append(errs, errors.New("REGISTRY_CYCLE_INTERVAL must be greater than zero"))
	}

	if cfg.VisibilityProbeInterval <= 0 {
		errs = append(errs, errors.New("VISIBILITY_PROBE_INTERVAL must be greater than zero"))
	}

	if cfg.LoadBalancerRequestInterval <= 0 {
		errs = append(errs, errors.New("LOAD_BALANCER_REQUEST_INTERVAL must be greater than zero"))
	}

	return cfg, errors.Join(errs...)
}
