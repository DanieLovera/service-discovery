package config

import (
	"errors"
	"os"

	"tpiii.local/daniel-tpiii/internal/logging"
)

const defaultLBHTTPAddress = "0.0.0.0:8080"

type LoadBalancer struct {
	HTTPAddress       string
	RegistryAddresses []string
	Log               logging.Config
}

func LoadLoadBalancer() (LoadBalancer, error) {
	cfg := LoadBalancer{
		HTTPAddress:       stringFromEnvOrDefault("LB_HTTP_ADDRESS", defaultLBHTTPAddress),
		RegistryAddresses: splitCSV(os.Getenv("LB_REGISTRY_ADDRESSES")),
		Log: logging.Config{
			Level:  stringFromEnvOrDefault("LB_LOG_LEVEL", "info"),
			Format: stringFromEnvOrDefault("LB_LOG_FORMAT", "json"),
		},
	}
	if len(cfg.RegistryAddresses) == 0 {
		return cfg, errors.New("LB_REGISTRY_ADDRESSES is required")
	}
	return cfg, nil
}
