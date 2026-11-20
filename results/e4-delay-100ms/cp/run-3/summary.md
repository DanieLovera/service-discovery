# e4-delay-100ms — cp — run 3

## Faults

| t (s) | action | target |
|---|---|---|
| 30.5 | degrade | delay=100ms loss=0% |
| 91.2 | restore |  |

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
| before.Register.p50 | 39.49 |
| before.Register.p95 | 54.44 |
| before.Register.p99 | 91.36 |
| before.Heartbeat.p50 | 7.26 |
| before.Heartbeat.p95 | 14.44 |
| before.Heartbeat.p99 | 27.66 |
| before.Lookup.p50 | 7.46 |
| before.Lookup.p95 | 14.27 |
| before.Lookup.p99 | 27.09 |
| before.Deregister.p50 | 31.67 |
| before.Deregister.p95 | 51.49 |
| before.Deregister.p99 | 68.80 |
| during.Register.p50 | 562.92 |
| during.Register.p95 | 655.40 |
| during.Register.p99 | 756.87 |
| during.Heartbeat.p50 | 341.03 |
| during.Heartbeat.p95 | 422.51 |
| during.Heartbeat.p99 | 454.73 |
| during.Lookup.p50 | 333.09 |
| during.Lookup.p95 | 426.85 |
| during.Lookup.p99 | 465.49 |
| during.Deregister.p50 | 574.28 |
| during.Deregister.p95 | 671.84 |
| during.Deregister.p99 | 751.42 |
| after.Register.p50 | 38.91 |
| after.Register.p95 | 59.07 |
| after.Register.p99 | 248.89 |
| after.Heartbeat.p50 | 7.51 |
| after.Heartbeat.p95 | 18.92 |
| after.Heartbeat.p99 | 66.99 |
| after.Lookup.p50 | 7.92 |
| after.Lookup.p95 | 28.21 |
| after.Lookup.p99 | 66.45 |
| after.Deregister.p50 | 29.49 |
| after.Deregister.p95 | 53.67 |
| after.Deregister.p99 | 80.06 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 84.50 |
| other_nodes.p50 | 89.78 |
| other_nodes.p99 | 122.55 |
| other_nodes.max | 122.79 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 61 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 4 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 0.34 |

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
| before.p99 | 13.51 |
| during.p99 | 9.86 |
| after.p99 | 12.51 |

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
| -5.1 | registry-1 |
