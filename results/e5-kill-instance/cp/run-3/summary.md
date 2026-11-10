# e5-kill-instance — cp — run 3

## Faults

| t (s) | action | target |
|---|---|---|
| 30.8 | kill | service-c-3 |
| 62.9 | start | service-c-3 |

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
| before.Register.p50 | 40.40 |
| before.Register.p95 | 55.02 |
| before.Register.p99 | 83.75 |
| before.Heartbeat.p50 | 7.19 |
| before.Heartbeat.p95 | 14.51 |
| before.Heartbeat.p99 | 18.95 |
| before.Lookup.p50 | 7.94 |
| before.Lookup.p95 | 16.41 |
| before.Lookup.p99 | 25.10 |
| before.Deregister.p50 | 30.82 |
| before.Deregister.p95 | 46.32 |
| before.Deregister.p99 | 70.81 |
| during.Register.p50 | 37.87 |
| during.Register.p95 | 76.84 |
| during.Register.p99 | 280.60 |
| during.Heartbeat.p50 | 7.05 |
| during.Heartbeat.p95 | 14.06 |
| during.Heartbeat.p99 | 23.94 |
| during.Lookup.p50 | 7.64 |
| during.Lookup.p95 | 19.14 |
| during.Lookup.p99 | 40.91 |
| during.Deregister.p50 | 30.13 |
| during.Deregister.p95 | 55.04 |
| during.Deregister.p99 | 100.98 |
| after.Register.p50 | 38.11 |
| after.Register.p95 | 72.75 |
| after.Register.p99 | 100.95 |
| after.Heartbeat.p50 | 7.15 |
| after.Heartbeat.p95 | 11.63 |
| after.Heartbeat.p99 | 14.18 |
| after.Lookup.p50 | 7.40 |
| after.Lookup.p95 | 14.76 |
| after.Lookup.p99 | 32.37 |
| after.Deregister.p50 | 29.19 |
| after.Deregister.p95 | 51.08 |
| after.Deregister.p99 | 138.88 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 87.22 |
| other_nodes.p50 | 88.90 |
| other_nodes.p99 | 130.45 |
| other_nodes.max | 130.47 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 33 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 4 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 0.03 |

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
| during | 105 |
| after | 0 |

## lb_latency_ms

| metric | value |
|---|---|
| before.p99 | 12.65 |
| during.p99 | 15.12 |
| after.p99 | 11.03 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 10.43 |

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
| -5.5 | registry-3 |
