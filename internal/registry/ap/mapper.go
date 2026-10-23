package ap

import (
	"fmt"
	"time"

	gossippb "tpiii.local/daniel-tpiii/gen/gossip"
	"tpiii.local/daniel-tpiii/internal/registry"
)

func instanceStateToProto(state instanceState) *gossippb.InstanceState {
	return &gossippb.InstanceState{
		ServiceName:         state.Instance.ID.ServiceName,
		InstanceId:          state.Instance.ID.InstanceID,
		Address:             state.Instance.Address,
		Weight:              int32(state.Instance.Weight),
		HeartbeatIntervalMs: state.Instance.HeartbeatInterval.Milliseconds(),
		Status:              statusToProto(state.Instance.Status),
		Version: &gossippb.Version{
			LamportTimestamp: state.Version.LamportTimestamp,
			NodeId:           state.Version.NodeID,
		},
	}
}

func instanceStateFromProto(state *gossippb.InstanceState) (instanceState, error) {
	status, err := statusFromProto(state.GetStatus())
	if err != nil {
		return instanceState{}, err
	}

	return instanceState{
		Instance: registry.ServiceInstance{
			ID: registry.ServiceInstanceID{
				ServiceName: state.GetServiceName(),
				InstanceID:  state.GetInstanceId(),
			},
			Address:           state.GetAddress(),
			Weight:            int(state.GetWeight()),
			HeartbeatInterval: time.Duration(state.GetHeartbeatIntervalMs()) * time.Millisecond,
			Status:            status,
		},
		Version: version{
			LamportTimestamp: state.GetVersion().GetLamportTimestamp(),
			NodeID:           state.GetVersion().GetNodeId(),
		},
	}, nil
}

func statusToProto(value registry.ServiceInstanceStatus) gossippb.InstanceStatus {
	switch value {
	case registry.ServiceInstanceStatusHealthy:
		return gossippb.InstanceStatus_INSTANCE_STATUS_HEALTHY
	case registry.ServiceInstanceStatusExpired:
		return gossippb.InstanceStatus_INSTANCE_STATUS_EXPIRED
	case registry.ServiceInstanceStatusDeleted:
		return gossippb.InstanceStatus_INSTANCE_STATUS_DELETED
	default:
		return gossippb.InstanceStatus_INSTANCE_STATUS_UNSPECIFIED
	}
}

func statusFromProto(value gossippb.InstanceStatus) (registry.ServiceInstanceStatus, error) {
	switch value {
	case gossippb.InstanceStatus_INSTANCE_STATUS_HEALTHY:
		return registry.ServiceInstanceStatusHealthy, nil
	case gossippb.InstanceStatus_INSTANCE_STATUS_EXPIRED:
		return registry.ServiceInstanceStatusExpired, nil
	case gossippb.InstanceStatus_INSTANCE_STATUS_DELETED:
		return registry.ServiceInstanceStatusDeleted, nil
	default:
		return "", fmt.Errorf("unknown instance status %v", value)
	}
}
