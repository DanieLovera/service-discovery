package registry

import (
	"context"
	"slices"
	"sync"
)

const watchSubscriberBufferSize = 100

type Subscriber struct {
	events chan StateChange
}

func (s *Subscriber) Events() <-chan StateChange {
	return s.events
}

type WatchManager struct {
	backend Backend

	mu          sync.Mutex
	subscribers map[string][]*Subscriber
}

func NewWatchManager(backend Backend) *WatchManager {
	manager := &WatchManager{
		backend:     backend,
		subscribers: make(map[string][]*Subscriber),
	}

	backend.RegisterStateChangeHandler(manager.onStateChange)

	return manager
}

func (w *WatchManager) Watch(ctx context.Context, serviceName string) ([]ServiceInstance, *Subscriber, func(), error) {
	subscriber := &Subscriber{
		events: make(chan StateChange, watchSubscriberBufferSize),
	}

	w.subscribe(serviceName, subscriber)

	cleanup := func() {
		w.removeSubscriber(serviceName, subscriber)
	}

	snapshot, err := w.backend.Lookup(ctx, ServiceInstanceID{
		ServiceName: serviceName,
	})
	if err != nil {
		cleanup()
		return nil, nil, nil, err
	}

	return snapshot, subscriber, cleanup, nil
}

func (w *WatchManager) subscribe(serviceName string, subscriber *Subscriber) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.subscribers[serviceName] = append(w.subscribers[serviceName], subscriber)
}

func (w *WatchManager) removeSubscriber(serviceName string, subscriber *Subscriber) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.subscribers[serviceName] = slices.DeleteFunc(
		w.subscribers[serviceName],
		func(s *Subscriber) bool {
			return s == subscriber
		},
	)
}

func (w *WatchManager) onStateChange(change StateChange) {
	serviceName := change.Instance.ID.ServiceName

	w.mu.Lock()
	defer w.mu.Unlock()

	w.subscribers[serviceName] = slices.DeleteFunc(
		w.subscribers[serviceName],
		func(s *Subscriber) bool {
			select {
			case s.events <- change:
				return false
			default:
				close(s.events)
				return true
			}
		},
	)
}
