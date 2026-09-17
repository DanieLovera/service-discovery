package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	registry *prometheus.Registry
}

func New(service, instance string) *Metrics {
	registry := prometheus.NewRegistry()

	serviceInfo := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tpiii_service_info",
			Help: "Static information about a running TPIII service instance.",
		},
		[]string{"service", "instance"},
	)

	serviceInfo.WithLabelValues(service, instance).Set(1)

	registry.MustRegister(
		serviceInfo,
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)

	return &Metrics{
		registry: registry,
	}
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}
