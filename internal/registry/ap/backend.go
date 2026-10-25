package ap

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"sync"
	"time"

	"google.golang.org/grpc"

	appb "tpiii.local/daniel-tpiii/gen/ap"
	"tpiii.local/daniel-tpiii/internal/registry"
)

// Compile-time check that *Backend implements registry.Backend
var _ registry.Backend = (*Backend)(nil)

type Config struct {
	NodeID              string
	BindAddress         string
	Members             []registry.ClusterMember
	GossipFanout        int
	AntiEntropyInterval time.Duration
	SWIM                SWIMConfig
	Logger              *slog.Logger
}

type Backend struct {
	store      *stateStore
	transport  *udpTransport
	gossiper   *gossiper
	reconciler *reconciler
	detector   *failureDetector

	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewBackend(config Config) *Backend {
	peers := slices.DeleteFunc(slices.Clone(config.Members), func(member registry.ClusterMember) bool {
		return member.ID == config.NodeID
	})

	store := newStateStore(config.NodeID)
	transport := newUDPTransport(config.BindAddress, config.Logger)
	membership := newMembership(
		config.NodeID,
		peers,
		config.SWIM.SuspectTimeout,
		config.SWIM.RetransmitMultiplier,
		config.Logger,
	)

	return &Backend{
		store:      store,
		transport:  transport,
		gossiper:   newGossiper(config.NodeID, membership, config.GossipFanout, store, transport, config.Logger),
		reconciler: newReconciler(config.NodeID, membership, config.AntiEntropyInterval, store, config.Logger),
		detector:   newFailureDetector(config.NodeID, peers, config.SWIM, membership, transport, config.Logger),
	}
}

func (b *Backend) RegisterGRPC(registrar grpc.ServiceRegistrar) {
	appb.RegisterAntiEntropyServer(registrar, &antiEntropyServer{reconciler: b.reconciler})
}

func (b *Backend) Start(ctx context.Context) error {
	if err := b.transport.listen(b.handle); err != nil {
		return err
	}

	b.reconciler.initialSync(ctx)

	ctx, b.cancel = context.WithCancel(ctx)
	b.wg.Go(func() {
		b.reconciler.run(ctx)
	})
	b.wg.Go(func() {
		b.detector.run(ctx)
	})

	return nil
}

func (b *Backend) Shutdown(ctx context.Context) error {
	if b.cancel != nil {
		b.cancel()
	}

	b.wg.Wait()

	return errors.Join(b.reconciler.close(), b.transport.close())
}

func (b *Backend) Register(ctx context.Context, params registry.RegisterParams) error {
	return b.publish(b.store.Register(params))
}

func (b *Backend) Update(ctx context.Context, params registry.UpdateParams) error {
	return b.publish(b.store.Update(params))
}

func (b *Backend) Deregister(ctx context.Context, id registry.ServiceInstanceID) error {
	return b.publish(b.store.Deregister(id))
}

func (b *Backend) Lookup(ctx context.Context, id registry.ServiceInstanceID) ([]registry.ServiceInstance, error) {
	return b.store.Lookup(id)
}

func (b *Backend) Expire(ctx context.Context, id registry.ServiceInstanceID) error {
	return b.publish(b.store.Expire(id))
}

func (b *Backend) Recover(ctx context.Context, id registry.ServiceInstanceID) error {
	return b.publish(b.store.Recover(id))
}

func (b *Backend) RegisterStateChangeHandler(handler registry.StateChangeHandler) {
	b.store.RegisterStateChangeHandler(handler)
}

func (b *Backend) publish(state instanceState, changed bool, err error) error {
	if err != nil || !changed {
		return err
	}

	b.gossiper.broadcast(state, "")

	return nil
}

func (b *Backend) handle(packet *appb.Packet) {
	switch payload := packet.GetPayload().(type) {
	case *appb.Packet_StateUpdate:
		b.gossiper.receive(packet.GetSenderId(), payload.StateUpdate)
	case *appb.Packet_Ping:
		b.detector.handlePing(packet.GetSenderId(), payload.Ping)
	case *appb.Packet_PingReq:
		b.detector.handlePingReq(packet.GetSenderId(), payload.PingReq)
	case *appb.Packet_Ack:
		b.detector.handleAck(payload.Ack)
	}
}
