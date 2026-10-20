package cp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	forwardingpb "tpiii.local/daniel-tpiii/gen/forwarding"
	"tpiii.local/daniel-tpiii/internal/registry"
)

type forwardingServer struct {
	forwardingpb.UnimplementedForwardingServer

	backend *Backend
}

func (s *forwardingServer) Apply(ctx context.Context, req *forwardingpb.ApplyRequest) (*forwardingpb.ApplyResponse, error) {
	status, err := toStatus(s.backend.applyLocal(ctx, req.GetCommand()))
	if err != nil {
		return nil, err
	}

	return &forwardingpb.ApplyResponse{Status: status}, nil
}

func (s *forwardingServer) Lookup(ctx context.Context, req *forwardingpb.LookupRequest) (*forwardingpb.LookupResponse, error) {
	id := registry.ServiceInstanceID{ServiceName: req.GetServiceName(), InstanceID: req.GetInstanceId()}

	instances, lookupErr := s.backend.lookupLocal(ctx, id)

	status, err := toStatus(lookupErr)
	if err != nil {
		return nil, err
	}

	if status != forwardingpb.Status_STATUS_OK {
		return &forwardingpb.LookupResponse{Status: status}, nil
	}

	data, err := json.Marshal(instances)
	if err != nil {
		return nil, fmt.Errorf("encode instances: %w", err)
	}

	return &forwardingpb.LookupResponse{Status: status, Instances: data}, nil
}

type forwarder struct {
	nodeID string

	mu    sync.Mutex
	conns map[string]*grpc.ClientConn
}

func newForwarder(nodeID string) *forwarder {
	return &forwarder{
		nodeID: nodeID,
		conns:  make(map[string]*grpc.ClientConn),
	}
}

func (f *forwarder) apply(ctx context.Context, leader registry.ClusterMember, command []byte) error {
	client, err := f.client(leader)
	if err != nil {
		return err
	}

	response, err := client.Apply(ctx, &forwardingpb.ApplyRequest{SenderId: f.nodeID, Command: command})
	if err != nil {
		return err
	}

	return fromStatus(response.GetStatus(), registry.ErrServiceInstanceNotFound)
}

func (f *forwarder) lookup(
	ctx context.Context,
	leader registry.ClusterMember,
	id registry.ServiceInstanceID,
) ([]registry.ServiceInstance, error) {
	client, err := f.client(leader)
	if err != nil {
		return nil, err
	}

	request := &forwardingpb.LookupRequest{SenderId: f.nodeID, ServiceName: id.ServiceName}
	if id.InstanceID != "" {
		request.InstanceId = &id.InstanceID
	}

	response, err := client.Lookup(ctx, request)
	if err != nil {
		return nil, err
	}

	notFound := registry.ErrServiceNotFound
	if id.InstanceID != "" {
		notFound = registry.ErrServiceInstanceNotFound
	}

	if err := fromStatus(response.GetStatus(), notFound); err != nil {
		return nil, err
	}

	var instances []registry.ServiceInstance
	if err := json.Unmarshal(response.GetInstances(), &instances); err != nil {
		return nil, fmt.Errorf("decode instances: %w", err)
	}

	return instances, nil
}

func (f *forwarder) client(member registry.ClusterMember) (forwardingpb.ForwardingClient, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	conn, ok := f.conns[member.ID]
	if !ok {
		var err error
		conn, err = grpc.NewClient(member.GRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("create forwarding connection to %s: %w", member.GRPCAddress, err)
		}

		f.conns[member.ID] = conn
	}

	return forwardingpb.NewForwardingClient(conn), nil
}

func (f *forwarder) close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	var errs []error
	for _, conn := range f.conns {
		errs = append(errs, conn.Close())
	}

	return errors.Join(errs...)
}

func toStatus(err error) (forwardingpb.Status, error) {
	switch {
	case err == nil:
		return forwardingpb.Status_STATUS_OK, nil
	case errors.Is(err, registry.ErrServiceInstanceAlreadyExists):
		return forwardingpb.Status_STATUS_ALREADY_EXISTS, nil
	case errors.Is(err, registry.ErrServiceInstanceNotFound), errors.Is(err, registry.ErrServiceNotFound):
		return forwardingpb.Status_STATUS_NOT_FOUND, nil
	case errors.Is(err, registry.ErrServiceInstanceDeleted):
		return forwardingpb.Status_STATUS_DELETED, nil
	case errors.Is(err, registry.ErrNotLeader):
		return forwardingpb.Status_STATUS_NOT_LEADER, nil
	default:
		return 0, err
	}
}

func fromStatus(status forwardingpb.Status, notFound error) error {
	switch status {
	case forwardingpb.Status_STATUS_OK:
		return nil
	case forwardingpb.Status_STATUS_ALREADY_EXISTS:
		return registry.ErrServiceInstanceAlreadyExists
	case forwardingpb.Status_STATUS_NOT_FOUND:
		return notFound
	case forwardingpb.Status_STATUS_DELETED:
		return registry.ErrServiceInstanceDeleted
	case forwardingpb.Status_STATUS_NOT_LEADER:
		return registry.ErrNotLeader
	default:
		return fmt.Errorf("unknown forwarding status %v", status)
	}
}
