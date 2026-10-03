package discovery

import (
	"context"
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

	wg sync.WaitGroup
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
	for serviceName, servicePool := range m.pools {
		client, err := registry.NewClient(m.addresses, m.keepalive)
		if err != nil {
			return err
		}

		m.wg.Go(func() {
			m.watch(ctx, client, serviceName, servicePool)
		})
	}

	return nil
}

func (m *Manager) Wait() {
	m.wg.Wait()
}

func (m *Manager) watch(ctx context.Context, registry *registry.Client, serviceName string, servicePool *loadbalancer.ServicePool) {
	defer func() {
		if err := registry.Close(); err != nil {
			m.logger.Error("Failed to close registry client", "service", serviceName, "error", err)
		}
	}()

	for {
		err := m.syncPool(ctx, registry, serviceName, servicePool)
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

func (m *Manager) syncPool(ctx context.Context, registry *registry.Client, serviceName string, servicePool *loadbalancer.ServicePool) error {
	snapshot, subscription, err := registry.Watch(ctx, serviceName)
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
