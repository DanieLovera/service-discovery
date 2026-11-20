package registry

import (
	"fmt"
	"slices"
	"testing"
	"time"
)

func receiveChange(t *testing.T, subscriber *Subscriber) StateChange {
	t.Helper()

	select {
	case change, ok := <-subscriber.Events():
		if !ok {
			t.Fatal("subscriber channel closed, want an event")
		}
		return change
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for an event")
		return StateChange{}
	}
}

func assertNoChange(t *testing.T, subscriber *Subscriber) {
	t.Helper()

	select {
	case change, ok := <-subscriber.Events():
		if !ok {
			t.Fatal("subscriber channel closed, want it open and empty")
		}
		t.Fatalf("received %s for %s, want no event", change.Type, change.Instance.ID.InstanceID)
	default:
	}
}

func instanceIDs(instances []ServiceInstance) []string {
	ids := make([]string, 0, len(instances))
	for _, instance := range instances {
		ids = append(ids, instance.ID.InstanceID)
	}

	slices.Sort(ids)

	return ids
}

func TestWatchStartsWithSnapshotOfHealthyInstances(t *testing.T) {
	backend := newFakeBackend()
	manager := NewWatchManager(backend)

	registerInstance(backend, "service-a", "a1", time.Second)
	expired := registerInstance(backend, "service-a", "a2", time.Second)
	registerInstance(backend, "service-b", "b1", time.Second)
	_ = backend.Expire(t.Context(), expired)

	snapshot, _, cleanup, err := manager.Watch(t.Context(), "service-a")
	if err != nil {
		t.Fatalf("Watch() error = %v, want nil", err)
	}
	defer cleanup()

	if got, want := instanceIDs(snapshot), []string{"a1"}; !slices.Equal(got, want) {
		t.Errorf("snapshot = %v, want %v", got, want)
	}
}

func TestWatchDeliversOnlyChangesOfWatchedService(t *testing.T) {
	backend := newFakeBackend()
	manager := NewWatchManager(backend)

	_, subscriber, cleanup, err := manager.Watch(t.Context(), "service-a")
	if err != nil {
		t.Fatalf("Watch() error = %v, want nil", err)
	}
	defer cleanup()

	registerInstance(backend, "service-b", "b1", time.Second)
	assertNoChange(t, subscriber)

	id := registerInstance(backend, "service-a", "a1", time.Second)
	_ = backend.Expire(t.Context(), id)
	_ = backend.Recover(t.Context(), id)
	_ = backend.Deregister(t.Context(), id)

	want := []StateChangeType{StateChangeRegistered, StateChangeExpired, StateChangeHealthy, StateChangeDeleted}
	for _, wantType := range want {
		change := receiveChange(t, subscriber)
		if change.Type != wantType || change.Instance.ID != id {
			t.Errorf("event = %s for %s, want %s for a1", change.Type, change.Instance.ID.InstanceID, wantType)
		}
	}
}

func TestWatchCleanupStopsDelivery(t *testing.T) {
	backend := newFakeBackend()
	manager := NewWatchManager(backend)

	_, subscriber, cleanup, err := manager.Watch(t.Context(), "service-a")
	if err != nil {
		t.Fatalf("Watch() error = %v, want nil", err)
	}

	cleanup()
	registerInstance(backend, "service-a", "a1", time.Second)

	assertNoChange(t, subscriber)
}

func TestWatchClosesSubscriberThatFallsBehind(t *testing.T) {
	backend := newFakeBackend()
	manager := NewWatchManager(backend)

	_, subscriber, cleanup, err := manager.Watch(t.Context(), "service-a")
	if err != nil {
		t.Fatalf("Watch() error = %v, want nil", err)
	}
	defer cleanup()

	for i := range watchSubscriberBufferSize + 1 {
		registerInstance(backend, "service-a", fmt.Sprintf("a%d", i), time.Second)
	}

	for range watchSubscriberBufferSize {
		receiveChange(t, subscriber)
	}

	if _, ok := <-subscriber.Events(); ok {
		t.Error("subscriber channel still open after overflowing its buffer, want it closed")
	}

	manager.mu.Lock()
	remaining := len(manager.subscribers["service-a"])
	manager.mu.Unlock()

	if remaining != 0 {
		t.Errorf("subscribers left = %d, want 0", remaining)
	}
}

func TestWatchSendsSnapshotOfEachServiceOnRestore(t *testing.T) {
	backend := newFakeBackend()
	manager := NewWatchManager(backend)

	_, subscriberA, cleanupA, _ := manager.Watch(t.Context(), "service-a")
	defer cleanupA()

	_, subscriberB, cleanupB, _ := manager.Watch(t.Context(), "service-b")
	defer cleanupB()

	backend.restore([]ServiceInstance{
		{ID: ServiceInstanceID{ServiceName: "service-a", InstanceID: "a1"}, Status: ServiceInstanceStatusHealthy},
		{ID: ServiceInstanceID{ServiceName: "service-a", InstanceID: "a2"}, Status: ServiceInstanceStatusExpired},
		{ID: ServiceInstanceID{ServiceName: "service-b", InstanceID: "b1"}, Status: ServiceInstanceStatusHealthy},
	})

	tests := []struct {
		subscriber *Subscriber
		want       []string
	}{
		{subscriberA, []string{"a1"}},
		{subscriberB, []string{"b1"}},
	}

	for _, tt := range tests {
		change := receiveChange(t, tt.subscriber)
		if change.Type != StateChangeRestored {
			t.Errorf("event type = %s, want %s", change.Type, StateChangeRestored)
		}

		if got := instanceIDs(change.Instances); !slices.Equal(got, tt.want) {
			t.Errorf("restored instances = %v, want healthy instances of the service %v", got, tt.want)
		}
	}
}

func TestWatchFailsWhenSnapshotLookupFails(t *testing.T) {
	backend := newFakeBackend()
	backend.lookupErr = ErrUnavailable
	manager := NewWatchManager(backend)

	if _, _, _, err := manager.Watch(t.Context(), "service-a"); err == nil {
		t.Fatal("Watch() error = nil, want the lookup error")
	}

	manager.mu.Lock()
	remaining := len(manager.subscribers["service-a"])
	manager.mu.Unlock()

	if remaining != 0 {
		t.Errorf("subscribers left = %d, want 0 after a failed Watch()", remaining)
	}
}
