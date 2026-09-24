package grpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	registrypb "tpiii.local/daniel-tpiii/gen/registry"
	"tpiii.local/daniel-tpiii/internal/registry"
)

func serviceInstanceToProto(instance registry.ServiceInstance) *registrypb.ServiceInstance {
	return &registrypb.ServiceInstance{
		ServiceName:         instance.ID.ServiceName,
		InstanceId:          instance.ID.InstanceID,
		Address:             instance.Address,
		Weight:              int32(instance.Weight),
		HeartbeatIntervalMs: instance.HeartbeatInterval.Milliseconds(),
		Status:              statusToProto(instance.Status),
	}
}

func statusToProto(value registry.ServiceInstanceStatus) registrypb.ServiceInstanceStatus {
	switch value {
	case registry.ServiceInstanceStatusHealthy:
		return registrypb.ServiceInstanceStatus_SERVICE_INSTANCE_STATUS_HEALTHY
	case registry.ServiceInstanceStatusExpired:
		return registrypb.ServiceInstanceStatus_SERVICE_INSTANCE_STATUS_EXPIRED
	case registry.ServiceInstanceStatusDeleted:
		return registrypb.ServiceInstanceStatus_SERVICE_INSTANCE_STATUS_DELETED
	default:
		return registrypb.ServiceInstanceStatus_SERVICE_INSTANCE_STATUS_UNSPECIFIED
	}
}

func grpcError(err error) error {
	switch {
	case errors.Is(err, registry.ErrServiceInstanceAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, registry.ErrServiceInstanceNotFound), errors.Is(err, registry.ErrServiceNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, registry.ErrServiceInstanceDeleted):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
