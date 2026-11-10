# e4-delay-250ms — ap — run 1

## Faults

| t (s) | action | target |
|---|---|---|
| 30.9 | degrade | delay=250ms loss=0% |
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
| before.Register.p50 | 5.00 |
| before.Register.p95 | 10.73 |
| before.Register.p99 | 13.82 |
| before.Heartbeat.p50 | 3.70 |
| before.Heartbeat.p95 | 7.28 |
| before.Heartbeat.p99 | 9.70 |
| before.Lookup.p50 | 3.68 |
| before.Lookup.p95 | 6.90 |
| before.Lookup.p99 | 8.70 |
| before.Deregister.p50 | 3.99 |
| before.Deregister.p95 | 7.48 |
| before.Deregister.p99 | 10.02 |
| during.Register.p50 | 253.23 |
| during.Register.p95 | 255.32 |
| during.Register.p99 | 256.52 |
| during.Heartbeat.p50 | 252.48 |
| during.Heartbeat.p95 | 254.47 |
| during.Heartbeat.p99 | 255.92 |
| during.Lookup.p50 | 252.90 |
| during.Lookup.p95 | 254.89 |
| during.Lookup.p99 | 256.15 |
| during.Deregister.p50 | 252.50 |
| during.Deregister.p95 | 253.92 |
| during.Deregister.p99 | 255.12 |
| after.Register.p50 | 4.29 |
| after.Register.p95 | 8.98 |
| after.Register.p99 | 11.19 |
| after.Heartbeat.p50 | 3.55 |
| after.Heartbeat.p95 | 6.74 |
| after.Heartbeat.p99 | 9.51 |
| after.Lookup.p50 | 3.72 |
| after.Lookup.p95 | 6.81 |
| after.Lookup.p99 | 8.58 |
| after.Deregister.p50 | 3.99 |
| after.Deregister.p95 | 6.89 |
| after.Deregister.p99 | 8.08 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.68 |
| other_nodes.p50 | 4.11 |
| other_nodes.p99 | 546.80 |
| other_nodes.max | 754.74 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 60 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 3 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 1.23 |

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
| before.p99 | 15.86 |
| during.p99 | 10.72 |
| after.p99 | 13.52 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 0.00 |

## max_suspect_members

| metric | value |
|---|---|
| before | 0.00 |
| during | 10.00 |
| after | 3.00 |

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
