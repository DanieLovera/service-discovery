package ap

import (
	"errors"
	"maps"
	"slices"
	"testing"
	"time"

	"tpiii.local/daniel-tpiii/internal/registry"
)

type changeRecorder struct {
	changes []registry.StateChange
}

func (r *changeRecorder) record(change registry.StateChange) {
	r.changes = append(r.changes, change)
}

func instanceID(name string) registry.ServiceInstanceID {
	return registry.ServiceInstanceID{ServiceName: "service-a", InstanceID: name}
}

func registerParams(name string) registry.RegisterParams {
	return registry.RegisterParams{
		ID:                instanceID(name),
		Address:           name + ":8080",
		Weight:            1,
		HeartbeatInterval: time.Second,
	}
}

func healthyInstance(name string) registry.ServiceInstance {
	return registry.ServiceInstance{
		ID:                instanceID(name),
		Address:           name + ":8080",
		Weight:            1,
		HeartbeatInterval: time.Second,
		Status:            registry.ServiceInstanceStatusHealthy,
	}
}

func lookupInstance(t *testing.T, store *stateStore, name string) registry.ServiceInstance {
	t.Helper()

	instances, err := store.Lookup(instanceID(name))
	if err != nil {
		t.Fatalf("Lookup(%s) error = %v, want nil", name, err)
	}

	return instances[0]
}

func instanceNames(instances []registry.ServiceInstance) []string {
	names := make([]string, 0, len(instances))
	for _, instance := range instances {
		names = append(names, instance.ID.InstanceID)
	}

	slices.Sort(names)

	return names
}

func TestStateStoreLocalWritesAdvanceLamportClock(t *testing.T) {
	store := newStateStore("node-1")

	registered, changed, err := store.Register(registerParams("a1"))
	if err != nil || !changed {
		t.Fatalf("Register() = changed %v, error %v; want changed true, error nil", changed, err)
	}

	if want := (version{LamportTimestamp: 1, NodeID: "node-1"}); registered.Version != want {
		t.Errorf("Register() version = %+v, want %+v", registered.Version, want)
	}

	expired, _, err := store.Expire(instanceID("a1"))
	if err != nil {
		t.Fatalf("Expire() error = %v, want nil", err)
	}

	if !expired.Version.newerThan(registered.Version) {
		t.Errorf("Expire() version %+v is not newer than Register() version %+v", expired.Version, registered.Version)
	}
}

func TestStateStoreLocalWriteRules(t *testing.T) {
	store := newStateStore("node-1")
	recorder := &changeRecorder{}
	store.RegisterStateChangeHandler(recorder.record)

	if _, _, err := store.Register(registerParams("a1")); err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}

	if _, _, err := store.Register(registerParams("a1")); !errors.Is(err, registry.ErrServiceInstanceAlreadyExists) {
		t.Errorf("duplicate Register() error = %v, want %v", err, registry.ErrServiceInstanceAlreadyExists)
	}

	unchanged := registry.UpdateParams{ID: instanceID("a1"), Address: "a1:8080", Weight: 1, HeartbeatInterval: time.Second}
	if _, changed, err := store.Update(unchanged); err != nil || changed {
		t.Errorf("Update() with same values = changed %v, error %v; want changed false, error nil", changed, err)
	}

	_, changed, err := store.Expire(instanceID("missing"))
	if !errors.Is(err, registry.ErrServiceInstanceNotFound) || changed {
		t.Errorf("Expire(missing) = changed %v, error %v; want %v", changed, err, registry.ErrServiceInstanceNotFound)
	}

	if _, _, err := store.Deregister(instanceID("a1")); err != nil {
		t.Fatalf("Deregister() error = %v, want nil", err)
	}

	if _, changed, err := store.Deregister(instanceID("a1")); err != nil || changed {
		t.Errorf("second Deregister() = changed %v, error %v; want changed false, error nil", changed, err)
	}

	updated := registry.UpdateParams{ID: instanceID("a1"), Address: "a1:9090", Weight: 2, HeartbeatInterval: time.Second}
	if _, _, err := store.Update(updated); !errors.Is(err, registry.ErrServiceInstanceDeleted) {
		t.Errorf("Update(deleted) error = %v, want %v", err, registry.ErrServiceInstanceDeleted)
	}

	if _, _, err := store.Register(registerParams("a1")); err != nil {
		t.Errorf("Register() over a deleted instance error = %v, want nil", err)
	}

	if status := lookupInstance(t, store, "a1").Status; status != registry.ServiceInstanceStatusHealthy {
		t.Errorf("status after registering again = %s, want %s", status, registry.ServiceInstanceStatusHealthy)
	}

	var types []registry.StateChangeType
	for _, change := range recorder.changes {
		types = append(types, change.Type)
	}

	want := []registry.StateChangeType{
		registry.StateChangeRegistered,
		registry.StateChangeDeleted,
		registry.StateChangeRegistered,
	}
	if !slices.Equal(types, want) {
		t.Errorf("notified changes = %v, want %v", types, want)
	}
}

