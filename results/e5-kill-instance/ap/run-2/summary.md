# e5-kill-instance — ap — run 2

## Faults

| t (s) | action | target |
|---|---|---|
| 30.9 | kill | service-c-3 |
| 62.3 | start | service-c-3 |

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
| before.Register.p50 | 4.84 |
| before.Register.p95 | 11.34 |
| before.Register.p99 | 14.19 |
| before.Heartbeat.p50 | 3.92 |
| before.Heartbeat.p95 | 7.99 |
| before.Heartbeat.p99 | 11.31 |
| before.Lookup.p50 | 3.75 |
| before.Lookup.p95 | 7.68 |
| before.Lookup.p99 | 9.54 |
| before.Deregister.p50 | 4.09 |
| before.Deregister.p95 | 8.46 |
| before.Deregister.p99 | 12.70 |
| during.Register.p50 | 4.34 |
| during.Register.p95 | 10.12 |
| during.Register.p99 | 15.76 |
| during.Heartbeat.p50 | 3.75 |
| during.Heartbeat.p95 | 7.08 |
| during.Heartbeat.p99 | 9.42 |
| during.Lookup.p50 | 3.90 |
| during.Lookup.p95 | 7.41 |
| during.Lookup.p99 | 9.03 |
| during.Deregister.p50 | 4.14 |
| during.Deregister.p95 | 8.04 |
| during.Deregister.p99 | 11.43 |
| after.Register.p50 | 4.27 |
| after.Register.p95 | 9.95 |
| after.Register.p99 | 13.22 |
| after.Heartbeat.p50 | 3.75 |
| after.Heartbeat.p95 | 7.01 |
| after.Heartbeat.p99 | 8.92 |
| after.Lookup.p50 | 3.92 |
| after.Lookup.p95 | 7.02 |
| after.Lookup.p99 | 9.07 |
| after.Deregister.p50 | 4.02 |
| after.Deregister.p95 | 7.86 |
| after.Deregister.p99 | 9.96 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.58 |
| other_nodes.p50 | 4.17 |
| other_nodes.p99 | 8.65 |
| other_nodes.max | 9.12 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 31 |

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
| before | 2 |
| during | 116 |
| after | 0 |

## lb_latency_ms

| metric | value |
|---|---|
| before.p99 | 16.56 |
| during.p99 | 19.76 |
| after.p99 | 14.58 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 11.54 |

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
| during | 8.00 |
| after | 8.00 |

## min_lb_pool_instances

| metric | value |
|---|---|
| before | 3.00 |
| during | 2.00 |
| after | 2.00 |
