# e4-delay-250ms — ap — run 3

## Faults

| t (s) | action | target |
|---|---|---|
| 30.8 | degrade | delay=250ms loss=0% |
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
| before.Register.p50 | 4.50 |
| before.Register.p95 | 9.84 |
| before.Register.p99 | 21.63 |
| before.Heartbeat.p50 | 3.83 |
| before.Heartbeat.p95 | 8.22 |
| before.Heartbeat.p99 | 9.78 |
| before.Lookup.p50 | 3.78 |
| before.Lookup.p95 | 7.17 |
| before.Lookup.p99 | 9.11 |
| before.Deregister.p50 | 4.16 |
| before.Deregister.p95 | 7.87 |
| before.Deregister.p99 | 10.27 |
| during.Register.p50 | 253.28 |
| during.Register.p95 | 255.55 |
| during.Register.p99 | 257.09 |
| during.Heartbeat.p50 | 252.34 |
| during.Heartbeat.p95 | 254.24 |
| during.Heartbeat.p99 | 255.66 |
| during.Lookup.p50 | 252.81 |
| during.Lookup.p95 | 254.79 |
| during.Lookup.p99 | 256.40 |
| during.Deregister.p50 | 252.50 |
| during.Deregister.p95 | 254.04 |
| during.Deregister.p99 | 255.23 |
| after.Register.p50 | 4.25 |
| after.Register.p95 | 9.11 |
| after.Register.p99 | 12.05 |
| after.Heartbeat.p50 | 3.62 |
| after.Heartbeat.p95 | 7.07 |
| after.Heartbeat.p99 | 9.12 |
| after.Lookup.p50 | 3.77 |
| after.Lookup.p95 | 6.81 |
| after.Lookup.p99 | 8.07 |
| after.Deregister.p50 | 3.99 |
| after.Deregister.p95 | 7.32 |
| after.Deregister.p99 | 9.48 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.61 |
| other_nodes.p50 | 4.26 |
| other_nodes.p99 | 8.16 |
| other_nodes.max | 12.61 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 61 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 8 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 5.64 |

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
| during.p99 | 10.89 |
| after.p99 | 13.52 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 0.00 |

## max_suspect_members

| metric | value |
|---|---|
| before | 0.00 |
| during | 13.00 |
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
