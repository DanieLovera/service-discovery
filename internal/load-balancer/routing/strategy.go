package routing

import (
	"fmt"

	loadbalancer "tpiii.local/daniel-tpiii/internal/load-balancer"
)

type Strategy interface {
	Pick(instances []loadbalancer.ServiceInstance) (loadbalancer.ServiceInstance, func(), bool)
}

func NewStrategy(name string) (Strategy, error) {
	switch name {
	case "round-robin":
		return NewRoundRobin(), nil
	case "least-connections":
		return NewLeastConnections(), nil
	case "weighted":
		return NewWeighted(), nil
	default:
		return nil, fmt.Errorf("unsupported routing strategy %q", name)
	}
}

func noop() {}
