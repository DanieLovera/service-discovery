package grpc

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type ServerMetrics struct {
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
}

func NewServerMetrics(registerer prometheus.Registerer) *ServerMetrics {
	m := &ServerMetrics{
		requests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "registry_grpc_requests_total",
			Help: "gRPC calls handled by the registry node, by method and status code.",
		}, []string{"method", "code"}),
		duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "registry_grpc_request_duration_seconds",
			Help:    "Duration of the unary gRPC calls handled by the registry node.",
			Buckets: prometheus.ExponentialBuckets(0.0005, 2, 15),
		}, []string{"method"}),
	}

	registerer.MustRegister(m.requests, m.duration)

	return m
}

func (m *ServerMetrics) unaryInterceptor(
	ctx context.Context,
	request any,
	info *googlegrpc.UnaryServerInfo,
	handler googlegrpc.UnaryHandler,
) (any, error) {
	start := time.Now()
	response, err := handler(ctx, request)

	m.requests.WithLabelValues(info.FullMethod, statusCode(err)).Inc()
	m.duration.WithLabelValues(info.FullMethod).Observe(time.Since(start).Seconds())

	return response, err
}

func (m *ServerMetrics) streamInterceptor(
	server any,
	stream googlegrpc.ServerStream,
	info *googlegrpc.StreamServerInfo,
	handler googlegrpc.StreamHandler,
) error {
	err := handler(server, stream)
	m.requests.WithLabelValues(info.FullMethod, statusCode(err)).Inc()

	return err
}

func statusCode(err error) string {
	if s, ok := status.FromError(err); ok {
		return s.Code().String()
	}

	return status.FromContextError(err).Code().String()
}
