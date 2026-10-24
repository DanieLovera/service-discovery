package ap

import (
	"fmt"
	"time"

	appb "tpiii.local/daniel-tpiii/gen/ap"
	"tpiii.local/daniel-tpiii/internal/registry"
)

func instanceStateToProto(state instanceState) *appb.InstanceState {
	return &appb.InstanceState{
		ServiceName:         state.Instance.ID.ServiceName,
		InstanceId:          state.Instance.ID.InstanceID,
		Address:             state.Instance.Address,
		Weight:              int32(state.Instance.Weight),
		HeartbeatIntervalMs: state.Instance.HeartbeatInterval.Milliseconds(),
		Status:              statusToProto(state.Instance.Status),
		Version:             versionToProto(state.Version),
	}
}

func instanceStateFromProto(state *appb.InstanceState) (instanceState, error) {
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
		Version: versionFromProto(state.GetVersion()),
	}, nil
}

func instanceStatesToProto(states []instanceState) []*appb.InstanceState {
	protoStates := make([]*appb.InstanceState, 0, len(states))
	for _, state := range states {
		protoStates = append(protoStates, instanceStateToProto(state))
	}

	return protoStates
}

func digestToProto(digest map[registry.ServiceInstanceID]version) []*appb.DigestEntry {
	entries := make([]*appb.DigestEntry, 0, len(digest))
	for id, value := range digest {
		entries = append(entries, &appb.DigestEntry{
			ServiceName: id.ServiceName,
			InstanceId:  id.InstanceID,
			Version:     versionToProto(value),
		})
	}

	return entries
}

func digestFromProto(entries []*appb.DigestEntry) map[registry.ServiceInstanceID]version {
	digest := make(map[registry.ServiceInstanceID]version, len(entries))
	for _, entry := range entries {
		id := registry.ServiceInstanceID{ServiceName: entry.GetServiceName(), InstanceID: entry.GetInstanceId()}
		digest[id] = versionFromProto(entry.GetVersion())
	}

	return digest
}

func instanceKeysToProto(ids []registry.ServiceInstanceID) []*appb.InstanceKey {
	keys := make([]*appb.InstanceKey, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, &appb.InstanceKey{ServiceName: id.ServiceName, InstanceId: id.InstanceID})
	}

	return keys
}

func instanceKeysFromProto(keys []*appb.InstanceKey) []registry.ServiceInstanceID {
	ids := make([]registry.ServiceInstanceID, 0, len(keys))
	for _, key := range keys {
		ids = append(ids, registry.ServiceInstanceID{ServiceName: key.GetServiceName(), InstanceID: key.GetInstanceId()})
	}

	return ids
}

func versionToProto(value version) *appb.Version {
	return &appb.Version{
		LamportTimestamp: value.LamportTimestamp,
		NodeId:           value.NodeID,
	}
}

func versionFromProto(value *appb.Version) version {
	return version{
		LamportTimestamp: value.GetLamportTimestamp(),
		NodeID:           value.GetNodeId(),
	}
}

func statusToProto(value registry.ServiceInstanceStatus) appb.InstanceStatus {
	switch value {
	case registry.ServiceInstanceStatusHealthy:
		return appb.InstanceStatus_INSTANCE_STATUS_HEALTHY
	case registry.ServiceInstanceStatusExpired:
		return appb.InstanceStatus_INSTANCE_STATUS_EXPIRED
	case registry.ServiceInstanceStatusDeleted:
		return appb.InstanceStatus_INSTANCE_STATUS_DELETED
	default:
		return appb.InstanceStatus_INSTANCE_STATUS_UNSPECIFIED
	}
}

func statusFromProto(value appb.InstanceStatus) (registry.ServiceInstanceStatus, error) {
	switch value {
	case appb.InstanceStatus_INSTANCE_STATUS_HEALTHY:
		return registry.ServiceInstanceStatusHealthy, nil
	case appb.InstanceStatus_INSTANCE_STATUS_EXPIRED:
		return registry.ServiceInstanceStatusExpired, nil
	case appb.InstanceStatus_INSTANCE_STATUS_DELETED:
		return registry.ServiceInstanceStatusDeleted, nil
	default:
		return "", fmt.Errorf("unknown instance status %v", value)
	}
}
