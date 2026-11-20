package integration

import (
	"slices"
	"sync"
	"testing"
)

func TestAPRemainsAvailableWithoutMajority(t *testing.T) {
	c := newCluster(t, "ap", 3)

	c.stop(0)
	c.stop(1)
	survivor := c.nodes[2]

	if err := survivor.register("service-a", "a1", defaultHeartbeatInterval); err != nil {
		t.Fatalf("Register() on the only running node error = %v, want nil", err)
	}

	instances, err := survivor.lookup("service-a")
	if err != nil || instances["a1"] == nil {
		t.Errorf("Lookup() on the only running node = %v, %v; want a1", instances, err)
	}
}

func TestAPRestartedNodeSyncsStateBeforeServing(t *testing.T) {
	c := newCluster(t, "ap", 3)

	c.stop(2)

	if err := c.nodes[0].register("service-a", "a1", defaultHeartbeatInterval); err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}
	c.waitVisible("service-a", "a1")

	c.start(2)

	instances, err := c.nodes[2].lookup("service-a")
	if err != nil || instances["a1"] == nil {
		t.Errorf("Lookup() on the restarted node = %v, %v; want a1 without waiting", instances, err)
	}

	w := c.nodes[2].watch(t, "service-a")
	if got := sortedIDs(w.snapshot); !slices.Equal(got, []string{"a1"}) {
		t.Errorf("watch snapshot on the restarted node = %v, want [a1]", got)
	}
}

func TestAPConcurrentUpdatesConverge(t *testing.T) {
	c := newCluster(t, "ap", 3)

	if err := c.nodes[0].register("service-a", "a1", defaultHeartbeatInterval); err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}
	c.waitVisible("service-a", "a1")

	var wg sync.WaitGroup
	errs := make([]error, 2)
	for i, weight := range []int32{2, 3} {
		wg.Go(func() {
			errs[i] = c.nodes[i].update("service-a", "a1", weight)
		})
	}
	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Fatalf("Update() on %s error = %v, want nil", c.nodes[i].params.NodeID, err)
		}
	}

	eventually(t, "every node converges to the same weight", func() bool {
		weights := make(map[int32]bool)
		converged := c.everyRunningNode(func(node *registryNode) bool {
			instances, err := node.lookup("service-a")
			if err != nil || instances["a1"] == nil {
				return false
			}

			weights[instances["a1"].GetWeight()] = true
			return true
		})

		return converged && len(weights) == 1 && (weights[2] || weights[3])
	})
}

func TestAPFailureDetectorConfirmsAndReadmitsNode(t *testing.T) {
	c := newCluster(t, "ap", 3)
	observer := c.nodes[0]

	c.stop(2)

	eventually(t, "the observer confirms the stopped node as failed", func() bool {
		value, ok := observer.metric(`registry_swim_members{status="CONFIRM"}`)
		return ok && value == 1
	})

	c.start(2)

	eventually(t, "the observer sees every peer alive again", func() bool {
		value, ok := observer.metric(`registry_swim_members{status="ALIVE"}`)
		return ok && value == 2
	})

	if err := observer.register("service-a", "a1", defaultHeartbeatInterval); err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}
	c.waitVisible("service-a", "a1")
}
