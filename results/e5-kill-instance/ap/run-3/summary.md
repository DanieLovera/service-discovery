# e5-kill-instance — ap — run 3

## Faults

| t (s) | action | target |
|---|---|---|
| 30.0 | kill | service-c-3 |
| 61.3 | start | service-c-3 |

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
| before.Register.p50 | 4.75 |
| before.Register.p95 | 10.40 |
| before.Register.p99 | 12.32 |
| before.Heartbeat.p50 | 3.75 |
| before.Heartbeat.p95 | 7.10 |
| before.Heartbeat.p99 | 9.76 |
| before.Lookup.p50 | 3.80 |
| before.Lookup.p95 | 6.93 |
| before.Lookup.p99 | 8.75 |
| before.Deregister.p50 | 4.50 |
| before.Deregister.p95 | 7.99 |
| before.Deregister.p99 | 10.09 |
| during.Register.p50 | 4.90 |
| during.Register.p95 | 10.52 |
| during.Register.p99 | 13.17 |
| during.Heartbeat.p50 | 3.79 |
| during.Heartbeat.p95 | 7.76 |
| during.Heartbeat.p99 | 10.30 |
| during.Lookup.p50 | 3.93 |
| during.Lookup.p95 | 7.75 |
| during.Lookup.p99 | 10.17 |
| during.Deregister.p50 | 4.33 |
| during.Deregister.p95 | 8.13 |
| during.Deregister.p99 | 10.78 |
| after.Register.p50 | 4.78 |
| after.Register.p95 | 10.40 |
| after.Register.p99 | 14.64 |
| after.Heartbeat.p50 | 3.50 |
| after.Heartbeat.p95 | 7.19 |
| after.Heartbeat.p99 | 9.20 |
| after.Lookup.p50 | 3.76 |
| after.Lookup.p95 | 7.38 |
| after.Lookup.p99 | 9.51 |
| after.Deregister.p50 | 4.22 |
| after.Deregister.p95 | 8.01 |
| after.Deregister.p99 | 10.42 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.71 |
| other_nodes.p50 | 3.97 |
| other_nodes.p99 | 5.68 |
| other_nodes.max | 6.52 |

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
| before | 0 |
| during | 131 |
| after | 0 |

## lb_latency_ms

| metric | value |
|---|---|
| before.p99 | 16.23 |
| during.p99 | 20.20 |
| after.p99 | 14.01 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 13.01 |

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
