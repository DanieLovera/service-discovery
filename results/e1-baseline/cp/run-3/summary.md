# e1-baseline — cp — run 3

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
| steady.Register.p50 | 39.83 |
| steady.Register.p95 | 73.17 |
| steady.Register.p99 | 232.86 |
| steady.Heartbeat.p50 | 7.40 |
| steady.Heartbeat.p95 | 12.94 |
| steady.Heartbeat.p99 | 22.91 |
| steady.Lookup.p50 | 7.46 |
| steady.Lookup.p95 | 16.55 |
| steady.Lookup.p99 | 26.13 |
| steady.Deregister.p50 | 29.49 |
| steady.Deregister.p95 | 65.19 |
| steady.Deregister.p99 | 106.00 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 78.64 |
| other_nodes.p50 | 85.58 |
| other_nodes.p99 | 140.83 |
| other_nodes.max | 203.40 |

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
| steady.p99 | 12.48 |

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
| -5.5 | registry-4 |
