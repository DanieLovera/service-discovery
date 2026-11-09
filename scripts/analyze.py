import json
import sys
from pathlib import Path

import matplotlib

matplotlib.use("Agg")

import matplotlib.pyplot as plt  # noqa: E402
import pandas as pd  # noqa: E402

WARMUP_S = 10
END_MARGIN_S = 5
TRANSIENT_CODES = {"Unavailable", "DeadlineExceeded"}
METHODS = ["Register", "Heartbeat", "Lookup", "Deregister"]
COLUMNS = {
    "operations.csv": ["timestamp_ms", "node", "method", "code", "latency_ms"],
    "probes.csv": ["timestamp_ms", "probe_id", "write_node", "code"],
    "sightings.csv": ["timestamp_ms", "probe_id", "node", "source"],
    "load_balancer.csv": ["timestamp_ms", "service", "code", "instance_id", "latency_ms"],
    "faults.csv": ["timestamp_ms", "action", "target"],
    "metrics.csv": ["timestamp_ms", "metric", "instance", "label", "value"],
}


def read_csv(run_dir, name):
    path = run_dir / name
    if not path.exists() or path.stat().st_size == 0:
        return pd.DataFrame(columns=COLUMNS[name])

    return pd.read_csv(path, keep_default_na=False, na_values=[""])


def load_run(run_dir):
    frames = {name: read_csv(run_dir, name) for name in COLUMNS}
    t0 = frames["operations.csv"]["timestamp_ms"].min()

    for frame in frames.values():
        frame["t"] = (frame["timestamp_ms"].astype(float) - t0) / 1000

    experiment, backend, repetition = identify(run_dir)

    return {
        "dir": run_dir,
        "experiment": experiment,
        "backend": backend,
        "repetition": repetition,
        "ops": frames["operations.csv"],
        "probes": frames["probes.csv"],
        "sightings": frames["sightings.csv"],
        "lb": frames["load_balancer.csv"],
        "faults": frames["faults.csv"],
        "metrics": frames["metrics.csv"],
        "end": frames["operations.csv"]["t"].max(),
    }


def identify(run_dir):
    parts = run_dir.resolve().parts
    if len(parts) >= 3 and parts[-2] in ("cp", "ap") and parts[-1].startswith("run-"):
        return parts[-3], parts[-2], parts[-1].removeprefix("run-")

    return run_dir.name, "", ""


def phase_windows(run):
    faults = run["faults"]
    if len(faults) >= 2:
        start, stop = faults["t"].iloc[0], faults["t"].iloc[1]
        return {"before": (WARMUP_S, start), "during": (start, stop), "after": (stop, run["end"])}

    return {"steady": (WARMUP_S, run["end"])}


def fault_window(run):
    faults = run["faults"]
    if len(faults) >= 2:
        return faults["t"].iloc[0], faults["t"].iloc[1]

    return None


def node_groups(run):
    nodes = sorted(run["ops"]["node"].unique())
    faults = run["faults"]
    if faults.empty:
        return {"all": nodes}

    action, target = faults["action"].iloc[0], str(faults["target"].iloc[0])

    if action == "partition":
        side_a, side_b = (side.split() for side in target.split("|"))
        return {"minority": side_a, "majority": side_b}

    if action == "kill" and target in nodes:
        return {"killed": [target], "others": [node for node in nodes if node != target]}

    return {"all": nodes}


def between(frame, window, column="t"):
    start, stop = window
    return frame[(frame[column] >= start) & (frame[column] < stop)]


def availability(run, windows, groups):
    ops = run["ops"]
    result = {}

    for phase, window in windows.items():
        in_phase = between(ops, window)

        for group, nodes in groups.items():
            subset = in_phase[in_phase["node"].isin(nodes)]
            if len(subset):
                result[f"availability_pct.{phase}.{group}"] = 100 * (~subset["code"].isin(TRANSIENT_CODES)).mean()

    return result


