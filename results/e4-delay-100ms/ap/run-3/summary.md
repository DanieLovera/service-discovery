# e4-delay-100ms — ap — run 3

## Faults

| t (s) | action | target |
|---|---|---|
| 30.9 | degrade | delay=100ms loss=0% |
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
| before.Register.p50 | 4.83 |
| before.Register.p95 | 10.85 |
| before.Register.p99 | 14.41 |
| before.Heartbeat.p50 | 3.91 |
| before.Heartbeat.p95 | 7.80 |
| before.Heartbeat.p99 | 9.55 |
| before.Lookup.p50 | 3.96 |
| before.Lookup.p95 | 7.24 |
| before.Lookup.p99 | 9.73 |
| before.Deregister.p50 | 4.23 |
| before.Deregister.p95 | 7.99 |
| before.Deregister.p99 | 10.59 |
| during.Register.p50 | 103.86 |
| during.Register.p95 | 106.86 |
| during.Register.p99 | 108.79 |
| during.Heartbeat.p50 | 102.79 |
| during.Heartbeat.p95 | 105.36 |
| during.Heartbeat.p99 | 107.15 |
| during.Lookup.p50 | 102.82 |
| during.Lookup.p95 | 105.46 |
| during.Lookup.p99 | 107.90 |
| during.Deregister.p50 | 102.86 |
| during.Deregister.p95 | 105.12 |
| during.Deregister.p99 | 106.76 |
| after.Register.p50 | 4.42 |
| after.Register.p95 | 9.86 |
| after.Register.p99 | 13.46 |
| after.Heartbeat.p50 | 3.48 |
| after.Heartbeat.p95 | 7.35 |
| after.Heartbeat.p99 | 9.06 |
| after.Lookup.p50 | 3.73 |
| after.Lookup.p95 | 6.64 |
| after.Lookup.p99 | 8.12 |
| after.Deregister.p50 | 4.08 |
| after.Deregister.p95 | 7.47 |
| after.Deregister.p99 | 9.51 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.98 |
| other_nodes.p50 | 4.36 |
| other_nodes.p99 | 221.72 |
| other_nodes.max | 304.73 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 60 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 2 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 0.06 |

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
| before.p99 | 15.68 |
| during.p99 | 12.35 |
| after.p99 | 13.52 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 0.00 |

## max_suspect_members

| metric | value |
|---|---|
| before | 0.00 |
| during | 0.00 |
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
