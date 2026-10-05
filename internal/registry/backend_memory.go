package registry

import (
	"context"
	"sync"
)

type BackendMemory struct {
	mu        sync.RWMutex
	instances map[ServiceInstanceID]ServiceInstance
	handlers  []StateChangeHandler
}

func NewBackendMemory() *BackendMemory {
	return &BackendMemory{
		instances: make(map[ServiceInstanceID]ServiceInstance),
	}
}

func (b *BackendMemory) Register(ctx context.Context, params RegisterParams) error {
	b.mu.Lock()

	if _, exists := b.instances[params.ID]; exists {
		b.mu.Unlock()
		return ErrServiceInstanceAlreadyExists
	}

	instance := ServiceInstance{
		ID:                params.ID,
		Address:           params.Address,
		Weight:            params.Weight,
		HeartbeatInterval: params.HeartbeatInterval,
		Status:            ServiceInstanceStatusHealthy,
	}

	b.instances[params.ID] = instance
	b.mu.Unlock()

	b.notify(StateChange{
		Type:     StateChangeRegistered,
		Instance: instance,
	})

	return nil
}

func (b *BackendMemory) Update(ctx context.Context, params UpdateParams) error {
	b.mu.Lock()

	instance, exists := b.instances[params.ID]
	if !exists {
		b.mu.Unlock()
		return ErrServiceInstanceNotFound
	}

	if instance.Status == ServiceInstanceStatusDeleted {
		b.mu.Unlock()
		return ErrServiceInstanceDeleted
	}

	instance.Address = params.Address
	instance.Weight = params.Weight
	instance.HeartbeatInterval = params.HeartbeatInterval

	b.instances[params.ID] = instance
	b.mu.Unlock()

	b.notify(StateChange{
		Type:     StateChangeUpdated,
		Instance: instance,
	})

	return nil
}

func (b *BackendMemory) Deregister(ctx context.Context, id ServiceInstanceID) error {
	b.mu.Lock()

	instance, exists := b.instances[id]
	if !exists {
		b.mu.Unlock()
		return ErrServiceInstanceNotFound
	}

	if instance.Status == ServiceInstanceStatusDeleted {
		b.mu.Unlock()
		return ErrServiceInstanceDeleted
	}

	instance.Status = ServiceInstanceStatusDeleted
	b.instances[id] = instance
	b.mu.Unlock()

	b.notify(StateChange{
		Type:     StateChangeDeleted,
		Instance: instance,
	})

	return nil
}

func (b *BackendMemory) Lookup(
	ctx context.Context,
	id ServiceInstanceID,
) ([]ServiceInstance, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if id.InstanceID != "" {
		instance, exists := b.instances[id]
		if !exists {
			return nil, ErrServiceInstanceNotFound
		}

		return []ServiceInstance{instance}, nil
	}

	instances := make([]ServiceInstance, 0)

	for _, instance := range b.instances {
		if instance.ID.ServiceName == id.ServiceName {
			instances = append(instances, instance)
		}
	}

	if len(instances) == 0 {
		return nil, ErrServiceNotFound
	}

	return instances, nil
}

func (b *BackendMemory) Expire(ctx context.Context, id ServiceInstanceID) error {
	b.mu.Lock()

	instance, exists := b.instances[id]
	if !exists {
		b.mu.Unlock()
		return ErrServiceInstanceNotFound
	}

	if instance.Status == ServiceInstanceStatusDeleted {
		b.mu.Unlock()
		return ErrServiceInstanceDeleted
	}

	instance.Status = ServiceInstanceStatusExpired
	b.instances[id] = instance
	b.mu.Unlock()

	b.notify(StateChange{
		Type:     StateChangeExpired,
		Instance: instance,
	})

	return nil
}

func (b *BackendMemory) Recover(ctx context.Context, id ServiceInstanceID) error {
	b.mu.Lock()

	instance, exists := b.instances[id]
	if !exists {
		b.mu.Unlock()
		return ErrServiceInstanceNotFound
	}

	if instance.Status == ServiceInstanceStatusDeleted {
		b.mu.Unlock()
		return ErrServiceInstanceDeleted
	}

	instance.Status = ServiceInstanceStatusHealthy
	b.instances[id] = instance
	b.mu.Unlock()

	b.notify(StateChange{
		Type:     StateChangeHealthy,
		Instance: instance,
	})

	return nil
}

func (b *BackendMemory) RegisterStateChangeHandler(handler StateChangeHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers = append(b.handlers, handler)
}

func (b *BackendMemory) notify(change StateChange) {
	b.mu.RLock()
	handlers := append([]StateChangeHandler(nil), b.handlers...)
	b.mu.RUnlock()

	for _, handler := range handlers {
		handler(change)
	}
}
