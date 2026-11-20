# e4-loss-10pct — cp — run 1

## Faults

| t (s) | action | target |
|---|---|---|
| 30.9 | degrade | delay=0ms loss=10% |
| 90.8 | restore |  |

## Node groups

| group | nodes |
|---|---|
| all | registry-1 registry-2 registry-3 registry-4 registry-5 |

## availability_pct

| metric | value |
|---|---|
| before.all | 100.00 |
| during.all | 99.88 |
| after.all | 100.00 |

## outage_after_fault_s

| metric | value |
|---|---|
| all | 52.94 |

## outage_after_restore_s

| metric | value |
|---|---|
| all | 0.00 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 40.58 |
| before.Register.p95 | 68.03 |
| before.Register.p99 | 89.99 |
| before.Heartbeat.p50 | 7.63 |
| before.Heartbeat.p95 | 14.50 |
| before.Heartbeat.p99 | 22.73 |
| before.Lookup.p50 | 7.63 |
| before.Lookup.p95 | 16.65 |
| before.Lookup.p99 | 30.10 |
| before.Deregister.p50 | 31.78 |
| before.Deregister.p95 | 62.10 |
| before.Deregister.p99 | 130.88 |
| during.Register.p50 | 65.42 |
| during.Register.p95 | 439.17 |
| during.Register.p99 | 717.85 |
| during.Heartbeat.p50 | 21.77 |
| during.Heartbeat.p95 | 238.81 |
| during.Heartbeat.p99 | 506.40 |
| during.Lookup.p50 | 23.69 |
| during.Lookup.p95 | 239.78 |
| during.Lookup.p99 | 518.34 |
| during.Deregister.p50 | 41.47 |
| during.Deregister.p95 | 294.48 |
| during.Deregister.p99 | 645.03 |
| after.Register.p50 | 38.51 |
| after.Register.p95 | 78.11 |
| after.Register.p99 | 188.01 |
| after.Heartbeat.p50 | 7.35 |
| after.Heartbeat.p95 | 13.04 |
| after.Heartbeat.p99 | 44.24 |
| after.Lookup.p50 | 7.72 |
| after.Lookup.p95 | 18.36 |
| after.Lookup.p99 | 42.28 |
| after.Deregister.p50 | 30.44 |
| after.Deregister.p95 | 59.50 |
| after.Deregister.p99 | 92.74 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 92.42 |
| other_nodes.p50 | 87.76 |
| other_nodes.p99 | 172.45 |
| other_nodes.max | 175.45 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 60 |

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
| during | 0 |
| after | 0 |

## lb_latency_ms

| metric | value |
|---|---|
| before.p99 | 14.11 |
| during.p99 | 11.02 |
| after.p99 | 11.42 |

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
| after | 8.00 |

## min_lb_pool_instances

| metric | value |
|---|---|
| before | 3.00 |
| during | 3.00 |
| after | 3.00 |

## Raft leader changes

| t (s) | leader |
|---|---|
| -5.4 | registry-3 |
