package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"tpiii.local/daniel-tpiii/internal/logging"
)

const (
	defaultRegistryGRPCAddress          = "0.0.0.0:7000"
	defaultRegistryObservabilityAddress = "0.0.0.0:9100"
)

type Registry struct {
	NodeID               string
	Backend              string
	GRPCAddress          string
	ObservabilityAddress string
	ClusterMembers       []string
	Log                  logging.Config
}

func LoadRegistry() (Registry, error) {
	cfg := Registry{
		NodeID:               strings.TrimSpace(os.Getenv("REGISTRY_INSTANCE_ID")),
		Backend:              strings.ToLower(stringFromEnvOrDefault("REGISTRY_BACKEND", "ap")),
		GRPCAddress:          stringFromEnvOrDefault("REGISTRY_GRPC_ADDRESS", defaultRegistryGRPCAddress),
		ObservabilityAddress: stringFromEnvOrDefault("REGISTRY_OBSERVABILITY_ADDRESS", defaultRegistryObservabilityAddress),
		ClusterMembers:       splitCSV(os.Getenv("REGISTRY_CLUSTER_MEMBERS")),
		Log: logging.Config{
			Level:  stringFromEnvOrDefault("REGISTRY_LOG_LEVEL", "info"),
			Format: stringFromEnvOrDefault("REGISTRY_LOG_FORMAT", "json"),
		},
	}
	var errs []error
	if cfg.NodeID == "" {
		errs = append(errs, errors.New("REGISTRY_INSTANCE_ID is required"))
	}
	if cfg.Backend != "cp" && cfg.Backend != "ap" {
		errs = append(errs, fmt.Errorf("REGISTRY_BACKEND must be one of cp, ap: %q", cfg.Backend))
	}
	if len(cfg.ClusterMembers) == 0 {
		errs = append(errs, errors.New("REGISTRY_CLUSTER_MEMBERS is required"))
	}
	return cfg, errors.Join(errs...)
}