def outages(run, groups):
    window = fault_window(run)
    if window is None:
        return {}

    start, stop = window
    failed = run["ops"][run["ops"]["code"].isin(TRANSIENT_CODES)]
    result = {}

    for group, nodes in groups.items():
        group_failed = failed[failed["node"].isin(nodes)]
        during = between(group_failed, (start, stop))
        after = group_failed[group_failed["t"] >= stop]

        result[f"outage_after_fault_s.{group}"] = during["t"].max() - start if len(during) else 0.0
        result[f"outage_after_restore_s.{group}"] = after["t"].max() - stop if len(after) else 0.0

    return result


def latency(run, windows):
    ok = run["ops"][run["ops"]["code"] == "OK"]
    result = {}

    for phase, window in windows.items():
        in_phase = between(ok, window)

        for method in METHODS:
            values = in_phase.loc[in_phase["method"] == method, "latency_ms"]
            if len(values):
                for quantile in (50, 95, 99):
                    result[f"latency_ms.{phase}.{method}.p{quantile}"] = values.quantile(quantile / 100)

    return result


def visibility_pairs(run):
    probes = run["probes"][run["probes"]["code"] == "OK"][["probe_id", "write_node", "t"]]
    sightings = run["sightings"][["probe_id", "node", "source", "t"]]

    pairs = sightings.merge(probes, on="probe_id", suffixes=("_seen", "_sent"))
    pairs["delay_ms"] = (pairs["t_seen"] - pairs["t_sent"]) * 1000

    return pairs


def visibility(run, pairs, windows):
    baseline = windows.get("before") or windows.get("steady")
    sent = between(pairs, baseline, "t_sent")
    own = sent.loc[sent["node"] == sent["write_node"], "delay_ms"]
    others = sent.loc[sent["node"] != sent["write_node"], "delay_ms"]

    result = {}
    if len(own):
        result["visibility_ms.write_node.p50"] = own.median()
    if len(others):
        result["visibility_ms.other_nodes.p50"] = others.median()
        result["visibility_ms.other_nodes.p99"] = others.quantile(0.99)
        result["visibility_ms.other_nodes.max"] = others.max()

    return result


def convergence(run, pairs):
    window = fault_window(run)
    if window is None:
        return {}

    start, stop = window
    nodes = sorted(run["ops"]["node"].unique())
    probes = between(run["probes"][run["probes"]["code"] == "OK"], (start, stop))
    seen = {(row.probe_id, row.node): row.t_seen for row in pairs.itertuples()}

    pending_after_restore, never_seen = [], 0
    for probe in probes.itertuples():
        for node in nodes:
            if node == probe.write_node:
                continue

            t_seen = seen.get((probe.probe_id, node))
            if t_seen is None:
                never_seen += 1
            elif t_seen >= stop:
                pending_after_restore.append(t_seen - stop)

    return {
        "probes_written_during_fault": len(probes),
        "pairs_pending_at_restore": len(pending_after_restore) + never_seen,
        "convergence_after_restore_s": max(pending_after_restore, default=0.0),
        "pairs_never_seen": never_seen,
    }


def unreplicated_writes(run, pairs):
    probes = run["probes"]
    acknowledged = probes[(probes["code"] == "OK") & (probes["t"] < run["end"] - END_MARGIN_S)]
    replicated = set(pairs.loc[pairs["node"] != pairs["write_node"], "probe_id"])

    return {"acknowledged_writes_never_replicated": int((~acknowledged["probe_id"].isin(replicated)).sum())}


def load_balancer(run, windows):
    lb = run["lb"]
    errors = lb[lb["code"] != 200]
    result = {}

    for phase, window in windows.items():
        in_phase = between(lb, window)
        if len(in_phase):
            result[f"lb_errors.{phase}"] = int((in_phase["code"] != 200).sum())
            result[f"lb_latency_ms.{phase}.p99"] = in_phase["latency_ms"].quantile(0.99)

    window = fault_window(run)
    if window is not None:
        after_fault = errors[errors["t"] >= window[0]]
        result["lb_last_error_after_fault_s"] = after_fault["t"].max() - window[0] if len(after_fault) else 0.0

    return result


