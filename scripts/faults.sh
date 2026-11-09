#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROJECT="${COMPOSE_PROJECT_NAME:-$(basename "$ROOT")}"
INJECTOR=fault-injector
REGISTRIES=(registry-1 registry-2 registry-3 registry-4 registry-5)
INSPECT_FORMAT='{{index .Config.Labels "com.docker.compose.service"}} {{.Id}} {{.State.Pid}} '
INSPECT_FORMAT+='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}'

declare -A IDS PIDS IPS

usage() {
	cat >&2 <<'EOF'
Usage: scripts/faults.sh <action> [arguments]

  kill <service>              Kill a container's main process with SIGKILL
  start <service>             Start a stopped container
  leader                      Print the current Raft leader (CP backend)
  partition <nodes> <nodes>   Drop traffic between two comma-separated groups of registry nodes
  heal                        Remove the partition
  degrade <delay> <loss>      Add latency and packet loss to every registry node, e.g. 100ms 5%
  restore                     Remove the added latency and packet loss

Faults run inside the fault-injector container:
  docker compose --profile experiments up -d fault-injector

Each fault is appended to $FAULTS_FILE when it is set.
EOF
}

record() {
	local micros
	micros=$(date +%s%6N)

	local row="${micros:0:-3}.${micros: -3},$1,${2:-}"
	echo "$row"

	if [[ -n "${FAULTS_FILE:-}" ]]; then
		[[ -s "$FAULTS_FILE" ]] || echo "timestamp_ms,action,target" >"$FAULTS_FILE"
		echo "$row" >>"$FAULTS_FILE"
	fi
}

resolve() {
	local -a ids
	local service id pid ip

	mapfile -t ids < <(docker ps -a -q --filter "label=com.docker.compose.project=$PROJECT")

	while read -r service id pid ip; do
		IDS[$service]=$id
		PIDS[$service]=$pid
		IPS[$service]=$ip
	done < <(docker inspect -f "$INSPECT_FORMAT" "${ids[@]}")
}

require_injector() {
	if [[ ${PIDS[$INJECTOR]:-0} == 0 ]]; then
		echo "$INJECTOR is not running, start it with: docker compose --profile experiments up -d $INJECTOR" >&2
		exit 1
	fi
}

apply_and_record() {
	local script=$1 action=$2 target=${3:-}

	docker exec "${IDS[$INJECTOR]}" sh -c "${script}echo applied" | {
		read -r _
		record "$action" "$target"
		cat >/dev/null
	}
}

in_each_registry() {
	local action=$1 target=$2
	shift 2

	local script="" node command

	for node in "${REGISTRIES[@]}"; do
		for command in "$@"; do
			script+="nsenter -t ${PIDS[$node]} -n $command && "
		done
	done

	apply_and_record "$script" "$action" "$target"
}

kill_service() {
	local pid=${PIDS[$1]:-0}

	if [[ $pid == 0 ]]; then
		echo "$1 is not running" >&2
		return 1
	fi

	apply_and_record "kill -KILL $pid && " kill "$1"

	if ! timeout 15 docker wait "${IDS[$1]}" >/dev/null; then
		echo "$1 did not stop" >&2
		return 1
	fi
}

start_service() {
	docker start "${IDS[$1]:?unknown service $1}" >/dev/null
	record start "$1"
}

leader() {
	local script

	script="for registry in ${REGISTRIES[*]}; do "
	script+="curl -s -m 1 http://\$registry:10000/metrics | grep -q '^registry_raft_state{state=\"LEADER\"} 1\$' "
	script+="&& echo \$registry && exit 0; done; exit 1"

	if ! docker exec "${IDS[$INJECTOR]}" sh -c "$script"; then
		echo "no leader found" >&2
		return 1
	fi
}

partition() {
	local -a side_a side_b
	local script=""

	IFS=, read -ra side_a <<<"$1"
	IFS=, read -ra side_b <<<"$2"

	script+=$(drop_rules side_a side_b)
	script+=$(drop_rules side_b side_a)

	apply_and_record "$script" partition "${side_a[*]} | ${side_b[*]}"
}

drop_rules() {
	local -n nodes=$1 peers=$2
	local node peer pid ip

	for node in "${nodes[@]}"; do
		pid=${PIDS[$node]}

		for peer in "${peers[@]}"; do
			ip=${IPS[$peer]}
			printf 'nsenter -t %s -n iptables -A INPUT -s %s -j DROP && ' "$pid" "$ip"
			printf 'nsenter -t %s -n iptables -A OUTPUT -d %s -j DROP && ' "$pid" "$ip"
		done
	done
}

heal() {
	in_each_registry heal "" "iptables -F INPUT" "iptables -F OUTPUT"
}

degrade() {
	in_each_registry degrade "delay=$1 loss=$2" "tc qdisc add dev eth0 root netem delay $1 loss $2"
}

restore() {
	in_each_registry restore "" "tc qdisc del dev eth0 root"
}

main() {
	local action=${1:-}
	shift || true

	case "$action" in
	kill | leader | partition | heal | degrade | restore)
		resolve
		require_injector
		;;
	start)
		resolve
		;;
	*)
		usage
		exit 1
		;;
	esac

	case "$action" in
	kill) kill_service "${1:?service required}" ;;
	start) start_service "${1:?service required}" ;;
	leader) leader ;;
	partition) partition "${1:?first group required}" "${2:?second group required}" ;;
	heal) heal ;;
	degrade) degrade "${1:?delay required}" "${2:?loss required}" ;;
	restore) restore ;;
	esac
}

main "$@"
