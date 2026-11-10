package registry

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc/status"

	registrypb "tpiii.local/daniel-tpiii/gen/registry"
	"tpiii.local/daniel-tpiii/internal/load-generator/results"
)

const (
	visibilityServiceName       = "visibility-probe"
	visibilityInstanceAddress   = "visibility-probe:0"
	visibilityHeartbeatInterval = time.Hour
	watchRetryInterval          = 100 * time.Millisecond
)

var (
	ProbesHeader    = []string{"timestamp_ms", "probe_id", "write_node", "code"}
	SightingsHeader = []string{"timestamp_ms", "probe_id", "node", "source"}
)

type Visibility struct {
	nodes           []*Node
	interval        time.Duration
	timeout         time.Duration
	probesWriter    *results.Writer
	sightingsWriter *results.Writer
}

func NewVisibility(
	nodes []*Node,
	interval time.Duration,
	timeout time.Duration,
	probesWriter *results.Writer,
	sightingsWriter *results.Writer,
) *Visibility {
	return &Visibility{
		nodes:           nodes,
		interval:        interval,
		timeout:         timeout,
		probesWriter:    probesWriter,
		sightingsWriter: sightingsWriter,
	}
}

func (v *Visibility) Run(ctx context.Context) {
	ready := make(chan struct{}, len(v.nodes))

	var watchers sync.WaitGroup
	for _, node := range v.nodes {
		watchers.Go(func() {
			v.watch(ctx, node, ready)
		})
	}

	for range v.nodes {
		select {
		case <-ready:
		case <-ctx.Done():
		}
	}

	v.sendProbes(ctx)
	watchers.Wait()
}

func (v *Visibility) watch(ctx context.Context, node *Node, ready chan<- struct{}) {
	seen := make(map[string]bool)

	var once sync.Once
	signalReady := func() {
		once.Do(func() {
			ready <- struct{}{}
		})
	}

	for ctx.Err() == nil {
		v.receive(ctx, node, seen, signalReady)

		select {
		case <-time.After(watchRetryInterval):
		case <-ctx.Done():
		}
	}
}

func (v *Visibility) receive(ctx context.Context, node *Node, seen map[string]bool, onSnapshot func()) {
	stream, err := node.client.Watch(ctx, &registrypb.WatchRequest{
		ServiceName: visibilityServiceName,
	})
	if err != nil {
		return
	}

	for {
		response, err := stream.Recv()
		if err != nil {
			return
		}

		receivedAt := time.Now()

		if snapshot := response.GetSnapshot(); snapshot != nil {
			for _, instance := range snapshot.GetInstances() {
				v.recordSighting(node.Name, instance.GetInstanceId(), receivedAt, "snapshot", seen)
			}

			onSnapshot()
			continue
		}

		instance := response.GetEvent().GetInstance()
		if instance.GetStatus() == registrypb.ServiceInstanceStatus_SERVICE_INSTANCE_STATUS_HEALTHY {
			v.recordSighting(node.Name, instance.GetInstanceId(), receivedAt, "event", seen)
		}
	}
}

func (v *Visibility) recordSighting(nodeName, probeID string, at time.Time, source string, seen map[string]bool) {
	if seen[probeID] {
		return
	}

	seen[probeID] = true
	v.sightingsWriter.Write(results.FormatTimestamp(at), probeID, nodeName, source)
}

func (v *Visibility) sendProbes(ctx context.Context) {
	ticker := time.NewTicker(v.interval)
	defer ticker.Stop()

	var wg sync.WaitGroup
	defer wg.Wait()

	for sequence := 0; ; sequence++ {
		select {
		case <-ticker.C:
			node := v.nodes[sequence%len(v.nodes)]
			probeID := fmt.Sprintf("probe-%d", sequence)
			wg.Go(func() {
				v.sendProbe(ctx, node, probeID)
			})

		case <-ctx.Done():
			return
		}
	}
}

func (v *Visibility) sendProbe(ctx context.Context, node *Node, probeID string) {
	sentAt := time.Now()
	err := v.register(ctx, node, probeID)

	if ctx.Err() != nil {
		return
	}

	v.probesWriter.Write(results.FormatTimestamp(sentAt), probeID, node.Name, status.Code(err).String())
}

func (v *Visibility) register(ctx context.Context, node *Node, instanceID string) error {
	ctx, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	_, err := node.client.Register(ctx, &registrypb.RegisterRequest{
		ServiceName:         visibilityServiceName,
		InstanceId:          instanceID,
		Address:             visibilityInstanceAddress,
		Weight:              1,
		HeartbeatIntervalMs: visibilityHeartbeatInterval.Milliseconds(),
	})

	return err
}
