# e2-kill-node — ap — run 3

## Faults

| t (s) | action | target |
|---|---|---|
| 31.1 | kill | registry-3 |
| 62.0 | start | registry-3 |

## Node groups

| group | nodes |
|---|---|
| killed | registry-3 |
| others | registry-1 registry-2 registry-4 registry-5 |

## availability_pct

| metric | value |
|---|---|
| before.killed | 99.52 |
| before.others | 100.00 |
| during.killed | 0.00 |
| during.others | 100.00 |
| after.killed | 99.52 |
| after.others | 100.00 |

## outage_after_fault_s

| metric | value |
|---|---|
| killed | 30.93 |
| others | 0.00 |

## outage_after_restore_s

| metric | value |
|---|---|
| killed | 1.00 |
| others | 0.00 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 4.88 |
| before.Register.p95 | 10.75 |
| before.Register.p99 | 13.93 |
| before.Heartbeat.p50 | 3.62 |
| before.Heartbeat.p95 | 6.80 |
| before.Heartbeat.p99 | 8.78 |
| before.Lookup.p50 | 3.69 |
| before.Lookup.p95 | 6.92 |
| before.Lookup.p99 | 8.59 |
| before.Deregister.p50 | 4.09 |
| before.Deregister.p95 | 7.53 |
| before.Deregister.p99 | 9.81 |
| during.Register.p50 | 4.29 |
| during.Register.p95 | 8.84 |
| during.Register.p99 | 14.98 |
| during.Heartbeat.p50 | 3.03 |
| during.Heartbeat.p95 | 5.79 |
| during.Heartbeat.p99 | 8.43 |
| during.Lookup.p50 | 3.03 |
| during.Lookup.p95 | 6.00 |
| during.Lookup.p99 | 7.85 |
| during.Deregister.p50 | 3.40 |
| during.Deregister.p95 | 6.56 |
| during.Deregister.p99 | 17.72 |
| after.Register.p50 | 4.40 |
| after.Register.p95 | 9.06 |
| after.Register.p99 | 12.82 |
| after.Heartbeat.p50 | 3.56 |
| after.Heartbeat.p95 | 7.02 |
| after.Heartbeat.p99 | 9.48 |
| after.Lookup.p50 | 3.74 |
| after.Lookup.p95 | 6.73 |
| after.Lookup.p99 | 8.92 |
| after.Deregister.p50 | 3.98 |
| after.Deregister.p95 | 7.07 |
| after.Deregister.p99 | 8.58 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.50 |
| other_nodes.p50 | 4.09 |
| other_nodes.p99 | 5516.97 |
| other_nodes.max | 32422.38 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 25 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 25 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 1.29 |

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
| before.p99 | 15.30 |
| during.p99 | 13.36 |
| after.p99 | 13.68 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 0.00 |

## max_suspect_members

| metric | value |
|---|---|
| before | 0.00 |
| during | 4.00 |
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
