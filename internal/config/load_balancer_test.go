package config

import (
	"strings"
	"testing"
)

func TestLoadLoadBalancer(t *testing.T) {
	t.Setenv("LB_REGISTRY_ADDRESSES", "registry-1:7000,registry-2:7000")

	cfg, err := LoadLoadBalancer()
	if err != nil {
		t.Fatalf("LoadLoadBalancer() error = %v, want nil", err)
	}

	if cfg.HTTPAddress != defaultLBHTTPAddress {
		t.Errorf("HTTPAddress = %q, want %q", cfg.HTTPAddress, defaultLBHTTPAddress)
	}

	if cfg.ObservabilityAddress != defaultLBObservabilityAddress {
		t.Errorf("ObservabilityAddress = %q, want %q", cfg.ObservabilityAddress, defaultLBObservabilityAddress)
	}

	if len(cfg.RegistryAddresses) != 2 {
		t.Errorf("len(RegistryAddresses) = %d, want %d", len(cfg.RegistryAddresses), 2)
	}
}

func TestLoadLoadBalancerRequiresRegistryAddresses(t *testing.T) {
	t.Setenv("LB_REGISTRY_ADDRESSES", "")

	_, err := LoadLoadBalancer()
	if err == nil {
		t.Fatal("LoadLoadBalancer() error = nil, want registry addresses validation error")
	}

	want := "LB_REGISTRY_ADDRESSES is required"
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("LoadLoadBalancer() error = %q, want it to contain %q", err, want)
	}
}
