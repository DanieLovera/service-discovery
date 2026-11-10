# e4-loss-10pct — ap — run 1

## Faults

| t (s) | action | target |
|---|---|---|
| 31.3 | degrade | delay=0ms loss=10% |
| 91.3 | restore |  |

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
| before.Register.p50 | 4.31 |
| before.Register.p95 | 10.06 |
| before.Register.p99 | 15.07 |
| before.Heartbeat.p50 | 4.01 |
| before.Heartbeat.p95 | 8.54 |
| before.Heartbeat.p99 | 10.13 |
| before.Lookup.p50 | 3.73 |
| before.Lookup.p95 | 7.55 |
| before.Lookup.p99 | 12.25 |
| before.Deregister.p50 | 3.92 |
| before.Deregister.p95 | 8.40 |
| before.Deregister.p99 | 10.96 |
| during.Register.p50 | 5.12 |
| during.Register.p95 | 214.21 |
| during.Register.p99 | 417.26 |
| during.Heartbeat.p50 | 3.64 |
| during.Heartbeat.p95 | 199.15 |
| during.Heartbeat.p99 | 218.33 |
| during.Lookup.p50 | 3.50 |
| during.Lookup.p95 | 197.25 |
| during.Lookup.p99 | 399.72 |
| during.Deregister.p50 | 3.80 |
| during.Deregister.p95 | 192.60 |
| during.Deregister.p99 | 231.70 |
| after.Register.p50 | 4.38 |
| after.Register.p95 | 9.11 |
| after.Register.p99 | 12.56 |
| after.Heartbeat.p50 | 3.42 |
| after.Heartbeat.p95 | 6.99 |
| after.Heartbeat.p99 | 8.88 |
| after.Lookup.p50 | 3.81 |
| after.Lookup.p95 | 6.66 |
| after.Lookup.p99 | 7.86 |
| after.Deregister.p50 | 3.87 |
| after.Deregister.p95 | 6.61 |
| after.Deregister.p99 | 8.17 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.74 |
| other_nodes.p50 | 4.34 |
| other_nodes.p99 | 9.20 |
| other_nodes.max | 9.24 |

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
| before.p99 | 14.69 |
| during.p99 | 13.38 |
| after.p99 | 13.20 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 0.00 |

## max_suspect_members

| metric | value |
|---|---|
| before | 0.00 |
| during | 3.00 |
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
