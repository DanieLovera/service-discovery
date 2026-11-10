# e2-kill-follower — cp — run 3

## Faults

| t (s) | action | target |
|---|---|---|
| 30.8 | kill | registry-1 |
| 62.3 | start | registry-1 |

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
| after.killed | 99.42 |
| after.others | 100.00 |

## outage_after_fault_s

| metric | value |
|---|---|
| killed | 31.36 |
| others | 0.00 |

## outage_after_restore_s

| metric | value |
|---|---|
| killed | 1.13 |
| others | 0.00 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 40.44 |
| before.Register.p95 | 61.97 |
| before.Register.p99 | 83.67 |
| before.Heartbeat.p50 | 7.21 |
| before.Heartbeat.p95 | 14.73 |
| before.Heartbeat.p99 | 20.03 |
| before.Lookup.p50 | 7.66 |
| before.Lookup.p95 | 16.34 |
| before.Lookup.p99 | 25.50 |
| before.Deregister.p50 | 32.11 |
| before.Deregister.p95 | 55.71 |
| before.Deregister.p99 | 67.15 |
| during.Register.p50 | 33.82 |
| during.Register.p95 | 50.22 |
| during.Register.p99 | 70.66 |
| during.Heartbeat.p50 | 6.78 |
| during.Heartbeat.p95 | 11.61 |
| during.Heartbeat.p99 | 22.58 |
| during.Lookup.p50 | 6.57 |
| during.Lookup.p95 | 11.24 |
| during.Lookup.p99 | 18.60 |
| during.Deregister.p50 | 24.09 |
| during.Deregister.p95 | 37.22 |
| during.Deregister.p99 | 53.52 |
| after.Register.p50 | 38.07 |
| after.Register.p95 | 55.35 |
| after.Register.p99 | 137.29 |
| after.Heartbeat.p50 | 7.09 |
| after.Heartbeat.p95 | 12.82 |
| after.Heartbeat.p99 | 19.47 |
| after.Lookup.p50 | 7.49 |
| after.Lookup.p95 | 16.25 |
| after.Lookup.p99 | 25.84 |
| after.Deregister.p50 | 28.13 |
| after.Deregister.p95 | 44.02 |
| after.Deregister.p99 | 65.66 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 42.89 |
| other_nodes.p50 | 44.28 |
| other_nodes.p99 | 81.55 |
| other_nodes.max | 82.82 |

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
| convergence_after_restore_s | 1.30 |

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
| before.p99 | 14.82 |
| during.p99 | 12.84 |
| after.p99 | 12.43 |

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
| -5.4 | registry-5 |
