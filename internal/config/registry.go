package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"tpiii.local/daniel-tpiii/internal/logging"
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
		NodeID:               strings.TrimSpace(os.Getenv("INSTANCE_ID")),
		Backend:              strings.ToLower(stringFromEnvOrDefault("BACKEND", "ap")),
		GRPCAddress:          stringFromEnvOrDefault("GRPC_ADDRESS", defaultRegistryGRPCAddress),
		ObservabilityAddress: stringFromEnvOrDefault("OBSERVABILITY_ADDRESS", defaultRegistryObservabilityAddress),
		ClusterMembers:       splitCSV(os.Getenv("CLUSTER_MEMBERS")),
		Log: logging.Config{
			Level:  stringFromEnvOrDefault("LOG_LEVEL", "info"),
			Format: stringFromEnvOrDefault("LOG_FORMAT", "json"),
		},
	}
	var errs []error
	if cfg.NodeID == "" {
		errs = append(errs, errors.New("INSTANCE_ID is required"))
	}
	if cfg.Backend != "cp" && cfg.Backend != "ap" {
		errs = append(errs, fmt.Errorf("BACKEND must be one of cp, ap: %q", cfg.Backend))
	}
	if len(cfg.ClusterMembers) == 0 {
		errs = append(errs, errors.New("CLUSTER_MEMBERS is required"))
	}
	return cfg, errors.Join(errs...)
}
