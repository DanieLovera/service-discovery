# e4-delay-100ms — ap — run 2

## Faults

| t (s) | action | target |
|---|---|---|
| 30.9 | degrade | delay=100ms loss=0% |
| 90.8 | restore |  |

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
| before.Register.p50 | 4.71 |
| before.Register.p95 | 10.32 |
| before.Register.p99 | 14.36 |
| before.Heartbeat.p50 | 3.62 |
| before.Heartbeat.p95 | 7.55 |
| before.Heartbeat.p99 | 9.18 |
| before.Lookup.p50 | 3.55 |
| before.Lookup.p95 | 7.16 |
| before.Lookup.p99 | 9.49 |
| before.Deregister.p50 | 4.29 |
| before.Deregister.p95 | 8.26 |
| before.Deregister.p99 | 10.15 |
| during.Register.p50 | 103.75 |
| during.Register.p95 | 106.94 |
| during.Register.p99 | 109.21 |
| during.Heartbeat.p50 | 102.80 |
| during.Heartbeat.p95 | 105.64 |
| during.Heartbeat.p99 | 108.85 |
| during.Lookup.p50 | 102.91 |
| during.Lookup.p95 | 106.24 |
| during.Lookup.p99 | 109.11 |
| during.Deregister.p50 | 102.98 |
| during.Deregister.p95 | 105.89 |
| during.Deregister.p99 | 107.63 |
| after.Register.p50 | 4.45 |
| after.Register.p95 | 9.41 |
| after.Register.p99 | 12.10 |
| after.Heartbeat.p50 | 3.62 |
| after.Heartbeat.p95 | 7.60 |
| after.Heartbeat.p99 | 9.96 |
| after.Lookup.p50 | 4.07 |
| after.Lookup.p95 | 7.42 |
| after.Lookup.p99 | 9.49 |
| after.Deregister.p50 | 4.29 |
| after.Deregister.p95 | 7.77 |
| after.Deregister.p99 | 9.71 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.85 |
| other_nodes.p50 | 4.19 |
| other_nodes.p99 | 221.48 |
| other_nodes.max | 305.10 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 59 |

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
| before.p99 | 14.41 |
| during.p99 | 13.06 |
| after.p99 | 15.06 |

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
