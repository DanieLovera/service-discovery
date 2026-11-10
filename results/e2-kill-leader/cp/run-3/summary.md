# e2-kill-leader — cp — run 3

## Faults

| t (s) | action | target |
|---|---|---|
| 30.8 | kill | registry-3 |
| 62.3 | start | registry-3 |

## Node groups

| group | nodes |
|---|---|
| killed | registry-3 |
| others | registry-1 registry-2 registry-4 registry-5 |

## availability_pct

| metric | value |
|---|---|
| before.killed | 100.00 |
| before.others | 100.00 |
| during.killed | 0.00 |
| during.others | 98.77 |
| after.killed | 99.42 |
| after.others | 100.00 |

## outage_after_fault_s

| metric | value |
|---|---|
| killed | 31.40 |
| others | 1.60 |

## outage_after_restore_s

| metric | value |
|---|---|
| killed | 1.09 |
| others | 0.00 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 39.88 |
| before.Register.p95 | 54.57 |
| before.Register.p99 | 68.70 |
| before.Heartbeat.p50 | 7.88 |
| before.Heartbeat.p95 | 14.78 |
| before.Heartbeat.p99 | 41.64 |
| before.Lookup.p50 | 7.74 |
| before.Lookup.p95 | 17.57 |
| before.Lookup.p99 | 37.21 |
| before.Deregister.p50 | 30.52 |
| before.Deregister.p95 | 50.99 |
| before.Deregister.p99 | 74.33 |
| during.Register.p50 | 36.27 |
| during.Register.p95 | 51.25 |
| during.Register.p99 | 75.78 |
| during.Heartbeat.p50 | 6.62 |
| during.Heartbeat.p95 | 11.65 |
| during.Heartbeat.p99 | 19.76 |
| during.Lookup.p50 | 6.68 |
| during.Lookup.p95 | 12.70 |
| during.Lookup.p99 | 21.98 |
| during.Deregister.p50 | 24.02 |
| during.Deregister.p95 | 36.87 |
| during.Deregister.p99 | 60.65 |
| after.Register.p50 | 39.38 |
| after.Register.p95 | 68.37 |
| after.Register.p99 | 163.68 |
| after.Heartbeat.p50 | 7.11 |
| after.Heartbeat.p95 | 13.66 |
| after.Heartbeat.p99 | 24.29 |
| after.Lookup.p50 | 7.59 |
| after.Lookup.p95 | 17.82 |
| after.Lookup.p99 | 28.31 |
| after.Deregister.p50 | 28.61 |
| after.Deregister.p95 | 50.44 |
| after.Deregister.p99 | 76.05 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 81.95 |
| other_nodes.p50 | 84.52 |
| other_nodes.p99 | 123.51 |
| other_nodes.max | 125.65 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 24 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 24 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 1.22 |

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
| before.p99 | 13.32 |
| during.p99 | 12.28 |
| after.p99 | 12.46 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 0.00 |

## leader_changes

| metric | value |
|---|---|
| leader_changes | 1 |

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
| -6.0 | registry-3 |
| 34.0 | registry-1 |
