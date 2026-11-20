package discovery

import (
	"context"
	"log/slog"
	"maps"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	registrypb "tpiii.local/daniel-tpiii/gen/registry"
	loadbalancer "tpiii.local/daniel-tpiii/internal/load-balancer"
	"tpiii.local/daniel-tpiii/internal/load-balancer/registry"
)

type fakeRegistry struct {
	registrypb.UnimplementedRegistryServer

	responses chan *registrypb.WatchResponse
}

func (f *fakeRegistry) Watch(_ *registrypb.WatchRequest, stream registrypb.Registry_WatchServer) error {
	for {
		select {
		case response := <-f.responses:
			if response == nil {
				return status.Error(codes.Unavailable, "stream closed by the test")
			}

			if err := stream.Send(response); err != nil {
				return err
			}

		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}

func startManager(t *testing.T) (*fakeRegistry, *loadbalancer.ServicePool) {
	t.Helper()

	fake := &fakeRegistry{responses: make(chan *registrypb.WatchResponse, 10)}

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

	pool := loadbalancer.NewServicePool()
	manager := NewManager(
		[]string{listener.Addr().String()},
		map[string]*loadbalancer.ServicePool{"service-a": pool},
		50*time.Millisecond,
		registry.KeepaliveConfig{Time: 10 * time.Second, Timeout: 5 * time.Second},
		slog.New(slog.DiscardHandler),
	)

	if err := manager.Start(context.Background()); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = manager.Shutdown(ctx)
	})

	return fake, pool
}

func instance(id string, weight int32, instanceStatus registrypb.ServiceInstanceStatus) *registrypb.ServiceInstance {
	return &registrypb.ServiceInstance{
		ServiceName:         "service-a",
		InstanceId:          id,
		Address:             id + ":8080",
		Weight:              weight,
		HeartbeatIntervalMs: 1000,
		Status:              instanceStatus,
	}
}

func healthy(id string, weight int32) *registrypb.ServiceInstance {
	return instance(id, weight, registrypb.ServiceInstanceStatus_SERVICE_INSTANCE_STATUS_HEALTHY)
}

func snapshot(instances ...*registrypb.ServiceInstance) *registrypb.WatchResponse {
	return &registrypb.WatchResponse{
		Payload: &registrypb.WatchResponse_Snapshot{Snapshot: &registrypb.WatchSnapshot{Instances: instances}},
	}
}

func event(eventType registrypb.StateChangeType, instance *registrypb.ServiceInstance) *registrypb.WatchResponse {
	return &registrypb.WatchResponse{
		Payload: &registrypb.WatchResponse_Event{Event: &registrypb.WatchEvent{Type: eventType, Instance: instance}},
	}
}

func poolWeights(pool *loadbalancer.ServicePool) map[string]int {
	weights := make(map[string]int)
	for _, instance := range pool.Instances() {
		weights[instance.ID] = instance.Weight
	}

	return weights
}

func waitForPool(t *testing.T, pool *loadbalancer.ServicePool, want map[string]int) {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	for !maps.Equal(poolWeights(pool), want) {
		if time.Now().After(deadline) {
			t.Fatalf("pool = %v, want %v", poolWeights(pool), want)
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func TestManagerAppliesWatchEventsToPool(t *testing.T) {
	fake, pool := startManager(t)

	steps := []struct {
		name     string
		response *registrypb.WatchResponse
		want     map[string]int
	}{
		{
			name:     "initial snapshot",
			response: snapshot(healthy("a1", 1), healthy("a2", 1)),
			want:     map[string]int{"a1": 1, "a2": 1},
		},
		{
			name: "expired instance is removed",
			response: event(registrypb.StateChangeType_STATE_CHANGE_TYPE_EXPIRED,
				instance("a1", 1, registrypb.ServiceInstanceStatus_SERVICE_INSTANCE_STATUS_EXPIRED)),
			want: map[string]int{"a2": 1},
		},
		{
			name:     "recovered instance is added back",
			response: event(registrypb.StateChangeType_STATE_CHANGE_TYPE_HEALTHY, healthy("a1", 1)),
			want:     map[string]int{"a1": 1, "a2": 1},
		},
		{
			name:     "updated instance changes its weight",
			response: event(registrypb.StateChangeType_STATE_CHANGE_TYPE_UPDATED, healthy("a2", 5)),
			want:     map[string]int{"a1": 1, "a2": 5},
		},
		{
			name: "deleted instance is removed",
			response: event(registrypb.StateChangeType_STATE_CHANGE_TYPE_DELETED,
				instance("a1", 1, registrypb.ServiceInstanceStatus_SERVICE_INSTANCE_STATUS_DELETED)),
			want: map[string]int{"a2": 5},
		},
		{
			name:     "snapshot in the middle of the stream replaces the pool",
			response: snapshot(healthy("a3", 2)),
			want:     map[string]int{"a3": 2},
		},
	}

	for _, step := range steps {
		fake.responses <- step.response
		waitForPool(t, pool, step.want)
	}
}

func TestManagerKeepsPoolWhileDisconnectedAndReplacesItOnReconnect(t *testing.T) {
	fake, pool := startManager(t)

	fake.responses <- snapshot(healthy("a1", 1), healthy("a2", 1))
	waitForPool(t, pool, map[string]int{"a1": 1, "a2": 1})

	fake.responses <- nil
	time.Sleep(200 * time.Millisecond)

	if got := poolWeights(pool); !maps.Equal(got, map[string]int{"a1": 1, "a2": 1}) {
		t.Fatalf("pool while disconnected = %v, want the last known instances", got)
	}

	fake.responses <- snapshot(healthy("a2", 1))
	waitForPool(t, pool, map[string]int{"a2": 1})
}
