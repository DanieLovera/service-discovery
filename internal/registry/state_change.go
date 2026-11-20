package registry

type StateChangeType string

const (
	StateChangeRegistered StateChangeType = "REGISTERED"
	StateChangeUpdated    StateChangeType = "UPDATED"
	StateChangeHealthy    StateChangeType = "HEALTHY"
	StateChangeExpired    StateChangeType = "EXPIRED"
	StateChangeDeleted    StateChangeType = "DELETED"
	StateChangeRestored   StateChangeType = "RESTORED"
)

type StateChange struct {
	Type      StateChangeType
	Instance  ServiceInstance
	Instances []ServiceInstance
}