func TestStateStoreMergeKeepsNewestVersion(t *testing.T) {
	current := version{LamportTimestamp: 5, NodeID: "node-2"}

	tests := []struct {
		name        string
		remote      version
		wantApplied bool
	}{
		{name: "higher timestamp wins", remote: version{LamportTimestamp: 6, NodeID: "node-1"}, wantApplied: true},
		{name: "lower timestamp loses", remote: version{LamportTimestamp: 4, NodeID: "node-3"}, wantApplied: false},
		{name: "tie broken by higher node id", remote: version{LamportTimestamp: 5, NodeID: "node-3"}, wantApplied: true},
		{name: "tie lost by lower node id", remote: version{LamportTimestamp: 5, NodeID: "node-1"}, wantApplied: false},
		{name: "same version is ignored", remote: current, wantApplied: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newStateStore("node-9")
			store.Merge(instanceState{Instance: healthyInstance("a1"), Version: current})

			remote := healthyInstance("a1")
			remote.Address = "remote:8080"

			applied := store.Merge(instanceState{Instance: remote, Version: tt.remote})
			if applied != tt.wantApplied {
				t.Fatalf("Merge() = %v, want %v", applied, tt.wantApplied)
			}

			wantAddress := "a1:8080"
			if tt.wantApplied {
				wantAddress = "remote:8080"
			}

			if address := lookupInstance(t, store, "a1").Address; address != wantAddress {
				t.Errorf("address after Merge() = %q, want %q", address, wantAddress)
			}
		})
	}
}

func TestStateStoreLocalWriteAfterMergeWinsOverRemoteState(t *testing.T) {
	remoteState := instanceState{
		Instance: healthyInstance("a1"),
		Version:  version{LamportTimestamp: 10, NodeID: "node-9"},
	}

	local := newStateStore("node-1")
	local.Merge(remoteState)

	deleted, _, err := local.Deregister(instanceID("a1"))
	if err != nil {
		t.Fatalf("Deregister() error = %v, want nil", err)
	}

	if deleted.Version.LamportTimestamp != 11 {
		t.Errorf("Deregister() timestamp = %d, want 11 (clock catches up with merged states)",
			deleted.Version.LamportTimestamp)
	}

	remote := newStateStore("node-9")
	remote.Merge(remoteState)

	if !remote.Merge(deleted) {
		t.Fatal("remote Merge() of the newer deregister = false, want true")
	}

	if status := lookupInstance(t, remote, "a1").Status; status != registry.ServiceInstanceStatusDeleted {
		t.Errorf("remote status = %s, want %s", status, registry.ServiceInstanceStatusDeleted)
	}
}

func TestStateStoreMergeNotifiesChangeType(t *testing.T) {
	withStatus := func(status registry.ServiceInstanceStatus) registry.ServiceInstance {
		instance := healthyInstance("a1")
		instance.Status = status
		return instance
	}

	reweighted := healthyInstance("a1")
	reweighted.Weight = 5

	tests := []struct {
		name     string
		previous *registry.ServiceInstance
		next     registry.ServiceInstance
		want     []registry.StateChangeType
	}{
		{
			name: "unknown healthy instance is registered",
			next: withStatus(registry.ServiceInstanceStatusHealthy),
			want: []registry.StateChangeType{registry.StateChangeRegistered},
		},
		{
			name: "unknown deleted instance is deleted",
			next: withStatus(registry.ServiceInstanceStatusDeleted),
			want: []registry.StateChangeType{registry.StateChangeDeleted},
		},
		{
			name:     "healthy to expired",
			previous: new(withStatus(registry.ServiceInstanceStatusHealthy)),
			next:     withStatus(registry.ServiceInstanceStatusExpired),
			want:     []registry.StateChangeType{registry.StateChangeExpired},
		},
		{
			name:     "expired to healthy",
			previous: new(withStatus(registry.ServiceInstanceStatusExpired)),
			next:     withStatus(registry.ServiceInstanceStatusHealthy),
			want:     []registry.StateChangeType{registry.StateChangeHealthy},
		},
		{
			name:     "deleted to healthy is registered again",
			previous: new(withStatus(registry.ServiceInstanceStatusDeleted)),
			next:     withStatus(registry.ServiceInstanceStatusHealthy),
			want:     []registry.StateChangeType{registry.StateChangeRegistered},
		},
		{
			name:     "healthy with new data is updated",
			previous: new(withStatus(registry.ServiceInstanceStatusHealthy)),
			next:     reweighted,
			want:     []registry.StateChangeType{registry.StateChangeUpdated},
		},
		{
			name:     "same instance with newer version is not notified",
			previous: new(withStatus(registry.ServiceInstanceStatusHealthy)),
			next:     withStatus(registry.ServiceInstanceStatusHealthy),
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newStateStore("node-1")
			if tt.previous != nil {
				store.Merge(instanceState{Instance: *tt.previous, Version: version{LamportTimestamp: 1, NodeID: "node-2"}})
			}

			recorder := &changeRecorder{}
			store.RegisterStateChangeHandler(recorder.record)

			if !store.Merge(instanceState{Instance: tt.next, Version: version{LamportTimestamp: 2, NodeID: "node-2"}}) {
				t.Fatal("Merge() = false, want true")
			}

			var types []registry.StateChangeType
			for _, change := range recorder.changes {
				types = append(types, change.Type)
			}

			if !slices.Equal(types, tt.want) {
				t.Errorf("notified changes = %v, want %v", types, tt.want)
			}
		})
	}
}

