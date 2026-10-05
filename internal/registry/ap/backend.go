package ap

import (
	"context"

	"tpiii.local/daniel-tpiii/internal/registry"
)

// Compile-time check that *Backend implements registry.Backend
var _ registry.Backend = (*Backend)(nil)

type Backend struct{}

func NewBackend() *Backend {
	return &Backend{}
}

func (b *Backend) Register(ctx context.Context, params registry.RegisterParams) error {
	return nil
}

func (b *Backend) Update(ctx context.Context, params registry.UpdateParams) error {
	return nil
}

func (b *Backend) Deregister(ctx context.Context, id registry.ServiceInstanceID) error {
	return nil
}

func (b *Backend) Lookup(ctx context.Context, id registry.ServiceInstanceID) ([]registry.ServiceInstance, error) {
	return nil, nil
}

func (b *Backend) Expire(ctx context.Context, id registry.ServiceInstanceID) error {
	return nil
}

func (b *Backend) Recover(ctx context.Context, id registry.ServiceInstanceID) error {
	return nil
}

func (b *Backend) RegisterStateChangeHandler(handler registry.StateChangeHandler) {}