def leader_changes(run):
    raft = run["metrics"][(run["metrics"]["metric"] == "registry_raft_state") & (run["metrics"]["label"] == "LEADER")]
    leaders = raft[raft["value"].astype(float) == 1].groupby("t")["instance"].first()

    changes, previous = [], None
    for t, leader in leaders.items():
        if leader != previous:
            changes.append((t, leader))
            previous = leader

    return changes


def cluster_state(run, windows):
    metrics = run["metrics"]
    if metrics.empty:
        return {}

    metrics = metrics.assign(value=metrics["value"].astype(float))
    result = {}

    changes = leader_changes(run)
    if changes:
        result["leader_changes"] = len(changes) - 1

    for phase, window in windows.items():
        in_phase = between(metrics, window)

        for status in ("SUSPECT", "CONFIRM"):
            swim = in_phase[(in_phase["metric"] == "registry_swim_members") & (in_phase["label"] == status)]
            if len(swim):
                result[f"max_{status.lower()}_members.{phase}"] = swim.groupby("t")["value"].sum().max()

        healthy = in_phase[(in_phase["metric"] == "registry_instances") & (in_phase["label"] == "HEALTHY")]
        if len(healthy):
            result[f"min_healthy_instances.{phase}"] = healthy["value"].min()

        pool = in_phase[in_phase["metric"] == "lb_pool_instances"]
        if len(pool):
            result[f"min_lb_pool_instances.{phase}"] = pool["value"].min()

    return result


def mark_faults(ax, run):
    faults = run["faults"]

    for row in faults.itertuples():
        ax.axvline(row.t, color="black", linestyle="--", linewidth=0.8)
        ax.text(row.t, 1.01, row.action, transform=ax.get_xaxis_transform(), ha="center", fontsize=8)

    window = fault_window(run)
    if window is not None:
        ax.axvspan(*window, color="grey", alpha=0.12)


def save(fig, path):
    fig.savefig(path, dpi=120, bbox_inches="tight")
    plt.close(fig)


def title(run, text):
    return f"{text} — {run['experiment']} {run['backend']} run {run['repetition']}".strip()


def plot_availability(run):
    ops = run["ops"].assign(second=run["ops"]["t"].astype(int))
    ops["answered"] = ~ops["code"].isin(TRANSIENT_CODES)
    per_second = ops.groupby(["second", "node"])["answered"].mean().unstack("node") * 100

    fig, ax = plt.subplots(figsize=(11, 4))
    per_second.plot(ax=ax, linewidth=1.2)
    mark_faults(ax, run)
    ax.set(xlabel="Time since load start (s)", ylabel="Answered operations (%)", ylim=(-5, 105))
    ax.set_title(title(run, "Availability by node"), pad=16)
    ax.legend(loc="lower left", fontsize=8)
    save(fig, run["dir"] / "availability.png")


def plot_latency(run):
    ok = run["ops"][run["ops"]["code"] == "OK"]
    ok = ok.assign(second=ok["t"].astype(int))

    fig, axes = plt.subplots(2, 1, figsize=(11, 6), sharex=True)
    for ax, quantile in zip(axes, (0.5, 0.99)):
        per_second = ok.groupby(["second", "method"])["latency_ms"].quantile(quantile).unstack("method")
        per_second.reindex(columns=[m for m in METHODS if m in per_second]).plot(ax=ax, linewidth=1.2)
        mark_faults(ax, run)
        ax.set(ylabel=f"p{int(quantile * 100)} latency (ms)", yscale="log")
        ax.legend(loc="upper left", fontsize=8)

    axes[0].set_title(title(run, "Operation latency per second"), pad=16)
    axes[1].set_xlabel("Time since load start (s)")
    save(fig, run["dir"] / "latency.png")


