# e1-baseline — ap — run 2

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
| steady.Register.p50 | 4.63 |
| steady.Register.p95 | 9.68 |
| steady.Register.p99 | 12.76 |
| steady.Heartbeat.p50 | 3.60 |
| steady.Heartbeat.p95 | 7.10 |
| steady.Heartbeat.p99 | 10.19 |
| steady.Lookup.p50 | 3.85 |
| steady.Lookup.p95 | 7.13 |
| steady.Lookup.p99 | 9.47 |
| steady.Deregister.p50 | 4.07 |
| steady.Deregister.p95 | 7.66 |
| steady.Deregister.p99 | 10.27 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.72 |
| other_nodes.p50 | 4.08 |
| other_nodes.p99 | 6.54 |
| other_nodes.max | 7.46 |

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
| steady.p99 | 14.61 |

## max_suspect_members

| metric | value |
|---|---|
| steady | 0.00 |

## max_confirm_members

| metric | value |
|---|---|
| steady | 0.00 |

## min_healthy_instances

| metric | value |
|---|---|
| steady | 9.00 |

## min_lb_pool_instances

| metric | value |
|---|---|
| steady | 3.00 |
