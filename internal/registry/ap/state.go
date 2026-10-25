package ap

import "tpiii.local/daniel-tpiii/internal/registry"

type version struct {
	LamportTimestamp uint64
	NodeID           string
}

func (v version) newerThan(other version) bool {
	if v.LamportTimestamp == other.LamportTimestamp {
		return v.NodeID > other.NodeID
	}

	return v.LamportTimestamp > other.LamportTimestamp
}

type instanceState struct {
	Instance registry.ServiceInstance
	Version  version
}
