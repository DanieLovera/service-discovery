package routing

import (
	"sync"

	loadbalancer "tpiii.local/daniel-tpiii/internal/load-balancer"
)

type LeastConnections struct {
	mu     sync.Mutex
	active map[string]int
	offset int
}

func NewLeastConnections() *LeastConnections {
	return &LeastConnections{
		active: make(map[string]int),
	}
}

func (l *LeastConnections) Pick(instances []loadbalancer.ServiceInstance) (loadbalancer.ServiceInstance, func(), bool) {
	if len(instances) == 0 {
		return loadbalancer.ServiceInstance{}, noop, false
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Rotating the starting point spreads ties instead of always favoring the first instance
	l.offset = (l.offset + 1) % len(instances)

	selected := instances[l.offset]
	for i := 1; i < len(instances); i++ {
		candidate := instances[(l.offset+i)%len(instances)]
		if l.active[candidate.ID] < l.active[selected.ID] {
			selected = candidate
		}
	}

	l.active[selected.ID]++

	return selected, func() { l.release(selected.ID) }, true
}

func (l *LeastConnections) release(instanceID string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.active[instanceID]--
	if l.active[instanceID] <= 0 {
		delete(l.active, instanceID)
	}
}
