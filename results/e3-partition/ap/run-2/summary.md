# e3-partition — ap — run 2

## Faults

| t (s) | action | target |
|---|---|---|
| 30.6 | partition | registry-1 registry-2 | registry-3 registry-4 registry-5 |
| 91.7 | heal |  |

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
| before.Register.p50 | 4.79 |
| before.Register.p95 | 10.30 |
| before.Register.p99 | 14.34 |
| before.Heartbeat.p50 | 3.72 |
| before.Heartbeat.p95 | 7.46 |
| before.Heartbeat.p99 | 11.01 |
| before.Lookup.p50 | 3.65 |
| before.Lookup.p95 | 6.84 |
| before.Lookup.p99 | 8.53 |
| before.Deregister.p50 | 4.24 |
| before.Deregister.p95 | 8.01 |
| before.Deregister.p99 | 11.29 |
| during.Register.p50 | 3.96 |
| during.Register.p95 | 7.70 |
| during.Register.p99 | 9.85 |
| during.Heartbeat.p50 | 3.42 |
| during.Heartbeat.p95 | 6.43 |
| during.Heartbeat.p99 | 8.12 |
| during.Lookup.p50 | 3.80 |
| during.Lookup.p95 | 6.31 |
| during.Lookup.p99 | 8.17 |
| during.Deregister.p50 | 3.40 |
| during.Deregister.p95 | 5.82 |
| during.Deregister.p99 | 7.28 |
| after.Register.p50 | 4.34 |
| after.Register.p95 | 8.95 |
| after.Register.p99 | 11.43 |
| after.Heartbeat.p50 | 3.52 |
| after.Heartbeat.p95 | 6.69 |
| after.Heartbeat.p99 | 8.49 |
| after.Lookup.p50 | 3.62 |
| after.Lookup.p95 | 6.52 |
| after.Lookup.p99 | 8.66 |
| after.Deregister.p50 | 3.94 |
| after.Deregister.p95 | 7.20 |
| after.Deregister.p99 | 9.42 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.55 |
| other_nodes.p50 | 3.79 |
| other_nodes.p99 | 11.50 |
| other_nodes.max | 14.11 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 61 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 147 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 15.76 |

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
| before.p99 | 15.57 |
| during.p99 | 13.40 |
| after.p99 | 13.73 |

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
| after | 8.00 |

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
