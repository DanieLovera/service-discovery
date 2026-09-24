package registry

import (
	"context"
	"time"
)

type RegisterParams struct {
	ID                ServiceInstanceID
	Address           string
	Weight            int
	HeartbeatInterval time.Duration
}

type UpdateParams struct {
	ID                ServiceInstanceID
	Address           string
	Weight            int
	HeartbeatInterval time.Duration
}

type RegistryManager struct {
	backend Backend
}

func NewRegistryManager(backend Backend) *RegistryManager {
	return &RegistryManager{backend: backend}
}

func (m *RegistryManager) Register(ctx context.Context, params RegisterParams) error {
	return m.backend.Register(ctx, params)
}

func (m *RegistryManager) Update(ctx context.Context, params UpdateParams) error {
	return m.backend.Update(ctx, params)
}

func (m *RegistryManager) Deregister(ctx context.Context, id ServiceInstanceID) error {
	return m.backend.Deregister(ctx, id)
}

func (m *RegistryManager) Lookup(ctx context.Context, id ServiceInstanceID) ([]ServiceInstance, error) {
	instances, err := m.backend.Lookup(ctx, id)
	if err != nil {
		return nil, err
	}

	healthy := make([]ServiceInstance, 0, len(instances))
	for _, instance := range instances {
		if instance.Status == ServiceInstanceStatusHealthy {
			healthy = append(healthy, instance)
		}
	}

	return healthy, nil
}
