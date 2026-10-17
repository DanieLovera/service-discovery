package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"tpiii.local/daniel-tpiii/internal/logging"
)

const (
	defaultTTLMultiplier = 3
	defaultRaftDataDir   = "/data/raft"
)

type Registry struct {
	NodeID               string
	Backend              string
	GRPCAddress          string
	ObservabilityAddress string
	ClusterMembers       []string
	TTLMultiplier        int
	RaftBindAddress      string
	RaftAdvertiseAddress string
	RaftDataDir          string
	Log                  logging.Config
}

func LoadRegistry() (Registry, error) {
	ttlMultiplier, err := intFromEnvOrDefault("HEARTBEAT_TTL_MULTIPLIER", defaultTTLMultiplier)
	if err != nil {
		return Registry{}, err
	}

	cfg := Registry{
		NodeID:               strings.TrimSpace(os.Getenv("INSTANCE_ID")),
		Backend:              strings.ToLower(stringFromEnvOrDefault("BACKEND", "ap")),
		GRPCAddress:          stringFromEnvOrDefault("GRPC_ADDRESS", defaultRegistryGRPCAddress),
		ObservabilityAddress: stringFromEnvOrDefault("OBSERVABILITY_ADDRESS", defaultRegistryObservabilityAddress),
		ClusterMembers:       splitCSV(os.Getenv("CLUSTER_MEMBERS")),
		TTLMultiplier:        ttlMultiplier,
		RaftBindAddress:      stringFromEnvOrDefault("RAFT_BIND_ADDRESS", defaultRegistryRaftAddress),
		RaftAdvertiseAddress: strings.TrimSpace(os.Getenv("RAFT_ADVERTISE_ADDRESS")),
		RaftDataDir:          stringFromEnvOrDefault("RAFT_DATA_DIR", defaultRaftDataDir),
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
	if cfg.TTLMultiplier <= 0 {
		errs = append(errs, errors.New("HEARTBEAT_TTL_MULTIPLIER must be greater than zero"))
	}
	if cfg.Backend == "cp" && cfg.RaftAdvertiseAddress == "" {
		errs = append(errs, errors.New("RAFT_ADVERTISE_ADDRESS is required for the cp backend"))
	}
	return cfg, errors.Join(errs...)
}
