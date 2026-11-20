package integration

import (
	"fmt"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCPAcknowledgedWritesAreReadFromEveryNode(t *testing.T) {
	c := newCluster(t, "cp", 3)

	for i, writer := range c.nodes {
		instanceID := fmt.Sprintf("a%d", i+1)
		if err := writer.register("service-a", instanceID, defaultHeartbeatInterval); err != nil {
			t.Fatalf("Register() on %s error = %v, want nil", writer.params.NodeID, err)
		}

		for _, reader := range c.nodes {
			instances, err := reader.lookup("service-a")
			if err != nil || instances[instanceID] == nil {
				t.Errorf("Lookup() on %s right after Register() on %s = %v, %v; want %s without waiting",
					reader.params.NodeID, writer.params.NodeID, instances, err, instanceID)
			}
		}
	}
}

func TestCPElectsNewLeaderAfterLeaderFailure(t *testing.T) {
	c := newCluster(t, "cp", 3)

	if err := c.nodes[0].register("service-a", "a1", defaultHeartbeatInterval); err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}

	oldLeader := c.leader()
	c.stop(oldLeader)

	survivor := c.nodes[(oldLeader+1)%len(c.nodes)]
	eventually(t, "a surviving node accepts writes", func() bool {
		err := survivor.register("service-a", "a2", defaultHeartbeatInterval)
		return err == nil || status.Code(err) == codes.AlreadyExists
	})

	if newLeader := c.leader(); newLeader == oldLeader {
		t.Errorf("leader after failure = %s, want a different node", c.nodes[newLeader].params.NodeID)
	}

	for _, node := range c.nodes {
		if node.process == nil {
			continue
		}

		instances, err := node.lookup("service-a")
		if err != nil || instances["a1"] == nil || instances["a2"] == nil {
			t.Errorf("Lookup() on %s = %v, %v; want a1 and a2", node.params.NodeID, instances, err)
		}
	}
}

func TestCPMinorityRejectsOperationsUntilMajorityReturns(t *testing.T) {
	c := newCluster(t, "cp", 3)

	if err := c.nodes[0].register("service-a", "a1", defaultHeartbeatInterval); err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}

	c.stop(0)
	c.stop(1)
	minority := c.nodes[2]

	if err := minority.register("service-a", "a2", defaultHeartbeatInterval); status.Code(err) != codes.Unavailable {
		t.Errorf("Register() without quorum error = %v, want code %s", err, codes.Unavailable)
	}

	if _, err := minority.lookup("service-a"); status.Code(err) != codes.Unavailable {
		t.Errorf("Lookup() without quorum error = %v, want code %s", err, codes.Unavailable)
	}

	c.start(1)

	eventually(t, "writes are accepted again once a majority is back", func() bool {
		err := minority.register("service-a", "a2", defaultHeartbeatInterval)
		return err == nil || status.Code(err) == codes.AlreadyExists
	})

	instances, err := c.nodes[1].lookup("service-a")
	if err != nil || instances["a1"] == nil || instances["a2"] == nil {
		t.Errorf("Lookup() on the restarted node = %v, %v; want a1 and a2", instances, err)
	}
}
