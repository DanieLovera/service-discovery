package observability

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestEndpoints(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		ready          bool
		expectedStatus int
	}{
		{
			name:           "health returns OK",
			path:           "/healthz",
			ready:          false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ready returns unavailable by default",
			path:           "/readyz",
			ready:          false,
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			name:           "ready returns OK when ready",
			path:           "/readyz",
			ready:          true,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := New(":0", http.NotFoundHandler(), testLogger())
			server.SetReady(tt.ready)

			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			response := httptest.NewRecorder()

			server.httpServer.Handler.ServeHTTP(response, request)
			if response.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, response.Code)
			}
		})
	}
}

func TestMetricsEndpointUsesProvidedHandler(t *testing.T) {
	metricsHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test_metric 1\n"))
	})

	server := New(":0", metricsHandler, testLogger())

	request := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	response := httptest.NewRecorder()

	server.httpServer.Handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if response.Body.String() != "test_metric 1\n" {
		t.Fatalf("unexpected body: %q", response.Body.String())
	}
}

func TestStartReturnsErrorWhenAddressIsAlreadyInUse(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	server := New(listener.Addr().String(), http.NotFoundHandler(), testLogger())

	err = server.Start()
	if err == nil {
		t.Fatal("expected Start to return an error")
	}
	if !strings.Contains(err.Error(), "serve observability HTTP") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestShutdownStopsRunningServer(t *testing.T) {
	address := freeAddress(t)
	server := New(address, http.NotFoundHandler(), testLogger())

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.Start()
	}()

	waitForServer(t, "http://"+address+"/healthz")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		t.Fatalf("shutdown: %v", err)
	}

	select {
	case err := <-serverErr:
		if err != nil {
			t.Fatalf("Start returned error after shutdown: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Start did not return after shutdown")
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func freeAddress(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	address := listener.Addr().String()

	if err := listener.Close(); err != nil {
		t.Fatalf("close listener: %v", err)
	}
	return address
}

func waitForServer(t *testing.T, url string) {
	t.Helper()

	deadline := time.Now().Add(time.Second)

	for time.Now().Before(deadline) {
		response, err := http.Get(url)
		if err == nil {
			_ = response.Body.Close()
			return
		}

		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("server did not start at %s", url)
}
