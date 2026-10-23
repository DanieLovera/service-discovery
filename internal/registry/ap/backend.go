package ap

import (
	"context"
	"log/slog"

	gossippb "tpiii.local/daniel-tpiii/gen/gossip"
	"tpiii.local/daniel-tpiii/internal/registry"
)

// Compile-time check that *Backend implements registry.Backend
var _ registry.Backend = (*Backend)(nil)

type Config struct {
	NodeID       string
	BindAddress  string
	Members      []registry.ClusterMember
	GossipFanout int
	Logger       *slog.Logger
}

type Backend struct {
	store     *stateStore
	transport *transport
	gossiper  *gossiper
}

func NewBackend(config Config) *Backend {
	store := newStateStore(config.NodeID)
	transport := newTransport(config.BindAddress, config.Logger)

	return &Backend{
		store:     store,
		transport: transport,
		gossiper:  newGossiper(config.NodeID, config.Members, config.GossipFanout, store, transport, config.Logger),
	}
}

func (b *Backend) Start(ctx context.Context) error {
	return b.transport.listen(b.handle)
}

func (b *Backend) Shutdown(ctx context.Context) error {
	return b.transport.close()
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

func (b *Backend) handle(message *gossippb.Message) {
	switch payload := message.GetPayload().(type) {
	case *gossippb.Message_StateUpdate:
		b.gossiper.receive(message.GetSenderId(), payload.StateUpdate)
	}
}