def plot_visibility(run, pairs):
    fig, ax = plt.subplots(figsize=(11, 4))

    for node, group in pairs.groupby("node"):
        ax.scatter(group["t_sent"], group["delay_ms"].clip(lower=0.1), s=10, label=node)

    mark_faults(ax, run)
    ax.set(xlabel="Probe write time (s)", ylabel="Time until visible (ms)", yscale="log")
    ax.set_title(title(run, "Probe visibility by node"), pad=16)
    ax.legend(loc="upper left", fontsize=8)
    save(fig, run["dir"] / "visibility.png")


def plot_load_balancer(run):
    lb = run["lb"].assign(second=run["lb"]["t"].astype(int))
    served = lb[lb["code"] == 200].groupby(["second", "instance_id"]).size().unstack("instance_id", fill_value=0)
    errors = lb[lb["code"] != 200].groupby(["second", "code"]).size().unstack("code", fill_value=0)

    fig, axes = plt.subplots(2, 1, figsize=(11, 6), sharex=True)
    served.plot(ax=axes[0], linewidth=1)
    axes[0].set(ylabel="Requests served per second")
    axes[0].legend(loc="upper left", fontsize=7, ncol=3)
    axes[0].set_title(title(run, "Load balancer traffic"), pad=16)

    if not errors.empty:
        errors.plot(ax=axes[1], linewidth=1.2, color="tab:red")
    axes[1].set(ylabel="Errors per second", xlabel="Time since load start (s)")

    for ax in axes:
        mark_faults(ax, run)

    save(fig, run["dir"] / "load_balancer.png")


def plot_cluster_state(run):
    metrics = run["metrics"]
    if metrics.empty:
        return

    metrics = metrics.assign(value=metrics["value"].astype(float))
    panels = []

    changes = leader_changes(run)
    if changes:
        panels.append(("leader", changes))

    swim = metrics[(metrics["metric"] == "registry_swim_members") & (metrics["label"] != "ALIVE")]
    if len(swim):
        panels.append(("swim", swim))

    healthy = metrics[(metrics["metric"] == "registry_instances") & (metrics["label"] == "HEALTHY")]
    if len(healthy):
        panels.append(("healthy", healthy))

    pool = metrics[metrics["metric"] == "lb_pool_instances"]
    if len(pool):
        panels.append(("pool", pool))

    fig, axes = plt.subplots(len(panels), 1, figsize=(11, 2.6 * len(panels)), sharex=True, squeeze=False)

    for ax, (kind, data) in zip(axes[:, 0], panels):
        if kind == "leader":
            times = [t for t, _ in data] + [run["end"]]
            leaders = [leader for _, leader in data]
            ax.step(times, leaders + leaders[-1:], where="post")
            ax.set(ylabel="Raft leader")
        elif kind == "swim":
            series = data.groupby(["t", "instance", "label"])["value"].sum().unstack(["instance", "label"])
            series.plot(ax=ax, drawstyle="steps-post", linewidth=1)
            ax.set(ylabel="Peers seen as")
            ax.legend(loc="upper left", fontsize=7, ncol=3)
        elif kind == "healthy":
            data.pivot_table(index="t", columns="instance", values="value").plot(ax=ax, drawstyle="steps-post")
            ax.set(ylabel="HEALTHY instances")
            ax.legend(loc="lower left", fontsize=7, ncol=5)
        else:
            data.pivot_table(index="t", columns="label", values="value").plot(ax=ax, drawstyle="steps-post")
            ax.set(ylabel="LB pool size")
            ax.legend(loc="lower left", fontsize=7)

        mark_faults(ax, run)

    axes[0, 0].set_title(title(run, "Cluster state (Prometheus, 5 s resolution)"), pad=16)
    axes[-1, 0].set_xlabel("Time since load start (s)")
    save(fig, run["dir"] / "cluster_state.png")


def markdown_table(headers, rows):
    lines = ["| " + " | ".join(headers) + " |", "|" + "---|" * len(headers)]
    lines += ["| " + " | ".join(str(cell) for cell in row) + " |" for row in rows]
    return "\n".join(lines)


def number(value, digits=2):
    if isinstance(value, float):
        return f"{value:.{digits}f}"

    return str(value)


