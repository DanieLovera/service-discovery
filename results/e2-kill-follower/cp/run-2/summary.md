# e2-kill-follower — cp — run 2

## Faults

| t (s) | action | target |
|---|---|---|
| 30.8 | kill | registry-1 |
| 62.3 | start | registry-1 |

## Node groups

| group | nodes |
|---|---|
| killed | registry-1 |
| others | registry-2 registry-3 registry-4 registry-5 |

## availability_pct

| metric | value |
|---|---|
| before.killed | 99.76 |
| before.others | 100.00 |
| during.killed | 0.00 |
| during.others | 100.00 |
| after.killed | 99.52 |
| after.others | 100.00 |

## outage_after_fault_s

| metric | value |
|---|---|
| killed | 31.39 |
| others | 0.00 |

## outage_after_restore_s

| metric | value |
|---|---|
| killed | 0.88 |
| others | 0.00 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 39.74 |
| before.Register.p95 | 59.35 |
| before.Register.p99 | 105.30 |
| before.Heartbeat.p50 | 7.07 |
| before.Heartbeat.p95 | 14.31 |
| before.Heartbeat.p99 | 17.68 |
| before.Lookup.p50 | 7.91 |
| before.Lookup.p95 | 19.90 |
| before.Lookup.p99 | 23.75 |
| before.Deregister.p50 | 30.33 |
| before.Deregister.p95 | 55.39 |
| before.Deregister.p99 | 95.29 |
| during.Register.p50 | 35.22 |
| during.Register.p95 | 59.87 |
| during.Register.p99 | 88.86 |
| during.Heartbeat.p50 | 6.12 |
| during.Heartbeat.p95 | 12.45 |
| during.Heartbeat.p99 | 17.30 |
| during.Lookup.p50 | 6.46 |
| during.Lookup.p95 | 18.18 |
| during.Lookup.p99 | 29.24 |
| during.Deregister.p50 | 24.78 |
| during.Deregister.p95 | 37.89 |
| during.Deregister.p99 | 54.30 |
| after.Register.p50 | 37.64 |
| after.Register.p95 | 79.62 |
| after.Register.p99 | 181.44 |
| after.Heartbeat.p50 | 6.92 |
| after.Heartbeat.p95 | 13.83 |
| after.Heartbeat.p99 | 19.36 |
| after.Lookup.p50 | 7.59 |
| after.Lookup.p95 | 20.38 |
| after.Lookup.p99 | 32.61 |
| after.Deregister.p50 | 28.93 |
| after.Deregister.p95 | 56.63 |
| after.Deregister.p99 | 110.12 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 89.07 |
| other_nodes.p50 | 88.91 |
| other_nodes.p99 | 128.36 |
| other_nodes.max | 130.48 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 25 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 25 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 0.99 |

## pairs_never_seen

| metric | value |
|---|---|
| pairs_never_seen | 0 |

## acknowledged_writes_never_replicated

| metric | value |
|---|---|
| acknowledged_writes_never_replicated | 0 |

## lb_errors

| metric | value |
|---|---|
| before | 0 |
| during | 0 |
| after | 0 |

## lb_latency_ms

| metric | value |
|---|---|
| before.p99 | 13.60 |
| during.p99 | 10.57 |
| after.p99 | 10.82 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 0.00 |

## leader_changes

| metric | value |
|---|---|
| leader_changes | 0 |

## min_healthy_instances

| metric | value |
|---|---|
| before | 9.00 |
| during | 9.00 |
| after | 9.00 |

## min_lb_pool_instances

| metric | value |
|---|---|
| before | 3.00 |
| during | 3.00 |
| after | 3.00 |

## Raft leader changes

| t (s) | leader |
|---|---|
| -5.3 | registry-3 |
