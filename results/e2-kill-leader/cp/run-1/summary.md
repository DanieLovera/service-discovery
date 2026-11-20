# e2-kill-leader — cp — run 1

## Faults

| t (s) | action | target |
|---|---|---|
| 31.5 | kill | registry-1 |
| 63.2 | start | registry-1 |

## Node groups

| group | nodes |
|---|---|
| killed | registry-1 |
| others | registry-2 registry-3 registry-4 registry-5 |

## availability_pct

| metric | value |
|---|---|
| before.killed | 100.00 |
| before.others | 99.77 |
| during.killed | 0.00 |
| during.others | 98.85 |
| after.killed | 99.61 |
| after.others | 100.00 |

## outage_after_fault_s

| metric | value |
|---|---|
| killed | 31.49 |
| others | 1.49 |

## outage_after_restore_s

| metric | value |
|---|---|
| killed | 0.61 |
| others | 0.00 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 42.13 |
| before.Register.p95 | 74.82 |
| before.Register.p99 | 92.51 |
| before.Heartbeat.p50 | 9.74 |
| before.Heartbeat.p95 | 20.81 |
| before.Heartbeat.p99 | 43.90 |
| before.Lookup.p50 | 8.59 |
| before.Lookup.p95 | 22.58 |
| before.Lookup.p99 | 34.15 |
| before.Deregister.p50 | 34.61 |
| before.Deregister.p95 | 64.54 |
| before.Deregister.p99 | 108.97 |
| during.Register.p50 | 37.18 |
| during.Register.p95 | 61.30 |
| during.Register.p99 | 76.09 |
| during.Heartbeat.p50 | 8.21 |
| during.Heartbeat.p95 | 14.64 |
| during.Heartbeat.p99 | 21.22 |
| during.Lookup.p50 | 7.32 |
| during.Lookup.p95 | 12.68 |
| during.Lookup.p99 | 19.18 |
| during.Deregister.p50 | 28.08 |
| during.Deregister.p95 | 49.86 |
| during.Deregister.p99 | 72.59 |
| after.Register.p50 | 38.17 |
| after.Register.p95 | 72.44 |
| after.Register.p99 | 137.31 |
| after.Heartbeat.p50 | 8.69 |
| after.Heartbeat.p95 | 15.21 |
| after.Heartbeat.p99 | 18.91 |
| after.Lookup.p50 | 8.04 |
| after.Lookup.p95 | 20.69 |
| after.Lookup.p99 | 32.82 |
| after.Deregister.p50 | 31.35 |
| after.Deregister.p95 | 58.46 |
| after.Deregister.p99 | 98.66 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 89.40 |
| other_nodes.p50 | 88.61 |
| other_nodes.p99 | 156.82 |
| other_nodes.max | 163.65 |

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
| convergence_after_restore_s | 0.87 |

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
| before.p99 | 19.70 |
| during.p99 | 12.46 |
| after.p99 | 12.58 |

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
| -5.1 | registry-1 |
| 39.9 | registry-2 |
