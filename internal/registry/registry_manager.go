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

func (r *RegistryManager) Register(ctx context.Context, params RegisterParams) error {
	return r.backend.Register(ctx, params)
}

func (r *RegistryManager) Update(ctx context.Context, params UpdateParams) error {
	return r.backend.Update(ctx, params)
}

func (r *RegistryManager) Deregister(ctx context.Context, id ServiceInstanceID) error {
	return r.backend.Deregister(ctx, id)
}

func (r *RegistryManager) Lookup(ctx context.Context, id ServiceInstanceID) ([]ServiceInstance, error) {
	instances, err := r.backend.Lookup(ctx, id)
	if err != nil {
		return nil, err
	}

	return healthyInstances(instances), nil
}

func healthyInstances(instances []ServiceInstance) []ServiceInstance {
	healthy := make([]ServiceInstance, 0, len(instances))
	for _, instance := range instances {
		if instance.Status == ServiceInstanceStatusHealthy {
			healthy = append(healthy, instance)
		}
	}

	return healthy
}
