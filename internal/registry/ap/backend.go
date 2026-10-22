package ap

import (
	"context"

	"tpiii.local/daniel-tpiii/internal/registry"
)

// Compile-time check that *Backend implements registry.Backend
var _ registry.Backend = (*Backend)(nil)

type Config struct {
	NodeID string
}

type Backend struct {
	store *stateStore
}

func NewBackend(config Config) *Backend {
	return &Backend{
		store: newStateStore(config.NodeID),
	}
}

func (b *Backend) Start(ctx context.Context) error {
	return nil
}

func (b *Backend) Shutdown(ctx context.Context) error {
	return nil
}

func (b *Backend) Register(ctx context.Context, params registry.RegisterParams) error {
	_, _, err := b.store.Register(params)
	return err
}

func (b *Backend) Update(ctx context.Context, params registry.UpdateParams) error {
	_, _, err := b.store.Update(params)
	return err
}

func (b *Backend) Deregister(ctx context.Context, id registry.ServiceInstanceID) error {
	_, _, err := b.store.Deregister(id)
	return err
}

func (b *Backend) Lookup(ctx context.Context, id registry.ServiceInstanceID) ([]registry.ServiceInstance, error) {
	return b.store.Lookup(id)
}

func (b *Backend) Expire(ctx context.Context, id registry.ServiceInstanceID) error {
	_, _, err := b.store.Expire(id)
	return err
}

func (b *Backend) Recover(ctx context.Context, id registry.ServiceInstanceID) error {
	_, _, err := b.store.Recover(id)
	return err
}

func (b *Backend) RegisterStateChangeHandler(handler registry.StateChangeHandler) {
	b.store.RegisterStateChangeHandler(handler)
}
