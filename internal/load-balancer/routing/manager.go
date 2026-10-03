package routing

import (
	"errors"

	loadbalancer "tpiii.local/daniel-tpiii/internal/load-balancer"
)

var (
	ErrUnknownService = errors.New("unknown service")
	ErrNoInstances    = errors.New("no healthy instances available")
)

type Route struct {
	Pool     *loadbalancer.ServicePool
	Strategy Strategy
}

type Manager struct {
	routes map[string]Route
}

func NewManager(routes map[string]Route) *Manager {
	return &Manager{routes: routes}
}

func (m *Manager) Pick(serviceName string) (loadbalancer.ServiceInstance, func(), error) {
	route, ok := m.routes[serviceName]
	if !ok {
		return loadbalancer.ServiceInstance{}, nil, ErrUnknownService
	}

	instance, done, ok := route.Strategy.Pick(route.Pool.Instances())
	if !ok {
		return loadbalancer.ServiceInstance{}, nil, ErrNoInstances
	}

	return instance, done, nil
}
