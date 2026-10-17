package cp

import (
	"time"

	"tpiii.local/daniel-tpiii/internal/registry"
)

type commandType string

const (
	commandRegister   commandType = "REGISTER"
	commandUpdate     commandType = "UPDATE"
	commandDeregister commandType = "DEREGISTER"
	commandExpire     commandType = "EXPIRE"
	commandRecover    commandType = "RECOVER"
)

type command struct {
	Type                commandType `json:"type"`
	ServiceName         string      `json:"service_name"`
	InstanceID          string      `json:"instance_id"`
	Address             string      `json:"address,omitempty"`
	Weight              int         `json:"weight,omitempty"`
	HeartbeatIntervalMs int64       `json:"heartbeat_interval_ms,omitempty"`
}

func newRegisterCommand(params registry.RegisterParams) command {
	return command{
		Type:                commandRegister,
		ServiceName:         params.ID.ServiceName,
		InstanceID:          params.ID.InstanceID,
		Address:             params.Address,
		Weight:              params.Weight,
		HeartbeatIntervalMs: params.HeartbeatInterval.Milliseconds(),
	}
}

func newUpdateCommand(params registry.UpdateParams) command {
	return command{
		Type:                commandUpdate,
		ServiceName:         params.ID.ServiceName,
		InstanceID:          params.ID.InstanceID,
		Address:             params.Address,
		Weight:              params.Weight,
		HeartbeatIntervalMs: params.HeartbeatInterval.Milliseconds(),
	}
}

func newInstanceCommand(commandType commandType, id registry.ServiceInstanceID) command {
	return command{
		Type:        commandType,
		ServiceName: id.ServiceName,
		InstanceID:  id.InstanceID,
	}
}

func (c command) registerParams() registry.RegisterParams {
	return registry.RegisterParams{
		ID:                c.id(),
		Address:           c.Address,
		Weight:            c.Weight,
		HeartbeatInterval: c.heartbeatInterval(),
	}
}

func (c command) updateParams() registry.UpdateParams {
	return registry.UpdateParams{
		ID:                c.id(),
		Address:           c.Address,
		Weight:            c.Weight,
		HeartbeatInterval: c.heartbeatInterval(),
	}
}

func (c command) id() registry.ServiceInstanceID {
	return registry.ServiceInstanceID{ServiceName: c.ServiceName, InstanceID: c.InstanceID}
}

func (c command) heartbeatInterval() time.Duration {
	return time.Duration(c.HeartbeatIntervalMs) * time.Millisecond
}
