package integration

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"

	registrypb "tpiii.local/daniel-tpiii/gen/registry"
	"tpiii.local/daniel-tpiii/internal/registry"
	"tpiii.local/daniel-tpiii/internal/registry/app"
)

const (
	startupTimeout           = 30 * time.Second
	convergeTimeout          = 15 * time.Second
	requestTimeout           = 5 * time.Second
	shutdownTimeout          = 2 * time.Second
	pollInterval             = 50 * time.Millisecond
	defaultHeartbeatInterval = time.Minute
)

var (
	backends      = []string{"cp", "ap"}
	discardLogger = slog.New(slog.DiscardHandler)
	httpClient    = &http.Client{Timeout: time.Second}
)

type process interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type runningProcess struct {
	name   string
	app    process
	cancel context.CancelFunc
	done   chan error
	once   sync.Once
}

func launch(name string, application process) *runningProcess {
	ctx, cancel := context.WithCancel(context.Background())

	p := &runningProcess{name: name, app: application, cancel: cancel, done: make(chan error, 1)}
	go func() {
		p.done <- application.Start(ctx)
	}()

	return p
}

func (p *runningProcess) waitReady(t *testing.T, observabilityAddress string) {
	t.Helper()

	url := "http://" + observabilityAddress + "/readyz"
	deadline := time.Now().Add(startupTimeout)

	for {
		select {
		case err := <-p.done:
			t.Fatalf("%s stopped during startup: %v", p.name, err)
		default:
		}

		if httpStatus(url) == http.StatusOK {
			return
		}

		if time.Now().After(deadline) {
			t.Fatalf("%s not ready after %s", p.name, startupTimeout)
		}

		time.Sleep(pollInterval)
	}
}

func (p *runningProcess) stop(t *testing.T) {
	p.once.Do(func() {
		p.cancel()

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := p.app.Shutdown(ctx); err != nil {
			t.Logf("%s shutdown: %v", p.name, err)
		}
	})
}

type registryNode struct {
	params  app.Params
	process *runningProcess
	client  registrypb.RegistryClient
}

type cluster struct {
	t     *testing.T
	nodes []*registryNode
}

func newCluster(t *testing.T, backend string, size int) *cluster {
	t.Helper()

	members := make([]registry.ClusterMember, size)
	for i := range members {
		members[i] = registry.ClusterMember{
			ID:          fmt.Sprintf("registry-%d", i+1),
			GRPCAddress: freeAddress(t),
			PeerAddress: freeAddress(t),
		}
	}

	c := &cluster{t: t}
	for _, member := range members {
		c.nodes = append(c.nodes, &registryNode{
			params: app.Params{
				NodeID:                   member.ID,
				Backend:                  backend,
				GRPCBindAddress:          member.GRPCAddress,
				ObservabilityBindAddress: freeAddress(t),
				ClusterMembers:           members,
				TTLMultiplier:            3,
				PeerBindAddress:          member.PeerAddress,
				RaftDataDir:              t.TempDir(),
				ForwardTimeout:           time.Second,
				ForwardMaxAttempts:       3,
				GossipFanout:             2,
				AntiEntropyInterval:      500 * time.Millisecond,
				SWIMProbeInterval:        200 * time.Millisecond,
				SWIMPingTimeout:          100 * time.Millisecond,
				SWIMPingReqTimeout:       200 * time.Millisecond,
				SWIMPingReqCount:         1,
				SWIMSuspectTimeout:       time.Second,
				SWIMRetransmitMultiplier: 3,
				Logger:                   discardLogger,
			},
			client: newRegistryClient(t, member.GRPCAddress),
		})
	}

	t.Cleanup(c.stopAll)

	for i := range c.nodes {
		c.launch(i)
	}

	for i := range c.nodes {
		c.nodes[i].process.waitReady(t, c.nodes[i].params.ObservabilityBindAddress)
	}

	return c
}

