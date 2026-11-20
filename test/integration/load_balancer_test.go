package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	lbapp "tpiii.local/daniel-tpiii/internal/load-balancer/app"
	mockapp "tpiii.local/daniel-tpiii/internal/mock-service/app"
)

func startMockService(t *testing.T, instanceID string, registryAddresses []string) *runningProcess {
	t.Helper()

	httpAddress := freeAddress(t)
	observabilityAddress := freeAddress(t)

	application, err := mockapp.New(mockapp.Params{
		ServiceName:              "service-a",
		InstanceID:               instanceID,
		HTTPBindAddress:          httpAddress,
		AdvertiseAddress:         httpAddress,
		ObservabilityBindAddress: observabilityAddress,
		RegistryAddresses:        registryAddresses,
		Weight:                   1,
		HeartbeatInterval:        time.Second,
		Logger:                   discardLogger,
	})
	if err != nil {
		t.Fatalf("create mock service %s: %v", instanceID, err)
	}

	process := launch(instanceID, application)
	t.Cleanup(func() {
		process.stop(t)
	})
	process.waitReady(t, observabilityAddress)

	return process
}

func startLoadBalancer(t *testing.T, registryAddresses []string) string {
	t.Helper()

	httpAddress := freeAddress(t)
	observabilityAddress := freeAddress(t)

	application, err := lbapp.New(lbapp.Params{
		HTTPBindAddress:          httpAddress,
		ObservabilityBindAddress: observabilityAddress,
		RegistryAddresses:        registryAddresses,
		RegistryRetryInterval:    100 * time.Millisecond,
		RegistryKeepaliveTime:    10 * time.Second,
		RegistryKeepaliveTimeout: 5 * time.Second,
		Services:                 map[string]string{"service-a": "round-robin"},
		Logger:                   discardLogger,
	})
	if err != nil {
		t.Fatalf("create load balancer: %v", err)
	}

	process := launch("load-balancer", application)
	t.Cleanup(func() {
		process.stop(t)
	})
	process.waitReady(t, observabilityAddress)

	return "http://" + httpAddress
}

func requestThroughLoadBalancer(url string) (int, string) {
	response, err := httpClient.Get(url)
	if err != nil {
		return 0, ""
	}
	defer response.Body.Close()

	var body struct {
		InstanceID string `json:"instance_id"`
	}
	_ = json.NewDecoder(response.Body).Decode(&body)

	return response.StatusCode, body.InstanceID
}

func servedInstances(loadBalancerURL string, requests int) (map[string]int, int) {
	served := make(map[string]int)
	failures := 0

	for range requests {
		code, instanceID := requestThroughLoadBalancer(loadBalancerURL + "/service-a/ping")
		if code != http.StatusOK {
			failures++
			continue
		}

		served[instanceID]++
	}

	return served, failures
}

func TestLoadBalancerRoutesOnlyToRegisteredInstances(t *testing.T) {
	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			c := newCluster(t, backend, 3)
			loadBalancer := startLoadBalancer(t, c.addresses())

			startMockService(t, "a1", c.addresses())
			a2 := startMockService(t, "a2", c.addresses())

			eventually(t, "the load balancer routes to both instances", func() bool {
				served, _ := servedInstances(loadBalancer, 10)
				return served["a1"] > 0 && served["a2"] > 0
			})

			a2.stop(t)

			eventually(t, "the load balancer only routes to the remaining instance", func() bool {
				served, failures := servedInstances(loadBalancer, 10)
				return failures == 0 && served["a1"] == 10
			})

			if code, _ := requestThroughLoadBalancer(loadBalancer + "/unknown-service/ping"); code != http.StatusNotFound {
				t.Errorf("request to an unknown service status = %d, want %d", code, http.StatusNotFound)
			}
		})
	}
}

func TestLoadBalancerKeepsRoutingWhenRegistryNodesFail(t *testing.T) {
	c := newCluster(t, "ap", 3)
	loadBalancer := startLoadBalancer(t, c.addresses(0, 1))

	startMockService(t, "a1", c.addresses(2))

	eventually(t, "the load balancer routes to a1", func() bool {
		served, _ := servedInstances(loadBalancer, 5)
		return served["a1"] > 0
	})

	failures := 0
	for i, stopped := range []int{0, 1} {
		c.stop(stopped)

		instanceID := fmt.Sprintf("a%d", i+2)
		startMockService(t, instanceID, c.addresses(2))

		eventually(t, "the load balancer learns about "+instanceID+" after a registry node failure", func() bool {
			served, failed := servedInstances(loadBalancer, 5)
			failures += failed
			return served[instanceID] > 0
		})

		c.start(stopped)
	}

	if failures != 0 {
		t.Errorf("failed requests while registry nodes were down = %d, want 0", failures)
	}
}
