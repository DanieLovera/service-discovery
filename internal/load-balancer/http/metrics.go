package http

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type ServerMetrics struct {
	requests *prometheus.CounterVec
}

func NewServerMetrics(registerer prometheus.Registerer) *ServerMetrics {
	m := &ServerMetrics{
		requests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "lb_requests_total",
			Help: "HTTP requests routed by the load balancer, by service, instance and status code.",
		}, []string{"service", "instance_id", "code"}),
	}

	registerer.MustRegister(m.requests)

	return m
}

func (m *ServerMetrics) observe(service, instance string, code int) {
	m.requests.WithLabelValues(service, instance, strconv.Itoa(code)).Inc()
}
