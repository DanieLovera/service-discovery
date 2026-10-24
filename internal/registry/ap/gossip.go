package ap

import (
	"log/slog"
	"math/rand/v2"
	"slices"

	appb "tpiii.local/daniel-tpiii/gen/ap"
	"tpiii.local/daniel-tpiii/internal/registry"
)

type gossiper struct {
	nodeID    string
	peers     []registry.ClusterMember
	fanout    int
	store     *stateStore
	transport *udpTransport
	logger    *slog.Logger
}

func newGossiper(
	nodeID string,
	peers []registry.ClusterMember,
	fanout int,
	store *stateStore,
	transport *udpTransport,
	logger *slog.Logger,
) *gossiper {
	return &gossiper{
		nodeID:    nodeID,
		peers:     peers,
		fanout:    fanout,
		store:     store,
		transport: transport,
		logger:    logger,
	}
}

func (g *gossiper) broadcast(state instanceState, exclude string) {
	packet := &appb.Packet{
		SenderId: g.nodeID,
		Payload: &appb.Packet_StateUpdate{
			StateUpdate: &appb.StateUpdate{
				States: []*appb.InstanceState{instanceStateToProto(state)},
			},
		},
	}

	for _, peer := range g.targets(exclude) {
		if err := g.transport.send(peer.PeerAddress, packet); err != nil {
			g.logger.Debug("Failed to gossip instance state", "peer", peer.ID, "error", err)
		}
	}
}

func (g *gossiper) receive(senderID string, update *appb.StateUpdate) {
	for _, protoState := range update.GetStates() {
		state, err := instanceStateFromProto(protoState)
		if err != nil {
			g.logger.Warn("Discarding invalid instance state", "sender", senderID, "error", err)
			continue
		}

		if g.store.Merge(state) {
			g.broadcast(state, senderID)
		}
	}
}

func (g *gossiper) targets(exclude string) []registry.ClusterMember {
	candidates := slices.DeleteFunc(slices.Clone(g.peers), func(peer registry.ClusterMember) bool {
		return peer.ID == exclude
	})

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	return candidates[:min(g.fanout, len(candidates))]
}
