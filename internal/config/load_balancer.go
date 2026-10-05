package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"tpiii.local/daniel-tpiii/internal/logging"
)

const (
	defaultRetryInterval    = time.Second
	defaultKeepaliveTime    = 10 * time.Second
	defaultKeepaliveTimeout = 5 * time.Second
	minKeepaliveTime        = 10 * time.Second
)

type LoadBalancer struct {
	HTTPAddress              string
	ObservabilityAddress     string
	RegistryAddresses        []string
	RegistryRetryInterval    time.Duration
	RegistryKeepaliveTime    time.Duration
	RegistryKeepaliveTimeout time.Duration
	Services                 map[string]string
	Log                      logging.Config
}

func LoadLoadBalancer() (LoadBalancer, error) {
	retryInterval, err := durationFromEnvOrDefault("REGISTRY_RETRY_INTERVAL", defaultRetryInterval)
	if err != nil {
		return LoadBalancer{}, err
	}

	keepaliveTime, err := durationFromEnvOrDefault("REGISTRY_KEEPALIVE_TIME", defaultKeepaliveTime)
	if err != nil {
		return LoadBalancer{}, err
	}

	keepaliveTimeout, err := durationFromEnvOrDefault("REGISTRY_KEEPALIVE_TIMEOUT", defaultKeepaliveTimeout)
	if err != nil {
		return LoadBalancer{}, err
	}

	services, err := parseServices(os.Getenv("SERVICES"))
	if err != nil {
		return LoadBalancer{}, err
	}

	cfg := LoadBalancer{
		HTTPAddress:              stringFromEnvOrDefault("HTTP_ADDRESS", defaultLBHTTPAddress),
		ObservabilityAddress:     stringFromEnvOrDefault("OBSERVABILITY_ADDRESS", defaultLBObservabilityAddress),
		RegistryAddresses:        splitCSV(os.Getenv("REGISTRY_ADDRESSES")),
		RegistryRetryInterval:    retryInterval,
		RegistryKeepaliveTime:    keepaliveTime,
		RegistryKeepaliveTimeout: keepaliveTimeout,
		Services:                 services,
		Log: logging.Config{
			Level:  stringFromEnvOrDefault("LOG_LEVEL", "info"),
			Format: stringFromEnvOrDefault("LOG_FORMAT", "json"),
		},
	}

	var errs []error

	if len(cfg.RegistryAddresses) == 0 {
		errs = append(errs, errors.New("REGISTRY_ADDRESSES is required"))
	}

	if len(cfg.Services) == 0 {
		errs = append(errs, errors.New("SERVICES is required"))
	}

	if cfg.RegistryRetryInterval <= 0 {
		errs = append(errs, errors.New("REGISTRY_RETRY_INTERVAL must be greater than zero"))
	}

	if cfg.RegistryKeepaliveTime < minKeepaliveTime {
		errs = append(errs, fmt.Errorf("REGISTRY_KEEPALIVE_TIME must be at least %s", minKeepaliveTime))
	}

	if cfg.RegistryKeepaliveTimeout <= 0 {
		errs = append(errs, errors.New("REGISTRY_KEEPALIVE_TIMEOUT must be greater than zero"))
	}

	return cfg, errors.Join(errs...)
}

func parseServices(raw string) (map[string]string, error) {
	services := make(map[string]string)

	for _, item := range splitCSV(raw) {
		name, strategy, ok := strings.Cut(item, ":")
		name, strategy = strings.TrimSpace(name), strings.TrimSpace(strategy)

		if !ok || name == "" || strategy == "" {
			return nil, fmt.Errorf("SERVICES entry %q must have the form service:strategy", item)
		}

		if _, exists := services[name]; exists {
			return nil, fmt.Errorf("SERVICES entry %q is duplicated", name)
		}

		services[name] = strategy
	}

	return services, nil
}
