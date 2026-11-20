# e4-loss-10pct — ap — run 3

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
| before.Register.p50 | 4.82 |
| before.Register.p95 | 10.50 |
| before.Register.p99 | 14.11 |
| before.Heartbeat.p50 | 3.84 |
| before.Heartbeat.p95 | 7.08 |
| before.Heartbeat.p99 | 9.83 |
| before.Lookup.p50 | 3.63 |
| before.Lookup.p95 | 6.57 |
| before.Lookup.p99 | 8.98 |
| before.Deregister.p50 | 4.08 |
| before.Deregister.p95 | 7.56 |
| before.Deregister.p99 | 9.58 |
| during.Register.p50 | 5.46 |
| during.Register.p95 | 213.25 |
| during.Register.p99 | 444.50 |
| during.Heartbeat.p50 | 3.57 |
| during.Heartbeat.p95 | 200.94 |
| during.Heartbeat.p99 | 248.96 |
| during.Lookup.p50 | 3.39 |
| during.Lookup.p95 | 191.89 |
| during.Lookup.p99 | 218.81 |
| during.Deregister.p50 | 3.89 |
| during.Deregister.p95 | 191.82 |
| during.Deregister.p99 | 232.17 |
| after.Register.p50 | 4.70 |
| after.Register.p95 | 10.13 |
| after.Register.p99 | 13.41 |
| after.Heartbeat.p50 | 3.56 |
| after.Heartbeat.p95 | 7.48 |
| after.Heartbeat.p99 | 8.74 |
| after.Lookup.p50 | 3.67 |
| after.Lookup.p95 | 6.49 |
| after.Lookup.p99 | 8.99 |
| after.Deregister.p50 | 4.05 |
| after.Deregister.p95 | 7.23 |
| after.Deregister.p99 | 9.07 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.58 |
| other_nodes.p50 | 3.92 |
| other_nodes.p99 | 5.29 |
| other_nodes.max | 6.25 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 60 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 0 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 0.00 |

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
| before.p99 | 15.73 |
| during.p99 | 13.98 |
| after.p99 | 13.77 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 0.00 |

## max_suspect_members

| metric | value |
|---|---|
| before | 0.00 |
| during | 1.00 |
| after | 0.00 |

## max_confirm_members

| metric | value |
|---|---|
| before | 0.00 |
| during | 0.00 |
| after | 0.00 |

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
