# e2-kill-node — ap — run 1

## Faults

| t (s) | action | target |
|---|---|---|
| 30.8 | kill | registry-3 |
| 62.2 | start | registry-3 |

## Node groups

| group | nodes |
|---|---|
| killed | registry-3 |
| others | registry-1 registry-2 registry-4 registry-5 |

## availability_pct

| metric | value |
|---|---|
| before.killed | 100.00 |
| before.others | 100.00 |
| during.killed | 0.00 |
| during.others | 100.00 |
| after.killed | 99.52 |
| after.others | 100.00 |

## outage_after_fault_s

| metric | value |
|---|---|
| killed | 31.43 |
| others | 0.00 |

## outage_after_restore_s

| metric | value |
|---|---|
| killed | 0.97 |
| others | 0.00 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 4.55 |
| before.Register.p95 | 9.72 |
| before.Register.p99 | 11.78 |
| before.Heartbeat.p50 | 3.71 |
| before.Heartbeat.p95 | 8.22 |
| before.Heartbeat.p99 | 11.15 |
| before.Lookup.p50 | 3.82 |
| before.Lookup.p95 | 6.92 |
| before.Lookup.p99 | 9.08 |
| before.Deregister.p50 | 3.98 |
| before.Deregister.p95 | 7.14 |
| before.Deregister.p99 | 9.68 |
| during.Register.p50 | 4.12 |
| during.Register.p95 | 7.96 |
| during.Register.p99 | 10.95 |
| during.Heartbeat.p50 | 2.95 |
| during.Heartbeat.p95 | 6.24 |
| during.Heartbeat.p99 | 7.87 |
| during.Lookup.p50 | 3.28 |
| during.Lookup.p95 | 6.17 |
| during.Lookup.p99 | 8.49 |
| during.Deregister.p50 | 3.52 |
| during.Deregister.p95 | 6.89 |
| during.Deregister.p99 | 10.77 |
| after.Register.p50 | 4.39 |
| after.Register.p95 | 9.00 |
| after.Register.p99 | 12.04 |
| after.Heartbeat.p50 | 3.54 |
| after.Heartbeat.p95 | 6.81 |
| after.Heartbeat.p99 | 8.90 |
| after.Lookup.p50 | 3.72 |
| after.Lookup.p95 | 6.96 |
| after.Lookup.p99 | 9.74 |
| after.Deregister.p50 | 4.13 |
| after.Deregister.p95 | 7.70 |
| after.Deregister.p99 | 9.97 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.59 |
| other_nodes.p50 | 4.00 |
| other_nodes.p99 | 5.18 |
| other_nodes.max | 5.33 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 26 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 26 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 1.05 |

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
| before.p99 | 14.60 |
| during.p99 | 12.71 |
| after.p99 | 14.22 |

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
| during | 4.00 |
| after | 3.00 |

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
