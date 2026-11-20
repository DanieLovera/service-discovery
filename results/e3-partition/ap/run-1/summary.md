# e3-partition — ap — run 1

## Faults

| t (s) | action | target |
|---|---|---|
| 31.2 | partition | registry-1 registry-2 | registry-3 registry-4 registry-5 |
| 90.9 | heal |  |

## Node groups

| group | nodes |
|---|---|
| minority | registry-1 registry-2 |
| majority | registry-3 registry-4 registry-5 |

## availability_pct

| metric | value |
|---|---|
| before.minority | 100.00 |
| before.majority | 100.00 |
| during.minority | 100.00 |
| during.majority | 100.00 |
| after.minority | 100.00 |
| after.majority | 100.00 |

## outage_after_fault_s

| metric | value |
|---|---|
| minority | 0.00 |
| majority | 0.00 |

## outage_after_restore_s

| metric | value |
|---|---|
| minority | 0.00 |
| majority | 0.00 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 4.89 |
| before.Register.p95 | 10.25 |
| before.Register.p99 | 14.53 |
| before.Heartbeat.p50 | 3.44 |
| before.Heartbeat.p95 | 7.16 |
| before.Heartbeat.p99 | 8.87 |
| before.Lookup.p50 | 3.79 |
| before.Lookup.p95 | 7.31 |
| before.Lookup.p99 | 9.24 |
| before.Deregister.p50 | 4.27 |
| before.Deregister.p95 | 7.76 |
| before.Deregister.p99 | 9.88 |
| during.Register.p50 | 3.89 |
| during.Register.p95 | 6.95 |
| during.Register.p99 | 8.61 |
| during.Heartbeat.p50 | 3.20 |
| during.Heartbeat.p95 | 5.94 |
| during.Heartbeat.p99 | 7.68 |
| during.Lookup.p50 | 3.58 |
| during.Lookup.p95 | 6.38 |
| during.Lookup.p99 | 7.72 |
| during.Deregister.p50 | 3.36 |
| during.Deregister.p95 | 6.03 |
| during.Deregister.p99 | 7.51 |
| after.Register.p50 | 4.47 |
| after.Register.p95 | 8.57 |
| after.Register.p99 | 11.13 |
| after.Heartbeat.p50 | 3.39 |
| after.Heartbeat.p95 | 6.60 |
| after.Heartbeat.p99 | 8.42 |
| after.Lookup.p50 | 3.51 |
| after.Lookup.p95 | 6.44 |
| after.Lookup.p99 | 8.49 |
| after.Deregister.p50 | 4.01 |
| after.Deregister.p95 | 7.20 |
| after.Deregister.p99 | 8.81 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.17 |
| other_nodes.p50 | 3.53 |
| other_nodes.p99 | 4.85 |
| other_nodes.max | 4.86 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 60 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 141 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 15.93 |

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
| before.p99 | 14.27 |
| during.p99 | 13.41 |
| after.p99 | 13.43 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 0.00 |

## max_suspect_members

| metric | value |
|---|---|
| before | 0.00 |
| during | 10.00 |
| after | 0.00 |

## max_confirm_members

| metric | value |
|---|---|
| before | 0.00 |
| during | 12.00 |
| after | 6.00 |

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
