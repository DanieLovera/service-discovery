# e3-partition — ap — run 3

## Faults

| t (s) | action | target |
|---|---|---|
| 31.1 | partition | registry-1 registry-2 | registry-3 registry-4 registry-5 |
| 90.4 | heal |  |

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
| before.Register.p50 | 4.69 |
| before.Register.p95 | 10.02 |
| before.Register.p99 | 19.31 |
| before.Heartbeat.p50 | 3.75 |
| before.Heartbeat.p95 | 7.69 |
| before.Heartbeat.p99 | 12.08 |
| before.Lookup.p50 | 3.85 |
| before.Lookup.p95 | 7.59 |
| before.Lookup.p99 | 10.42 |
| before.Deregister.p50 | 4.10 |
| before.Deregister.p95 | 7.96 |
| before.Deregister.p99 | 10.39 |
| during.Register.p50 | 4.01 |
| during.Register.p95 | 7.53 |
| during.Register.p99 | 9.53 |
| during.Heartbeat.p50 | 3.51 |
| during.Heartbeat.p95 | 6.43 |
| during.Heartbeat.p99 | 8.23 |
| during.Lookup.p50 | 3.74 |
| during.Lookup.p95 | 6.47 |
| during.Lookup.p99 | 8.12 |
| during.Deregister.p50 | 3.32 |
| during.Deregister.p95 | 5.88 |
| during.Deregister.p99 | 8.04 |
| after.Register.p50 | 4.63 |
| after.Register.p95 | 9.14 |
| after.Register.p99 | 11.59 |
| after.Heartbeat.p50 | 3.41 |
| after.Heartbeat.p95 | 6.36 |
| after.Heartbeat.p99 | 7.83 |
| after.Lookup.p50 | 3.66 |
| after.Lookup.p95 | 6.33 |
| after.Lookup.p99 | 8.23 |
| after.Deregister.p50 | 4.04 |
| after.Deregister.p95 | 7.26 |
| after.Deregister.p99 | 9.02 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.60 |
| other_nodes.p50 | 3.96 |
| other_nodes.p99 | 71520.92 |
| other_nodes.max | 71551.59 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 59 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 141 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 12.06 |

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
| before.p99 | 15.20 |
| during.p99 | 13.35 |
| after.p99 | 13.49 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 0.00 |

## max_suspect_members

| metric | value |
|---|---|
| before | 0.00 |
| during | 7.00 |
| after | 0.00 |

## max_confirm_members

| metric | value |
|---|---|
| before | 0.00 |
| during | 12.00 |
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
