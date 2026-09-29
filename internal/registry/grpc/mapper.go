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

func stateChangeTypeToProto(value registry.StateChangeType) registrypb.StateChangeType {
	switch value {
	case registry.StateChangeRegistered:
		return registrypb.StateChangeType_STATE_CHANGE_TYPE_REGISTERED
	case registry.StateChangeUpdated:
		return registrypb.StateChangeType_STATE_CHANGE_TYPE_UPDATED
	case registry.StateChangeHealthy:
		return registrypb.StateChangeType_STATE_CHANGE_TYPE_HEALTHY
	case registry.StateChangeExpired:
		return registrypb.StateChangeType_STATE_CHANGE_TYPE_EXPIRED
	case registry.StateChangeDeleted:
		return registrypb.StateChangeType_STATE_CHANGE_TYPE_DELETED
	default:
		return registrypb.StateChangeType_STATE_CHANGE_TYPE_UNSPECIFIED
	}
}

func watchSnapshotToProto(instances []registry.ServiceInstance) *registrypb.WatchResponse {
	protoInstances := make([]*registrypb.ServiceInstance, 0, len(instances))

	for _, instance := range instances {
		protoInstances = append(protoInstances, serviceInstanceToProto(instance))
	}

	return &registrypb.WatchResponse{
		Payload: &registrypb.WatchResponse_Snapshot{
			Snapshot: &registrypb.WatchSnapshot{
				Instances: protoInstances,
			},
		},
	}
}

func watchEventToProto(change registry.StateChange) *registrypb.WatchResponse {
	return &registrypb.WatchResponse{
		Payload: &registrypb.WatchResponse_Event{
			Event: &registrypb.WatchEvent{
				Type:     stateChangeTypeToProto(change.Type),
				Instance: serviceInstanceToProto(change.Instance),
			},
		},
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
