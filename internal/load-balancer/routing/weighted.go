package routing

import (
	"sync"

	loadbalancer "tpiii.local/daniel-tpiii/internal/load-balancer"
)

type Weighted struct {
	mu      sync.Mutex
	current map[string]int
}

func NewWeighted() *Weighted {
	return &Weighted{
		current: make(map[string]int),
	}
}

// Smooth Weighted Round Robin algorithm implementation
func (w *Weighted) Pick(instances []loadbalancer.ServiceInstance) (loadbalancer.ServiceInstance, func(), bool) {
	if len(instances) == 0 {
		return loadbalancer.ServiceInstance{}, noop, false
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	total := 0
	selected := instances[0]

	for _, instance := range instances {
		w.current[instance.ID] += instance.Weight
		total += instance.Weight

		if w.current[instance.ID] > w.current[selected.ID] {
			selected = instance
		}
	}

	w.current[selected.ID] -= total

	if len(w.current) > len(instances) {
		w.prune(instances)
	}

	return selected, noop, true
}

func (w *Weighted) prune(instances []loadbalancer.ServiceInstance) {
	present := make(map[string]struct{}, len(instances))
	for _, instance := range instances {
		present[instance.ID] = struct{}{}
	}

	for id := range w.current {
		if _, ok := present[id]; !ok {
			delete(w.current, id)
		}
	}
}
