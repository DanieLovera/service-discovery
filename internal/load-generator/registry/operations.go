package registry

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	registrypb "tpiii.local/daniel-tpiii/gen/registry"
	"tpiii.local/daniel-tpiii/internal/load-generator/results"
)

const (
	operationsServiceName       = "load-generator"
	operationsInstanceAddress   = "load-generator:0"
	operationsHeartbeatInterval = time.Minute
	instancesPerNode            = 100
)

var OperationsHeader = []string{"timestamp_ms", "node", "method", "code", "latency_ms"}

type Operations struct {
	node     *Node
	interval time.Duration
	timeout  time.Duration
	writer   *results.Writer
}

func NewOperations(node *Node, interval, timeout time.Duration, writer *results.Writer) *Operations {
	return &Operations{
		node:     node,
		interval: interval,
		timeout:  timeout,
		writer:   writer,
	}
}

func (o *Operations) Run(ctx context.Context) {
	ticker := time.NewTicker(o.interval)
	defer ticker.Stop()

	var wg sync.WaitGroup
	defer wg.Wait()

	for cycle := 0; ; cycle++ {
		select {
		case <-ticker.C:
			instanceID := fmt.Sprintf("%s-%d", o.node.Name, cycle%instancesPerNode)
			wg.Go(func() {
				o.runCycle(ctx, instanceID)
			})

		case <-ctx.Done():
			return
		}
	}
}

func (o *Operations) runCycle(ctx context.Context, instanceID string) {
	err := o.call(ctx, "Register", func(ctx context.Context) error {
		_, err := o.node.client.Register(ctx, &registrypb.RegisterRequest{
			ServiceName:         operationsServiceName,
			InstanceId:          instanceID,
			Address:             operationsInstanceAddress,
			Weight:              1,
			HeartbeatIntervalMs: operationsHeartbeatInterval.Milliseconds(),
		})
		return err
	})
	if ctx.Err() != nil || isTransient(err) {
		return
	}

	o.call(ctx, "Heartbeat", func(ctx context.Context) error {
		_, err := o.node.client.Heartbeat(ctx, &registrypb.HeartbeatRequest{
			ServiceName: operationsServiceName,
			InstanceId:  instanceID,
		})
		return err
	})

	o.call(ctx, "Lookup", func(ctx context.Context) error {
		_, err := o.node.client.Lookup(ctx, &registrypb.LookupRequest{
			ServiceName: operationsServiceName,
		})
		return err
	})

	o.call(ctx, "Deregister", func(ctx context.Context) error {
		_, err := o.node.client.Deregister(ctx, &registrypb.DeregisterRequest{
			ServiceName: operationsServiceName,
			InstanceId:  instanceID,
		})
		return err
	})
}

func (o *Operations) call(ctx context.Context, method string, operation func(context.Context) error) error {
	callCtx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	start := time.Now()
	err := operation(callCtx)
	latency := time.Since(start)

	if ctx.Err() != nil {
		return ctx.Err()
	}

	o.writer.Write(
		results.FormatTimestamp(start),
		o.node.Name,
		method,
		status.Code(err).String(),
		results.FormatMilliseconds(latency),
	)

	return err
}

func isTransient(err error) bool {
	switch status.Code(err) {
	case codes.Unavailable, codes.DeadlineExceeded:
		return true
	default:
		return false
	}
}
