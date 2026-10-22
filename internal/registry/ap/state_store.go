package ap

import (
	"sync"

	"tpiii.local/daniel-tpiii/internal/registry"
)

type stateStore struct {
	nodeID string

	mu       sync.RWMutex
	clock    uint64
	states   map[registry.ServiceInstanceID]instanceState
	handlers []registry.StateChangeHandler
}

func newStateStore(nodeID string) *stateStore {
	return &stateStore{
		nodeID: nodeID,
		states: make(map[registry.ServiceInstanceID]instanceState),
	}
}

func (s *stateStore) Register(params registry.RegisterParams) (instanceState, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, exists := s.states[params.ID]
	if exists && current.Instance.Status != registry.ServiceInstanceStatusDeleted {
		return instanceState{}, false, registry.ErrServiceInstanceAlreadyExists
	}

	instance := registry.ServiceInstance{
		ID:                params.ID,
		Address:           params.Address,
		Weight:            params.Weight,
		HeartbeatInterval: params.HeartbeatInterval,
		Status:            registry.ServiceInstanceStatusHealthy,
	}

	return s.apply(instance, registry.StateChangeRegistered), true, nil
}

func (s *stateStore) Update(params registry.UpdateParams) (instanceState, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, err := s.active(params.ID)
	if err != nil {
		return instanceState{}, false, err
	}

	updated := current
	updated.Address = params.Address
	updated.Weight = params.Weight
	updated.HeartbeatInterval = params.HeartbeatInterval

	if updated == current {
		return instanceState{}, false, nil
	}

	return s.apply(updated, registry.StateChangeUpdated), true, nil
}

func (s *stateStore) Deregister(id registry.ServiceInstanceID) (instanceState, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, exists := s.states[id]
	if !exists {
		return instanceState{}, false, registry.ErrServiceInstanceNotFound
	}

	if current.Instance.Status == registry.ServiceInstanceStatusDeleted {
		return instanceState{}, false, nil
	}

	deleted := current.Instance
	deleted.Status = registry.ServiceInstanceStatusDeleted

	return s.apply(deleted, registry.StateChangeDeleted), true, nil
}

func (s *stateStore) Expire(id registry.ServiceInstanceID) (instanceState, bool, error) {
	return s.setStatus(id, registry.ServiceInstanceStatusExpired, registry.StateChangeExpired)
}

func (s *stateStore) Recover(id registry.ServiceInstanceID) (instanceState, bool, error) {
	return s.setStatus(id, registry.ServiceInstanceStatusHealthy, registry.StateChangeHealthy)
}

func (s *stateStore) Merge(remote instanceState) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clock = max(s.clock, remote.Version.LamportTimestamp)

	current, exists := s.states[remote.Instance.ID]
	if exists && !remote.Version.newerThan(current.Version) {
		return false
	}

	s.states[remote.Instance.ID] = remote

	if exists && current.Instance == remote.Instance {
		return true
	}

	s.notify(registry.StateChange{
		Type:     mergeChangeType(current.Instance, exists, remote.Instance),
		Instance: remote.Instance,
	})

	return true
}

func (s *stateStore) Lookup(id registry.ServiceInstanceID) ([]registry.ServiceInstance, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if id.InstanceID != "" {
		state, exists := s.states[id]
		if !exists {
			return nil, registry.ErrServiceInstanceNotFound
		}

		return []registry.ServiceInstance{state.Instance}, nil
	}

	var instances []registry.ServiceInstance
	for _, state := range s.states {
		if state.Instance.ID.ServiceName == id.ServiceName {
			instances = append(instances, state.Instance)
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

func (s *stateStore) setStatus(
	id registry.ServiceInstanceID,
	status registry.ServiceInstanceStatus,
	changeType registry.StateChangeType,
) (instanceState, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, err := s.active(id)
	if err != nil {
		return instanceState{}, false, err
	}

	if current.Status == status {
		return instanceState{}, false, nil
	}

	current.Status = status

	return s.apply(current, changeType), true, nil
}

func (s *stateStore) active(id registry.ServiceInstanceID) (registry.ServiceInstance, error) {
	state, exists := s.states[id]
	if !exists {
		return registry.ServiceInstance{}, registry.ErrServiceInstanceNotFound
	}

	if state.Instance.Status == registry.ServiceInstanceStatusDeleted {
		return registry.ServiceInstance{}, registry.ErrServiceInstanceDeleted
	}

	return state.Instance, nil
}

func (s *stateStore) apply(instance registry.ServiceInstance, changeType registry.StateChangeType) instanceState {
	s.clock++

	state := instanceState{
		Instance: instance,
		Version:  version{LamportTimestamp: s.clock, NodeID: s.nodeID},
	}

	s.states[instance.ID] = state
	s.notify(registry.StateChange{Type: changeType, Instance: instance})

	return state
}

func (s *stateStore) notify(change registry.StateChange) {
	for _, handler := range s.handlers {
		handler(change)
	}
}

func mergeChangeType(
	previous registry.ServiceInstance,
	exists bool,
	next registry.ServiceInstance,
) registry.StateChangeType {
	switch next.Status {
	case registry.ServiceInstanceStatusDeleted:
		return registry.StateChangeDeleted
	case registry.ServiceInstanceStatusExpired:
		return registry.StateChangeExpired
	}

	switch {
	case !exists || previous.Status == registry.ServiceInstanceStatusDeleted:
		return registry.StateChangeRegistered
	case previous.Status == registry.ServiceInstanceStatusExpired:
		return registry.StateChangeHealthy
	default:
		return registry.StateChangeUpdated
	}
}
