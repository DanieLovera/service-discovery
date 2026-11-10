# e5-kill-instance — cp — run 2

## Faults

| t (s) | action | target |
|---|---|---|
| 30.9 | kill | service-c-3 |
| 62.5 | start | service-c-3 |

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
| before.Register.p50 | 41.85 |
| before.Register.p95 | 62.49 |
| before.Register.p99 | 85.38 |
| before.Heartbeat.p50 | 7.47 |
| before.Heartbeat.p95 | 16.94 |
| before.Heartbeat.p99 | 43.82 |
| before.Lookup.p50 | 7.58 |
| before.Lookup.p95 | 16.91 |
| before.Lookup.p99 | 28.23 |
| before.Deregister.p50 | 30.28 |
| before.Deregister.p95 | 56.96 |
| before.Deregister.p99 | 87.02 |
| during.Register.p50 | 38.43 |
| during.Register.p95 | 53.61 |
| during.Register.p99 | 71.40 |
| during.Heartbeat.p50 | 7.12 |
| during.Heartbeat.p95 | 13.16 |
| during.Heartbeat.p99 | 38.09 |
| during.Lookup.p50 | 7.77 |
| during.Lookup.p95 | 15.83 |
| during.Lookup.p99 | 26.43 |
| during.Deregister.p50 | 29.24 |
| during.Deregister.p95 | 48.91 |
| during.Deregister.p99 | 66.13 |
| after.Register.p50 | 39.10 |
| after.Register.p95 | 64.01 |
| after.Register.p99 | 106.54 |
| after.Heartbeat.p50 | 7.11 |
| after.Heartbeat.p95 | 13.25 |
| after.Heartbeat.p99 | 19.41 |
| after.Lookup.p50 | 7.68 |
| after.Lookup.p95 | 12.57 |
| after.Lookup.p99 | 23.20 |
| after.Deregister.p50 | 28.98 |
| after.Deregister.p95 | 64.46 |
| after.Deregister.p99 | 94.83 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 93.82 |
| other_nodes.p50 | 83.31 |
| other_nodes.p99 | 130.97 |
| other_nodes.max | 131.79 |

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
| during | 105 |
| after | 0 |

## lb_latency_ms

| metric | value |
|---|---|
| before.p99 | 13.31 |
| during.p99 | 16.69 |
| after.p99 | 12.49 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 10.42 |

## leader_changes

| metric | value |
|---|---|
| leader_changes | 0 |

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

## Raft leader changes

| t (s) | leader |
|---|---|
| -5.2 | registry-3 |
