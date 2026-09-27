package registry

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

const (
	healthManagerChangesBufferSize = 100
	healthManagerErrorsBufferSize  = 10
)

type HealthManager struct {
	backend Backend
	changes chan StateChange
	errs    chan error
	worker  *Worker[StateChange]
	logger  *slog.Logger

	mu        sync.Mutex
	deadlines map[ServiceInstanceID]time.Time
	wakeUp    chan struct{}
}

func NewHealthManager(backend Backend, logger *slog.Logger) *HealthManager {
	manager := &HealthManager{
		backend:   backend,
		changes:   make(chan StateChange, healthManagerChangesBufferSize),
		errs:      make(chan error, healthManagerErrorsBufferSize),
		logger:    logger,
		deadlines: make(map[ServiceInstanceID]time.Time),
		wakeUp:    make(chan struct{}, 1),
	}

	backend.SubscribeStateChanges(manager.onStateChange)
	manager.worker = New(manager.changes, manager.errs, manager.stateChangeHandler)

	return manager
}

func (h *HealthManager) Start(ctx context.Context) {
	h.worker.Start(ctx)
	go h.monitorExpirations(ctx)
	go h.monitorErrors(ctx)
}

func (h *HealthManager) Heartbeat(ctx context.Context, id ServiceInstanceID) error {
	instances, err := h.backend.Lookup(ctx, id)
	if err != nil {
		return err
	}

	if len(instances) == 0 {
		return ErrServiceInstanceNotFound
	}

	instance := instances[0]

	switch instance.Status {
	case ServiceInstanceStatusHealthy:
		h.renewTTL(id, instance.HeartbeatInterval)
		return nil

	case ServiceInstanceStatusExpired:
		if err := h.backend.Recover(ctx, id); err != nil {
			return err
		}

		h.renewTTL(id, instance.HeartbeatInterval)
		return nil

	case ServiceInstanceStatusDeleted:
		return ErrServiceInstanceDeleted

	default:
		return nil
	}
}

func (h *HealthManager) renewTTL(id ServiceInstanceID, ttl time.Duration) {
	h.mu.Lock()
	h.deadlines[id] = time.Now().Add(ttl)
	h.mu.Unlock()

	h.notifyScheduler()
}

func (h *HealthManager) notifyScheduler() {
	select {
	case h.wakeUp <- struct{}{}:
	default:
		// If the channel is full, we don't need to send another notification. Don't block the caller
	}
}

func (h *HealthManager) monitorExpirations(ctx context.Context) {
	timer := time.NewTimer(time.Hour)
	timer.Stop()
	defer timer.Stop()

	for {
		id, deadline, ok := h.nextExpiration()

		if !ok {
			select {
			case <-h.wakeUp:
				continue
			case <-ctx.Done():
				return
			}
		}

		timer.Reset(time.Until(deadline))

		select {
		case <-timer.C:
			if err := h.expire(ctx, id, deadline); err != nil {
				h.sendError(ctx, err)
			}

		case <-h.wakeUp:
			timer.Stop()

		case <-ctx.Done():
			return
		}
	}
}

func (h *HealthManager) nextExpiration() (ServiceInstanceID, time.Time, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	var (
		nextID       ServiceInstanceID
		nextDeadline time.Time
		found        bool
	)

	for id, deadline := range h.deadlines {
		if !found || deadline.Before(nextDeadline) {
			nextID = id
			nextDeadline = deadline
			found = true
		}
	}

	return nextID, nextDeadline, found
}

func (h *HealthManager) expire(ctx context.Context, id ServiceInstanceID, expectedDeadline time.Time) error {
	h.mu.Lock()

	deadline, ok := h.deadlines[id]
	if !ok || !deadline.Equal(expectedDeadline) {
		h.mu.Unlock()
		return nil
	}

	delete(h.deadlines, id)
	h.mu.Unlock()

	return h.backend.Expire(ctx, id)
}

func (h *HealthManager) sendError(ctx context.Context, err error) {
	select {
	case h.errs <- err:
	case <-ctx.Done():
	}
}

func (h *HealthManager) onStateChange(change StateChange) {
	h.changes <- change
}

func (h *HealthManager) stateChangeHandler(ctx context.Context, change StateChange) error {
	if change.Type != StateChangeExpired {
		return nil
	}

	id := change.Instance.ID
	if !h.hasValidTTL(id) {
		return nil
	}

	if err := h.backend.Recover(ctx, id); err != nil {
		return err
	}

	// A late Recover must not leave the instance HEALTHY after its local TTL has expired
	if !h.hasValidTTL(id) {
		return h.backend.Expire(ctx, id)
	}
	return nil
}

func (h *HealthManager) hasValidTTL(id ServiceInstanceID) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	deadline, ok := h.deadlines[id]
	return ok && time.Now().Before(deadline)
}

func (h *HealthManager) monitorErrors(ctx context.Context) {
	for {
		select {
		case err := <-h.errs:
			h.logger.Error("Failed to process health manager event", "error", err)
		case <-ctx.Done():
			return
		}
	}
}
