package registry

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"

	registrypb "tpiii.local/daniel-tpiii/gen/registry"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Instance struct {
	ServiceName       string
	InstanceID        string
	Address           string
	Weight            int
	HeartbeatInterval time.Duration
}

type Client struct {
	addresses      []string
	current        int
	requestTimeout time.Duration

	conn   *grpc.ClientConn
	client registrypb.RegistryClient

	mu sync.Mutex
}

func NewClient(addresses []string, requestTimeout time.Duration) (*Client, error) {
	if len(addresses) == 0 {
		return nil, errors.New("at least one registry address is required")
	}

	client := &Client{
		addresses:      slices.Clone(addresses),
		requestTimeout: requestTimeout,
	}

	if err := client.connect(0); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Register(ctx context.Context, instance Instance) error {
	return c.execute(ctx, func(ctx context.Context, client registrypb.RegistryClient) error {
		_, err := client.Register(ctx, &registrypb.RegisterRequest{
			ServiceName:         instance.ServiceName,
			InstanceId:          instance.InstanceID,
			Address:             instance.Address,
			Weight:              int32(instance.Weight),
			HeartbeatIntervalMs: instance.HeartbeatInterval.Milliseconds(),
		})
		return err
	})
}

func (c *Client) Heartbeat(ctx context.Context, instance Instance) error {
	return c.execute(ctx, func(ctx context.Context, client registrypb.RegistryClient) error {
		_, err := client.Heartbeat(ctx, &registrypb.HeartbeatRequest{
			ServiceName: instance.ServiceName,
			InstanceId:  instance.InstanceID,
		})
		return err
	})
}

func (c *Client) Deregister(ctx context.Context, instance Instance) error {
	return c.execute(ctx, func(ctx context.Context, client registrypb.RegistryClient) error {
		_, err := client.Deregister(ctx, &registrypb.DeregisterRequest{
			ServiceName: instance.ServiceName,
			InstanceId:  instance.InstanceID,
		})
		return err
	})
}

func (c *Client) execute(
	ctx context.Context,
	operation func(context.Context, registrypb.RegistryClient) error,
) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	start := c.current
	var errs []error

	for attempt := range len(c.addresses) {
		if err := ctx.Err(); err != nil {
			return err
		}

		index := (start + attempt) % len(c.addresses)

		if index != c.current {
			if err := c.connect(index); err != nil {
				errs = append(errs, err)
				continue
			}
		}

		attemptCtx, cancel := context.WithTimeout(ctx, c.requestTimeout)
		err := operation(attemptCtx, c.client)
		cancel()

		if err == nil {
			return nil
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if !isRetryable(err) {
			return err
		}

		errs = append(errs, err)
	}

	return fmt.Errorf("all registry nodes unavailable: %w", errors.Join(errs...))
}

func isRetryable(err error) bool {
	switch status.Code(err) {
	case codes.Unavailable, codes.DeadlineExceeded:
		return true
	default:
		return false
	}
}

func (c *Client) connect(index int) error {
	conn, err := grpc.NewClient(c.addresses[index], grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return nil
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("close registry connection: %w", err)
	}

	return nil
}