func (c *cluster) launch(i int) {
	node := c.nodes[i]

	application, err := app.New(node.params)
	if err != nil {
		c.t.Fatalf("create %s: %v", node.params.NodeID, err)
	}

	node.process = launch(node.params.NodeID, application)
}

func (c *cluster) start(i int) {
	c.t.Helper()

	c.launch(i)
	c.nodes[i].process.waitReady(c.t, c.nodes[i].params.ObservabilityBindAddress)
}

func (c *cluster) stop(i int) {
	node := c.nodes[i]
	if node.process == nil {
		return
	}

	node.process.stop(c.t)
	node.process = nil
}

func (c *cluster) stopAll() {
	for i := range c.nodes {
		c.stop(i)
	}
}

func (c *cluster) addresses(indexes ...int) []string {
	if len(indexes) == 0 {
		for i := range c.nodes {
			indexes = append(indexes, i)
		}
	}

	addresses := make([]string, 0, len(indexes))
	for _, i := range indexes {
		addresses = append(addresses, c.nodes[i].params.GRPCBindAddress)
	}

	return addresses
}

func (c *cluster) everyRunningNode(condition func(*registryNode) bool) bool {
	for _, node := range c.nodes {
		if node.process != nil && !condition(node) {
			return false
		}
	}

	return true
}

func (c *cluster) waitVisible(serviceName, instanceID string) {
	c.t.Helper()

	eventually(c.t, instanceID+" is visible on every running node", func() bool {
		return c.everyRunningNode(func(node *registryNode) bool {
			instances, err := node.lookup(serviceName)
			return err == nil && instances[instanceID] != nil
		})
	})
}

func (c *cluster) waitGone(serviceName, instanceID string) {
	c.t.Helper()

	eventually(c.t, instanceID+" is no longer returned by any running node", func() bool {
		return c.everyRunningNode(func(node *registryNode) bool {
			instances, err := node.lookup(serviceName)
			return err == nil && instances[instanceID] == nil
		})
	})
}

func (c *cluster) leader() int {
	c.t.Helper()

	leader := -1
	eventually(c.t, "a raft leader is elected", func() bool {
		for i, node := range c.nodes {
			if node.process == nil {
				continue
			}

			if value, ok := node.metric(`registry_raft_state{state="LEADER"}`); ok && value == 1 {
				leader = i
				return true
			}
		}

		return false
	})

	return leader
}

func (n *registryNode) register(serviceName, instanceID string, heartbeatInterval time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	_, err := n.client.Register(ctx, &registrypb.RegisterRequest{
		ServiceName:         serviceName,
		InstanceId:          instanceID,
		Address:             instanceID + ":8080",
		Weight:              1,
		HeartbeatIntervalMs: heartbeatInterval.Milliseconds(),
	})

	return err
}

func (n *registryNode) update(serviceName, instanceID string, weight int32) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	_, err := n.client.Update(ctx, &registrypb.UpdateRequest{
		ServiceName:         serviceName,
		InstanceId:          instanceID,
		Address:             instanceID + ":8080",
		Weight:              weight,
		HeartbeatIntervalMs: defaultHeartbeatInterval.Milliseconds(),
	})

	return err
}

func (n *registryNode) heartbeat(serviceName, instanceID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	_, err := n.client.Heartbeat(ctx, &registrypb.HeartbeatRequest{ServiceName: serviceName, InstanceId: instanceID})

	return err
}

func (n *registryNode) deregister(serviceName, instanceID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	_, err := n.client.Deregister(ctx, &registrypb.DeregisterRequest{ServiceName: serviceName, InstanceId: instanceID})

	return err
}

func (n *registryNode) lookup(serviceName string) (map[string]*registrypb.ServiceInstance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	response, err := n.client.Lookup(ctx, &registrypb.LookupRequest{ServiceName: serviceName})
	if err != nil {
		return nil, err
	}

	instances := make(map[string]*registrypb.ServiceInstance, len(response.GetInstances()))
	for _, instance := range response.GetInstances() {
		instances[instance.GetInstanceId()] = instance
	}

	return instances, nil
}

