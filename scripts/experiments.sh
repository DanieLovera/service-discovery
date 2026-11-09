#!/usr/bin/env bash
set -Eeuo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FAULTS="$ROOT/scripts/faults.sh"
PROMETHEUS_URL="http://localhost:8083"
READY_TIMEOUT=120
CLEAN_ATTEMPTS=3
REGISTRIES=(registry-1 registry-2 registry-3 registry-4 registry-5)
SERVICE_COUNT=3
CP_EXPERIMENTS=(
	e1-baseline e2-kill-leader e2-kill-follower e3-partition
	e4-delay-100ms e4-delay-250ms e4-loss-10pct e5-kill-instance
)
AP_EXPERIMENTS=(
	e1-baseline e2-kill-node e3-partition
	e4-delay-100ms e4-delay-250ms e4-loss-10pct e5-kill-instance
)

trap 'log "Failed at line $LINENO: $BASH_COMMAND"' ERR

usage() {
	cat >&2 <<'EOF'
Usage:
  scripts/experiments.sh run <experiment> <cp|ap> [repetition]
  scripts/experiments.sh all [repetitions]

Experiments (t measured from the load generator start):
  e1-baseline         no faults, 120 s
  e2-kill-leader      (cp) kill the leader at 30 s, start it at 60 s, 120 s
  e2-kill-follower    (cp) kill a follower at 30 s, start it at 60 s, 120 s
  e2-kill-node        (ap) kill registry-3 at 30 s, start it at 60 s, 120 s
  e3-partition        partition 2 vs 3 at 30 s (cp: leader in the minority), heal at 90 s, 150 s
  e4-delay-100ms      add 100 ms of latency at 30 s, restore at 90 s, 120 s
  e4-delay-250ms      add 250 ms of latency at 30 s, restore at 90 s, 120 s
  e4-loss-10pct       add 10% packet loss at 30 s, restore at 90 s, 120 s
  e5-kill-instance    kill service-c-3 at 30 s, start it at 60 s, 90 s

Results are written to results/<experiment>/<backend>/run-<repetition>/.
EOF
}

compose() {
	docker compose --project-directory "$ROOT" "$@"
}

log() {
	echo "[$(date +%H:%M:%S)] $*" >&2
}

now_ms() {
	date +%s%3N
}

fault() {
	"$FAULTS" "$@" >&2
}

experiments_for() {
	case "$1" in
	cp) echo "${CP_EXPERIMENTS[@]}" ;;
	ap) echo "${AP_EXPERIMENTS[@]}" ;;
	*) return 1 ;;
	esac
}

validate() {
	local experiment=$1 backend=$2 experiments
	local -a available

	if ! experiments=$(experiments_for "$backend"); then
		usage
		exit 1
	fi

	read -ra available <<<"$experiments"

	if [[ " ${available[*]} " != *" $experiment "* ]]; then
		echo "Unknown experiment $experiment for backend $backend" >&2
		exit 1
	fi
}

reset_environment() {
	local backend=$1

	log "Resetting environment with BACKEND=$backend"

	if ! clean_environment >"$RUN_DIR/setup.log" 2>&1; then
		log "Could not clean the environment, see $RUN_DIR/setup.log"
		return 1
	fi

	if ! start_environment "$backend" >>"$RUN_DIR/setup.log" 2>&1; then
		log "Could not start the environment, see $RUN_DIR/setup.log"
		return 1
	fi

	wait_until_ready "$backend"
}

clean_environment() {
	local attempt

	for ((attempt = 1; attempt <= CLEAN_ATTEMPTS; attempt++)); do
		if make -C "$ROOT" docker-clean; then
			return 0
		fi
	done

	return 1
}

start_environment() {
	BACKEND=$1 make -C "$ROOT" docker-up &&
		compose --profile experiments build load-generator &&
		BACKEND=$1 compose --profile experiments up -d fault-injector
}

wait_until_ready() {
	local backend=$1
	local deadline=$((SECONDS + READY_TIMEOUT))

	until is_ready "$backend"; do
		if ((SECONDS >= deadline)); then
			log "Environment not ready after ${READY_TIMEOUT}s"
			return 1
		fi

		sleep 1
	done
}

