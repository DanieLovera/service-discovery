package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerExposesMetrics(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "Go metrics",
			expected: "go_goroutines",
		},
		{
			name:     "process metrics",
			expected: "process_start_time_seconds",
		},
	}

	request := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	response := httptest.NewRecorder()
	metrics := New("registry", "registry-1")
	metrics.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	body := response.Body.String()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(body, tt.expected) {
				t.Errorf("expected metrics output to contain %q", tt.expected)
			}
		})
	}
}
