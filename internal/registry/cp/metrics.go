package cp

import (
	"strings"

	"github.com/hashicorp/raft"
	"github.com/prometheus/client_golang/prometheus"

	"tpiii.local/daniel-tpiii/internal/registry"
)

func (b *Backend) RegisterMetrics(registerer prometheus.Registerer) {
	registerer.MustRegister(registry.NewInstancesCollector(b.store.Instances))

	for _, state := range []raft.RaftState{raft.Follower, raft.Candidate, raft.Leader} {
		registerer.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Name:        "registry_raft_state",
			Help:        "Raft role of the registry node: 1 for its current state, 0 for the others.",
			ConstLabels: prometheus.Labels{"state": strings.ToUpper(state.String())},
		}, func() float64 {
			if b.raft.State() == state {
				return 1
			}

			return 0
		}))
	}
}
