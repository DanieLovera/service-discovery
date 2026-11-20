package http

import (
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/prometheus/client_golang/prometheus"

	loadbalancer "tpiii.local/daniel-tpiii/internal/load-balancer"
	"tpiii.local/daniel-tpiii/internal/load-balancer/routing"
)

type countingStrategy struct {
	picks    atomic.Int32
	releases atomic.Int32
}

func (s *countingStrategy) Pick(instances []loadbalancer.ServiceInstance) (loadbalancer.ServiceInstance, func(), bool) {
	if len(instances) == 0 {
		return loadbalancer.ServiceInstance{}, func() {}, false
	}

	s.picks.Add(1)

	return instances[0], func() { s.releases.Add(1) }, true
}

func newTestServer(strategy routing.Strategy, instances ...loadbalancer.ServiceInstance) *Server {
	pool := loadbalancer.NewServicePool()
	pool.Replace(instances)

	router := routing.NewManager(map[string]routing.Upstream{
		"service-a": {Pool: pool, Strategy: strategy},
	})

	return NewServer("127.0.0.1:0", router, NewServerMetrics(prometheus.NewRegistry()), slog.New(slog.DiscardHandler))
}

func serve(server *Server, target string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, target, nil))

	return recorder
}

func TestServerProxiesToInstanceWithoutServicePrefix(t *testing.T) {
	received := make(chan *http.Request, 1)
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received <- r
		_, _ = io.WriteString(w, "pong")
	}))
	defer backend.Close()

	strategy := &countingStrategy{}
	server := newTestServer(strategy, loadbalancer.ServiceInstance{ID: "a1", Address: backend.Listener.Addr().String()})

	response := serve(server, "/service-a/api/items?page=2")

	if response.Code != http.StatusOK || response.Body.String() != "pong" {
		t.Fatalf("response = %d %q, want %d %q", response.Code, response.Body.String(), http.StatusOK, "pong")
	}

	request := <-received
	if request.URL.Path != "/api/items" || request.URL.RawQuery != "page=2" {
		t.Errorf("instance received %s?%s, want /api/items?page=2", request.URL.Path, request.URL.RawQuery)
	}

	if host := request.Header.Get("X-Forwarded-Host"); host != "example.com" {
		t.Errorf("X-Forwarded-Host = %q, want %q", host, "example.com")
	}

	if releases := strategy.releases.Load(); releases != 1 {
		t.Errorf("strategy releases = %d, want 1 after the request finished", releases)
	}
}

func TestServerReturnsBadGatewayWhenInstanceIsDown(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	downAddress := listener.Addr().String()
	_ = listener.Close()

	strategy := &countingStrategy{}
	server := newTestServer(strategy, loadbalancer.ServiceInstance{ID: "a1", Address: downAddress})

	if response := serve(server, "/service-a/ping"); response.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", response.Code, http.StatusBadGateway)
	}

	if releases := strategy.releases.Load(); releases != 1 {
		t.Errorf("strategy releases = %d, want 1 after a failed request", releases)
	}
}

func TestServerReturnsServiceUnavailableWithoutInstances(t *testing.T) {
	strategy := &countingStrategy{}
	server := newTestServer(strategy)

	if response := serve(server, "/service-a/ping"); response.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want %d", response.Code, http.StatusServiceUnavailable)
	}

	if picks := strategy.picks.Load(); picks != 0 {
		t.Errorf("strategy picks = %d, want 0", picks)
	}
}
