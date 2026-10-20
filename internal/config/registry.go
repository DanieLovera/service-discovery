package config

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"tpiii.local/daniel-tpiii/internal/logging"
	"tpiii.local/daniel-tpiii/internal/registry"
)

const (
	defaultTTLMultiplier      = 3
	defaultRaftDataDir        = "/data/raft"
	defaultForwardTimeout     = 2 * time.Second
	defaultForwardMaxAttempts = 3
)

type Registry struct {
	NodeID               string
	Backend              string
	GRPCAddress          string
	ObservabilityAddress string
	ClusterMembers       []registry.ClusterMember
	TTLMultiplier        int
	RaftBindAddress      string
	RaftDataDir          string
	ForwardTimeout       time.Duration
	ForwardMaxAttempts   int
	Log                  logging.Config
}

func LoadRegistry() (Registry, error) {
	ttlMultiplier, err := intFromEnvOrDefault("HEARTBEAT_TTL_MULTIPLIER", defaultTTLMultiplier)
	if err != nil {
		return Registry{}, err
	}

	forwardTimeout, err := durationFromEnvOrDefault("FORWARD_TIMEOUT", defaultForwardTimeout)
	if err != nil {
		return Registry{}, err
	}

	forwardMaxAttempts, err := intFromEnvOrDefault("FORWARD_MAX_ATTEMPTS", defaultForwardMaxAttempts)
	if err != nil {
		return Registry{}, err
	}

	var errs []error

	clusterMembers, err := parseClusterMembers(os.Getenv("CLUSTER_MEMBERS"))
	if err != nil {
		errs = append(errs, err)
	}

	cfg := Registry{
		NodeID:               strings.TrimSpace(os.Getenv("INSTANCE_ID")),
		Backend:              strings.ToLower(stringFromEnvOrDefault("BACKEND", "ap")),
		GRPCAddress:          stringFromEnvOrDefault("GRPC_ADDRESS", defaultRegistryGRPCAddress),
		ObservabilityAddress: stringFromEnvOrDefault("OBSERVABILITY_ADDRESS", defaultRegistryObservabilityAddress),
		ClusterMembers:       clusterMembers,
		TTLMultiplier:        ttlMultiplier,
		RaftBindAddress:      stringFromEnvOrDefault("RAFT_BIND_ADDRESS", defaultRegistryRaftAddress),
		RaftDataDir:          stringFromEnvOrDefault("RAFT_DATA_DIR", defaultRaftDataDir),
		ForwardTimeout:       forwardTimeout,
		ForwardMaxAttempts:   forwardMaxAttempts,
		Log: logging.Config{
			Level:  stringFromEnvOrDefault("LOG_LEVEL", "info"),
			Format: stringFromEnvOrDefault("LOG_FORMAT", "json"),
		},
	}

	if cfg.NodeID == "" {
		errs = append(errs, errors.New("INSTANCE_ID is required"))
	}
	if cfg.Backend != "cp" && cfg.Backend != "ap" {
		errs = append(errs, fmt.Errorf("BACKEND must be one of cp, ap: %q", cfg.Backend))
	}
	if err == nil && len(cfg.ClusterMembers) == 0 {
		errs = append(errs, errors.New("CLUSTER_MEMBERS is required"))
	}
	if cfg.NodeID != "" && len(cfg.ClusterMembers) > 0 && !containsMember(cfg.ClusterMembers, cfg.NodeID) {
		errs = append(errs, fmt.Errorf("CLUSTER_MEMBERS must include INSTANCE_ID %q", cfg.NodeID))
	}
	if cfg.TTLMultiplier <= 0 {
		errs = append(errs, errors.New("HEARTBEAT_TTL_MULTIPLIER must be greater than zero"))
	}
	if cfg.ForwardTimeout <= 0 {
		errs = append(errs, errors.New("FORWARD_TIMEOUT must be greater than zero"))
	}
	if cfg.ForwardMaxAttempts <= 0 {
		errs = append(errs, errors.New("FORWARD_MAX_ATTEMPTS must be greater than zero"))
	}
	return cfg, errors.Join(errs...)
}

func parseClusterMembers(raw string) ([]registry.ClusterMember, error) {
	var members []registry.ClusterMember

	for _, item := range splitCSV(raw) {
		id, addresses, ok := strings.Cut(item, "=")
		grpcAddress, raftAddress, ok2 := strings.Cut(addresses, "|")

		member := registry.ClusterMember{
			ID:          strings.TrimSpace(id),
			GRPCAddress: strings.TrimSpace(grpcAddress),
			RaftAddress: strings.TrimSpace(raftAddress),
		}

		if !ok || !ok2 || member.ID == "" || member.GRPCAddress == "" || member.RaftAddress == "" {
			return nil, fmt.Errorf("CLUSTER_MEMBERS entry %q must have the form id=grpc_address|raft_address", item)
		}

		if containsMember(members, member.ID) {
			return nil, fmt.Errorf("CLUSTER_MEMBERS entry %q is duplicated", member.ID)
		}

		members = append(members, member)
	}

	return members, nil
}

func containsMember(members []registry.ClusterMember, id string) bool {
	return slices.ContainsFunc(members, func(member registry.ClusterMember) bool {
		return member.ID == id
	})
}
