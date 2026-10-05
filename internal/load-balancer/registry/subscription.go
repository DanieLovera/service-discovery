package registry

import (
	"google.golang.org/grpc"

	registrypb "tpiii.local/daniel-tpiii/gen/registry"
	loadbalancer "tpiii.local/daniel-tpiii/internal/load-balancer"
)

type WatchEvent struct {
	Instance loadbalancer.ServiceInstance
	Healthy  bool
}

type Subscription struct {
	stream grpc.ServerStreamingClient[registrypb.WatchResponse]
}

func (s *Subscription) WaitEvent() (WatchEvent, error) {
	response, err := s.stream.Recv()
	if err != nil {
		return WatchEvent{}, err
	}

	instance := response.GetEvent().GetInstance()

	return WatchEvent{
		Instance: serviceInstanceFromProto(instance),
		Healthy:  isHealthy(instance),
	}, nil
}

func isHealthy(instance *registrypb.ServiceInstance) bool {
	return instance.GetStatus() == registrypb.ServiceInstanceStatus_SERVICE_INSTANCE_STATUS_HEALTHY
}