func (n *registryNode) metric(name string) (float64, bool) {
	response, err := httpClient.Get("http://" + n.params.ObservabilityBindAddress + "/metrics")
	if err != nil {
		return 0, false
	}
	defer response.Body.Close()

	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		if value, found := strings.CutPrefix(scanner.Text(), name+" "); found {
			parsed, err := strconv.ParseFloat(value, 64)
			return parsed, err == nil
		}
	}

	return 0, false
}

type watcher struct {
	serviceName string
	snapshot    []*registrypb.ServiceInstance
	responses   chan *registrypb.WatchResponse
}

func (n *registryNode) watch(t *testing.T, serviceName string) *watcher {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	stream, err := n.client.Watch(ctx, &registrypb.WatchRequest{ServiceName: serviceName})
	if err != nil {
		t.Fatalf("Watch() error = %v, want nil", err)
	}

	first, err := stream.Recv()
	if err != nil {
		t.Fatalf("receive watch snapshot: %v", err)
	}

	if first.GetSnapshot() == nil {
		t.Fatalf("first watch message = %v, want a snapshot", first)
	}

	w := &watcher{
		serviceName: serviceName,
		snapshot:    first.GetSnapshot().GetInstances(),
		responses:   make(chan *registrypb.WatchResponse, 100),
	}

	go func() {
		defer close(w.responses)

		for {
			response, err := stream.Recv()
			if err != nil {
				return
			}

			w.responses <- response
		}
	}()

	return w
}

func (w *watcher) waitEvent(
	t *testing.T,
	eventType registrypb.StateChangeType,
	instanceID string,
) *registrypb.WatchEvent {
	t.Helper()

	timeout := time.After(convergeTimeout)

	for {
		select {
		case response, ok := <-w.responses:
			if !ok {
				t.Fatalf("watch stream closed while waiting for %s of %s", eventType, instanceID)
			}

			event := response.GetEvent()
			if service := event.GetInstance().GetServiceName(); event != nil && service != w.serviceName {
				t.Errorf("watch of %s received an event of %s", w.serviceName, service)
			}

			if event.GetType() == eventType && event.GetInstance().GetInstanceId() == instanceID {
				return event
			}

		case <-timeout:
			t.Fatalf("timed out after %s waiting for %s of %s", convergeTimeout, eventType, instanceID)
		}
	}
}

func newRegistryClient(t *testing.T, address string) registrypb.RegistryClient {
	t.Helper()

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  50 * time.Millisecond,
				Multiplier: 1.6,
				MaxDelay:   500 * time.Millisecond,
			},
			MinConnectTimeout: time.Second,
		}),
	)
	if err != nil {
		t.Fatalf("create registry client: %v", err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
	})

	return registrypb.NewRegistryClient(conn)
}

func freeAddress(t *testing.T) string {
	t.Helper()

	for range 20 {
		tcp, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("listen: %v", err)
		}

		address := tcp.Addr().String()
		udp, err := net.ListenPacket("udp", address)
		_ = tcp.Close()

		if err != nil {
			continue
		}

		_ = udp.Close()

		return address
	}

	t.Fatal("no free port available for both TCP and UDP")

	return ""
}

func httpStatus(url string) int {
	response, err := httpClient.Get(url)
	if err != nil {
		return 0
	}
	_ = response.Body.Close()

	return response.StatusCode
}

func eventually(t *testing.T, description string, condition func() bool) {
	t.Helper()

	deadline := time.Now().Add(convergeTimeout)
	for !condition() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out after %s waiting until %s", convergeTimeout, description)
		}

		time.Sleep(pollInterval)
	}
}

func sortedIDs(instances []*registrypb.ServiceInstance) []string {
	ids := make([]string, 0, len(instances))
	for _, instance := range instances {
		ids = append(ids, instance.GetInstanceId())
	}

	slices.Sort(ids)

	return ids
}
