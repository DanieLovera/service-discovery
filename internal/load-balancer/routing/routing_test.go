package routing

import (
	"errors"
	"sync"
	"testing"

	loadbalancer "tpiii.local/daniel-tpiii/internal/load-balancer"
)

func instances(ids ...string) []loadbalancer.ServiceInstance {
	result := make([]loadbalancer.ServiceInstance, 0, len(ids))
	for _, id := range ids {
		result = append(result, loadbalancer.ServiceInstance{ID: id, Address: id + ":8080", Weight: 1})
	}

	return result
}

func pickID(t *testing.T, strategy Strategy, pool []loadbalancer.ServiceInstance) (string, func()) {
	t.Helper()

	instance, done, ok := strategy.Pick(pool)
	if !ok {
		t.Fatal("Pick() ok = false, want true")
	}

	return instance.ID, done
}

func TestStrategiesRejectEmptyPool(t *testing.T) {
	for _, name := range []string{"round-robin", "least-connections", "weighted"} {
		strategy, err := NewStrategy(name)
		if err != nil {
			t.Fatalf("NewStrategy(%q) error = %v, want nil", name, err)
		}

		if _, _, ok := strategy.Pick(nil); ok {
			t.Errorf("%s: Pick(empty) ok = true, want false", name)
		}
	}

	if _, err := NewStrategy("random"); err == nil {
		t.Error("NewStrategy(\"random\") error = nil, want unsupported strategy error")
	}
}

func TestRoundRobinCyclesThroughInstances(t *testing.T) {
	strategy := NewRoundRobin()
	pool := instances("a", "b", "c")

	var got []string
	for range 6 {
		id, _ := pickID(t, strategy, pool)
		got = append(got, id)
	}

	want := []string{"a", "b", "c", "a", "b", "c"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("picks = %v, want %v", got, want)
		}
	}
}

func TestLeastConnectionsPrefersLeastBusyInstance(t *testing.T) {
	strategy := NewLeastConnections()
	pool := instances("a", "b", "c")

	first, _ := pickID(t, strategy, pool)
	second, releaseSecond := pickID(t, strategy, pool)
	third, _ := pickID(t, strategy, pool)

	if first == second || second == third || first == third {
		t.Fatalf("three concurrent picks = %s, %s, %s; want three different instances", first, second, third)
	}

	releaseSecond()

	if next, _ := pickID(t, strategy, pool); next != second {
		t.Errorf("pick after releasing %s = %s, want %s (the only idle instance)", second, next, second)
	}
}

func TestLeastConnectionsSpreadsTies(t *testing.T) {
	strategy := NewLeastConnections()
	pool := instances("a", "b", "c")

	counts := make(map[string]int)
	for range 30 {
		id, done := pickID(t, strategy, pool)
		counts[id]++
		done()
	}

	for _, id := range []string{"a", "b", "c"} {
		if counts[id] != 10 {
			t.Errorf("picks per instance = %v, want 10 each when every request finishes before the next", counts)
			break
		}
	}
}

func TestLeastConnectionsIsSafeForConcurrentUse(t *testing.T) {
	strategy := NewLeastConnections()
	pool := instances("a", "b", "c")

	var wg sync.WaitGroup
	for range 20 {
		wg.Go(func() {
			for range 100 {
				_, done, _ := strategy.Pick(pool)
				done()
			}
		})
	}
	wg.Wait()

	strategy.mu.Lock()
	defer strategy.mu.Unlock()

	if len(strategy.active) != 0 {
		t.Errorf("active connections after all requests finished = %v, want none", strategy.active)
	}
}

func TestWeightedFollowsWeightsSmoothly(t *testing.T) {
	strategy := NewWeighted()
	pool := []loadbalancer.ServiceInstance{
		{ID: "a", Weight: 5},
		{ID: "b", Weight: 1},
		{ID: "c", Weight: 1},
	}

	for cycle := range 10 {
		counts := make(map[string]int)
		for range 7 {
			id, _ := pickID(t, strategy, pool)
			counts[id]++
		}

		if counts["a"] != 5 || counts["b"] != 1 || counts["c"] != 1 {
			t.Fatalf("cycle %d picks = %v, want a=5 b=1 c=1 in every cycle of 7", cycle, counts)
		}
	}
}

func TestWeightedForgetsRemovedInstances(t *testing.T) {
	strategy := NewWeighted()

	pickID(t, strategy, instances("a", "b", "c"))

	counts := make(map[string]int)
	for range 1000 {
		id, _ := pickID(t, strategy, instances("a", "b"))
		counts[id]++
	}

	if counts["c"] != 0 || counts["a"]+counts["b"] != 1000 || max(counts["a"]-counts["b"], counts["b"]-counts["a"]) > 2 {
		t.Errorf("picks after removing c = %v, want only a and b, differing by at most 2", counts)
	}

	strategy.mu.Lock()
	defer strategy.mu.Unlock()

	if _, ok := strategy.current["c"]; ok {
		t.Error("weighted state still tracks the removed instance c")
	}
}

func TestManagerPick(t *testing.T) {
	pool := loadbalancer.NewServicePool()
	manager := NewManager(map[string]Upstream{
		"service-a": {Pool: pool, Strategy: NewRoundRobin()},
	})

	if _, _, err := manager.Pick("unknown"); !errors.Is(err, ErrUnknownService) {
		t.Errorf("Pick(unknown) error = %v, want %v", err, ErrUnknownService)
	}

	if _, _, err := manager.Pick("service-a"); !errors.Is(err, ErrNoInstances) {
		t.Errorf("Pick(empty pool) error = %v, want %v", err, ErrNoInstances)
	}

	pool.Replace(instances("a1"))

	instance, done, err := manager.Pick("service-a")
	if err != nil || instance.ID != "a1" {
		t.Fatalf("Pick() = %s, %v; want a1, nil", instance.ID, err)
	}
	done()
}
