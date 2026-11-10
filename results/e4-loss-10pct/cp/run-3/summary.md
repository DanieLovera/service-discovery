# e4-loss-10pct — cp — run 3

## Faults

| t (s) | action | target |
|---|---|---|
| 30.9 | degrade | delay=0ms loss=10% |
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
| before.Register.p50 | 40.33 |
| before.Register.p95 | 67.32 |
| before.Register.p99 | 91.93 |
| before.Heartbeat.p50 | 7.71 |
| before.Heartbeat.p95 | 14.19 |
| before.Heartbeat.p99 | 19.31 |
| before.Lookup.p50 | 7.87 |
| before.Lookup.p95 | 15.53 |
| before.Lookup.p99 | 29.64 |
| before.Deregister.p50 | 29.57 |
| before.Deregister.p95 | 48.67 |
| before.Deregister.p99 | 95.40 |
| during.Register.p50 | 61.06 |
| during.Register.p95 | 524.81 |
| during.Register.p99 | 1067.95 |
| during.Heartbeat.p50 | 21.06 |
| during.Heartbeat.p95 | 242.43 |
| during.Heartbeat.p99 | 526.29 |
| during.Lookup.p50 | 22.74 |
| during.Lookup.p95 | 242.71 |
| during.Lookup.p99 | 481.00 |
| during.Deregister.p50 | 39.17 |
| during.Deregister.p95 | 260.77 |
| during.Deregister.p99 | 656.45 |
| after.Register.p50 | 38.84 |
| after.Register.p95 | 94.81 |
| after.Register.p99 | 289.69 |
| after.Heartbeat.p50 | 7.38 |
| after.Heartbeat.p95 | 18.14 |
| after.Heartbeat.p99 | 56.42 |
| after.Lookup.p50 | 7.66 |
| after.Lookup.p95 | 18.15 |
| after.Lookup.p99 | 99.55 |
| after.Deregister.p50 | 28.45 |
| after.Deregister.p95 | 61.15 |
| after.Deregister.p99 | 132.56 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 104.52 |
| other_nodes.p50 | 92.31 |
| other_nodes.p99 | 158.74 |
| other_nodes.max | 250.73 |

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
| convergence_after_restore_s | 0.55 |

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
| before.p99 | 12.93 |
| during.p99 | 10.96 |
| after.p99 | 11.85 |

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
| -5.3 | registry-1 |
