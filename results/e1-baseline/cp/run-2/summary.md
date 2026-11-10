# e1-baseline — cp — run 2

## Node groups

| group | nodes |
|---|---|
| all | registry-1 registry-2 registry-3 registry-4 registry-5 |

## availability_pct

| metric | value |
|---|---|
| steady.all | 100.00 |

## latency_ms

| metric | value |
|---|---|
| steady.Register.p50 | 39.49 |
| steady.Register.p95 | 63.40 |
| steady.Register.p99 | 153.41 |
| steady.Heartbeat.p50 | 8.03 |
| steady.Heartbeat.p95 | 15.41 |
| steady.Heartbeat.p99 | 25.60 |
| steady.Lookup.p50 | 7.97 |
| steady.Lookup.p95 | 16.81 |
| steady.Lookup.p99 | 26.92 |
| steady.Deregister.p50 | 30.03 |
| steady.Deregister.p95 | 47.74 |
| steady.Deregister.p99 | 72.78 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 94.02 |
| other_nodes.p50 | 90.15 |
| other_nodes.p99 | 139.37 |
| other_nodes.max | 152.61 |

## acknowledged_writes_never_replicated

| metric | value |
|---|---|
| acknowledged_writes_never_replicated | 0 |

## lb_errors

| metric | value |
|---|---|
| steady | 0 |

## lb_latency_ms

| metric | value |
|---|---|
| steady.p99 | 14.80 |

## leader_changes

| metric | value |
|---|---|
| leader_changes | 0 |

## min_healthy_instances

| metric | value |
|---|---|
| steady | 9.00 |

## min_lb_pool_instances

| metric | value |
|---|---|
| steady | 3.00 |

## Raft leader changes

| t (s) | leader |
|---|---|
| -5.3 | registry-1 |
