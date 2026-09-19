package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"tpiii.local/daniel-tpiii/internal/logging"
)

const (
	defaultHeartbeatInterval = 5 * time.Second
	defaultWeight            = 1
)

type MockService struct {
	ServiceName          string
	InstanceID           string
	HTTPAddress          string
	ObservabilityAddress string
	RegistryAddresses    []string
	Weight               int
	HeartbeatInterval    time.Duration
	Log                  logging.Config
}

func LoadMockService() (MockService, error) {
	weight, err := intFromEnvOrDefault("WEIGHT", defaultWeight)
	if err != nil {
		return MockService{}, err
	}
	interval, err := durationFromEnvOrDefault("HEARTBEAT_INTERVAL", defaultHeartbeatInterval)
	if err != nil {
		return MockService{}, err
	}

	cfg := MockService{
		ServiceName:          strings.TrimSpace(os.Getenv("NAME")),
		InstanceID:           strings.TrimSpace(os.Getenv("INSTANCE_ID")),
		HTTPAddress:          stringFromEnvOrDefault("HTTP_ADDRESS", defaultMockHTTPAddress),
		ObservabilityAddress: stringFromEnvOrDefault("OBSERVABILITY_ADDRESS", defaultMockObservabilityAddress),
		RegistryAddresses:    splitCSV(os.Getenv("REGISTRY_ADDRESSES")),
		Weight:               weight,
		HeartbeatInterval:    interval,
		Log: logging.Config{
			Level:  stringFromEnvOrDefault("LOG_LEVEL", "info"),
			Format: stringFromEnvOrDefault("LOG_FORMAT", "json"),
		},
	}
	var errs []error
	if cfg.ServiceName == "" {
		errs = append(errs, errors.New("NAME is required"))
	}
	if cfg.InstanceID == "" {
		errs = append(errs, errors.New("INSTANCE_ID is required"))
	}
	if len(cfg.RegistryAddresses) == 0 {
		errs = append(errs, errors.New("REGISTRY_ADDRESSES is required"))
	}
	if cfg.Weight <= 0 {
		errs = append(errs, errors.New("WEIGHT must be greater than zero"))
	}

	if cfg.HeartbeatInterval <= 0 {
		errs = append(errs, errors.New("HEARTBEAT_INTERVAL must be greater than zero"))
	}
	return cfg, errors.Join(errs...)
}

func intFromEnvOrDefault(key string, fallback int) (int, error) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", key, err)
	}
	return v, nil
}

func durationFromEnvOrDefault(key string, fallback time.Duration) (time.Duration, error) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback, nil
	}
	v, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be a Go duration: %w", key, err)
	}
	return v, nil
}
