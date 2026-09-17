package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadMockService(t *testing.T) {
	t.Setenv("MS_NAME", "payments")
	t.Setenv("MS_INSTANCE_ID", "payments-1")
	t.Setenv("MS_HTTP_ADDRESS", "0.0.0.0:9100")
	t.Setenv("MS_REGISTRY_ADDRESSES", "registry-1:7000, registry-2:7000")
	t.Setenv("MS_WEIGHT", "3")
	t.Setenv("MS_HEARTBEAT_INTERVAL", "2s")
	t.Setenv("MS_LOG_LEVEL", "debug")
	t.Setenv("MS_LOG_FORMAT", "text")

	cfg, err := LoadMockService()
	if err != nil {
		t.Fatalf("LoadMockService() error = %v", err)
	}

	if cfg.ServiceName != "payments" {
		t.Errorf("ServiceName = %q, want %q", cfg.ServiceName, "payments")
	}

	if cfg.InstanceID != "payments-1" {
		t.Errorf("InstanceID = %q, want %q", cfg.InstanceID, "payments-1")
	}

	if cfg.HTTPAddress != "0.0.0.0:9100" {
		t.Errorf("HTTPAddress = %q, want %q", cfg.HTTPAddress, "0.0.0.0:9100")
	}

	wantRegistryAddresses := []string{
		"registry-1:7000",
		"registry-2:7000",
	}
	if len(cfg.RegistryAddresses) != len(wantRegistryAddresses) {
		t.Fatalf("len(RegistryAddresses) = %d, want %d", len(cfg.RegistryAddresses), len(wantRegistryAddresses))
	}

	for i, want := range wantRegistryAddresses {
		if cfg.RegistryAddresses[i] != want {
			t.Errorf("RegistryAddresses[%d] = %q, want %q", i, cfg.RegistryAddresses[i], want)
		}
	}

	if cfg.Weight != 3 {
		t.Errorf("Weight = %d, want %d", cfg.Weight, 3)
	}

	if cfg.HeartbeatInterval != 2*time.Second {
		t.Errorf("HeartbeatInterval = %s, want %s", cfg.HeartbeatInterval, 2*time.Second)
	}

	if cfg.Log.Level != "debug" {
		t.Errorf("Log.Level = %q, want %q", cfg.Log.Level, "debug")
	}

	if cfg.Log.Format != "text" {
		t.Errorf("Log.Format = %q, want %q", cfg.Log.Format, "text")
	}
}

func TestLoadMockServiceUsesDefaults(t *testing.T) {
	t.Setenv("MS_NAME", "payments")
	t.Setenv("MS_INSTANCE_ID", "payments-1")
	t.Setenv("MS_REGISTRY_ADDRESSES", "registry-1:7000")

	cfg, err := LoadMockService()
	if err != nil {
		t.Fatalf("LoadMockService() error = %v", err)
	}

	if cfg.HTTPAddress != defaultMockHTTPAddress {
		t.Errorf("HTTPAddress = %q, want %q", cfg.HTTPAddress, defaultMockHTTPAddress)
	}

	if cfg.ObservabilityAddress != defaultMockObservabilityAddress {
		t.Errorf("ObservabilityAddress = %q, want %q", cfg.ObservabilityAddress, defaultMockObservabilityAddress)
	}

	if cfg.Weight != defaultWeight {
		t.Errorf("Weight = %d, want %d", cfg.Weight, defaultWeight)
	}

	if cfg.HeartbeatInterval != defaultHeartbeatInterval {
		t.Errorf("HeartbeatInterval = %s, want %s", cfg.HeartbeatInterval, defaultHeartbeatInterval)
	}

	if cfg.Log.Level != "info" {
		t.Errorf("Log.Level = %q, want %q", cfg.Log.Level, "info")
	}

	if cfg.Log.Format != "json" {
		t.Errorf("Log.Format = %q, want %q", cfg.Log.Format, "json")
	}
}

func TestLoadMockServiceRequiresFields(t *testing.T) {
	tests := []struct {
		name      string
		envKey    string
		wantError string
	}{
		{
			name:      "missing service name",
			envKey:    "MS_NAME",
			wantError: "MS_NAME is required",
		},
		{
			name:      "missing instance id",
			envKey:    "MS_INSTANCE_ID",
			wantError: "MS_INSTANCE_ID is required",
		},
		{
			name:      "missing registry addresses",
			envKey:    "MS_REGISTRY_ADDRESSES",
			wantError: "MS_REGISTRY_ADDRESSES is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("MS_NAME", "payments")
			t.Setenv("MS_INSTANCE_ID", "payments-1")
			t.Setenv("MS_REGISTRY_ADDRESSES", "registry-1:7000")
			t.Setenv(tt.envKey, "")

			_, err := LoadMockService()
			if err == nil {
				t.Fatal("LoadMockService() error = nil, want error")
			}

			if !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("LoadMockService() error = %q, want it to contain %q", err, tt.wantError)
			}
		})
	}
}

func TestLoadMockServiceRejectsNonPositiveValues(t *testing.T) {
	tests := []struct {
		name      string
		envKey    string
		envValue  string
		wantError string
	}{
		{
			name:      "zero weight",
			envKey:    "MS_WEIGHT",
			envValue:  "0",
			wantError: "MS_WEIGHT must be greater than zero",
		},
		{
			name:      "negative weight",
			envKey:    "MS_WEIGHT",
			envValue:  "-1",
			wantError: "MS_WEIGHT must be greater than zero",
		},
		{
			name:      "zero heartbeat interval",
			envKey:    "MS_HEARTBEAT_INTERVAL",
			envValue:  "0s",
			wantError: "MS_HEARTBEAT_INTERVAL must be greater than zero",
		},
		{
			name:      "negative heartbeat interval",
			envKey:    "MS_HEARTBEAT_INTERVAL",
			envValue:  "-1s",
			wantError: "MS_HEARTBEAT_INTERVAL must be greater than zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("MS_NAME", "payments")
			t.Setenv("MS_INSTANCE_ID", "payments-1")
			t.Setenv("MS_REGISTRY_ADDRESSES", "registry-1:7000")
			t.Setenv(tt.envKey, tt.envValue)

			_, err := LoadMockService()
			if err == nil {
				t.Fatal("LoadMockService() error = nil, want error")
			}

			if !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("LoadMockService() error = %q, want it to contain %q", err, tt.wantError)
			}
		})
	}
}

func TestLoadMockServiceRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name      string
		envKey    string
		envValue  string
		wantError string
	}{
		{
			name:      "invalid weight",
			envKey:    "MS_WEIGHT",
			envValue:  "invalid",
			wantError: "MS_WEIGHT must be an integer",
		},
		{
			name:      "invalid heartbeat interval",
			envKey:    "MS_HEARTBEAT_INTERVAL",
			envValue:  "invalid",
			wantError: "MS_HEARTBEAT_INTERVAL must be a Go duration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("MS_NAME", "payments")
			t.Setenv("MS_INSTANCE_ID", "payments-1")
			t.Setenv("MS_REGISTRY_ADDRESSES", "registry-1:7000")
			t.Setenv(tt.envKey, tt.envValue)

			_, err := LoadMockService()
			if err == nil {
				t.Fatal("LoadMockService() error = nil, want error")
			}

			if !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("LoadMockService() error = %q, want it to contain %q", err, tt.wantError)
			}
		})
	}
}
