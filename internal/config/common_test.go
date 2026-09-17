package config

import (
	"reflect"
	"testing"
)

func TestStringFromEnvOrDefault(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		fallback string
		want     string
	}{
		{
			name:     "returns environment value",
			value:    "debug",
			fallback: "info",
			want:     "debug",
		},
		{
			name:     "returns fallback for empty value",
			value:    "",
			fallback: "info",
			want:     "info",
		},
		{
			name:     "returns fallback for whitespace-only value",
			value:    "   ",
			fallback: "info",
			want:     "info",
		},
		{
			name:     "trims environment value",
			value:    "  debug  ",
			fallback: "info",
			want:     "debug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TEST_CONFIG_VALUE", tt.value)

			got := stringFromEnvOrDefault("TEST_CONFIG_VALUE", tt.fallback)
			if got != tt.want {
				t.Fatalf("stringFromEnvOrDefault() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSplitCSV(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  []string
	}{
		{
			name:  "single value",
			value: "registry-1:7000",
			want:  []string{"registry-1:7000"},
		},
		{
			name:  "multiple values",
			value: "registry-1:7000,registry-2:7000",
			want:  []string{"registry-1:7000", "registry-2:7000"},
		},
		{
			name:  "trims whitespace",
			value: " registry-1:7000 , registry-2:7000 ",
			want:  []string{"registry-1:7000", "registry-2:7000"},
		},
		{
			name:  "ignores empty entries",
			value: "registry-1:7000,, ,registry-2:7000,",
			want:  []string{"registry-1:7000", "registry-2:7000"},
		},
		{
			name:  "empty value",
			value: "",
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitCSV(tt.value)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("splitCSV(%q) = %#v, want %#v", tt.value, got, tt.want)
			}
		})
	}
}