is_ready() {
	local backend=$1 registry metrics

	for registry in "${REGISTRIES[@]}"; do
		compose exec -T "$registry" wget -qO- http://localhost:10000/readyz >/dev/null 2>&1 || return 1
	done

	if [[ $backend == cp ]]; then
		"$FAULTS" leader >/dev/null 2>&1 || return 1
	fi

	metrics=$(compose exec -T load-balancer wget -qO- http://localhost:10010/metrics 2>/dev/null) || return 1

	[[ $(grep -c '^lb_pool_instances{service="[^"]*"} 3$' <<<"$metrics") -eq $SERVICE_COUNT ]]
}

start_load() {
	local duration=$1

	compose --profile experiments run --rm --user "$(id -u):$(id -g)" \
		-e DURATION="${duration}s" -e OUTPUT_DIR="/results/$RUN_PATH" \
		load-generator >"$RUN_DIR/load-generator.log" 2>&1 &
	LOAD_PID=$!

	until [[ -f "$RUN_DIR/operations.csv" ]]; do
		if ! kill -0 "$LOAD_PID" 2>/dev/null; then
			log "Load generator exited before starting, see $RUN_DIR/load-generator.log"
			return 1
		fi

		sleep 0.1
	done

	START_MS=$(now_ms)
	log "Load generator started for ${duration}s"
}

export_metrics() {
	if ! python3 -I "$ROOT/scripts/export_metrics.py" "$PROMETHEUS_URL" \
		"$((START_MS / 1000 - 5))" "$(($(now_ms) / 1000))" "$RUN_DIR/metrics.csv"; then
		log "Could not export Prometheus metrics"
		return 1
	fi
}

at() {
	local offset=$1
	shift

	local wait_ms=$((START_MS + offset * 1000 - $(now_ms)))
	if ((wait_ms > 0)); then
		sleep "$((wait_ms / 1000)).$(printf '%03d' $((wait_ms % 1000)))"
	fi

	"$@"
}

follower_of() {
	local leader=$1 registry

	for registry in "${REGISTRIES[@]}"; do
		if [[ $registry != "$leader" ]]; then
			echo "$registry"
			return
		fi
	done
}

crash_registry() {
	local target=$1

	start_load 120
	at 30 fault kill "$target"
	at 60 fault start "$target"
}

partition_cluster() {
	local minority=$1 majority=$2

	start_load 150
	at 30 fault partition "$minority" "$majority"
	at 90 fault heal
}

split_around_leader() {
	local leader=$1 registry
	local -a minority=("$leader") majority=()

	for registry in "${REGISTRIES[@]}"; do
		[[ $registry == "$leader" ]] && continue

		if ((${#minority[@]} < 2)); then
			minority+=("$registry")
		else
			majority+=("$registry")
		fi
	done

	(
		IFS=,
		echo "${minority[*]} ${majority[*]}"
	)
}

degrade_network() {
	start_load 120
	at 30 fault degrade "$1" "$2"
	at 90 fault restore
}

crash_instance() {
	start_load 90
	at 30 fault kill "$1"
	at 60 fault start "$1"
}

run_experiment() {
	local experiment=$1 backend=$2 repetition=${3:-1}
	local leader minority majority

	validate "$experiment" "$backend"

	RUN_PATH="$experiment/$backend/run-$repetition"
	RUN_DIR="$ROOT/results/$RUN_PATH"
	export FAULTS_FILE="$RUN_DIR/faults.csv"

	rm -rf "$RUN_DIR"
	mkdir -p "$RUN_DIR"

	log "Running $RUN_PATH"
	reset_environment "$backend"

	case "$experiment" in
	e1-baseline)
		start_load 120
		;;
	e2-kill-leader)
		crash_registry "$("$FAULTS" leader)"
		;;
	e2-kill-follower)
		crash_registry "$(follower_of "$("$FAULTS" leader)")"
		;;
	e2-kill-node)
		crash_registry registry-3
		;;
	e3-partition)
		if [[ $backend == cp ]]; then
			leader=$("$FAULTS" leader)
			read -r minority majority < <(split_around_leader "$leader")
		else
			minority="registry-1,registry-2"
			majority="registry-3,registry-4,registry-5"
		fi

		partition_cluster "$minority" "$majority"
		;;
	e4-delay-100ms)
		degrade_network 100ms 0%
		;;
	e4-delay-250ms)
		degrade_network 250ms 0%
		;;
	e4-loss-10pct)
		degrade_network 0ms 10%
		;;
	e5-kill-instance)
		crash_instance service-c-3
		;;
	esac

	wait "$LOAD_PID"
	export_metrics
	log "Finished $RUN_PATH"
}

run_all() {
	local repetitions=${1:-3} repetition backend experiment
	local -a failed=()

	for ((repetition = 1; repetition <= repetitions; repetition++)); do
		for backend in cp ap; do
			for experiment in $(experiments_for "$backend"); do
				if ! "$ROOT/scripts/experiments.sh" run "$experiment" "$backend" "$repetition"; then
					failed+=("$experiment/$backend/run-$repetition")
				fi
			done
		done
	done

	if ((${#failed[@]} > 0)); then
		log "Failed runs: ${failed[*]}"
		return 1
	fi

	log "All runs finished"
}

main() {
	local mode=${1:-}
	shift || true

	case "$mode" in
	run) run_experiment "${1:?experiment required}" "${2:?backend required}" "${3:-1}" ;;
	all) run_all "${1:-3}" ;;
	*)
		usage
		exit 1
		;;
	esac
}

main "$@"
