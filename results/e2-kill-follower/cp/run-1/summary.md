# e2-kill-follower — cp — run 1

## Faults

| t (s) | action | target |
|---|---|---|
| 30.9 | kill | registry-1 |
| 62.4 | start | registry-1 |

## Node groups

| group | nodes |
|---|---|
| killed | registry-1 |
| others | registry-2 registry-3 registry-4 registry-5 |

## availability_pct

| metric | value |
|---|---|
| before.killed | 99.76 |
| before.others | 100.00 |
| during.killed | 0.00 |
| during.others | 100.00 |
| after.killed | 99.71 |
| after.others | 100.00 |

## outage_after_fault_s

| metric | value |
|---|---|
| killed | 31.53 |
| others | 0.00 |

## outage_after_restore_s

| metric | value |
|---|---|
| killed | 0.56 |
| others | 0.00 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 42.11 |
| before.Register.p95 | 121.44 |
| before.Register.p99 | 145.93 |
| before.Heartbeat.p50 | 7.16 |
| before.Heartbeat.p95 | 14.25 |
| before.Heartbeat.p99 | 31.67 |
| before.Lookup.p50 | 7.39 |
| before.Lookup.p95 | 15.73 |
| before.Lookup.p99 | 42.82 |
| before.Deregister.p50 | 33.68 |
| before.Deregister.p95 | 102.06 |
| before.Deregister.p99 | 148.98 |
| during.Register.p50 | 34.81 |
| during.Register.p95 | 55.19 |
| during.Register.p99 | 72.31 |
| during.Heartbeat.p50 | 6.87 |
| during.Heartbeat.p95 | 13.26 |
| during.Heartbeat.p99 | 23.72 |
| during.Lookup.p50 | 6.84 |
| during.Lookup.p95 | 13.68 |
| during.Lookup.p99 | 24.13 |
| during.Deregister.p50 | 24.93 |
| during.Deregister.p95 | 55.46 |
| during.Deregister.p99 | 95.79 |
| after.Register.p50 | 38.41 |
| after.Register.p95 | 99.39 |
| after.Register.p99 | 250.89 |
| after.Heartbeat.p50 | 7.65 |
| after.Heartbeat.p95 | 13.94 |
| after.Heartbeat.p99 | 22.07 |
| after.Lookup.p50 | 7.57 |
| after.Lookup.p95 | 18.05 |
| after.Lookup.p99 | 43.11 |
| after.Deregister.p50 | 30.15 |
| after.Deregister.p95 | 84.04 |
| after.Deregister.p99 | 126.23 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 82.24 |
| other_nodes.p50 | 82.70 |
| other_nodes.p99 | 138.67 |
| other_nodes.max | 141.67 |

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
| convergence_after_restore_s | 0.76 |

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
| before.p99 | 13.11 |
| during.p99 | 12.45 |
| after.p99 | 13.75 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 0.00 |

## leader_changes

| metric | value |
|---|---|
| leader_changes | 0 |

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
