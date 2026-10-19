package cp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"

	"tpiii.local/daniel-tpiii/internal/registry"
)

const (
	raftLogFile        = "raft.db"
	snapshotsRetained  = 2
	transportMaxPool   = 3
	transportIOTimeout = 10 * time.Second
)

// Compile-time check that *Backend implements registry.Backend
var _ registry.Backend = (*Backend)(nil)

type Config struct {
	NodeID           string
	BindAddress      string
	AdvertiseAddress string
	DataDir          string
}

type Backend struct {
	config Config
	store  *stateStore

	raft      *raft.Raft
	boltStore *raftboltdb.BoltStore
}

func NewBackend(config Config) *Backend {
	return &Backend{
		config: config,
		store:  newStateStore(),
	}
}

func (b *Backend) Start(ctx context.Context) error {
	if err := os.MkdirAll(b.config.DataDir, 0o750); err != nil {
		return fmt.Errorf("create raft data dir: %w", err)
	}

	boltStore, err := raftboltdb.NewBoltStore(filepath.Join(b.config.DataDir, raftLogFile))
	if err != nil {
		return fmt.Errorf("open raft log store: %w", err)
	}

	snapshots, err := raft.NewFileSnapshotStore(b.config.DataDir, snapshotsRetained, os.Stderr)
	if err != nil {
		_ = boltStore.Close()
		return fmt.Errorf("open raft snapshot store: %w", err)
	}

	advertise, err := net.ResolveTCPAddr("tcp", b.config.AdvertiseAddress)
	if err != nil {
		_ = boltStore.Close()
		return fmt.Errorf("resolve raft advertise address: %w", err)
	}

	transport, err := raft.NewTCPTransport(b.config.BindAddress, advertise, transportMaxPool, transportIOTimeout, os.Stderr)
	if err != nil {
		_ = boltStore.Close()
		return fmt.Errorf("create raft transport: %w", err)
	}

	hasState, err := raft.HasExistingState(boltStore, boltStore, snapshots)
	if err != nil {
		_ = transport.Close()
		_ = boltStore.Close()
		return fmt.Errorf("check raft state: %w", err)
	}

	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(b.config.NodeID)
	raftConfig.Logger = hclog.New(&hclog.LoggerOptions{
		Name:       "raft",
		JSONFormat: true,
		Output:     os.Stderr,
	})

	node, err := raft.NewRaft(raftConfig, newFSM(b.store), boltStore, boltStore, snapshots, transport)
	if err != nil {
		_ = transport.Close()
		_ = boltStore.Close()
		return fmt.Errorf("create raft node: %w", err)
	}

	b.raft = node
	b.boltStore = boltStore

	if !hasState {
		bootstrap := raft.Configuration{
			Servers: []raft.Server{{ID: raftConfig.LocalID, Address: transport.LocalAddr()}},
		}

		if err := node.BootstrapCluster(bootstrap).Error(); err != nil {
			return fmt.Errorf("bootstrap raft cluster: %w", err)
		}
	}

	if err := b.waitForLeader(ctx); err != nil {
		return fmt.Errorf("wait for raft leader: %w", err)
	}

	return nil
}

func (b *Backend) waitForLeader(ctx context.Context) error {
	observations := make(chan raft.Observation, 1)
	observer := raft.NewObserver(observations, false, func(o *raft.Observation) bool {
		_, ok := o.Data.(raft.LeaderObservation)
		return ok
	})

	b.raft.RegisterObserver(observer)
	defer b.raft.DeregisterObserver(observer)

	for {
		if address, _ := b.raft.LeaderWithID(); address != "" {
			return nil
		}

		select {
		case <-observations:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Backend) Shutdown(ctx context.Context) error {
	if b.raft == nil {
		return nil
	}

	done := make(chan error, 1)
	go func() {
		done <- b.raft.Shutdown().Error()
	}()

	var err error
	select {
	case err = <-done:
	case <-ctx.Done():
		err = ctx.Err()
	}

	return errors.Join(err, b.boltStore.Close())
}

func (b *Backend) Register(ctx context.Context, params registry.RegisterParams) error {
	return b.apply(ctx, newRegisterCommand(params))
}

func (b *Backend) Update(ctx context.Context, params registry.UpdateParams) error {
	return b.apply(ctx, newUpdateCommand(params))
}

func (b *Backend) Deregister(ctx context.Context, id registry.ServiceInstanceID) error {
	return b.apply(ctx, newInstanceCommand(commandDeregister, id))
}

func (b *Backend) Expire(ctx context.Context, id registry.ServiceInstanceID) error {
	return b.apply(ctx, newInstanceCommand(commandExpire, id))
}

func (b *Backend) Recover(ctx context.Context, id registry.ServiceInstanceID) error {
	return b.apply(ctx, newInstanceCommand(commandRecover, id))
}

func (b *Backend) Lookup(ctx context.Context, id registry.ServiceInstanceID) ([]registry.ServiceInstance, error) {
	if err := b.consistentRead(ctx); err != nil {
		return nil, err
	}

	return b.store.Lookup(id)
}

func (b *Backend) RegisterStateChangeHandler(handler registry.StateChangeHandler) {
	b.store.RegisterStateChangeHandler(handler)
}

func (b *Backend) apply(ctx context.Context, cmd command) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("encode command: %w", err)
	}

	timeout, err := remaining(ctx)
	if err != nil {
		return err
	}

	future := b.raft.Apply(data, timeout)
	if err := future.Error(); err != nil {
		return raftError(err)
	}

	if err, ok := future.Response().(error); ok {
		return err
	}

	return nil
}

func (b *Backend) consistentRead(ctx context.Context) error {
	readIndex := b.raft.CommitIndex()

	if err := b.raft.VerifyLeader().Error(); err != nil {
		return raftError(err)
	}

	if b.raft.AppliedIndex() >= readIndex {
		return nil
	}

	timeout, err := remaining(ctx)
	if err != nil {
		return err
	}

	if err := b.raft.Barrier(timeout).Error(); err != nil {
		return raftError(err)
	}

	return nil
}

func remaining(ctx context.Context) (time.Duration, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return 0, nil
	}

	return time.Until(deadline), nil
}

func raftError(err error) error {
	if errors.Is(err, raft.ErrNotLeader) || errors.Is(err, raft.ErrLeadershipLost) {
		return fmt.Errorf("%w: %w", registry.ErrNotLeader, err)
	}

	return err
}
