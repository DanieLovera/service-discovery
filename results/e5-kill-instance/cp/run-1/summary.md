# e5-kill-instance — cp — run 1

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
| before.Register.p50 | 40.63 |
| before.Register.p95 | 68.12 |
| before.Register.p99 | 115.04 |
| before.Heartbeat.p50 | 7.38 |
| before.Heartbeat.p95 | 14.63 |
| before.Heartbeat.p99 | 18.82 |
| before.Lookup.p50 | 7.66 |
| before.Lookup.p95 | 13.43 |
| before.Lookup.p99 | 25.15 |
| before.Deregister.p50 | 30.75 |
| before.Deregister.p95 | 54.97 |
| before.Deregister.p99 | 75.90 |
| during.Register.p50 | 38.49 |
| during.Register.p95 | 77.37 |
| during.Register.p99 | 104.40 |
| during.Heartbeat.p50 | 7.19 |
| during.Heartbeat.p95 | 12.75 |
| during.Heartbeat.p99 | 17.72 |
| during.Lookup.p50 | 7.55 |
| during.Lookup.p95 | 17.32 |
| during.Lookup.p99 | 33.71 |
| during.Deregister.p50 | 29.31 |
| during.Deregister.p95 | 68.12 |
| during.Deregister.p99 | 104.67 |
| after.Register.p50 | 39.47 |
| after.Register.p95 | 66.32 |
| after.Register.p99 | 111.57 |
| after.Heartbeat.p50 | 7.26 |
| after.Heartbeat.p95 | 13.22 |
| after.Heartbeat.p99 | 18.37 |
| after.Lookup.p50 | 7.34 |
| after.Lookup.p95 | 14.12 |
| after.Lookup.p99 | 23.11 |
| after.Deregister.p50 | 28.16 |
| after.Deregister.p95 | 65.17 |
| after.Deregister.p99 | 99.07 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 77.52 |
| other_nodes.p50 | 86.59 |
| other_nodes.p99 | 134.60 |
| other_nodes.max | 136.11 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 32 |

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
| before | 1 |
| during | 105 |
| after | 0 |

## lb_latency_ms

| metric | value |
|---|---|
| before.p99 | 12.32 |
| during.p99 | 15.28 |
| after.p99 | 11.52 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 10.50 |

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
| -6.0 | registry-2 |
