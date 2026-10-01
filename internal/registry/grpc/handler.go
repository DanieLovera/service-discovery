package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	registrypb "tpiii.local/daniel-tpiii/gen/registry"
	"tpiii.local/daniel-tpiii/internal/registry"
)

type Handler struct {
	registrypb.UnimplementedRegistryServer

	registryManager *registry.RegistryManager
	healthManager   *registry.HealthManager
	watchManager    *registry.WatchManager
}

func NewHandler(
	registryManager *registry.RegistryManager,
	healthManager *registry.HealthManager,
	watchManager *registry.WatchManager,
) *Handler {
	return &Handler{
		registryManager: registryManager,
		healthManager:   healthManager,
		watchManager:    watchManager,
	}
}

func (h *Handler) Register(ctx context.Context, req *registrypb.RegisterRequest) (*registrypb.RegisterResponse, error) {
	err := h.registryManager.Register(ctx, registry.RegisterParams{
		ID: registry.ServiceInstanceID{
			ServiceName: req.GetServiceName(),
			InstanceID:  req.GetInstanceId(),
		},
		Address:           req.GetAddress(),
		Weight:            int(req.GetWeight()),
		HeartbeatInterval: time.Duration(req.GetHeartbeatIntervalMs()) * time.Millisecond,
	})
	if err != nil {
		return nil, grpcError(err)
	}

	return &registrypb.RegisterResponse{}, nil
}

func (h *Handler) Update(ctx context.Context, req *registrypb.UpdateRequest) (*registrypb.UpdateResponse, error) {
	err := h.registryManager.Update(ctx, registry.UpdateParams{
		ID: registry.ServiceInstanceID{
			ServiceName: req.GetServiceName(),
			InstanceID:  req.GetInstanceId(),
		},
		Address:           req.GetAddress(),
		Weight:            int(req.GetWeight()),
		HeartbeatInterval: time.Duration(req.GetHeartbeatIntervalMs()) * time.Millisecond,
	})
	if err != nil {
		return nil, grpcError(err)
	}

	return &registrypb.UpdateResponse{}, nil
}

func (h *Handler) Deregister(ctx context.Context, req *registrypb.DeregisterRequest) (*registrypb.DeregisterResponse, error) {
	err := h.registryManager.Deregister(ctx, registry.ServiceInstanceID{
		ServiceName: req.GetServiceName(),
		InstanceID:  req.GetInstanceId(),
	})
	if err != nil {
		return nil, grpcError(err)
	}

	return &registrypb.DeregisterResponse{}, nil
}

func (h *Handler) Heartbeat(ctx context.Context, req *registrypb.HeartbeatRequest) (*registrypb.HeartbeatResponse, error) {
	err := h.healthManager.Heartbeat(ctx, registry.ServiceInstanceID{
		ServiceName: req.GetServiceName(),
		InstanceID:  req.GetInstanceId(),
	})
	if err != nil {
		return nil, grpcError(err)
	}

	return &registrypb.HeartbeatResponse{}, nil
}

func (h *Handler) Lookup(ctx context.Context, req *registrypb.LookupRequest) (*registrypb.LookupResponse, error) {
	id := registry.ServiceInstanceID{ServiceName: req.GetServiceName()}
	if req.InstanceId != nil {
		id.InstanceID = req.GetInstanceId()
	}

	instances, err := h.registryManager.Lookup(ctx, id)
	if err != nil {
		return nil, grpcError(err)
	}

	response := &registrypb.LookupResponse{
		Instances: make([]*registrypb.ServiceInstance, 0, len(instances)),
	}
	for _, instance := range instances {
		response.Instances = append(response.Instances, serviceInstanceToProto(instance))
	}

	return response, nil
}

func (h *Handler) Watch(req *registrypb.WatchRequest, stream registrypb.Registry_WatchServer) error {
	snapshot, subscriber, cleanup, err := h.watchManager.Watch(stream.Context(), req.GetServiceName())
	if err != nil {
		return grpcError(err)
	}
	defer cleanup()

	if err := stream.Send(watchSnapshotToProto(snapshot)); err != nil {
		return err
	}

	for {
		select {
		case change, ok := <-subscriber.Events():
			if !ok {
				return status.Error(codes.ResourceExhausted, "watch subscriber fell behind")
			}

			if err := stream.Send(watchEventToProto(change)); err != nil {
				return err
			}

		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}
