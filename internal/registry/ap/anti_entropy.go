package ap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	appb "tpiii.local/daniel-tpiii/gen/ap"
	"tpiii.local/daniel-tpiii/internal/registry"
)

type reconciler struct {
	nodeID     string
	membership *membership
	interval   time.Duration
	store      *stateStore
	logger     *slog.Logger

	mu    sync.Mutex
	conns map[string]*grpc.ClientConn
}

func newReconciler(
	nodeID string,
	membership *membership,
	interval time.Duration,
	store *stateStore,
	logger *slog.Logger,
) *reconciler {
	return &reconciler{
		nodeID:     nodeID,
		membership: membership,
		interval:   interval,
		store:      store,
		logger:     logger,
		conns:      make(map[string]*grpc.ClientConn),
	}
}

func (r *reconciler) initialSync(ctx context.Context) {
	peers := r.membership.alivePeers()
	shuffle(peers)

	for _, peer := range peers {
		if ctx.Err() != nil {
			return
		}

		if err := r.reconcile(ctx, peer); err != nil {
			r.logger.Debug("Initial anti-entropy sync failed", "peer", peer.ID, "error", err)
			continue
		}

		r.logger.Info("Initial anti-entropy sync completed", "peer", peer.ID)
		return
	}

	r.logger.Info("Initial anti-entropy sync skipped, no peer available")
}

func (r *reconciler) run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			peers := r.membership.alivePeers()
			if len(peers) == 0 {
				continue
			}

			peer := peers[rand.IntN(len(peers))]
			if err := r.reconcile(ctx, peer); err != nil {
				r.logger.Debug("Anti-entropy round aborted", "peer", peer.ID, "error", err)
			}
		}
	}
}

func (r *reconciler) reconcile(ctx context.Context, peer registry.ClusterMember) error {
	ctx, cancel := context.WithTimeout(ctx, r.interval)
	defer cancel()

	client, err := r.client(peer)
	if err != nil {
		return err
	}

	diff, err := client.Compare(ctx, &appb.Digest{
		SenderId: r.nodeID,
		Entries:  digestToProto(r.store.Digest()),
	})
	if err != nil {
		return fmt.Errorf("compare digest: %w", err)
	}

	r.mergeStates(peer.ID, diff.GetNewerStates())

	requested := instanceKeysFromProto(diff.GetRequestedInstances())
	if len(requested) == 0 {
		return nil
	}

	_, err = client.Merge(ctx, &appb.StateTransfer{
		SenderId: r.nodeID,
		States:   instanceStatesToProto(r.store.States(requested)),
	})
	if err != nil {
		return fmt.Errorf("merge states: %w", err)
	}

	return nil
}

func (r *reconciler) mergeStates(senderID string, states []*appb.InstanceState) {
	for _, protoState := range states {
		state, err := instanceStateFromProto(protoState)
		if err != nil {
			r.logger.Warn("Discarding invalid instance state", "sender", senderID, "error", err)
			continue
		}

		r.store.Merge(state)
	}
}

func (r *reconciler) client(peer registry.ClusterMember) (appb.AntiEntropyClient, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	conn, ok := r.conns[peer.ID]
	if !ok {
		var err error
		conn, err = grpc.NewClient(peer.GRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("create anti-entropy connection to %s: %w", peer.GRPCAddress, err)
		}

		r.conns[peer.ID] = conn
	}

	return appb.NewAntiEntropyClient(conn), nil
}

func (r *reconciler) close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var errs []error
	for _, conn := range r.conns {
		errs = append(errs, conn.Close())
	}

	return errors.Join(errs...)
}

type antiEntropyServer struct {
	appb.UnimplementedAntiEntropyServer

	reconciler *reconciler
}

func (s *antiEntropyServer) Compare(ctx context.Context, digest *appb.Digest) (*appb.Diff, error) {
	newer, requested := s.reconciler.store.Diff(digestFromProto(digest.GetEntries()))

	return &appb.Diff{
		SenderId:           s.reconciler.nodeID,
		NewerStates:        instanceStatesToProto(newer),
		RequestedInstances: instanceKeysToProto(requested),
	}, nil
}

func (s *antiEntropyServer) Merge(
	ctx context.Context,
	transfer *appb.StateTransfer,
) (*appb.MergeResponse, error) {
	s.reconciler.mergeStates(transfer.GetSenderId(), transfer.GetStates())

	return &appb.MergeResponse{}, nil
}
