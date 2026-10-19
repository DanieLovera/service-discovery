package registry

type ClusterMember struct {
	ID          string
	GRPCAddress string
	RaftAddress string
}
