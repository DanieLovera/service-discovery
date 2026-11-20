# e5-kill-instance — ap — run 1

## Faults

| t (s) | action | target |
|---|---|---|
| 30.7 | kill | service-c-3 |
| 62.2 | start | service-c-3 |

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
| before.Register.p95 | 10.33 |
| before.Register.p99 | 12.83 |
| before.Heartbeat.p50 | 3.54 |
| before.Heartbeat.p95 | 6.79 |
| before.Heartbeat.p99 | 10.87 |
| before.Lookup.p50 | 3.67 |
| before.Lookup.p95 | 7.33 |
| before.Lookup.p99 | 9.70 |
| before.Deregister.p50 | 4.31 |
| before.Deregister.p95 | 8.47 |
| before.Deregister.p99 | 10.43 |
| during.Register.p50 | 4.77 |
| during.Register.p95 | 10.38 |
| during.Register.p99 | 15.74 |
| during.Heartbeat.p50 | 3.67 |
| during.Heartbeat.p95 | 7.39 |
| during.Heartbeat.p99 | 13.99 |
| during.Lookup.p50 | 3.88 |
| during.Lookup.p95 | 7.49 |
| during.Lookup.p99 | 9.55 |
| during.Deregister.p50 | 4.29 |
| during.Deregister.p95 | 8.12 |
| during.Deregister.p99 | 10.83 |
| after.Register.p50 | 4.77 |
| after.Register.p95 | 10.11 |
| after.Register.p99 | 14.70 |
| after.Heartbeat.p50 | 3.48 |
| after.Heartbeat.p95 | 7.06 |
| after.Heartbeat.p99 | 8.82 |
| after.Lookup.p50 | 4.00 |
| after.Lookup.p95 | 7.82 |
| after.Lookup.p99 | 11.17 |
| after.Deregister.p50 | 4.19 |
| after.Deregister.p95 | 8.19 |
| after.Deregister.p99 | 10.49 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.57 |
| other_nodes.p50 | 3.80 |
| other_nodes.p99 | 5.20 |
| other_nodes.max | 5.25 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 32 |

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
| before | 1 |
| during | 118 |
| after | 0 |

## lb_latency_ms

| metric | value |
|---|---|
| before.p99 | 15.13 |
| during.p99 | 19.66 |
| after.p99 | 15.69 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 11.73 |

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
