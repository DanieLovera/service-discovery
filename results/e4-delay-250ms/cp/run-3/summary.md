# e4-delay-250ms — cp — run 3

## Faults

| t (s) | action | target |
|---|---|---|
| 31.5 | degrade | delay=250ms loss=0% |
| 91.4 | restore |  |

## Node groups

| group | nodes |
|---|---|
| all | registry-1 registry-2 registry-3 registry-4 registry-5 |

## availability_pct

| metric | value |
|---|---|
| before.all | 99.77 |
| during.all | 88.49 |
| after.all | 95.66 |

## outage_after_fault_s

| metric | value |
|---|---|
| all | 59.93 |

## outage_after_restore_s

| metric | value |
|---|---|
| all | 1.97 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 40.74 |
| before.Register.p95 | 68.99 |
| before.Register.p99 | 141.22 |
| before.Heartbeat.p50 | 7.54 |
| before.Heartbeat.p95 | 13.13 |
| before.Heartbeat.p99 | 45.73 |
| before.Lookup.p50 | 7.74 |
| before.Lookup.p95 | 16.66 |
| before.Lookup.p99 | 23.06 |
| before.Deregister.p50 | 30.90 |
| before.Deregister.p95 | 64.36 |
| before.Deregister.p99 | 117.40 |
| during.Register.p50 | 1441.04 |
| during.Register.p95 | 1927.35 |
| during.Register.p99 | 1979.65 |
| during.Heartbeat.p50 | 944.30 |
| during.Heartbeat.p95 | 1420.72 |
| during.Heartbeat.p99 | 1512.14 |
| during.Lookup.p50 | 1005.27 |
| during.Lookup.p95 | 1562.37 |
| during.Lookup.p99 | 1919.11 |
| during.Deregister.p50 | 1499.01 |
| during.Deregister.p95 | 1994.13 |
| during.Deregister.p99 | 2000.06 |
| after.Register.p50 | 40.03 |
| after.Register.p95 | 100.23 |
| after.Register.p99 | 698.51 |
| after.Heartbeat.p50 | 7.42 |
| after.Heartbeat.p95 | 15.40 |
| after.Heartbeat.p99 | 30.56 |
| after.Lookup.p50 | 10.54 |
| after.Lookup.p95 | 19.32 |
| after.Lookup.p99 | 27.57 |
| after.Deregister.p50 | 29.75 |
| after.Deregister.p95 | 73.47 |
| after.Deregister.p99 | 127.57 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 93.71 |
| other_nodes.p50 | 90.46 |
| other_nodes.p99 | 143.98 |
| other_nodes.max | 158.22 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 47 |

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
| before.p99 | 13.57 |
| during.p99 | 10.06 |
| after.p99 | 13.11 |

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
| -5.2 | registry-1 |
