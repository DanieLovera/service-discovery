package cp

import (
	"maps"
	"slices"
	"sync"

	"tpiii.local/daniel-tpiii/internal/registry"
)

type stateStore struct {
	mu        sync.RWMutex
	instances map[registry.ServiceInstanceID]registry.ServiceInstance
	handlers  []registry.StateChangeHandler
}

func newStateStore() *stateStore {
	return &stateStore{
		instances: make(map[registry.ServiceInstanceID]registry.ServiceInstance),
	}
}

func (s *stateStore) Register(params registry.RegisterParams) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, exists := s.instances[params.ID]
	if exists && current.Status != registry.ServiceInstanceStatusDeleted {
		return registry.ErrServiceInstanceAlreadyExists
	}

	instance := registry.ServiceInstance{
		ID:                params.ID,
		Address:           params.Address,
		Weight:            params.Weight,
		HeartbeatInterval: params.HeartbeatInterval,
		Status:            registry.ServiceInstanceStatusHealthy,
	}

	s.instances[instance.ID] = instance
	s.notify(registry.StateChange{Type: registry.StateChangeRegistered, Instance: instance})

	return nil
}

func (s *stateStore) Update(params registry.UpdateParams) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, err := s.active(params.ID)
	if err != nil {
		return err
	}

	updated := current
	updated.Address = params.Address
	updated.Weight = params.Weight
	updated.HeartbeatInterval = params.HeartbeatInterval

	if updated == current {
		return nil
	}

	s.instances[updated.ID] = updated
	s.notify(registry.StateChange{Type: registry.StateChangeUpdated, Instance: updated})

	return nil
}

func (s *stateStore) Deregister(id registry.ServiceInstanceID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, exists := s.instances[id]
	if !exists {
		return registry.ErrServiceInstanceNotFound
	}

	if current.Status == registry.ServiceInstanceStatusDeleted {
		return nil
	}

	current.Status = registry.ServiceInstanceStatusDeleted
	s.instances[id] = current
	s.notify(registry.StateChange{Type: registry.StateChangeDeleted, Instance: current})

	return nil
}

func (s *stateStore) Expire(id registry.ServiceInstanceID) error {
	return s.setStatus(id, registry.ServiceInstanceStatusExpired, registry.StateChangeExpired)
}

func (s *stateStore) Recover(id registry.ServiceInstanceID) error {
	return s.setStatus(id, registry.ServiceInstanceStatusHealthy, registry.StateChangeHealthy)
}

func (s *stateStore) Lookup(id registry.ServiceInstanceID) ([]registry.ServiceInstance, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if id.InstanceID != "" {
		instance, exists := s.instances[id]
		if !exists {
			return nil, registry.ErrServiceInstanceNotFound
		}

		return []registry.ServiceInstance{instance}, nil
	}

	var instances []registry.ServiceInstance
	for _, instance := range s.instances {
		if instance.ID.ServiceName == id.ServiceName {
			instances = append(instances, instance)
		}
	}

	if len(instances) == 0 {
		return nil, registry.ErrServiceNotFound
	}

	return instances, nil
}

func (s *stateStore) RegisterStateChangeHandler(handler registry.StateChangeHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.handlers = append(s.handlers, handler)
}

func (s *stateStore) Instances() []registry.ServiceInstance {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return slices.Collect(maps.Values(s.instances))
}

func (s *stateStore) Replace(instances []registry.ServiceInstance) {
	restored := make(map[registry.ServiceInstanceID]registry.ServiceInstance, len(instances))
	for _, instance := range instances {
		restored[instance.ID] = instance
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.instances = restored
	s.notify(registry.StateChange{Type: registry.StateChangeRestored, Instances: instances})
}

func (s *stateStore) setStatus(
	id registry.ServiceInstanceID,
	status registry.ServiceInstanceStatus,
	changeType registry.StateChangeType,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, err := s.active(id)
	if err != nil {
		return err
	}

	if current.Status == status {
		return nil
	}

	current.Status = status
	s.instances[id] = current
	s.notify(registry.StateChange{Type: changeType, Instance: current})

	return nil
}

func (s *stateStore) active(id registry.ServiceInstanceID) (registry.ServiceInstance, error) {
	instance, exists := s.instances[id]
	if !exists {
		return registry.ServiceInstance{}, registry.ErrServiceInstanceNotFound
	}

	if instance.Status == registry.ServiceInstanceStatusDeleted {
		return registry.ServiceInstance{}, registry.ErrServiceInstanceDeleted
	}

	return instance, nil
}

func (s *stateStore) notify(change registry.StateChange) {
	for _, handler := range s.handlers {
		handler(change)
	}
}
