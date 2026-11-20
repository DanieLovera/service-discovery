# e4-delay-250ms — ap — run 2

## Faults

| t (s) | action | target |
|---|---|---|
| 30.5 | degrade | delay=250ms loss=0% |
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
| before.Register.p50 | 4.69 |
| before.Register.p95 | 10.58 |
| before.Register.p99 | 12.61 |
| before.Heartbeat.p50 | 3.77 |
| before.Heartbeat.p95 | 8.21 |
| before.Heartbeat.p99 | 10.72 |
| before.Lookup.p50 | 3.94 |
| before.Lookup.p95 | 7.69 |
| before.Lookup.p99 | 11.64 |
| before.Deregister.p50 | 4.52 |
| before.Deregister.p95 | 8.83 |
| before.Deregister.p99 | 11.59 |
| during.Register.p50 | 253.51 |
| during.Register.p95 | 255.77 |
| during.Register.p99 | 257.22 |
| during.Heartbeat.p50 | 252.38 |
| during.Heartbeat.p95 | 254.09 |
| during.Heartbeat.p99 | 255.23 |
| during.Lookup.p50 | 252.74 |
| during.Lookup.p95 | 254.48 |
| during.Lookup.p99 | 255.46 |
| during.Deregister.p50 | 252.45 |
| during.Deregister.p95 | 253.96 |
| during.Deregister.p99 | 255.02 |
| after.Register.p50 | 4.30 |
| after.Register.p95 | 9.42 |
| after.Register.p99 | 12.59 |
| after.Heartbeat.p50 | 3.56 |
| after.Heartbeat.p95 | 7.35 |
| after.Heartbeat.p99 | 9.28 |
| after.Lookup.p50 | 3.60 |
| after.Lookup.p95 | 7.29 |
| after.Lookup.p99 | 9.24 |
| after.Deregister.p50 | 3.92 |
| after.Deregister.p95 | 7.31 |
| after.Deregister.p99 | 9.19 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.84 |
| other_nodes.p50 | 4.30 |
| other_nodes.p99 | 5.68 |
| other_nodes.max | 6.36 |

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
| convergence_after_restore_s | 0.44 |

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
| before.p99 | 16.40 |
| during.p99 | 10.69 |
| after.p99 | 14.80 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 0.00 |

## max_suspect_members

| metric | value |
|---|---|
| before | 0.00 |
| during | 11.00 |
| after | 6.00 |

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
