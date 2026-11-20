# e4-delay-100ms — cp — run 1

## Faults

| t (s) | action | target |
|---|---|---|
| 31.0 | degrade | delay=100ms loss=0% |
| 90.5 | restore |  |

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
| before.Register.p50 | 39.09 |
| before.Register.p95 | 72.24 |
| before.Register.p99 | 106.57 |
| before.Heartbeat.p50 | 7.48 |
| before.Heartbeat.p95 | 13.30 |
| before.Heartbeat.p99 | 19.09 |
| before.Lookup.p50 | 7.62 |
| before.Lookup.p95 | 16.20 |
| before.Lookup.p99 | 22.29 |
| before.Deregister.p50 | 31.44 |
| before.Deregister.p95 | 70.48 |
| before.Deregister.p99 | 102.57 |
| during.Register.p50 | 565.17 |
| during.Register.p95 | 663.12 |
| during.Register.p99 | 842.47 |
| during.Heartbeat.p50 | 334.64 |
| during.Heartbeat.p95 | 424.99 |
| during.Heartbeat.p99 | 507.72 |
| during.Lookup.p50 | 337.78 |
| during.Lookup.p95 | 429.82 |
| during.Lookup.p99 | 546.22 |
| during.Deregister.p50 | 577.48 |
| during.Deregister.p95 | 682.05 |
| during.Deregister.p99 | 804.20 |
| after.Register.p50 | 41.38 |
| after.Register.p95 | 88.22 |
| after.Register.p99 | 219.19 |
| after.Heartbeat.p50 | 8.17 |
| after.Heartbeat.p95 | 29.47 |
| after.Heartbeat.p99 | 93.16 |
| after.Lookup.p50 | 7.89 |
| after.Lookup.p95 | 39.34 |
| after.Lookup.p99 | 112.38 |
| after.Deregister.p50 | 29.77 |
| after.Deregister.p95 | 202.19 |
| after.Deregister.p99 | 283.33 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 98.64 |
| other_nodes.p50 | 80.45 |
| other_nodes.p99 | 309.02 |
| other_nodes.max | 317.78 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 59 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 3 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 0.32 |

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
| before.p99 | 11.25 |
| during.p99 | 9.78 |
| after.p99 | 14.79 |

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
| -5.9 | registry-3 |
