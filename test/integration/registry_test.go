package integration

import (
	"slices"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	registrypb "tpiii.local/daniel-tpiii/gen/registry"
)

func TestRegistryInstanceLifecycle(t *testing.T) {
	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			c := newCluster(t, backend, 3)

			if err := c.nodes[0].register("service-a", "a1", defaultHeartbeatInterval); err != nil {
				t.Fatalf("Register() error = %v, want nil", err)
			}
			c.waitVisible("service-a", "a1")

			if err := c.nodes[1].register("service-a", "a1", defaultHeartbeatInterval); status.Code(err) != codes.AlreadyExists {
				t.Errorf("duplicate Register() error = %v, want code %s", err, codes.AlreadyExists)
			}

			if err := c.nodes[1].update("service-a", "a1", 5); err != nil {
				t.Fatalf("Update() error = %v, want nil", err)
			}

			eventually(t, "every node reports the new weight", func() bool {
				return c.everyRunningNode(func(node *registryNode) bool {
					instances, err := node.lookup("service-a")
					return err == nil && instances["a1"].GetWeight() == 5
				})
			})

			if err := c.nodes[0].heartbeat("service-a", "unknown"); status.Code(err) != codes.NotFound {
				t.Errorf("Heartbeat(unknown) error = %v, want code %s", err, codes.NotFound)
			}

			if instances, err := c.nodes[2].lookup("unknown-service"); err != nil || len(instances) != 0 {
				t.Errorf("Lookup(unknown service) = %v, %v; want an empty list and nil error", instances, err)
			}

			if err := c.nodes[2].deregister("service-a", "a1"); err != nil {
				t.Fatalf("Deregister() error = %v, want nil", err)
			}
			c.waitGone("service-a", "a1")

			eventually(t, "every node rejects heartbeats of the deleted instance", func() bool {
				return c.everyRunningNode(func(node *registryNode) bool {
					return status.Code(node.heartbeat("service-a", "a1")) == codes.FailedPrecondition
				})
			})
		})
	}
}

func TestRegistryExpiresAndRecoversInstancesByHeartbeat(t *testing.T) {
	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			c := newCluster(t, backend, 3)
			w := c.nodes[2].watch(t, "service-a")

			if err := c.nodes[0].register("service-a", "a1", 200*time.Millisecond); err != nil {
				t.Fatalf("Register() error = %v, want nil", err)
			}
			w.waitEvent(t, registrypb.StateChangeType_STATE_CHANGE_TYPE_REGISTERED, "a1")

			if err := c.nodes[0].heartbeat("service-a", "a1"); err != nil {
				t.Fatalf("Heartbeat() error = %v, want nil", err)
			}

			w.waitEvent(t, registrypb.StateChangeType_STATE_CHANGE_TYPE_EXPIRED, "a1")
			c.waitGone("service-a", "a1")

			if err := c.nodes[1].heartbeat("service-a", "a1"); err != nil {
				t.Fatalf("Heartbeat(expired) on another node error = %v, want nil", err)
			}

			w.waitEvent(t, registrypb.StateChangeType_STATE_CHANGE_TYPE_HEALTHY, "a1")
			c.waitVisible("service-a", "a1")

			w.waitEvent(t, registrypb.StateChangeType_STATE_CHANGE_TYPE_EXPIRED, "a1")
		})
	}
}

func TestRegistryWatchStreamsChangesOfService(t *testing.T) {
	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			c := newCluster(t, backend, 3)

			if err := c.nodes[0].register("service-a", "a1", defaultHeartbeatInterval); err != nil {
				t.Fatalf("Register() error = %v, want nil", err)
			}
			c.waitVisible("service-a", "a1")

			w := c.nodes[2].watch(t, "service-a")
			if got, want := sortedIDs(w.snapshot), []string{"a1"}; !slices.Equal(got, want) {
				t.Errorf("watch snapshot = %v, want %v", got, want)
			}

			if err := c.nodes[0].register("service-b", "b1", defaultHeartbeatInterval); err != nil {
				t.Fatalf("Register(service-b) error = %v, want nil", err)
			}

			if err := c.nodes[1].register("service-a", "a2", defaultHeartbeatInterval); err != nil {
				t.Fatalf("Register(a2) error = %v, want nil", err)
			}
			w.waitEvent(t, registrypb.StateChangeType_STATE_CHANGE_TYPE_REGISTERED, "a2")

			c.waitVisible("service-a", "a2")
			if err := c.nodes[0].update("service-a", "a2", 3); err != nil {
				t.Fatalf("Update() error = %v, want nil", err)
			}

			event := w.waitEvent(t, registrypb.StateChangeType_STATE_CHANGE_TYPE_UPDATED, "a2")
			if weight := event.GetInstance().GetWeight(); weight != 3 {
				t.Errorf("updated weight = %d, want 3", weight)
			}

			if err := c.nodes[1].deregister("service-a", "a1"); err != nil {
				t.Fatalf("Deregister() error = %v, want nil", err)
			}
			w.waitEvent(t, registrypb.StateChangeType_STATE_CHANGE_TYPE_DELETED, "a1")
		})
	}
}
