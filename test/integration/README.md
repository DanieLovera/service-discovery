# Integration tests

These tests start real project components inside the test process: 3-node registry clusters for both backends
(CP and AP), mock services and the load balancer, all on free localhost ports. No Docker is required.

- `registry_test.go`: instance lifecycle, TTL expiration and recovery by heartbeat, and Watch events, on both
  backends.
- `cp_test.go`: writes read from every node right away, leader failover, and a minority that refuses operations
  until the majority returns.
- `ap_test.go`: availability without a majority, state sync of a restarted node before it serves, convergence of
  concurrent updates, and SWIM failure detection.
- `load_balancer_test.go`: end-to-end routing through the load balancer to mock services, and routing that
  continues while registry nodes fail.
- `cluster_test.go`: shared helpers to start, stop and query the processes.

Run only these tests with `make test-integration`. They also run with `make test-race` and `make quality`, so CI
runs them on every push.
