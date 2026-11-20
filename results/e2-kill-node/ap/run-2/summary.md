# e2-kill-node — ap — run 2

## Faults

| t (s) | action | target |
|---|---|---|
| 30.8 | kill | registry-3 |
| 62.3 | start | registry-3 |

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
| after.killed | 99.32 |
| after.others | 100.00 |

## outage_after_fault_s

| metric | value |
|---|---|
| killed | 31.40 |
| others | 0.00 |

## outage_after_restore_s

| metric | value |
|---|---|
| killed | 1.31 |
| others | 0.00 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 5.13 |
| before.Register.p95 | 10.51 |
| before.Register.p99 | 14.28 |
| before.Heartbeat.p50 | 3.88 |
| before.Heartbeat.p95 | 7.57 |
| before.Heartbeat.p99 | 9.73 |
| before.Lookup.p50 | 3.90 |
| before.Lookup.p95 | 6.76 |
| before.Lookup.p99 | 9.08 |
| before.Deregister.p50 | 4.25 |
| before.Deregister.p95 | 8.12 |
| before.Deregister.p99 | 12.25 |
| during.Register.p50 | 4.68 |
| during.Register.p95 | 14.82 |
| during.Register.p99 | 24.88 |
| during.Heartbeat.p50 | 2.97 |
| during.Heartbeat.p95 | 6.16 |
| during.Heartbeat.p99 | 8.02 |
| during.Lookup.p50 | 3.08 |
| during.Lookup.p95 | 6.16 |
| during.Lookup.p99 | 10.62 |
| during.Deregister.p50 | 3.58 |
| during.Deregister.p95 | 14.74 |
| during.Deregister.p99 | 29.44 |
| after.Register.p50 | 4.64 |
| after.Register.p95 | 9.83 |
| after.Register.p99 | 13.09 |
| after.Heartbeat.p50 | 3.49 |
| after.Heartbeat.p95 | 6.70 |
| after.Heartbeat.p99 | 8.78 |
| after.Lookup.p50 | 3.64 |
| after.Lookup.p95 | 6.47 |
| after.Lookup.p99 | 8.70 |
| after.Deregister.p50 | 4.17 |
| after.Deregister.p95 | 7.96 |
| after.Deregister.p99 | 10.23 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.54 |
| other_nodes.p50 | 3.89 |
| other_nodes.p99 | 7.66 |
| other_nodes.max | 7.66 |

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
| convergence_after_restore_s | 1.41 |

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
| before.p99 | 15.19 |
| during.p99 | 14.17 |
| after.p99 | 14.33 |

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
