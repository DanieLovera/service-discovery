# e4-delay-100ms — ap — run 1

## Faults

| t (s) | action | target |
|---|---|---|
| 31.1 | degrade | delay=100ms loss=0% |
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
| before.Register.p50 | 4.70 |
| before.Register.p95 | 10.78 |
| before.Register.p99 | 14.34 |
| before.Heartbeat.p50 | 3.69 |
| before.Heartbeat.p95 | 7.19 |
| before.Heartbeat.p99 | 10.62 |
| before.Lookup.p50 | 3.77 |
| before.Lookup.p95 | 7.04 |
| before.Lookup.p99 | 9.83 |
| before.Deregister.p50 | 4.19 |
| before.Deregister.p95 | 8.02 |
| before.Deregister.p99 | 11.08 |
| during.Register.p50 | 103.89 |
| during.Register.p95 | 106.57 |
| during.Register.p99 | 108.63 |
| during.Heartbeat.p50 | 102.80 |
| during.Heartbeat.p95 | 105.18 |
| during.Heartbeat.p99 | 106.70 |
| during.Lookup.p50 | 102.78 |
| during.Lookup.p95 | 104.97 |
| during.Lookup.p99 | 106.95 |
| during.Deregister.p50 | 102.70 |
| during.Deregister.p95 | 104.69 |
| during.Deregister.p99 | 106.44 |
| after.Register.p50 | 4.72 |
| after.Register.p95 | 9.97 |
| after.Register.p99 | 11.98 |
| after.Heartbeat.p50 | 3.73 |
| after.Heartbeat.p95 | 7.33 |
| after.Heartbeat.p99 | 9.41 |
| after.Lookup.p50 | 3.65 |
| after.Lookup.p95 | 7.26 |
| after.Lookup.p99 | 8.89 |
| after.Deregister.p50 | 4.03 |
| after.Deregister.p95 | 7.37 |
| after.Deregister.p99 | 8.28 |

## visibility_ms

| metric | value |
|---|---|
| write_node.p50 | 2.59 |
| other_nodes.p50 | 4.08 |
| other_nodes.p99 | 5.31 |
| other_nodes.max | 6.23 |

## probes_written_during_fault

| metric | value |
|---|---|
| probes_written_during_fault | 60 |

## pairs_pending_at_restore

| metric | value |
|---|---|
| pairs_pending_at_restore | 2 |

## convergence_after_restore_s

| metric | value |
|---|---|
| convergence_after_restore_s | 0.06 |

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
| before.p99 | 16.27 |
| during.p99 | 11.97 |
| after.p99 | 15.16 |

## lb_last_error_after_fault_s

| metric | value |
|---|---|
| lb_last_error_after_fault_s | 0.00 |

## max_suspect_members

| metric | value |
|---|---|
| before | 0.00 |
| during | 0.00 |
| after | 0.00 |

## max_confirm_members

| metric | value |
|---|---|
| before | 0.00 |
| during | 0.00 |
| after | 0.00 |

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
