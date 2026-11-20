package registry

import (
	registrypb "tpiii.local/daniel-tpiii/gen/registry"
	loadbalancer "tpiii.local/daniel-tpiii/internal/load-balancer"
)

func serviceInstancesFromProto(instances []*registrypb.ServiceInstance) []loadbalancer.ServiceInstance {
	result := make([]loadbalancer.ServiceInstance, 0, len(instances))
	for _, instance := range instances {
		result = append(result, serviceInstanceFromProto(instance))
	}

	return result
}

func serviceInstanceFromProto(instance *registrypb.ServiceInstance) loadbalancer.ServiceInstance {
	return loadbalancer.ServiceInstance{
		ID:      instance.GetInstanceId(),
		Address: instance.GetAddress(),
		Weight:  int(instance.GetWeight()),
	}
}
