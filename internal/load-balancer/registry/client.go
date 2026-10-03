package registry

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	registrypb "tpiii.local/daniel-tpiii/gen/registry"
	loadbalancer "tpiii.local/daniel-tpiii/internal/load-balancer"
)

type KeepaliveConfig struct {
	Time    time.Duration
	Timeout time.Duration
}

type Client struct {
	addresses []string
	current   int
	keepalive KeepaliveConfig

	conn   *grpc.ClientConn
	client registrypb.RegistryClient
}

func NewClient(addresses []string, keepaliveConfig KeepaliveConfig) (*Client, error) {
	if len(addresses) == 0 {
		return nil, errors.New("at least one registry address is required")
	}

	client := &Client{
		addresses: slices.Clone(addresses),
		keepalive: keepaliveConfig,
	}

	if err := client.connect(0); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Watch(ctx context.Context, serviceName string) ([]loadbalancer.ServiceInstance, *Subscription, error) {
	stream, err := c.client.Watch(ctx, &registrypb.WatchRequest{ServiceName: serviceName})
	if err != nil {
		return nil, nil, c.failover(ctx, err)
	}

	response, err := stream.Recv()
	if err != nil {
		return nil, nil, c.failover(ctx, err)
	}

	snapshot := response.GetSnapshot()
	if snapshot == nil {
		return nil, nil, c.failover(ctx, errors.New("watch stream did not start with a snapshot"))
	}

	return serviceInstancesFromProto(snapshot.GetInstances()), &Subscription{ctx: ctx, client: c, stream: stream}, nil
}

func (c *Client) failover(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	next := (c.current + 1) % len(c.addresses)
	if connErr := c.connect(next); connErr != nil {
		return errors.Join(err, connErr)
	}

	return err
}

func (c *Client) connect(index int) error {
	conn, err := grpc.NewClient(
		c.addresses[index],
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                c.keepalive.Time,
			Timeout:             c.keepalive.Timeout,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return fmt.Errorf("create registry connection to %s: %w", c.addresses[index], err)
	}

	if c.conn != nil {
		_ = c.conn.Close()
	}

	c.conn = conn
	c.client = registrypb.NewRegistryClient(conn)
	c.current = index

	return nil
}

func (c *Client) Close() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("close registry connection: %w", err)
	}

	return nil
}
