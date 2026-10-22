package config

import (
	"strings"
	"testing"
)

func TestLoadLoadBalancer(t *testing.T) {
	t.Setenv("REGISTRY_ADDRESSES", "registry-1:7000,registry-2:7000")
	t.Setenv("SERVICES", "mock-service:round-robin")

	cfg, err := LoadLoadBalancer()
	if err != nil {
		t.Fatalf("LoadLoadBalancer() error = %v, want nil", err)
	}

	if cfg.HTTPBindAddress != defaultLBHTTPBindAddress {
		t.Errorf("HTTPBindAddress = %q, want %q", cfg.HTTPBindAddress, defaultLBHTTPBindAddress)
	}

	if cfg.ObservabilityBindAddress != defaultLBObservabilityBindAddress {
		t.Errorf("ObservabilityBindAddress = %q, want %q", cfg.ObservabilityBindAddress, defaultLBObservabilityBindAddress)
	}

	if len(cfg.RegistryAddresses) != 2 {
		t.Errorf("len(RegistryAddresses) = %d, want %d", len(cfg.RegistryAddresses), 2)
	}
}

func TestLoadLoadBalancerRequiresRegistryAddresses(t *testing.T) {
	t.Setenv("REGISTRY_ADDRESSES", "")

	_, err := LoadLoadBalancer()
	if err == nil {
		t.Fatal("LoadLoadBalancer() error = nil, want registry addresses validation error")
	}

	want := "REGISTRY_ADDRESSES is required"
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("LoadLoadBalancer() error = %q, want it to contain %q", err, want)
	}
}
