# e1-baseline — cp — run 1

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
| steady.Register.p50 | 39.34 |
| steady.Register.p95 | 62.27 |
| steady.Register.p99 | 183.95 |
| steady.Heartbeat.p50 | 8.81 |
| steady.Heartbeat.p95 | 16.60 |
| steady.Heartbeat.p99 | 28.86 |
| steady.Lookup.p50 | 8.09 |
| steady.Lookup.p95 | 18.63 |
| steady.Lookup.p99 | 32.14 |
| steady.Deregister.p50 | 31.61 |
| steady.Deregister.p95 | 57.83 |
| steady.Deregister.p99 | 99.51 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 86.70 |
| other_nodes.p50 | 85.06 |
| other_nodes.p99 | 207.72 |
| other_nodes.max | 302.18 |

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
| steady.p99 | 13.82 |

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
| -5.1 | registry-1 |
