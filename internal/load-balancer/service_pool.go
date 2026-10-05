package loadbalancer

import (
	"slices"
	"sync"
)

type ServiceInstance struct {
	ID      string
	Address string
	Weight  int
}

type ServicePool struct {
	mu        sync.RWMutex
	instances []ServiceInstance
}

func NewServicePool() *ServicePool {
	return &ServicePool{}
}

func (p *ServicePool) Replace(instances []ServiceInstance) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.instances = slices.Clone(instances)
}

func (p *ServicePool) Upsert(instance ServiceInstance) {
	p.mu.Lock()
	defer p.mu.Unlock()

	instances := slices.Clone(p.instances)

	index := slices.IndexFunc(instances, func(i ServiceInstance) bool {
		return i.ID == instance.ID
	})
	if index >= 0 {
		instances[index] = instance
	} else {
		instances = append(instances, instance)
	}

	p.instances = instances
}

func (p *ServicePool) Remove(instanceID string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.instances = slices.DeleteFunc(slices.Clone(p.instances), func(i ServiceInstance) bool {
		return i.ID == instanceID
	})
}

// The returned slice is shared and must not be modified
func (p *ServicePool) Instances() []ServiceInstance {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.instances
}
