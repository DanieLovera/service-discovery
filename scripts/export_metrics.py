import csv
import json
import sys
import urllib.parse
import urllib.request

STEP = "5s"
QUERIES = [
    ("registry_raft_state", "registry_raft_state", "state"),
    ("registry_swim_members", "registry_swim_members", "status"),
    ("registry_instances", 'sum by (instance, status) (registry_instances{service=~"service-.*"})', "status"),
    ("lb_pool_instances", "lb_pool_instances", "service"),
    ("up", "up", "job"),
]


def query_range(prometheus_url, expr, start, end):
    params = urllib.parse.urlencode({"query": expr, "start": start, "end": end, "step": STEP})

    with urllib.request.urlopen(f"{prometheus_url}/api/v1/query_range?{params}", timeout=30) as response:
        body = json.load(response)

    if body.get("status") != "success":
        raise RuntimeError(f"query {expr!r} failed: {body}")

    return body["data"]["result"]


def main():
    if len(sys.argv) != 5:
        sys.exit("Usage: export_metrics.py <prometheus-url> <start-seconds> <end-seconds> <output.csv>")

    prometheus_url, start, end, output = sys.argv[1:]

    with open(output, "w", newline="") as file:
        writer = csv.writer(file, lineterminator="\n")
        writer.writerow(["timestamp_ms", "metric", "instance", "label", "value"])

        for metric, expr, label in QUERIES:
            for series in query_range(prometheus_url, expr, start, end):
                instance = series["metric"].get("instance", "").split(":")[0]
                label_value = series["metric"].get(label, "")

                for timestamp, value in series["values"]:
                    writer.writerow([f"{float(timestamp) * 1000:.3f}", metric, instance, label_value, value])


if __name__ == "__main__":
    main()
