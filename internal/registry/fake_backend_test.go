package registry

import (
	"context"
	"sync"
	"time"
)

type fakeBackend struct {
	mu        sync.Mutex
	instances map[ServiceInstanceID]ServiceInstance
	handlers  []StateChangeHandler
	lookupErr error
	expires   int
	recovers  int
}

func newFakeBackend() *fakeBackend {
	return &fakeBackend{instances: make(map[ServiceInstanceID]ServiceInstance)}
}

func (b *fakeBackend) Start(context.Context) error    { return nil }
func (b *fakeBackend) Shutdown(context.Context) error { return nil }

func (b *fakeBackend) Register(_ context.Context, params RegisterParams) error {
	instance := ServiceInstance{
		ID:                params.ID,
		Address:           params.Address,
		Weight:            params.Weight,
		HeartbeatInterval: params.HeartbeatInterval,
		Status:            ServiceInstanceStatusHealthy,
	}

	b.mu.Lock()
	b.instances[params.ID] = instance
	b.mu.Unlock()

	b.notify(StateChange{Type: StateChangeRegistered, Instance: instance})

	return nil
}

func (b *fakeBackend) Update(context.Context, UpdateParams) error { return nil }

func (b *fakeBackend) Deregister(_ context.Context, id ServiceInstanceID) error {
	return b.setStatus(id, ServiceInstanceStatusDeleted, StateChangeDeleted)
}

func (b *fakeBackend) Expire(_ context.Context, id ServiceInstanceID) error {
	b.mu.Lock()
	b.expires++
	b.mu.Unlock()

	return b.setStatus(id, ServiceInstanceStatusExpired, StateChangeExpired)
}

func (b *fakeBackend) Recover(_ context.Context, id ServiceInstanceID) error {
	b.mu.Lock()
	b.recovers++
	b.mu.Unlock()

	return b.setStatus(id, ServiceInstanceStatusHealthy, StateChangeHealthy)
}

func (b *fakeBackend) Lookup(_ context.Context, id ServiceInstanceID) ([]ServiceInstance, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.lookupErr != nil {
		return nil, b.lookupErr
	}

	if id.InstanceID != "" {
		instance, ok := b.instances[id]
		if !ok {
			return nil, ErrServiceInstanceNotFound
		}

		return []ServiceInstance{instance}, nil
	}

	var instances []ServiceInstance
	for _, instance := range b.instances {
		if instance.ID.ServiceName == id.ServiceName && instance.Status == ServiceInstanceStatusHealthy {
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

func (b *fakeBackend) RegisterStateChangeHandler(handler StateChangeHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers = append(b.handlers, handler)
}

func (b *fakeBackend) restore(instances []ServiceInstance) {
	b.mu.Lock()
	b.instances = make(map[ServiceInstanceID]ServiceInstance, len(instances))
	for _, instance := range instances {
		b.instances[instance.ID] = instance
	}
	b.mu.Unlock()

	b.notify(StateChange{Type: StateChangeRestored, Instances: instances})
}

func (b *fakeBackend) status(id ServiceInstanceID) ServiceInstanceStatus {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.instances[id].Status
}

func (b *fakeBackend) calls() (expires, recovers int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.expires, b.recovers
}

func (b *fakeBackend) setStatus(id ServiceInstanceID, status ServiceInstanceStatus, changeType StateChangeType) error {
	b.mu.Lock()

	instance, ok := b.instances[id]
	switch {
	case !ok:
		b.mu.Unlock()
		return ErrServiceInstanceNotFound
	case instance.Status == ServiceInstanceStatusDeleted:
		b.mu.Unlock()
		return ErrServiceInstanceDeleted
	case instance.Status == status:
		b.mu.Unlock()
		return nil
	}

	instance.Status = status
	b.instances[id] = instance
	b.mu.Unlock()

	b.notify(StateChange{Type: changeType, Instance: instance})

	return nil
}

func (b *fakeBackend) notify(change StateChange) {
	b.mu.Lock()
	handlers := b.handlers
	b.mu.Unlock()

	for _, handler := range handlers {
		handler(change)
	}
}

func registerInstance(b *fakeBackend, service, instanceID string, heartbeatInterval time.Duration) ServiceInstanceID {
	id := ServiceInstanceID{ServiceName: service, InstanceID: instanceID}

	_ = b.Register(context.Background(), RegisterParams{
		ID:                id,
		Address:           instanceID + ":8080",
		Weight:            1,
		HeartbeatInterval: heartbeatInterval,
	})

	return id
}
