package routing

import (
	"sync/atomic"

	loadbalancer "tpiii.local/daniel-tpiii/internal/load-balancer"
)

type RoundRobin struct {
	next atomic.Uint64
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{}
}

func (r *RoundRobin) Pick(instances []loadbalancer.ServiceInstance) (loadbalancer.ServiceInstance, func(), bool) {
	if len(instances) == 0 {
		return loadbalancer.ServiceInstance{}, noop, false
	}

	i := (r.next.Add(1) - 1) % uint64(len(instances))
	return instances[i], noop, true
}
