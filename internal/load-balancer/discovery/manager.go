package discovery

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	loadbalancer "tpiii.local/daniel-tpiii/internal/load-balancer"
	"tpiii.local/daniel-tpiii/internal/load-balancer/registry"
)

type Manager struct {
	addresses     []string
	pools         map[string]*loadbalancer.ServicePool
	retryInterval time.Duration
	keepalive     registry.KeepaliveConfig
	logger        *slog.Logger

	client *registry.Client
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewManager(
	addresses []string,
	pools map[string]*loadbalancer.ServicePool,
	retryInterval time.Duration,
	keepalive registry.KeepaliveConfig,
	logger *slog.Logger,
) *Manager {
	return &Manager{
		addresses:     addresses,
		pools:         pools,
		retryInterval: retryInterval,
		keepalive:     keepalive,
		logger:        logger,
	}
}

func (m *Manager) Start(ctx context.Context) error {
	client, err := registry.NewClient(m.addresses, m.keepalive)
	if err != nil {
		return err
	}

	m.client = client
	ctx, m.cancel = context.WithCancel(ctx)

	for serviceName, servicePool := range m.pools {
		m.wg.Go(func() {
			m.watch(ctx, serviceName, servicePool)
		})
	}

	return nil
}

func (m *Manager) Shutdown(ctx context.Context) error {
	if m.cancel == nil {
		return nil
	}

	m.cancel()

	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()

	var err error
	select {
	case <-done:
	case <-ctx.Done():
		err = ctx.Err()
	}

	return errors.Join(err, m.client.Close())
}

func (m *Manager) watch(ctx context.Context, serviceName string, servicePool *loadbalancer.ServicePool) {
	for {
		err := m.syncPool(ctx, serviceName, servicePool)
		if ctx.Err() != nil {
			return
		}

		m.logger.Warn("Watch stream failed, retrying", "service", serviceName, "error", err)

		select {
		case <-time.After(m.retryInterval):
		case <-ctx.Done():
			return
		}
	}
}

func (m *Manager) syncPool(ctx context.Context, serviceName string, servicePool *loadbalancer.ServicePool) error {
	snapshot, subscription, err := m.client.Watch(ctx, serviceName)
	if err != nil {
		return err
	}

	servicePool.Replace(snapshot)

	for {
		event, err := subscription.WaitEvent()
		if err != nil {
			return err
		}

		if event.Healthy {
			servicePool.Upsert(event.Instance)
		} else {
			servicePool.Remove(event.Instance.ID)
		}
	}
}