def write_summary(run, summary, groups):
    faults = run["faults"]
    sections = [f"# {run['experiment']} — {run['backend']} — run {run['repetition']}".strip()]

    if not faults.empty:
        rows = [
            (number(row.t, 1), row.action, row.target if isinstance(row.target, str) else "")
            for row in faults.itertuples()
        ]
        sections += ["## Faults", markdown_table(["t (s)", "action", "target"], rows)]

    sections += [
        "## Node groups",
        markdown_table(["group", "nodes"], [(group, " ".join(nodes)) for group, nodes in groups.items()]),
    ]

    by_prefix = {}
    for key, value in summary.items():
        prefix, _, rest = key.partition(".")
        by_prefix.setdefault(prefix, []).append((rest or prefix, number(value)))

    for prefix, rows in by_prefix.items():
        if prefix in ("experiment", "backend", "repetition"):
            continue

        sections += [f"## {prefix}", markdown_table(["metric", "value"], rows)]

    changes = leader_changes(run)
    if changes:
        rows = [(number(t, 1), leader) for t, leader in changes]
        sections += ["## Raft leader changes", markdown_table(["t (s)", "leader"], rows)]

    (run["dir"] / "summary.md").write_text("\n\n".join(sections) + "\n")


def analyze_run(run_dir):
    run = load_run(run_dir)
    windows = phase_windows(run)
    groups = node_groups(run)
    pairs = visibility_pairs(run)

    summary = {"experiment": run["experiment"], "backend": run["backend"], "repetition": run["repetition"]}
    summary |= availability(run, windows, groups)
    summary |= outages(run, groups)
    summary |= latency(run, windows)
    summary |= visibility(run, pairs, windows)
    summary |= convergence(run, pairs)
    summary |= unreplicated_writes(run, pairs)
    summary |= load_balancer(run, windows)
    summary |= cluster_state(run, windows)

    summary = {key: (float(value) if hasattr(value, "item") else value) for key, value in summary.items()}
    (run_dir / "summary.json").write_text(json.dumps(summary, indent=2) + "\n")
    write_summary(run, summary, groups)

    plot_availability(run)
    plot_latency(run)
    plot_visibility(run, pairs)
    plot_load_balancer(run)
    plot_cluster_state(run)

    return summary


def aggregate(results_dir, summaries):
    frame = pd.DataFrame(summaries)
    metrics = [column for column in frame.columns if column not in ("experiment", "backend", "repetition")]

    long = frame.melt(id_vars=["experiment", "backend"], value_vars=metrics, var_name="metric").dropna()
    long["value"] = long["value"].astype(float)
    table = long.groupby(["experiment", "backend", "metric"])["value"].agg(["mean", "std", "count"]).reset_index()
    table.to_csv(results_dir / "summary.csv", index=False)

    sections = ["# Experiment summary (mean ± std over runs)"]
    for (experiment, backend), group in table.groupby(["experiment", "backend"]):
        rows = [
            (row.metric, f"{row.mean:.2f} ± {0 if pd.isna(row.std) else row.std:.2f}", int(row.count))
            for row in group.itertuples()
        ]
        sections += [f"## {experiment} — {backend}", markdown_table(["metric", "mean ± std", "runs"], rows)]

    (results_dir / "summary.md").write_text("\n\n".join(sections) + "\n")


def main():
    if len(sys.argv) != 2:
        sys.exit("Usage: analyze.py <run directory | results directory>")

    target = Path(sys.argv[1])

    if (target / "operations.csv").exists():
        analyze_run(target)
        print(f"Analyzed {target}")
        return

    run_dirs = sorted(path.parent for path in target.rglob("operations.csv"))
    if not run_dirs:
        sys.exit(f"No runs found under {target}")

    summaries = []
    for run_dir in run_dirs:
        summaries.append(analyze_run(run_dir))
        print(f"Analyzed {run_dir}")

    aggregate(target, summaries)
    print(f"Wrote {target / 'summary.md'} and {target / 'summary.csv'}")


if __name__ == "__main__":
    main()
