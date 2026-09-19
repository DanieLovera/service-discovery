package config

import (
	"errors"
	"os"

	"tpiii.local/daniel-tpiii/internal/logging"
)

type LoadBalancer struct {
	HTTPAddress          string
	ObservabilityAddress string
	RegistryAddresses    []string
	Log                  logging.Config
}

func LoadLoadBalancer() (LoadBalancer, error) {
	cfg := LoadBalancer{
		HTTPAddress:          stringFromEnvOrDefault("HTTP_ADDRESS", defaultLBHTTPAddress),
		ObservabilityAddress: stringFromEnvOrDefault("OBSERVABILITY_ADDRESS", defaultLBObservabilityAddress),
		RegistryAddresses:    splitCSV(os.Getenv("REGISTRY_ADDRESSES")),
		Log: logging.Config{
			Level:  stringFromEnvOrDefault("LOG_LEVEL", "info"),
			Format: stringFromEnvOrDefault("LOG_FORMAT", "json"),
		},
	}
	if len(cfg.RegistryAddresses) == 0 {
		return cfg, errors.New("REGISTRY_ADDRESSES is required")
	}
	return cfg, nil
}
