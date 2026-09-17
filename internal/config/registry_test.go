package config

import (
	"strings"
	"testing"
)

func TestLoadRegistry(t *testing.T) {
	t.Setenv("REGISTRY_INSTANCE_ID", "registry-1")
	t.Setenv("REGISTRY_BACKEND", "ap")
	t.Setenv("REGISTRY_CLUSTER_MEMBERS", "registry-1:7000, registry-2:7000")

	cfg, err := LoadRegistry()
	if err != nil {
		t.Fatalf("LoadRegistry() error = %v, want nil", err)
	}

	if cfg.NodeID != "registry-1" {
		t.Errorf("NodeID = %q, want %q", cfg.NodeID, "registry-1")
	}

	if len(cfg.ClusterMembers) != 2 {
		t.Errorf("len(ClusterMembers) = %d, want %d", len(cfg.ClusterMembers), 2)
	}
}

func TestLoadRegistryRejectsInvalidBackend(t *testing.T) {
	t.Setenv("REGISTRY_INSTANCE_ID", "registry-1")
	t.Setenv("REGISTRY_BACKEND", "invalid_backend")
	t.Setenv("REGISTRY_CLUSTER_MEMBERS", "registry-1:7000")

	_, err := LoadRegistry()
	if err == nil {
		t.Fatal("LoadRegistry() error = nil, want backend validation error")
	}

	want := "REGISTRY_BACKEND must be one of cp, ap"
	if !strings.Contains(err.Error(), want) {
		t.Errorf("LoadRegistry() error = %q, want it to contain %q", err, want)
	}
}

func TestLoadRegistryRequiresFields(t *testing.T) {
	t.Setenv("REGISTRY_INSTANCE_ID", "")
	t.Setenv("REGISTRY_CLUSTER_MEMBERS", "")

	_, err := LoadRegistry()
	if err == nil {
		t.Fatal("LoadRegistry() error = nil, want required field errors")
	}

	wantErrors := []string{
		"REGISTRY_INSTANCE_ID is required",
		"REGISTRY_CLUSTER_MEMBERS is required",
	}

	for _, want := range wantErrors {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("LoadRegistry() error = %q, want it to contain %q", err, want)
		}
	}
}
