package ap

import (
	"github.com/prometheus/client_golang/prometheus"

	"tpiii.local/daniel-tpiii/internal/registry"
)

func (b *Backend) RegisterMetrics(registerer prometheus.Registerer) {
	registerer.MustRegister(registry.NewInstancesCollector(b.store.Instances))

	for _, status := range []memberStatus{memberAlive, memberSuspect, memberConfirmed} {
		registerer.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Name:        "registry_swim_members",
			Help:        "Peers seen by the SWIM failure detector of the registry node, by status.",
			ConstLabels: prometheus.Labels{"status": string(status)},
		}, func() float64 {
			return float64(b.membership.count(status))
		}))
	}
}
