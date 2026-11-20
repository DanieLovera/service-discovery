package cp

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"slices"
	"testing"
	"time"

	"github.com/hashicorp/raft"

	"tpiii.local/daniel-tpiii/internal/registry"
)

type memorySink struct {
	bytes.Buffer
	closed bool
}

func (s *memorySink) ID() string    { return "test-snapshot" }
func (s *memorySink) Cancel() error { return nil }
func (s *memorySink) Close() error  { s.closed = true; return nil }

func instanceID(name string) registry.ServiceInstanceID {
	return registry.ServiceInstanceID{ServiceName: "service-a", InstanceID: name}
}

func registerCommand(name string) command {
	return newRegisterCommand(registry.RegisterParams{
		ID:                instanceID(name),
		Address:           name + ":8080",
		Weight:            1,
		HeartbeatInterval: time.Second,
	})
}

func updateCommand(name string, weight int) command {
	return newUpdateCommand(registry.UpdateParams{
		ID:                instanceID(name),
		Address:           name + ":8080",
		Weight:            weight,
		HeartbeatInterval: time.Second,
	})
}

func applyCommand(t *testing.T, f *fsm, cmd command) error {
	t.Helper()

	data, err := json.Marshal(cmd)
	if err != nil {
		t.Fatalf("encode command: %v", err)
	}

	return applyData(t, f, data)
}

func applyData(t *testing.T, f *fsm, data []byte) error {
	t.Helper()

	result := f.Apply(&raft.Log{Data: data})
	if result == nil {
		return nil
	}

	err, ok := result.(error)
	if !ok {
		t.Fatalf("Apply() returned %T, want error or nil", result)
	}

	return err
}

