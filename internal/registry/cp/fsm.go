package cp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/hashicorp/raft"

	"tpiii.local/daniel-tpiii/internal/registry"
)

const (
	resultOK            = "OK"
	resultAlreadyExists = "ALREADY_EXISTS"
	resultNotFound      = "NOT_FOUND"
	resultDeleted       = "DELETED"
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

	if cmd.RequestID == "" {
		return f.apply(cmd)
	}

	if result, ok := f.store.RequestResult(cmd.RequestID); ok {
		return resultError(result)
	}

	err := f.apply(cmd)
	f.store.RecordRequest(cmd.RequestID, resultCode(err))

	return err
}

func (f *fsm) apply(cmd command) error {
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
	return &fsmSnapshot{
		state: snapshotState{
			Instances:         f.store.Instances(),
			ProcessedRequests: f.store.ProcessedRequests(),
		},
	}, nil
}

func (f *fsm) Restore(snapshot io.ReadCloser) error {
	defer snapshot.Close()

	var state snapshotState
	if err := json.NewDecoder(snapshot).Decode(&state); err != nil {
		return fmt.Errorf("decode snapshot: %w", err)
	}

	f.store.Replace(state.Instances, state.ProcessedRequests)

	return nil
}

type snapshotState struct {
	Instances         []registry.ServiceInstance `json:"instances"`
	ProcessedRequests map[string]string          `json:"processed_requests"`
}

type fsmSnapshot struct {
	state snapshotState
}

func (s *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	if err := json.NewEncoder(sink).Encode(s.state); err != nil {
		_ = sink.Cancel()
		return fmt.Errorf("encode snapshot: %w", err)
	}

	return sink.Close()
}

func (s *fsmSnapshot) Release() {}

func resultCode(err error) string {
	switch {
	case err == nil:
		return resultOK
	case errors.Is(err, registry.ErrServiceInstanceAlreadyExists):
		return resultAlreadyExists
	case errors.Is(err, registry.ErrServiceInstanceNotFound):
		return resultNotFound
	case errors.Is(err, registry.ErrServiceInstanceDeleted):
		return resultDeleted
	default:
		return err.Error()
	}
}

func resultError(code string) error {
	switch code {
	case resultOK:
		return nil
	case resultAlreadyExists:
		return registry.ErrServiceInstanceAlreadyExists
	case resultNotFound:
		return registry.ErrServiceInstanceNotFound
	case resultDeleted:
		return registry.ErrServiceInstanceDeleted
	default:
		return errors.New(code)
	}
}
