# e4-delay-250ms — cp — run 1

## Faults

| t (s) | action | target |
|---|---|---|
| 30.8 | degrade | delay=250ms loss=0% |
| 90.9 | restore |  |

## Node groups

| group | nodes |
|---|---|
| all | registry-1 registry-2 registry-3 registry-4 registry-5 |

## availability_pct

| metric | value |
|---|---|
| before.all | 99.76 |
| during.all | 91.96 |
| after.all | 95.09 |

## outage_after_fault_s

| metric | value |
|---|---|
| all | 59.95 |

## outage_after_restore_s

| metric | value |
|---|---|
| all | 1.70 |

## latency_ms

| metric | value |
|---|---|
| before.Register.p50 | 40.06 |
| before.Register.p95 | 56.11 |
| before.Register.p99 | 69.39 |
| before.Heartbeat.p50 | 7.72 |
| before.Heartbeat.p95 | 14.30 |
| before.Heartbeat.p99 | 21.73 |
| before.Lookup.p50 | 7.65 |
| before.Lookup.p95 | 16.11 |
| before.Lookup.p99 | 27.36 |
| before.Deregister.p50 | 30.59 |
| before.Deregister.p95 | 47.68 |
| before.Deregister.p99 | 76.78 |
| during.Register.p50 | 1462.38 |
| during.Register.p95 | 1910.45 |
| during.Register.p99 | 1973.54 |
| during.Heartbeat.p50 | 934.18 |
| during.Heartbeat.p95 | 1141.19 |
| during.Heartbeat.p99 | 1631.29 |
| during.Lookup.p50 | 1014.49 |
| during.Lookup.p95 | 1530.83 |
| during.Lookup.p99 | 1894.09 |
| during.Deregister.p50 | 1515.19 |
| during.Deregister.p95 | 1947.73 |
| during.Deregister.p99 | 2000.32 |
| after.Register.p50 | 38.05 |
| after.Register.p95 | 74.65 |
| after.Register.p99 | 517.58 |
| after.Heartbeat.p50 | 7.38 |
| after.Heartbeat.p95 | 19.49 |
| after.Heartbeat.p99 | 45.16 |
| after.Lookup.p50 | 11.28 |
| after.Lookup.p95 | 23.73 |
| after.Lookup.p99 | 52.52 |
| after.Deregister.p50 | 28.98 |
| after.Deregister.p95 | 61.59 |
| after.Deregister.p99 | 95.94 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 87.29 |
| other_nodes.p50 | 86.99 |
| other_nodes.p99 | 127.89 |
| other_nodes.max | 130.37 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 54 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 4 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 1.94 |

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
| before.p99 | 13.90 |
| during.p99 | 11.71 |
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
| during | 7.00 |
| after | 9.00 |

## min_lb_pool_instances

| metric | value |
|---|---|
| before | 3.00 |
| during | 1.00 |
| after | 2.00 |

## Raft leader changes

| t (s) | leader |
|---|---|
| -5.8 | registry-2 |
