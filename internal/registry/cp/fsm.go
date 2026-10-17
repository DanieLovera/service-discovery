package cp

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/raft"

	"tpiii.local/daniel-tpiii/internal/registry"
)

type fsm struct {
	store *stateStore
}

func newFSM(store *stateStore) *fsm {
	return &fsm{store: store}
}

func (f *fsm) Apply(log *raft.Log) any {
	var cmd command
	if err := json.Unmarshal(log.Data, &cmd); err != nil {
		return fmt.Errorf("decode command: %w", err)
	}

	switch cmd.Type {
	case commandRegister:
		return f.store.Register(cmd.registerParams())
	case commandUpdate:
		return f.store.Update(cmd.updateParams())
	case commandDeregister:
		return f.store.Deregister(cmd.id())
	case commandExpire:
		return f.store.Expire(cmd.id())
	case commandRecover:
		return f.store.Recover(cmd.id())
	default:
		return fmt.Errorf("unknown command type %q", cmd.Type)
	}
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	return &fsmSnapshot{instances: f.store.Instances()}, nil
}

func (f *fsm) Restore(snapshot io.ReadCloser) error {
	defer snapshot.Close()

	var instances []registry.ServiceInstance
	if err := json.NewDecoder(snapshot).Decode(&instances); err != nil {
		return fmt.Errorf("decode snapshot: %w", err)
	}

	f.store.Replace(instances)

	return nil
}

type fsmSnapshot struct {
	instances []registry.ServiceInstance
}

func (s *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	if err := json.NewEncoder(sink).Encode(s.instances); err != nil {
		_ = sink.Cancel()
		return fmt.Errorf("encode snapshot: %w", err)
	}

	return sink.Close()
}

func (s *fsmSnapshot) Release() {}
