package registry

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"
)

const (
	testTTLMultiplier     = 3
	shortHeartbeat        = 20 * time.Millisecond
	longHeartbeat         = time.Second
	conditionTimeout      = 2 * time.Second
	conditionPollInterval = 5 * time.Millisecond
)

func startHealthManager(t *testing.T, backend *fakeBackend) *HealthManager {
	t.Helper()

	manager := NewHealthManager(backend, testTTLMultiplier, slog.New(slog.DiscardHandler))
	manager.Start(t.Context())

	return manager
}

func waitUntil(t *testing.T, description string, condition func() bool) {
	t.Helper()

	deadline := time.Now().Add(conditionTimeout)
	for !condition() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out after %s waiting until %s", conditionTimeout, description)
		}

		time.Sleep(conditionPollInterval)
	}
}

func TestHeartbeatRejectsUnknownAndDeletedInstances(t *testing.T) {
	backend := newFakeBackend()
	manager := startHealthManager(t, backend)

	unknown := ServiceInstanceID{ServiceName: "service-a", InstanceID: "unknown"}
	if err := manager.Heartbeat(t.Context(), unknown); !errors.Is(err, ErrServiceInstanceNotFound) {
		t.Errorf("Heartbeat(unknown) error = %v, want %v", err, ErrServiceInstanceNotFound)
	}

	id := registerInstance(backend, "service-a", "a1", longHeartbeat)
	_ = backend.Deregister(t.Context(), id)

	if err := manager.Heartbeat(t.Context(), id); !errors.Is(err, ErrServiceInstanceDeleted) {
		t.Errorf("Heartbeat(deleted) error = %v, want %v", err, ErrServiceInstanceDeleted)
	}
}

func TestInstanceWithoutHeartbeatsIsNotExpired(t *testing.T) {
	backend := newFakeBackend()
	startHealthManager(t, backend)

	id := registerInstance(backend, "service-a", "a1", shortHeartbeat)

	time.Sleep(5 * testTTLMultiplier * shortHeartbeat)

	if status := backend.status(id); status != ServiceInstanceStatusHealthy {
		t.Errorf("status = %s, want %s (the TTL starts with the first heartbeat)", status, ServiceInstanceStatusHealthy)
	}
}

func TestHeartbeatsKeepInstanceHealthy(t *testing.T) {
	backend := newFakeBackend()
	manager := startHealthManager(t, backend)

	heartbeatInterval := 50 * time.Millisecond
	id := registerInstance(backend, "service-a", "a1", heartbeatInterval)

	for range 20 {
		if err := manager.Heartbeat(t.Context(), id); err != nil {
			t.Fatalf("Heartbeat() error = %v, want nil", err)
		}

		time.Sleep(heartbeatInterval / 2)
	}

	if expires, _ := backend.calls(); expires != 0 {
		t.Errorf("Expire() calls = %d, want 0 while heartbeats arrive", expires)
	}
}

func TestMissedHeartbeatsExpireInstanceAfterTTL(t *testing.T) {
	backend := newFakeBackend()
	manager := startHealthManager(t, backend)

	id := registerInstance(backend, "service-a", "a1", shortHeartbeat)

	start := time.Now()
	if err := manager.Heartbeat(t.Context(), id); err != nil {
		t.Fatalf("Heartbeat() error = %v, want nil", err)
	}

	waitUntil(t, "the instance expires", func() bool {
		return backend.status(id) == ServiceInstanceStatusExpired
	})

	ttl := testTTLMultiplier * shortHeartbeat
	if elapsed := time.Since(start); elapsed < ttl {
		t.Errorf("instance expired after %s, want at least the TTL %s", elapsed, ttl)
	}
}

func TestHeartbeatRecoversExpiredInstance(t *testing.T) {
	backend := newFakeBackend()
	manager := startHealthManager(t, backend)

	id := registerInstance(backend, "service-a", "a1", shortHeartbeat)
	_ = manager.Heartbeat(t.Context(), id)

	waitUntil(t, "the instance expires", func() bool {
		return backend.status(id) == ServiceInstanceStatusExpired
	})

	if err := manager.Heartbeat(t.Context(), id); err != nil {
		t.Fatalf("Heartbeat(expired) error = %v, want nil", err)
	}

	if status := backend.status(id); status != ServiceInstanceStatusHealthy {
		t.Errorf("status after heartbeat = %s, want %s", status, ServiceInstanceStatusHealthy)
	}

	waitUntil(t, "the recovered instance expires again without heartbeats", func() bool {
		return backend.status(id) == ServiceInstanceStatusExpired
	})
}

func TestDeregisterCancelsPendingExpiration(t *testing.T) {
	backend := newFakeBackend()
	manager := startHealthManager(t, backend)

	id := registerInstance(backend, "service-a", "a1", shortHeartbeat)
	_ = manager.Heartbeat(t.Context(), id)
	_ = backend.Deregister(t.Context(), id)

	time.Sleep(3 * testTTLMultiplier * shortHeartbeat)

	if expires, _ := backend.calls(); expires != 0 {
		t.Errorf("Expire() calls = %d, want 0 after Deregister", expires)
	}
}

func TestClusterExpirationIsRecoveredWhileLocalTTLIsValid(t *testing.T) {
	backend := newFakeBackend()
	manager := startHealthManager(t, backend)

	id := registerInstance(backend, "service-a", "a1", longHeartbeat)
	_ = manager.Heartbeat(t.Context(), id)

	_ = backend.Expire(context.Background(), id)

	waitUntil(t, "the node holding a valid TTL recovers the instance", func() bool {
		return backend.status(id) == ServiceInstanceStatusHealthy
	})
}

func TestClusterExpirationIsKeptWithoutLocalTTL(t *testing.T) {
	backend := newFakeBackend()
	startHealthManager(t, backend)

	id := registerInstance(backend, "service-a", "a1", longHeartbeat)
	_ = backend.Expire(context.Background(), id)

	time.Sleep(100 * time.Millisecond)

	if _, recovers := backend.calls(); recovers != 0 {
		t.Errorf("Recover() calls = %d, want 0 when this node holds no TTL", recovers)
	}
}

func TestRestoreAppliesExpiredAndDeletedInstances(t *testing.T) {
	backend := newFakeBackend()
	manager := startHealthManager(t, backend)

	expired := registerInstance(backend, "service-a", "a1", longHeartbeat)
	deleted := registerInstance(backend, "service-a", "a2", longHeartbeat)
	_ = manager.Heartbeat(t.Context(), expired)
	_ = manager.Heartbeat(t.Context(), deleted)

	backend.restore([]ServiceInstance{
		{ID: expired, HeartbeatInterval: longHeartbeat, Status: ServiceInstanceStatusExpired},
		{ID: deleted, HeartbeatInterval: longHeartbeat, Status: ServiceInstanceStatusDeleted},
	})

	waitUntil(t, "the restored expired instance with a valid TTL is recovered", func() bool {
		return backend.status(expired) == ServiceInstanceStatusHealthy
	})

	if manager.hasValidTTL(deleted) {
		t.Error("the restored deleted instance still has a TTL")
	}
}
