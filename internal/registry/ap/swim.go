package ap

import (
	"context"
	"log/slog"
	"slices"
	"sync"
	"time"

	appb "tpiii.local/daniel-tpiii/gen/ap"
	"tpiii.local/daniel-tpiii/internal/registry"
)

type SWIMConfig struct {
	ProbeInterval        time.Duration
	PingTimeout          time.Duration
	PingReqTimeout       time.Duration
	PingReqCount         int
	SuspectTimeout       time.Duration
	RetransmitMultiplier int
}

type failureDetector struct {
	nodeID     string
	peers      []registry.ClusterMember
	config     SWIMConfig
	membership *membership
	transport  *udpTransport
	logger     *slog.Logger

	mu          sync.Mutex
	nextProbeID uint64
	pending     map[uint64]func()
}

func newFailureDetector(
	nodeID string,
	peers []registry.ClusterMember,
	config SWIMConfig,
	membership *membership,
	transport *udpTransport,
	logger *slog.Logger,
) *failureDetector {
	return &failureDetector{
		nodeID:     nodeID,
		peers:      peers,
		config:     config,
		membership: membership,
		transport:  transport,
		logger:     logger,
		pending:    make(map[uint64]func()),
	}
}

func (d *failureDetector) run(ctx context.Context) {
	if len(d.peers) == 0 {
		return
	}

	ticker := time.NewTicker(d.config.ProbeInterval)
	defer ticker.Stop()

	var order []registry.ClusterMember

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if len(order) == 0 {
				order = slices.Clone(d.peers)
				shuffle(order)
			}

			target := order[0]
			order = order[1:]

			d.probe(ctx, target)
		}
	}
}

func (d *failureDetector) probe(ctx context.Context, target registry.ClusterMember) {
	if d.ping(ctx, target) {
		return
	}

	if ctx.Err() != nil || d.membership.status(target.ID) != memberAlive {
		return
	}

	if d.pingReq(ctx, target) {
		return
	}

	if ctx.Err() != nil {
		return
	}

	d.logger.Debug("Probe failed", "peer", target.ID)
	d.membership.suspect(target.ID)
}

func (d *failureDetector) ping(ctx context.Context, target registry.ClusterMember) bool {
	probeID, acked := d.await()
	defer d.forget(probeID)

	d.send(target, d.pingPacket(probeID, target.ID))

	return wait(ctx, acked, d.config.PingTimeout)
}

func (d *failureDetector) pingReq(ctx context.Context, target registry.ClusterMember) bool {
	intermediaries := slices.DeleteFunc(d.membership.alivePeers(), func(peer registry.ClusterMember) bool {
		return peer.ID == target.ID
	})
	if len(intermediaries) == 0 {
		return false
	}

	shuffle(intermediaries)
	intermediaries = intermediaries[:min(d.config.PingReqCount, len(intermediaries))]

	probeID, acked := d.await()
	defer d.forget(probeID)

	for _, intermediary := range intermediaries {
		d.send(intermediary, d.pingReqPacket(probeID, target.ID, intermediary.ID))
	}

	return wait(ctx, acked, d.config.PingReqTimeout)
}

func (d *failureDetector) handlePing(senderID string, ping *appb.Ping) {
	d.membership.merge(membershipUpdatesFromProto(ping.GetMembershipUpdates()))

	sender, ok := d.membership.peer(senderID)
	if !ok {
		return
	}

	d.send(sender, d.ackPacket(ping.GetProbeId(), d.nodeID, sender.ID))
}

func (d *failureDetector) handlePingReq(senderID string, pingReq *appb.PingReq) {
	d.membership.merge(membershipUpdatesFromProto(pingReq.GetMembershipUpdates()))

	requester, ok := d.membership.peer(senderID)
	if !ok {
		return
	}

	target, ok := d.membership.peer(pingReq.GetTargetNodeId())
	if !ok {
		return
	}

	requestedProbeID := pingReq.GetProbeId()
	probeID := d.register(func() {
		d.send(requester, d.ackPacket(requestedProbeID, target.ID, requester.ID))
	})

	time.AfterFunc(d.config.PingTimeout, func() {
		d.forget(probeID)
	})

	d.send(target, d.pingPacket(probeID, target.ID))
}

func (d *failureDetector) handleAck(ack *appb.Ack) {
	d.membership.merge(membershipUpdatesFromProto(ack.GetMembershipUpdates()))

	d.mu.Lock()
	onAck, ok := d.pending[ack.GetProbeId()]
	delete(d.pending, ack.GetProbeId())
	d.mu.Unlock()

	if ok {
		onAck()
	}
}

func (d *failureDetector) await() (uint64, <-chan struct{}) {
	acked := make(chan struct{})
	probeID := d.register(func() {
		close(acked)
	})

	return probeID, acked
}

func (d *failureDetector) register(onAck func()) uint64 {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.nextProbeID++
	d.pending[d.nextProbeID] = onAck

	return d.nextProbeID
}

func (d *failureDetector) forget(probeID uint64) {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.pending, probeID)
}

func (d *failureDetector) send(peer registry.ClusterMember, packet *appb.Packet) {
	if err := d.transport.send(peer.PeerAddress, packet); err != nil {
		d.logger.Debug("Failed to send SWIM packet", "peer", peer.ID, "error", err)
	}
}

func (d *failureDetector) pingPacket(probeID uint64, targetID string) *appb.Packet {
	return &appb.Packet{
		SenderId: d.nodeID,
		Payload: &appb.Packet_Ping{
			Ping: &appb.Ping{
				ProbeId:           probeID,
				MembershipUpdates: membershipUpdatesToProto(d.membership.updatesFor(targetID)),
			},
		},
	}
}

func (d *failureDetector) pingReqPacket(probeID uint64, targetID, intermediaryID string) *appb.Packet {
	return &appb.Packet{
		SenderId: d.nodeID,
		Payload: &appb.Packet_PingReq{
			PingReq: &appb.PingReq{
				ProbeId:           probeID,
				TargetNodeId:      targetID,
				MembershipUpdates: membershipUpdatesToProto(d.membership.updatesFor(intermediaryID)),
			},
		},
	}
}

func (d *failureDetector) ackPacket(probeID uint64, responderID, destinationID string) *appb.Packet {
	return &appb.Packet{
		SenderId: d.nodeID,
		Payload: &appb.Packet_Ack{
			Ack: &appb.Ack{
				ProbeId:           probeID,
				ResponderId:       responderID,
				MembershipUpdates: membershipUpdatesToProto(d.membership.updatesFor(destinationID)),
			},
		},
	}
}

func wait(ctx context.Context, acked <-chan struct{}, timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-acked:
		return true
	case <-timer.C:
		return false
	case <-ctx.Done():
		return false
	}
}
