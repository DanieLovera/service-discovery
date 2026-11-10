# e4-loss-10pct — ap — run 2

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
| before.Register.p50 | 4.79 |
| before.Register.p95 | 9.72 |
| before.Register.p99 | 12.25 |
| before.Heartbeat.p50 | 3.57 |
| before.Heartbeat.p95 | 7.57 |
| before.Heartbeat.p99 | 9.50 |
| before.Lookup.p50 | 3.77 |
| before.Lookup.p95 | 6.95 |
| before.Lookup.p99 | 8.86 |
| before.Deregister.p50 | 4.38 |
| before.Deregister.p95 | 7.81 |
| before.Deregister.p99 | 10.54 |
| during.Register.p50 | 5.29 |
| during.Register.p95 | 217.07 |
| during.Register.p99 | 860.00 |
| during.Heartbeat.p50 | 3.46 |
| during.Heartbeat.p95 | 199.82 |
| during.Heartbeat.p99 | 236.40 |
| during.Lookup.p50 | 3.37 |
| during.Lookup.p95 | 195.57 |
| during.Lookup.p99 | 265.27 |
| during.Deregister.p50 | 3.98 |
| during.Deregister.p95 | 192.80 |
| during.Deregister.p99 | 233.09 |
| after.Register.p50 | 4.28 |
| after.Register.p95 | 9.01 |
| after.Register.p99 | 12.34 |
| after.Heartbeat.p50 | 3.59 |
| after.Heartbeat.p95 | 6.92 |
| after.Heartbeat.p99 | 9.42 |
| after.Lookup.p50 | 3.66 |
| after.Lookup.p95 | 6.85 |
| after.Lookup.p99 | 8.84 |
| after.Deregister.p50 | 4.14 |
| after.Deregister.p95 | 7.28 |
| after.Deregister.p99 | 8.75 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.65 |
| other_nodes.p50 | 4.13 |
| other_nodes.p99 | 7.08 |
| other_nodes.max | 7.29 |

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
| before.p99 | 15.20 |
| during.p99 | 13.52 |
| after.p99 | 14.48 |

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
