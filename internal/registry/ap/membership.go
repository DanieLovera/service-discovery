package ap

import (
	"log/slog"
	"math"
	"math/rand/v2"
	"sync"
	"time"

	"tpiii.local/daniel-tpiii/internal/registry"
)

type memberStatus string

const (
	memberAlive     memberStatus = "ALIVE"
	memberSuspect   memberStatus = "SUSPECT"
	memberConfirmed memberStatus = "CONFIRM"
)

type memberState struct {
	status      memberStatus
	incarnation uint64
}

type membershipUpdate struct {
	Status      memberStatus
	NodeID      string
	Incarnation uint64
}

func (u membershipUpdate) supersedes(current memberState) bool {
	switch u.Status {
	case memberAlive:
		return u.Incarnation > current.incarnation
	case memberSuspect:
		switch current.status {
		case memberAlive:
			return u.Incarnation >= current.incarnation
		case memberSuspect:
			return u.Incarnation > current.incarnation
		default:
			return false
		}
	case memberConfirmed:
		if current.status == memberConfirmed {
			return u.Incarnation > current.incarnation
		}

		return u.Incarnation >= current.incarnation
	default:
		return false
	}
}

type queuedUpdate struct {
	update        membershipUpdate
	transmissions int
}

type membership struct {
	nodeID          string
	peers           []registry.ClusterMember
	suspectTimeout  time.Duration
	retransmitLimit int
	logger          *slog.Logger

	mu          sync.Mutex
	incarnation uint64
	states      map[string]memberState
	queue       map[string]queuedUpdate
}

func newMembership(
	nodeID string,
	peers []registry.ClusterMember,
	suspectTimeout time.Duration,
	retransmitMultiplier int,
	logger *slog.Logger,
) *membership {
	states := make(map[string]memberState, len(peers))
	for _, peer := range peers {
		states[peer.ID] = memberState{status: memberAlive}
	}

	clusterSize := len(peers) + 1

	return &membership{
		nodeID:          nodeID,
		peers:           peers,
		suspectTimeout:  suspectTimeout,
		retransmitLimit: retransmitMultiplier * int(math.Ceil(math.Log10(float64(clusterSize+1)))),
		logger:          logger,
		states:          states,
		queue:           make(map[string]queuedUpdate),
	}
}

func (m *membership) peer(nodeID string) (registry.ClusterMember, bool) {
	for _, peer := range m.peers {
		if peer.ID == nodeID {
			return peer, true
		}
	}

	return registry.ClusterMember{}, false
}

func (m *membership) alivePeers() []registry.ClusterMember {
	m.mu.Lock()
	defer m.mu.Unlock()

	var alive []registry.ClusterMember
	for _, peer := range m.peers {
		if m.states[peer.ID].status == memberAlive {
			alive = append(alive, peer)
		}
	}

	return alive
}

func (m *membership) status(nodeID string) memberStatus {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.states[nodeID].status
}

func (m *membership) suspect(nodeID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	current, exists := m.states[nodeID]
	if !exists || current.status != memberAlive {
		return
	}

	m.apply(membershipUpdate{Status: memberSuspect, NodeID: nodeID, Incarnation: current.incarnation})
}

func (m *membership) merge(updates []membershipUpdate) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, update := range updates {
		if update.NodeID == m.nodeID {
			m.refute(update)
			continue
		}

		m.apply(update)
	}
}

func (m *membership) updatesFor(targetID string) []membershipUpdate {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, targetQueued := m.queue[targetID]

	updates := make([]membershipUpdate, 0, len(m.queue)+1)
	for nodeID, queued := range m.queue {
		updates = append(updates, queued.update)

		queued.transmissions++
		if queued.transmissions >= m.retransmitLimit {
			delete(m.queue, nodeID)
		} else {
			m.queue[nodeID] = queued
		}
	}

	if target, exists := m.states[targetID]; exists && target.status != memberAlive && !targetQueued {
		updates = append(updates, membershipUpdate{Status: target.status, NodeID: targetID, Incarnation: target.incarnation})
	}

	return updates
}

func (m *membership) refute(update membershipUpdate) {
	if update.Status == memberAlive {
		return
	}

	if update.Incarnation >= m.incarnation {
		m.incarnation = update.Incarnation + 1
		m.logger.Info("Refuting membership update about self", "status", update.Status, "incarnation", m.incarnation)
	}

	m.enqueue(membershipUpdate{Status: memberAlive, NodeID: m.nodeID, Incarnation: m.incarnation})
}

func (m *membership) apply(update membershipUpdate) {
	current, exists := m.states[update.NodeID]
	if !exists || !update.supersedes(current) {
		return
	}

	m.states[update.NodeID] = memberState{status: update.Status, incarnation: update.Incarnation}
	m.enqueue(update)

	if current.status != update.Status {
		m.logger.Info("Member status changed",
			"peer", update.NodeID,
			"status", update.Status,
			"incarnation", update.Incarnation,
		)
	}

	if update.Status == memberSuspect {
		time.AfterFunc(m.suspectTimeout, func() {
			m.confirm(update.NodeID, update.Incarnation)
		})
	}
}

func (m *membership) confirm(nodeID string, incarnation uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.states[nodeID] != (memberState{status: memberSuspect, incarnation: incarnation}) {
		return
	}

	m.apply(membershipUpdate{Status: memberConfirmed, NodeID: nodeID, Incarnation: incarnation})
}

func (m *membership) enqueue(update membershipUpdate) {
	m.queue[update.NodeID] = queuedUpdate{update: update}
}

func shuffle(members []registry.ClusterMember) {
	rand.Shuffle(len(members), func(i, j int) {
		members[i], members[j] = members[j], members[i]
	})
}