func TestStateStoreAntiEntropyRoundConverges(t *testing.T) {
	a := newStateStore("node-a")
	b := newStateStore("node-b")

	shared, _, _ := a.Register(registerParams("shared"))
	b.Merge(shared)

	if _, _, err := b.Expire(instanceID("shared")); err != nil {
		t.Fatalf("Expire() error = %v, want nil", err)
	}

	a.Register(registerParams("only-a"))
	b.Register(registerParams("only-b"))

	newer, requested := b.Diff(a.Digest())
	for _, state := range newer {
		a.Merge(state)
	}

	for _, state := range a.States(requested) {
		b.Merge(state)
	}

	if !maps.Equal(a.Digest(), b.Digest()) {
		t.Fatalf("digests differ after one round:\na = %v\nb = %v", a.Digest(), b.Digest())
	}

	if newer, requested := b.Diff(a.Digest()); len(newer) != 0 || len(requested) != 0 {
		t.Errorf("second Diff() = %d newer, %d requested; want 0 and 0", len(newer), len(requested))
	}

	if status := lookupInstance(t, a, "shared").Status; status != registry.ServiceInstanceStatusExpired {
		t.Errorf("shared instance status on a = %s, want %s", status, registry.ServiceInstanceStatusExpired)
	}

	want := []string{"only-a", "only-b", "shared"}
	if got := instanceNames(a.Instances()); !slices.Equal(got, want) {
		t.Errorf("instances on a = %v, want %v", got, want)
	}
}

func TestStateStoreLookup(t *testing.T) {
	store := newStateStore("node-1")

	for _, name := range []string{"a1", "a2", "a3"} {
		store.Register(registerParams(name))
	}

	store.Expire(instanceID("a2"))
	store.Deregister(instanceID("a3"))
	store.Register(registry.RegisterParams{
		ID:                registry.ServiceInstanceID{ServiceName: "service-b", InstanceID: "b1"},
		Address:           "b1:8080",
		HeartbeatInterval: time.Second,
	})

	instances, err := store.Lookup(registry.ServiceInstanceID{ServiceName: "service-a"})
	if err != nil {
		t.Fatalf("service Lookup() error = %v, want nil", err)
	}

	if got, want := instanceNames(instances), []string{"a1"}; !slices.Equal(got, want) {
		t.Errorf("service Lookup() = %v, want only healthy instances %v", got, want)
	}

	if status := lookupInstance(t, store, "a2").Status; status != registry.ServiceInstanceStatusExpired {
		t.Errorf("instance Lookup() status = %s, want %s", status, registry.ServiceInstanceStatusExpired)
	}

	if _, err := store.Lookup(instanceID("missing")); !errors.Is(err, registry.ErrServiceInstanceNotFound) {
		t.Errorf("Lookup(missing) error = %v, want %v", err, registry.ErrServiceInstanceNotFound)
	}

	instances, err = store.Lookup(registry.ServiceInstanceID{ServiceName: "unknown"})
	if err != nil || len(instances) != 0 {
		t.Errorf("Lookup(unknown service) = %v, %v; want empty list and nil error", instances, err)
	}
}
