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
	defaultTTLMultiplier            = 3
	defaultRaftDataDir              = "/data/raft"
	defaultForwardTimeout           = 2 * time.Second
	defaultForwardMaxAttempts       = 3
	defaultGossipFanout             = 3
	defaultAntiEntropyInterval      = 5 * time.Second
	defaultSWIMProbeInterval        = time.Second
	defaultSWIMPingTimeout          = 300 * time.Millisecond
	defaultSWIMPingReqTimeout       = 600 * time.Millisecond
	defaultSWIMPingReqCount         = 2
	defaultSWIMSuspectTimeout       = 5 * time.Second
	defaultSWIMRetransmitMultiplier = 4
)

type Registry struct {
	NodeID                   string
	Backend                  string
	GRPCBindAddress          string
	ObservabilityBindAddress string
	ClusterMembers           []registry.ClusterMember
	TTLMultiplier            int
	PeerBindAddress          string
	RaftDataDir              string
	ForwardTimeout           time.Duration
	ForwardMaxAttempts       int
	GossipFanout             int
	AntiEntropyInterval      time.Duration
	SWIMProbeInterval        time.Duration
	SWIMPingTimeout          time.Duration
	SWIMPingReqTimeout       time.Duration
	SWIMPingReqCount         int
	SWIMSuspectTimeout       time.Duration
	SWIMRetransmitMultiplier int
	Log                      logging.Config
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

	gossipFanout, err := intFromEnvOrDefault("GOSSIP_FANOUT", defaultGossipFanout)
	if err != nil {
		return Registry{}, err
	}

	antiEntropyInterval, err := durationFromEnvOrDefault("ANTI_ENTROPY_INTERVAL", defaultAntiEntropyInterval)
	if err != nil {
		return Registry{}, err
	}

	swimProbeInterval, err := durationFromEnvOrDefault("SWIM_PROBE_INTERVAL", defaultSWIMProbeInterval)
	if err != nil {
		return Registry{}, err
	}

	swimPingTimeout, err := durationFromEnvOrDefault("SWIM_PING_TIMEOUT", defaultSWIMPingTimeout)
	if err != nil {
		return Registry{}, err
	}

	swimPingReqTimeout, err := durationFromEnvOrDefault("SWIM_PING_REQ_TIMEOUT", defaultSWIMPingReqTimeout)
	if err != nil {
		return Registry{}, err
	}

	swimPingReqCount, err := intFromEnvOrDefault("SWIM_PING_REQ_COUNT", defaultSWIMPingReqCount)
	if err != nil {
		return Registry{}, err
	}

	swimSuspectTimeout, err := durationFromEnvOrDefault("SWIM_SUSPECT_TIMEOUT", defaultSWIMSuspectTimeout)
	if err != nil {
		return Registry{}, err
	}

	swimRetransmitMultiplier, err := intFromEnvOrDefault("SWIM_RETRANSMIT_MULTIPLIER", defaultSWIMRetransmitMultiplier)
	if err != nil {
		return Registry{}, err
	}

	var errs []error

	clusterMembers, err := parseClusterMembers(os.Getenv("CLUSTER_MEMBERS"))
	if err != nil {
		errs = append(errs, err)
	}

	cfg := Registry{
		NodeID:                   strings.TrimSpace(os.Getenv("INSTANCE_ID")),
		Backend:                  strings.ToLower(stringFromEnvOrDefault("BACKEND", "ap")),
		GRPCBindAddress:          stringFromEnvOrDefault("GRPC_BIND_ADDRESS", defaultRegistryGRPCBindAddress),
		ObservabilityBindAddress: stringFromEnvOrDefault("OBSERVABILITY_BIND_ADDRESS", defaultRegistryObservabilityBindAddress),
		ClusterMembers:           clusterMembers,
		TTLMultiplier:            ttlMultiplier,
		PeerBindAddress:          stringFromEnvOrDefault("PEER_BIND_ADDRESS", defaultRegistryPeerBindAddress),
		RaftDataDir:              stringFromEnvOrDefault("RAFT_DATA_DIR", defaultRaftDataDir),
		ForwardTimeout:           forwardTimeout,
		ForwardMaxAttempts:       forwardMaxAttempts,
		GossipFanout:             gossipFanout,
		AntiEntropyInterval:      antiEntropyInterval,
		SWIMProbeInterval:        swimProbeInterval,
		SWIMPingTimeout:          swimPingTimeout,
		SWIMPingReqTimeout:       swimPingReqTimeout,
		SWIMPingReqCount:         swimPingReqCount,
		SWIMSuspectTimeout:       swimSuspectTimeout,
		SWIMRetransmitMultiplier: swimRetransmitMultiplier,
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
	if cfg.GossipFanout <= 0 {
		errs = append(errs, errors.New("GOSSIP_FANOUT must be greater than zero"))
	}
	if cfg.AntiEntropyInterval <= 0 {
		errs = append(errs, errors.New("ANTI_ENTROPY_INTERVAL must be greater than zero"))
	}
	if cfg.SWIMProbeInterval <= 0 {
		errs = append(errs, errors.New("SWIM_PROBE_INTERVAL must be greater than zero"))
	}
	if cfg.SWIMPingTimeout <= 0 {
		errs = append(errs, errors.New("SWIM_PING_TIMEOUT must be greater than zero"))
	}
	if cfg.SWIMPingReqTimeout <= 0 {
		errs = append(errs, errors.New("SWIM_PING_REQ_TIMEOUT must be greater than zero"))
	}
	if cfg.SWIMPingReqCount <= 0 {
		errs = append(errs, errors.New("SWIM_PING_REQ_COUNT must be greater than zero"))
	}
	if cfg.SWIMSuspectTimeout <= 0 {
		errs = append(errs, errors.New("SWIM_SUSPECT_TIMEOUT must be greater than zero"))
	}
	if cfg.SWIMRetransmitMultiplier <= 0 {
		errs = append(errs, errors.New("SWIM_RETRANSMIT_MULTIPLIER must be greater than zero"))
	}
	if cfg.SWIMPingTimeout+cfg.SWIMPingReqTimeout > cfg.SWIMProbeInterval {
		errs = append(errs, errors.New("SWIM_PING_TIMEOUT plus SWIM_PING_REQ_TIMEOUT must not exceed SWIM_PROBE_INTERVAL"))
	}
	return cfg, errors.Join(errs...)
}

func parseClusterMembers(raw string) ([]registry.ClusterMember, error) {
	var members []registry.ClusterMember

	for _, item := range splitCSV(raw) {
		id, addresses, ok := strings.Cut(item, "=")
		grpcAddress, peerAddress, ok2 := strings.Cut(addresses, "|")

		member := registry.ClusterMember{
			ID:          strings.TrimSpace(id),
			GRPCAddress: strings.TrimSpace(grpcAddress),
			PeerAddress: strings.TrimSpace(peerAddress),
		}

		if !ok || !ok2 || member.ID == "" || member.GRPCAddress == "" || member.PeerAddress == "" {
			return nil, fmt.Errorf("CLUSTER_MEMBERS entry %q must have the form id=grpc_address|peer_address", item)
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
