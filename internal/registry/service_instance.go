package registry

import "time"

type ServiceInstanceStatus string

const (
	ServiceInstanceStatusHealthy ServiceInstanceStatus = "HEALTHY"
	ServiceInstanceStatusExpired ServiceInstanceStatus = "EXPIRED"
	ServiceInstanceStatusDeleted ServiceInstanceStatus = "DELETED"
)

type ServiceInstanceID struct {
	ServiceName string
	InstanceID  string
}

type ServiceInstance struct {
	ID                ServiceInstanceID
	Address           string
	Weight            int
	HeartbeatInterval time.Duration
	Status            ServiceInstanceStatus
}
