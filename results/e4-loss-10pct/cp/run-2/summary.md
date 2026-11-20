# e4-loss-10pct — cp — run 2

## Faults

| t (s) | action | target |
|---|---|---|
| 31.0 | degrade | delay=0ms loss=10% |
| 90.9 | restore |  |

## Node groups

| group | nodes |
|---|---|
| all | registry-1 registry-2 registry-3 registry-4 registry-5 |

## availability_pct

| metric | value |
|---|---|
| before.all | 100.00 |
| during.all | 99.83 |
| after.all | 100.00 |

## outage_after_fault_s

| metric | value |
|---|---|
| all | 6.24 |

## outage_after_restore_s

| metric | value |
|---|---|
| all | 0.00 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 40.74 |
| before.Register.p95 | 73.72 |
| before.Register.p99 | 99.61 |
| before.Heartbeat.p50 | 7.76 |
| before.Heartbeat.p95 | 15.93 |
| before.Heartbeat.p99 | 44.05 |
| before.Lookup.p50 | 7.69 |
| before.Lookup.p95 | 19.33 |
| before.Lookup.p99 | 49.27 |
| before.Deregister.p50 | 30.73 |
| before.Deregister.p95 | 60.58 |
| before.Deregister.p99 | 90.92 |
| during.Register.p50 | 68.90 |
| during.Register.p95 | 479.67 |
| during.Register.p99 | 907.09 |
| during.Heartbeat.p50 | 23.56 |
| during.Heartbeat.p95 | 240.73 |
| during.Heartbeat.p99 | 457.23 |
| during.Lookup.p50 | 23.55 |
| during.Lookup.p95 | 241.04 |
| during.Lookup.p99 | 646.53 |
| during.Deregister.p50 | 43.20 |
| during.Deregister.p95 | 263.09 |
| during.Deregister.p99 | 629.62 |
| after.Register.p50 | 39.26 |
| after.Register.p95 | 75.27 |
| after.Register.p99 | 201.77 |
| after.Heartbeat.p50 | 7.14 |
| after.Heartbeat.p95 | 14.11 |
| after.Heartbeat.p99 | 23.47 |
| after.Lookup.p50 | 7.45 |
| after.Lookup.p95 | 16.86 |
| after.Lookup.p99 | 33.39 |
| after.Deregister.p50 | 28.92 |
| after.Deregister.p95 | 61.65 |
| after.Deregister.p99 | 98.54 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 86.64 |
| other_nodes.p50 | 84.38 |
| other_nodes.p99 | 133.33 |
| other_nodes.max | 137.32 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 59 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 4 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 0.48 |

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
| before.p99 | 13.79 |
| during.p99 | 10.89 |
| after.p99 | 11.82 |

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
| -5.8 | registry-4 |
