package registry

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"

	registrypb "tpiii.local/daniel-tpiii/gen/registry"
	loadbalancer "tpiii.local/daniel-tpiii/internal/load-balancer"
)

const (
	resolverScheme = "registry"
	serviceConfig  = `{"loadBalancingConfig": [{"round_robin": {}}]}`
)

type KeepaliveConfig struct {
	Time    time.Duration
	Timeout time.Duration
}

type Client struct {
	conn   *grpc.ClientConn
	client registrypb.RegistryClient
}

func NewClient(addresses []string, keepaliveConfig KeepaliveConfig) (*Client, error) {
	if len(addresses) == 0 {
		return nil, errors.New("at least one registry address is required")
	}

	resolverAddresses := make([]resolver.Address, 0, len(addresses))
	for _, address := range addresses {
		resolverAddresses = append(resolverAddresses, resolver.Address{Addr: address})
	}

	builder := manual.NewBuilderWithScheme(resolverScheme)
	builder.InitialState(resolver.State{Addresses: resolverAddresses})

	conn, err := grpc.NewClient(
		resolverScheme+":///",
		grpc.WithResolvers(builder),
		grpc.WithDefaultServiceConfig(serviceConfig),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                keepaliveConfig.Time,
			Timeout:             keepaliveConfig.Timeout,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("create registry connection: %w", err)
	}

	return &Client{
		conn:   conn,
		client: registrypb.NewRegistryClient(conn),
	}, nil
}

func (c *Client) Watch(ctx context.Context, serviceName string) ([]loadbalancer.ServiceInstance, *Subscription, error) {
	stream, err := c.client.Watch(ctx, &registrypb.WatchRequest{ServiceName: serviceName})
	if err != nil {
		return nil, nil, err
	}

	response, err := stream.Recv()
	if err != nil {
		return nil, nil, err
	}

	snapshot := response.GetSnapshot()
	if snapshot == nil {
		return nil, nil, errors.New("watch stream did not start with a snapshot")
	}

	return serviceInstancesFromProto(snapshot.GetInstances()), &Subscription{stream: stream}, nil
}

func (c *Client) Close() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("close registry connection: %w", err)
	}

	return nil
}
