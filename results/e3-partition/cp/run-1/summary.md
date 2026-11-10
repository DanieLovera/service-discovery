# e3-partition — cp — run 1

## Faults

| t (s) | action | target |
|---|---|---|
| 31.2 | partition | registry-2 registry-1 | registry-3 registry-4 registry-5 |
| 91.7 | heal |  |

## Node groups

| group | nodes |
|---|---|
| minority | registry-2 registry-1 |
| majority | registry-3 registry-4 registry-5 |

## availability_pct

| metric | value |
|---|---|
| before.minority | 99.53 |
| before.majority | 99.53 |
| during.minority | 15.25 |
| during.majority | 99.21 |
| after.minority | 100.00 |
| after.majority | 100.00 |

## outage_after_fault_s

| metric | value |
|---|---|
| minority | 57.75 |
| majority | 2.35 |

## outage_after_restore_s

| metric | value |
|---|---|
| minority | 0.00 |
| majority | 0.00 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 41.47 |
| before.Register.p95 | 81.39 |
| before.Register.p99 | 133.15 |
| before.Heartbeat.p50 | 7.51 |
| before.Heartbeat.p95 | 15.90 |
| before.Heartbeat.p99 | 24.84 |
| before.Lookup.p50 | 7.76 |
| before.Lookup.p95 | 17.96 |
| before.Lookup.p99 | 41.57 |
| before.Deregister.p50 | 31.44 |
| before.Deregister.p95 | 85.36 |
| before.Deregister.p99 | 101.50 |
| during.Register.p50 | 32.57 |
| during.Register.p95 | 67.72 |
| during.Register.p99 | 1155.34 |
| during.Heartbeat.p50 | 5.91 |
| during.Heartbeat.p95 | 11.04 |
| during.Heartbeat.p99 | 24.69 |
| during.Lookup.p50 | 5.69 |
| during.Lookup.p95 | 14.46 |
| during.Lookup.p99 | 31.03 |
| during.Deregister.p50 | 22.43 |
| during.Deregister.p95 | 59.35 |
| during.Deregister.p99 | 84.62 |
| after.Register.p50 | 37.98 |
| after.Register.p95 | 98.65 |
| after.Register.p99 | 342.29 |
| after.Heartbeat.p50 | 7.15 |
| after.Heartbeat.p95 | 11.95 |
| after.Heartbeat.p99 | 27.83 |
| after.Lookup.p50 | 7.25 |
| after.Lookup.p95 | 15.12 |
| after.Lookup.p99 | 47.36 |
| after.Deregister.p50 | 28.59 |
| after.Deregister.p95 | 83.26 |
| after.Deregister.p99 | 108.43 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 93.21 |
| other_nodes.p50 | 89.77 |
| other_nodes.p99 | 2857.36 |
| other_nodes.max | 2857.40 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 37 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 12 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 0.14 |

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
| before.p99 | 13.64 |
| during.p99 | 9.90 |
| after.p99 | 11.89 |

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
| -6.0 | registry-2 |
| 39.0 | registry-5 |
