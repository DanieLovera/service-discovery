package config

import (
	"os"
	"strings"
)

const (
	defaultRegistryGRPCAddress          = "0.0.0.0:7000"
	defaultLBHTTPAddress                = "0.0.0.0:8000"
	defaultMockHTTPAddress              = "0.0.0.0:9000"
	defaultRegistryObservabilityAddress = "0.0.0.0:10000"
	defaultLBObservabilityAddress       = "0.0.0.0:10010"
	defaultMockObservabilityAddress     = "0.0.0.0:10020"
)

func stringFromEnvOrDefault(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func splitCSV(value string) []string {
	raw := strings.Split(value, ",")
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		if item = strings.TrimSpace(item); item != "" {
			out = append(out, item)
		}
	}
	return out
}
