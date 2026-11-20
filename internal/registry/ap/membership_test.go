package ap

import (
	"log/slog"
	"slices"
	"testing"
	"time"

	"tpiii.local/daniel-tpiii/internal/registry"
)

func newTestMembership(suspectTimeout time.Duration) *membership {
	peers := []registry.ClusterMember{{ID: "node-2"}, {ID: "node-3"}, {ID: "node-4"}, {ID: "node-5"}}

	return newMembership("node-1", peers, suspectTimeout, 3, slog.New(slog.DiscardHandler))
}

func waitForStatus(t *testing.T, m *membership, nodeID string, want memberStatus) {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for m.status(nodeID) != want {
		if time.Now().After(deadline) {
			t.Fatalf("status of %s = %s, want %s", nodeID, m.status(nodeID), want)
		}

		time.Sleep(5 * time.Millisecond)
	}
}

func TestMembershipUpdateSupersedes(t *testing.T) {
	tests := []struct {
		current memberState
		update  membershipUpdate
		want    bool
	}{
		{memberState{memberAlive, 1}, membershipUpdate{Status: memberAlive, Incarnation: 2}, true},
		{memberState{memberAlive, 1}, membershipUpdate{Status: memberAlive, Incarnation: 1}, false},
		{memberState{memberSuspect, 1}, membershipUpdate{Status: memberAlive, Incarnation: 2}, true},
		{memberState{memberSuspect, 1}, membershipUpdate{Status: memberAlive, Incarnation: 1}, false},
		{memberState{memberConfirmed, 1}, membershipUpdate{Status: memberAlive, Incarnation: 2}, true},
		{memberState{memberConfirmed, 1}, membershipUpdate{Status: memberAlive, Incarnation: 1}, false},

		{memberState{memberAlive, 1}, membershipUpdate{Status: memberSuspect, Incarnation: 1}, true},
		{memberState{memberAlive, 1}, membershipUpdate{Status: memberSuspect, Incarnation: 0}, false},
		{memberState{memberSuspect, 1}, membershipUpdate{Status: memberSuspect, Incarnation: 2}, true},
		{memberState{memberSuspect, 1}, membershipUpdate{Status: memberSuspect, Incarnation: 1}, false},
		{memberState{memberConfirmed, 1}, membershipUpdate{Status: memberSuspect, Incarnation: 5}, false},

		{memberState{memberAlive, 1}, membershipUpdate{Status: memberConfirmed, Incarnation: 1}, true},
		{memberState{memberAlive, 1}, membershipUpdate{Status: memberConfirmed, Incarnation: 0}, false},
		{memberState{memberSuspect, 1}, membershipUpdate{Status: memberConfirmed, Incarnation: 1}, true},
		{memberState{memberConfirmed, 1}, membershipUpdate{Status: memberConfirmed, Incarnation: 2}, true},
		{memberState{memberConfirmed, 1}, membershipUpdate{Status: memberConfirmed, Incarnation: 1}, false},
	}

	for _, tt := range tests {
		got := tt.update.supersedes(tt.current)
		if got != tt.want {
			t.Errorf("%s(%d) over %s(%d) = %v, want %v",
				tt.update.Status, tt.update.Incarnation, tt.current.status, tt.current.incarnation, got, tt.want)
		}
	}
}

func TestMembershipConfirmsSuspectAfterTimeout(t *testing.T) {
	m := newTestMembership(50 * time.Millisecond)

	m.suspect("node-2")

	if status := m.status("node-2"); status != memberSuspect {
		t.Fatalf("status after suspect() = %s, want %s", status, memberSuspect)
	}

	if alive := len(m.alivePeers()); alive != 3 {
		t.Errorf("alive peers = %d, want 3 (suspects are excluded)", alive)
	}

	waitForStatus(t, m, "node-2", memberConfirmed)
}

func TestMembershipRefutedSuspicionIsNotConfirmed(t *testing.T) {
	m := newTestMembership(100 * time.Millisecond)

	m.suspect("node-2")
	m.merge([]membershipUpdate{{Status: memberAlive, NodeID: "node-2", Incarnation: 1}})

	if status := m.status("node-2"); status != memberAlive {
		t.Fatalf("status after refutation = %s, want %s", status, memberAlive)
	}

	time.Sleep(300 * time.Millisecond)

	if status := m.status("node-2"); status != memberAlive {
		t.Errorf("status after the suspect timeout = %s, want %s", status, memberAlive)
	}
}

func TestMembershipRefutesRumorsAboutItself(t *testing.T) {
	m := newTestMembership(time.Minute)

	tests := []struct {
		rumor           membershipUpdate
		wantIncarnation uint64
	}{
		{membershipUpdate{Status: memberSuspect, NodeID: "node-1", Incarnation: 0}, 1},
		{membershipUpdate{Status: memberConfirmed, NodeID: "node-1", Incarnation: 1}, 2},
		{membershipUpdate{Status: memberSuspect, NodeID: "node-1", Incarnation: 0}, 2},
	}

	for _, tt := range tests {
		m.merge([]membershipUpdate{tt.rumor})

		if m.incarnation != tt.wantIncarnation {
			t.Errorf("incarnation after %s(%d) = %d, want %d",
				tt.rumor.Status, tt.rumor.Incarnation, m.incarnation, tt.wantIncarnation)
		}

		want := membershipUpdate{Status: memberAlive, NodeID: "node-1", Incarnation: tt.wantIncarnation}
		if updates := m.updatesFor("node-2"); !slices.Contains(updates, want) {
			t.Errorf("updates after %s(%d) = %v, want them to contain %v",
				tt.rumor.Status, tt.rumor.Incarnation, updates, want)
		}
	}
}

func TestMembershipStopsPiggybackingAfterRetransmitLimit(t *testing.T) {
	m := newTestMembership(time.Minute)
	suspicion := membershipUpdate{Status: memberSuspect, NodeID: "node-3", Incarnation: 0}

	m.merge([]membershipUpdate{suspicion})

	for transmission := 1; transmission <= m.retransmitLimit; transmission++ {
		if updates := m.updatesFor("node-2"); !slices.Contains(updates, suspicion) {
			t.Fatalf("transmission %d updates = %v, want them to contain %v", transmission, updates, suspicion)
		}
	}

	if updates := m.updatesFor("node-2"); slices.Contains(updates, suspicion) {
		t.Errorf("updates after %d transmissions = %v, want the suspicion dropped", m.retransmitLimit, updates)
	}
}

func TestMembershipAlwaysTellsNonAliveTargetAboutItself(t *testing.T) {
	m := newTestMembership(time.Minute)
	suspicion := membershipUpdate{Status: memberSuspect, NodeID: "node-2", Incarnation: 0}

	m.merge([]membershipUpdate{suspicion})

	for range m.retransmitLimit {
		m.updatesFor("node-3")
	}

	if updates := m.updatesFor("node-3"); len(updates) != 0 {
		t.Errorf("updates for an alive peer = %v, want none once the queue is drained", updates)
	}

	if updates := m.updatesFor("node-2"); !slices.Contains(updates, suspicion) {
		t.Errorf("updates for the suspected peer = %v, want them to contain %v", updates, suspicion)
	}
}
