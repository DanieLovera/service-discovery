# e3-partition — cp — run 2

## Faults

| t (s) | action | target |
|---|---|---|
| 31.2 | partition | registry-3 registry-1 | registry-2 registry-4 registry-5 |
| 90.9 | heal |  |

## Node groups

| group | nodes |
|---|---|
| minority | registry-3 registry-1 |
| majority | registry-2 registry-4 registry-5 |

## availability_pct

| metric | value |
|---|---|
| before.minority | 99.41 |
| before.majority | 99.29 |
| during.minority | 2.35 |
| during.majority | 99.26 |
| after.minority | 100.00 |
| after.majority | 100.00 |

## outage_after_fault_s

| metric | value |
|---|---|
| minority | 58.19 |
| majority | 2.39 |

## outage_after_restore_s

| metric | value |
|---|---|
| minority | 0.00 |
| majority | 0.00 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 39.80 |
| before.Register.p95 | 55.77 |
| before.Register.p99 | 62.87 |
| before.Heartbeat.p50 | 7.97 |
| before.Heartbeat.p95 | 16.27 |
| before.Heartbeat.p99 | 22.41 |
| before.Lookup.p50 | 7.58 |
| before.Lookup.p95 | 16.19 |
| before.Lookup.p99 | 22.89 |
| before.Deregister.p50 | 28.48 |
| before.Deregister.p95 | 43.22 |
| before.Deregister.p99 | 54.25 |
| during.Register.p50 | 33.73 |
| during.Register.p95 | 62.11 |
| during.Register.p99 | 1234.91 |
| during.Heartbeat.p50 | 5.81 |
| during.Heartbeat.p95 | 11.85 |
| during.Heartbeat.p99 | 15.73 |
| during.Lookup.p50 | 5.63 |
| during.Lookup.p95 | 10.76 |
| during.Lookup.p99 | 19.30 |
| during.Deregister.p50 | 21.78 |
| during.Deregister.p95 | 34.76 |
| during.Deregister.p99 | 51.64 |
| after.Register.p50 | 39.69 |
| after.Register.p95 | 79.06 |
| after.Register.p99 | 210.71 |
| after.Heartbeat.p50 | 7.34 |
| after.Heartbeat.p95 | 16.20 |
| after.Heartbeat.p99 | 27.03 |
| after.Lookup.p50 | 7.46 |
| after.Lookup.p95 | 18.42 |
| after.Lookup.p99 | 38.79 |
| after.Deregister.p50 | 28.82 |
| after.Deregister.p95 | 58.30 |
| after.Deregister.p99 | 86.00 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 50.27 |
| other_nodes.p50 | 48.48 |
| other_nodes.p99 | 85.08 |
| other_nodes.max | 85.42 |

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
| convergence_after_restore_s | 1.17 |

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
| before.p99 | 13.21 |
| during.p99 | 10.50 |
| after.p99 | 12.43 |

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
| -5.7 | registry-3 |
| 39.3 | registry-5 |
