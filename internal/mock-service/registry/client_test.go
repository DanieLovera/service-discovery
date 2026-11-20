package registry

import (
	"context"
	"errors"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	registrypb "tpiii.local/daniel-tpiii/gen/registry"
)

type fakeRegistry struct {
	registrypb.UnimplementedRegistryServer

	code  codes.Code
	delay time.Duration
	calls atomic.Int32
}

func (f *fakeRegistry) Register(
	ctx context.Context,
	_ *registrypb.RegisterRequest,
) (*registrypb.RegisterResponse, error) {
	f.calls.Add(1)

	select {
	case <-time.After(f.delay):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if f.code != codes.OK {
		return nil, status.Error(f.code, "fake registry error")
	}

	return &registrypb.RegisterResponse{}, nil
}

func startFakeRegistry(t *testing.T, fake *fakeRegistry) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	server := grpc.NewServer()
	registrypb.RegisterRegistryServer(server, fake)

	go func() {
		_ = server.Serve(listener)
	}()
	t.Cleanup(server.Stop)

	return listener.Addr().String()
}

func unusedAddress(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	address := listener.Addr().String()
	_ = listener.Close()

	return address
}

func newTestClient(t *testing.T, requestTimeout time.Duration, addresses ...string) *Client {
	t.Helper()

	client, err := NewClient(addresses, requestTimeout, NewClientMetrics(prometheus.NewRegistry()))
	if err != nil {
		t.Fatalf("NewClient() error = %v, want nil", err)
	}
	t.Cleanup(func() {
		_ = client.Close()
	})

	return client
}

var testInstance = Instance{
	ServiceName:       "service-a",
	InstanceID:        "a1",
	Address:           "a1:8080",
	Weight:            1,
	HeartbeatInterval: time.Second,
}

func TestClientFailsOverAndStaysOnHealthyNode(t *testing.T) {
	unavailable := &fakeRegistry{code: codes.Unavailable}
	healthy := &fakeRegistry{}

	client := newTestClient(t, time.Second, startFakeRegistry(t, unavailable), startFakeRegistry(t, healthy))

	for range 2 {
		if err := client.Register(t.Context(), testInstance); err != nil {
			t.Fatalf("Register() error = %v, want nil", err)
		}
	}

	if calls := unavailable.calls.Load(); calls != 1 {
		t.Errorf("calls to the unavailable node = %d, want 1 (the client stays on the node that answered)", calls)
	}

	if calls := healthy.calls.Load(); calls != 2 {
		t.Errorf("calls to the healthy node = %d, want 2", calls)
	}
}

func TestClientFailsOverWhenNodeIsDown(t *testing.T) {
	healthy := &fakeRegistry{}

	client := newTestClient(t, time.Second, unusedAddress(t), startFakeRegistry(t, healthy))

	if err := client.Register(t.Context(), testInstance); err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}

	if calls := healthy.calls.Load(); calls != 1 {
		t.Errorf("calls to the healthy node = %d, want 1", calls)
	}
}

func TestClientFailsOverWhenNodeTimesOut(t *testing.T) {
	slow := &fakeRegistry{delay: time.Second}
	healthy := &fakeRegistry{}

	client := newTestClient(t, 100*time.Millisecond, startFakeRegistry(t, slow), startFakeRegistry(t, healthy))

	if err := client.Register(t.Context(), testInstance); err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}

	if calls := healthy.calls.Load(); calls != 1 {
		t.Errorf("calls to the healthy node = %d, want 1 after the per-attempt timeout", calls)
	}
}

func TestClientReturnsFunctionalErrorsWithoutFailover(t *testing.T) {
	for _, code := range []codes.Code{codes.NotFound, codes.AlreadyExists, codes.FailedPrecondition} {
		t.Run(code.String(), func(t *testing.T) {
			failing := &fakeRegistry{code: code}
			healthy := &fakeRegistry{}

			client := newTestClient(t, time.Second, startFakeRegistry(t, failing), startFakeRegistry(t, healthy))

			if err := client.Register(t.Context(), testInstance); status.Code(err) != code {
				t.Errorf("Register() error = %v, want code %s", err, code)
			}

			if calls := healthy.calls.Load(); calls != 0 {
				t.Errorf("calls to the second node = %d, want 0", calls)
			}
		})
	}
}

func TestClientReportsUnavailableWhenEveryNodeFails(t *testing.T) {
	first := &fakeRegistry{code: codes.Unavailable}
	second := &fakeRegistry{code: codes.Unavailable}

	client := newTestClient(t, time.Second, startFakeRegistry(t, first), startFakeRegistry(t, second))

	err := client.Register(t.Context(), testInstance)
	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("Register() error = %v, want %v", err, ErrUnavailable)
	}

	if first.calls.Load() != 1 || second.calls.Load() != 1 {
		t.Errorf("calls = %d and %d, want one attempt per node", first.calls.Load(), second.calls.Load())
	}
}
