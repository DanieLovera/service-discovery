# e4-delay-250ms — cp — run 2

## Faults

| t (s) | action | target |
|---|---|---|
| 30.9 | degrade | delay=250ms loss=0% |
| 90.8 | restore |  |

## Node groups

| group | nodes |
|---|---|
| all | registry-1 registry-2 registry-3 registry-4 registry-5 |

## availability_pct

| metric | value |
|---|---|
| before.all | 100.00 |
| during.all | 98.46 |
| after.all | 95.75 |

## outage_after_fault_s

| metric | value |
|---|---|
| all | 59.91 |

## outage_after_restore_s

| metric | value |
|---|---|
| all | 1.96 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 40.59 |
| before.Register.p95 | 92.25 |
| before.Register.p99 | 146.22 |
| before.Heartbeat.p50 | 7.79 |
| before.Heartbeat.p95 | 17.81 |
| before.Heartbeat.p99 | 42.46 |
| before.Lookup.p50 | 7.59 |
| before.Lookup.p95 | 14.05 |
| before.Lookup.p99 | 25.84 |
| before.Deregister.p50 | 31.99 |
| before.Deregister.p95 | 84.48 |
| before.Deregister.p99 | 140.73 |
| during.Register.p50 | 1384.47 |
| during.Register.p95 | 1654.74 |
| during.Register.p99 | 1738.06 |
| during.Heartbeat.p50 | 846.53 |
| during.Heartbeat.p95 | 1032.36 |
| during.Heartbeat.p99 | 1056.71 |
| during.Lookup.p50 | 830.91 |
| during.Lookup.p95 | 1038.13 |
| during.Lookup.p99 | 1063.32 |
| during.Deregister.p50 | 1398.54 |
| during.Deregister.p95 | 1642.93 |
| during.Deregister.p99 | 1695.82 |
| after.Register.p50 | 37.05 |
| after.Register.p95 | 95.78 |
| after.Register.p99 | 868.64 |
| after.Heartbeat.p50 | 7.76 |
| after.Heartbeat.p95 | 13.80 |
| after.Heartbeat.p99 | 24.64 |
| after.Lookup.p50 | 11.60 |
| after.Lookup.p95 | 22.03 |
| after.Lookup.p99 | 71.89 |
| after.Deregister.p50 | 28.57 |
| after.Deregister.p95 | 60.57 |
| after.Deregister.p99 | 126.93 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 90.39 |
| other_nodes.p50 | 90.10 |
| other_nodes.p99 | 1350.21 |
| other_nodes.max | 1567.98 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 58 |

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
| before.p99 | 12.66 |
| during.p99 | 9.78 |
| after.p99 | 13.61 |

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
| during | 8.00 |
| after | 8.00 |

## min_lb_pool_instances

| metric | value |
|---|---|
| before | 3.00 |
| during | 2.00 |
| after | 3.00 |

## Raft leader changes

| t (s) | leader |
|---|---|
| -5.8 | registry-5 |
