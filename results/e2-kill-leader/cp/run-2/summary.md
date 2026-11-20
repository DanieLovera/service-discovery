# e2-kill-leader — cp — run 2

## Faults

| t (s) | action | target |
|---|---|---|
| 31.7 | kill | registry-3 |
| 62.4 | start | registry-3 |

## Node groups

| group | nodes |
|---|---|
| killed | registry-3 |
| others | registry-1 registry-2 registry-4 registry-5 |

## availability_pct

| metric | value |
|---|---|
| before.killed | 98.81 |
| before.others | 98.81 |
| during.killed | 0.00 |
| during.others | 99.38 |
| after.killed | 99.42 |
| after.others | 100.00 |

## outage_after_fault_s

| metric | value |
|---|---|
| killed | 30.54 |
| others | 1.14 |

## outage_after_restore_s

| metric | value |
|---|---|
| killed | 1.01 |
| others | 0.00 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 40.20 |
| before.Register.p95 | 54.29 |
| before.Register.p99 | 61.14 |
| before.Heartbeat.p50 | 7.35 |
| before.Heartbeat.p95 | 13.41 |
| before.Heartbeat.p99 | 16.49 |
| before.Lookup.p50 | 7.73 |
| before.Lookup.p95 | 15.08 |
| before.Lookup.p99 | 26.37 |
| before.Deregister.p50 | 29.64 |
| before.Deregister.p95 | 47.04 |
| before.Deregister.p99 | 62.04 |
| during.Register.p50 | 36.00 |
| during.Register.p95 | 66.69 |
| during.Register.p99 | 372.53 |
| during.Heartbeat.p50 | 7.06 |
| during.Heartbeat.p95 | 12.16 |
| during.Heartbeat.p99 | 19.53 |
| during.Lookup.p50 | 6.82 |
| during.Lookup.p95 | 15.12 |
| during.Lookup.p99 | 25.17 |
| during.Deregister.p50 | 25.68 |
| during.Deregister.p95 | 41.62 |
| during.Deregister.p99 | 63.37 |
| after.Register.p50 | 38.90 |
| after.Register.p95 | 72.29 |
| after.Register.p99 | 352.16 |
| after.Heartbeat.p50 | 7.35 |
| after.Heartbeat.p95 | 13.72 |
| after.Heartbeat.p99 | 29.38 |
| after.Lookup.p50 | 7.73 |
| after.Lookup.p95 | 18.16 |
| after.Lookup.p99 | 32.88 |
| after.Deregister.p50 | 28.93 |
| after.Deregister.p95 | 61.24 |
| after.Deregister.p99 | 113.21 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 72.44 |
| other_nodes.p50 | 80.82 |
| other_nodes.p99 | 111.02 |
| other_nodes.max | 111.79 |

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
| convergence_after_restore_s | 1.19 |

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
| before.p99 | 12.99 |
| during.p99 | 13.15 |
| after.p99 | 13.48 |

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
| -5.2 | registry-3 |
| 39.8 | registry-5 |
