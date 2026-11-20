package registry

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type ClientMetrics struct {
	requests *prometheus.CounterVec
}

func NewClientMetrics(registerer prometheus.Registerer) *ClientMetrics {
	m := &ClientMetrics{
		requests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "mock_registry_requests_total",
			Help: "gRPC calls sent by the mock service to the registry, by method and status code.",
		}, []string{"method", "code"}),
	}

	registerer.MustRegister(m.requests)

	return m
}

func (m *ClientMetrics) unaryInterceptor(
	ctx context.Context,
	method string,
	request, reply any,
	conn *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	err := invoker(ctx, method, request, reply, conn, opts...)
	m.requests.WithLabelValues(method, status.Code(err).String()).Inc()

	return err
}
