package routing

import (
	"errors"

	loadbalancer "tpiii.local/daniel-tpiii/internal/load-balancer"
)

var (
	ErrUnknownService = errors.New("unknown service")
	ErrNoInstances    = errors.New("no healthy instances available")
)

type Upstream struct {
	Pool     *loadbalancer.ServicePool
	Strategy Strategy
}

type Manager struct {
	upstreams map[string]Upstream
}

func NewManager(upstreams map[string]Upstream) *Manager {
	return &Manager{upstreams: upstreams}
}

func (m *Manager) Pick(serviceName string) (loadbalancer.ServiceInstance, func(), error) {
	upstream, ok := m.upstreams[serviceName]
	if !ok {
		return loadbalancer.ServiceInstance{}, nil, ErrUnknownService
	}

	instance, done, ok := upstream.Strategy.Pick(upstream.Pool.Instances())
	if !ok {
		return loadbalancer.ServiceInstance{}, nil, ErrNoInstances
	}

	return instance, done, nil
}
