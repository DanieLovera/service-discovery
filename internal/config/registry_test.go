package config

import (
	"strings"
	"testing"
)

func TestLoadRegistry(t *testing.T) {
	t.Setenv("INSTANCE_ID", "registry-1")
	t.Setenv("BACKEND", "ap")
	t.Setenv("CLUSTER_MEMBERS", "registry-1:7000, registry-2:7000")

	cfg, err := LoadRegistry()
	if err != nil {
		t.Fatalf("LoadRegistry() error = %v, want nil", err)
	}

	if cfg.NodeID != "registry-1" {
		t.Errorf("NodeID = %q, want %q", cfg.NodeID, "registry-1")
	}

	if cfg.ObservabilityAddress != defaultRegistryObservabilityAddress {
		t.Errorf("ObservabilityAddress = %q, want %q", cfg.ObservabilityAddress, defaultRegistryObservabilityAddress)
	}

	if len(cfg.ClusterMembers) != 2 {
		t.Errorf("len(ClusterMembers) = %d, want %d", len(cfg.ClusterMembers), 2)
	}
}

func TestLoadRegistryRejectsInvalidBackend(t *testing.T) {
	t.Setenv("INSTANCE_ID", "registry-1")
	t.Setenv("BACKEND", "invalid_backend")
	t.Setenv("CLUSTER_MEMBERS", "registry-1:7000")

	_, err := LoadRegistry()
	if err == nil {
		t.Fatal("LoadRegistry() error = nil, want backend validation error")
	}

	want := "BACKEND must be one of cp, ap"
	if !strings.Contains(err.Error(), want) {
		t.Errorf("LoadRegistry() error = %q, want it to contain %q", err, want)
	}
}

func TestLoadRegistryRequiresFields(t *testing.T) {
	t.Setenv("INSTANCE_ID", "")
	t.Setenv("CLUSTER_MEMBERS", "")

	_, err := LoadRegistry()
	if err == nil {
		t.Fatal("LoadRegistry() error = nil, want required field errors")
	}

	wantErrors := []string{
		"INSTANCE_ID is required",
		"CLUSTER_MEMBERS is required",
	}

	for _, want := range wantErrors {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("LoadRegistry() error = %q, want it to contain %q", err, want)
		}
	}
}