func instanceOf(t *testing.T, store *stateStore, name string) registry.ServiceInstance {
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

func TestFSMAppliesInstanceLifecycle(t *testing.T) {
	store := newStateStore()
	f := newFSM(store)

	steps := []struct {
		name       string
		cmd        command
		wantStatus registry.ServiceInstanceStatus
		wantWeight int
	}{
		{"register", registerCommand("a1"), registry.ServiceInstanceStatusHealthy, 1},
		{"update", updateCommand("a1", 3), registry.ServiceInstanceStatusHealthy, 3},
		{"expire", newInstanceCommand(commandExpire, instanceID("a1")), registry.ServiceInstanceStatusExpired, 3},
		{"recover", newInstanceCommand(commandRecover, instanceID("a1")), registry.ServiceInstanceStatusHealthy, 3},
		{"deregister", newInstanceCommand(commandDeregister, instanceID("a1")), registry.ServiceInstanceStatusDeleted, 3},
		{"register again", registerCommand("a1"), registry.ServiceInstanceStatusHealthy, 1},
	}

	for _, step := range steps {
		if err := applyCommand(t, f, step.cmd); err != nil {
			t.Fatalf("%s: Apply() error = %v, want nil", step.name, err)
		}

		instance := instanceOf(t, store, "a1")
		if instance.Status != step.wantStatus || instance.Weight != step.wantWeight {
			t.Errorf("%s: instance = %s with weight %d, want %s with weight %d",
				step.name, instance.Status, instance.Weight, step.wantStatus, step.wantWeight)
		}
	}
}

func TestFSMReturnsDomainErrors(t *testing.T) {
	store := newStateStore()
	f := newFSM(store)

	applyCommand(t, f, registerCommand("a1"))
	applyCommand(t, f, registerCommand("deleted"))
	applyCommand(t, f, newInstanceCommand(commandDeregister, instanceID("deleted")))

	missing := instanceID("missing")
	deleted := instanceID("deleted")

	tests := []struct {
		name string
		cmd  command
		want error
	}{
		{"register existing", registerCommand("a1"), registry.ErrServiceInstanceAlreadyExists},
		{"update missing", updateCommand("missing", 2), registry.ErrServiceInstanceNotFound},
		{"deregister missing", newInstanceCommand(commandDeregister, missing), registry.ErrServiceInstanceNotFound},
		{"expire deleted", newInstanceCommand(commandExpire, deleted), registry.ErrServiceInstanceDeleted},
		{"recover deleted", newInstanceCommand(commandRecover, deleted), registry.ErrServiceInstanceDeleted},
	}

	for _, tt := range tests {
		if err := applyCommand(t, f, tt.cmd); !errors.Is(err, tt.want) {
			t.Errorf("%s: Apply() error = %v, want %v", tt.name, err, tt.want)
		}
	}

	if err := applyCommand(t, f, command{Type: "UNKNOWN", ServiceName: "service-a", InstanceID: "a1"}); err == nil {
		t.Error("Apply() of an unknown command type error = nil, want error")
	}

	if err := applyData(t, f, []byte("not json")); err == nil {
		t.Error("Apply() of invalid data error = nil, want error")
	}
}

func TestFSMDeduplicatesRequestIDs(t *testing.T) {
	store := newStateStore()
	f := newFSM(store)

	applyCommand(t, f, registerCommand("a1"))

	expire := newInstanceCommand(commandExpire, instanceID("a1")).withRequestID("expire-1")
	if err := applyCommand(t, f, expire); err != nil {
		t.Fatalf("Expire() error = %v, want nil", err)
	}

	recovery := newInstanceCommand(commandRecover, instanceID("a1")).withRequestID("recover-1")
	if err := applyCommand(t, f, recovery); err != nil {
		t.Fatalf("Recover() error = %v, want nil", err)
	}

	if err := applyCommand(t, f, expire); err != nil {
		t.Errorf("replayed Expire() error = %v, want the recorded result nil", err)
	}

	if status := instanceOf(t, store, "a1").Status; status != registry.ServiceInstanceStatusHealthy {
		t.Errorf("status after replayed Expire() = %s, want %s", status, registry.ServiceInstanceStatusHealthy)
	}

	update := updateCommand("a2", 5).withRequestID("update-1")
	if err := applyCommand(t, f, update); !errors.Is(err, registry.ErrServiceInstanceNotFound) {
		t.Fatalf("Update(missing) error = %v, want %v", err, registry.ErrServiceInstanceNotFound)
	}

	applyCommand(t, f, registerCommand("a2"))

	if err := applyCommand(t, f, update); !errors.Is(err, registry.ErrServiceInstanceNotFound) {
		t.Errorf("replayed Update() error = %v, want the recorded result %v", err, registry.ErrServiceInstanceNotFound)
	}

	if weight := instanceOf(t, store, "a2").Weight; weight != 1 {
		t.Errorf("weight after replayed Update() = %d, want 1", weight)
	}
}

func TestFSMSnapshotRestore(t *testing.T) {
	source := newStateStore()
	f := newFSM(source)

	for _, name := range []string{"a1", "a2", "a3"} {
		applyCommand(t, f, registerCommand(name))
	}

	applyCommand(t, f, newInstanceCommand(commandExpire, instanceID("a2")).withRequestID("expire-a2"))
	applyCommand(t, f, newInstanceCommand(commandDeregister, instanceID("a3")))

	snapshot, err := f.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot() error = %v, want nil", err)
	}

	sink := &memorySink{}
	if err := snapshot.Persist(sink); err != nil {
		t.Fatalf("Persist() error = %v, want nil", err)
	}

	if !sink.closed {
		t.Error("Persist() did not close the sink")
	}

	target := newStateStore()
	var restored []registry.StateChange
	target.RegisterStateChangeHandler(func(change registry.StateChange) {
		restored = append(restored, change)
	})

	if err := newFSM(target).Restore(io.NopCloser(&sink.Buffer)); err != nil {
		t.Fatalf("Restore() error = %v, want nil", err)
	}

	for _, name := range []string{"a1", "a2", "a3"} {
		if got, want := instanceOf(t, target, name), instanceOf(t, source, name); got != want {
			t.Errorf("restored %s = %+v, want %+v", name, got, want)
		}
	}

	if _, ok := target.RequestResult("expire-a2"); !ok {
		t.Error("processed request ids were not restored")
	}

	if len(restored) != 1 || restored[0].Type != registry.StateChangeRestored {
		t.Fatalf("notified changes = %+v, want a single %s change", restored, registry.StateChangeRestored)
	}

	if got, want := instanceNames(restored[0].Instances), []string{"a1", "a2", "a3"}; !slices.Equal(got, want) {
		t.Errorf("restored change instances = %v, want %v", got, want)
	}
}

func TestStateStoreLookup(t *testing.T) {
	store := newStateStore()
	f := newFSM(store)

	for _, name := range []string{"a1", "a2", "a3"} {
		applyCommand(t, f, registerCommand(name))
	}

	applyCommand(t, f, newInstanceCommand(commandExpire, instanceID("a2")))
	applyCommand(t, f, newInstanceCommand(commandDeregister, instanceID("a3")))

	instances, err := store.Lookup(registry.ServiceInstanceID{ServiceName: "service-a"})
	if err != nil {
		t.Fatalf("service Lookup() error = %v, want nil", err)
	}

	if got, want := instanceNames(instances), []string{"a1"}; !slices.Equal(got, want) {
		t.Errorf("service Lookup() = %v, want only healthy instances %v", got, want)
	}

	if status := instanceOf(t, store, "a3").Status; status != registry.ServiceInstanceStatusDeleted {
		t.Errorf("instance Lookup() status = %s, want %s", status, registry.ServiceInstanceStatusDeleted)
	}

	if _, err := store.Lookup(instanceID("missing")); !errors.Is(err, registry.ErrServiceInstanceNotFound) {
		t.Errorf("Lookup(missing) error = %v, want %v", err, registry.ErrServiceInstanceNotFound)
	}
}
