# e4-delay-100ms — cp — run 2

## Faults

| t (s) | action | target |
|---|---|---|
| 31.0 | degrade | delay=100ms loss=0% |
| 90.9 | restore |  |

## Node groups

| group | nodes |
|---|---|
| all | registry-1 registry-2 registry-3 registry-4 registry-5 |

## availability_pct

| metric | value |
|---|---|
| before.all | 100.00 |
| during.all | 100.00 |
| after.all | 100.00 |

## outage_after_fault_s

| metric | value |
|---|---|
| all | 0.00 |

## outage_after_restore_s

| metric | value |
|---|---|
| all | 0.00 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 39.09 |
| before.Register.p95 | 53.67 |
| before.Register.p99 | 61.08 |
| before.Heartbeat.p50 | 7.34 |
| before.Heartbeat.p95 | 12.90 |
| before.Heartbeat.p99 | 19.14 |
| before.Lookup.p50 | 7.71 |
| before.Lookup.p95 | 14.52 |
| before.Lookup.p99 | 22.39 |
| before.Deregister.p50 | 30.10 |
| before.Deregister.p95 | 47.38 |
| before.Deregister.p99 | 65.46 |
| during.Register.p50 | 563.42 |
| during.Register.p95 | 651.90 |
| during.Register.p99 | 739.32 |
| during.Heartbeat.p50 | 339.01 |
| during.Heartbeat.p95 | 425.35 |
| during.Heartbeat.p99 | 456.65 |
| during.Lookup.p50 | 333.30 |
| during.Lookup.p95 | 429.72 |
| during.Lookup.p99 | 453.32 |
| during.Deregister.p50 | 570.10 |
| during.Deregister.p95 | 661.55 |
| during.Deregister.p99 | 801.02 |
| after.Register.p50 | 41.29 |
| after.Register.p95 | 79.15 |
| after.Register.p99 | 257.50 |
| after.Heartbeat.p50 | 7.36 |
| after.Heartbeat.p95 | 14.96 |
| after.Heartbeat.p99 | 41.57 |
| after.Lookup.p50 | 7.61 |
| after.Lookup.p95 | 23.42 |
| after.Lookup.p99 | 63.13 |
| after.Deregister.p50 | 29.73 |
| after.Deregister.p95 | 50.82 |
| after.Deregister.p99 | 99.62 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 81.02 |
| other_nodes.p50 | 79.33 |
| other_nodes.p99 | 733.09 |
| other_nodes.max | 734.52 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 60 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 4 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 0.32 |

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
| before.p99 | 12.43 |
| during.p99 | 8.84 |
| after.p99 | 11.83 |

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
| -5.7 | registry-4 |
