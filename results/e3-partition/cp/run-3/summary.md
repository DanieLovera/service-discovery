# e3-partition — cp — run 3

## Faults

| t (s) | action | target |
|---|---|---|
| 31.3 | partition | registry-5 registry-1 | registry-2 registry-3 registry-4 |
| 91.0 | heal |  |

## Node groups

| group | nodes |
|---|---|
| minority | registry-5 registry-1 |
| majority | registry-2 registry-3 registry-4 |

## availability_pct

| metric | value |
|---|---|
| before.minority | 99.53 |
| before.majority | 99.53 |
| during.minority | 0.00 |
| during.majority | 99.52 |
| after.minority | 100.00 |
| after.majority | 100.00 |

## outage_after_fault_s

| metric | value |
|---|---|
| minority | 59.75 |
| majority | 1.35 |

## outage_after_restore_s

| metric | value |
|---|---|
| minority | 0.00 |
| majority | 0.00 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 42.64 |
| before.Register.p95 | 68.27 |
| before.Register.p99 | 89.19 |
| before.Heartbeat.p50 | 7.40 |
| before.Heartbeat.p95 | 15.49 |
| before.Heartbeat.p99 | 22.46 |
| before.Lookup.p50 | 7.74 |
| before.Lookup.p95 | 19.47 |
| before.Lookup.p99 | 29.46 |
| before.Deregister.p50 | 31.05 |
| before.Deregister.p95 | 57.32 |
| before.Deregister.p99 | 105.20 |
| during.Register.p50 | 34.74 |
| during.Register.p95 | 59.18 |
| during.Register.p99 | 123.81 |
| during.Heartbeat.p50 | 5.67 |
| during.Heartbeat.p95 | 10.87 |
| during.Heartbeat.p99 | 12.96 |
| during.Lookup.p50 | 5.59 |
| during.Lookup.p95 | 8.28 |
| during.Lookup.p99 | 14.30 |
| during.Deregister.p50 | 20.75 |
| during.Deregister.p95 | 34.39 |
| during.Deregister.p99 | 57.84 |
| after.Register.p50 | 38.93 |
| after.Register.p95 | 73.92 |
| after.Register.p99 | 620.56 |
| after.Heartbeat.p50 | 7.02 |
| after.Heartbeat.p95 | 13.65 |
| after.Heartbeat.p99 | 21.78 |
| after.Lookup.p50 | 7.21 |
| after.Lookup.p95 | 14.64 |
| after.Lookup.p99 | 21.78 |
| after.Deregister.p50 | 28.00 |
| after.Deregister.p95 | 49.91 |
| after.Deregister.p99 | 99.94 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 91.22 |
| other_nodes.p50 | 90.31 |
| other_nodes.p99 | 1917.82 |
| other_nodes.max | 1938.90 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 35 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 70 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 2.75 |

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
| before.p99 | 13.95 |
| during.p99 | 9.52 |
| after.p99 | 11.64 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 0.00 |

## leader_changes

| metric | value |
|---|---|
| leader_changes | 1 |

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

## Raft leader changes

| t (s) | leader |
|---|---|
| -5.9 | registry-5 |
| 39.1 | registry-2 |
