package registry

import "github.com/prometheus/client_golang/prometheus"

var instanceStatuses = []ServiceInstanceStatus{
	ServiceInstanceStatusHealthy,
	ServiceInstanceStatusExpired,
	ServiceInstanceStatusDeleted,
}

type instancesCollector struct {
	desc      *prometheus.Desc
	instances func() []ServiceInstance
}

func NewInstancesCollector(instances func() []ServiceInstance) prometheus.Collector {
	return &instancesCollector{
		desc: prometheus.NewDesc(
			"registry_instances",
			"Service instances in the local state of the registry node, by service and status.",
			[]string{"service", "status"},
			nil,
		),
		instances: instances,
	}
}

func (c *instancesCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- c.desc
}

func (c *instancesCollector) Collect(metrics chan<- prometheus.Metric) {
	counts := make(map[string]map[ServiceInstanceStatus]int)
	for _, instance := range c.instances() {
		service := instance.ID.ServiceName
		if counts[service] == nil {
			counts[service] = make(map[ServiceInstanceStatus]int)
		}

		counts[service][instance.Status]++
	}

	for service, byStatus := range counts {
		for _, status := range instanceStatuses {
			metrics <- prometheus.MustNewConstMetric(
				c.desc,
				prometheus.GaugeValue,
				float64(byStatus[status]),
				service,
				string(status),
			)
		}
	}
}
