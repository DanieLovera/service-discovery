# DANIEL_TPIII — Service Registry & Load Balancer

## Procesos

- `registry`: nodo del cluster de Service Registry. Expone la API gRPC (Register, Update, Deregister,
  Heartbeat, Lookup y Watch) y replica el estado con uno de dos backends: CP (Raft) o AP (Gossip, SWIM y
  Anti-Entropy).
- `load-balancer`: hace Watch de los servicios configurados y rutea `/{servicio}/{path}` hacia
  `http://{instancia}/{path}` con la estrategia de cada servicio.
- `mock-service`: instancia de prueba que se registra, envía heartbeats y responde las requests que le
  rutea el load balancer.

## Requisitos

- Go 1.26.x
- Docker con Docker Compose
- `staticcheck` para ejecutar `make quality`
- `protoc`, `protoc-gen-go` y `protoc-gen-go-grpc`, solo para regenerar código con `make proto`

```sh
go install honnef.co/go/tools/cmd/staticcheck@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Configuración

La configuración se recibe por variables de entorno. `docker-compose.yml` define los valores del entorno
local. Las duraciones usan el formato de `time.ParseDuration` de Go (`300ms`, `5s`, `1m`).

### Compartidas

| Variable | Default | Valores | Descripción |
| --- | --- | --- | --- |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warn`, `error` | Nivel mínimo de los logs |
| `LOG_FORMAT` | `json` | `json`, `text` | Formato de los logs |

### Registry

| Variable | Requerida | Default | Backend | Descripción |
| --- | --- | --- | --- | --- |
| `INSTANCE_ID` | sí | — | ambos | Id del nodo, debe estar en `CLUSTER_MEMBERS` |
| `BACKEND` | sí | — | — | Backend de consistencia: `cp` (Raft) o `ap` (Gossip + SWIM) |
| `CLUSTER_MEMBERS` | sí | — | ambos | Membership estático: `id=grpc_address\|peer_address` separados por coma |
| `GRPC_BIND_ADDRESS` | no | `0.0.0.0:7000` | ambos | Escucha de la API gRPC, Forwarding (CP) y Anti-Entropy (AP) |
| `PEER_BIND_ADDRESS` | no | `0.0.0.0:7001` | ambos | Escucha entre nodos: Raft/TCP (CP), SWIM y Gossip/UDP (AP) |
| `OBSERVABILITY_BIND_ADDRESS` | no | `0.0.0.0:10000` | ambos | Escucha de `/metrics`, `/healthz` y `/readyz` |
| `HEARTBEAT_TTL_MULTIPLIER` | no | `3` | ambos | TTL de cada instancia = `heartbeat_interval × multiplicador` |
| `RAFT_DATA_DIR` | no | `/data/raft` | CP | Directorio del log, el estado y los snapshots de Raft |
| `FORWARD_TIMEOUT` | no | `2s` | CP | Timeout de cada intento de reenvío al líder |
| `FORWARD_MAX_ATTEMPTS` | no | `3` | CP | Intentos de reenvío al líder antes de devolver error |
| `GOSSIP_FANOUT` | no | `3` | AP | Máximo de nodos ALIVE a los que se propaga cada versión nueva |
| `ANTI_ENTROPY_INTERVAL` | no | `5s` | AP | Tiempo entre rondas de anti-entropy, también es el timeout de cada ronda |
| `SWIM_PROBE_INTERVAL` | no | `1s` | AP | Período de SWIM: se chequea un peer por período |
| `SWIM_PING_TIMEOUT` | no | `300ms` | AP | Espera máxima del ACK a un PING directo |
| `SWIM_PING_REQ_TIMEOUT` | no | `600ms` | AP | Espera máxima del ACK indirecto luego de enviar los PING_REQ |
| `SWIM_PING_REQ_COUNT` | no | `2` | AP | Cantidad de intermediarios que reciben un PING_REQ |
| `SWIM_SUSPECT_TIMEOUT` | no | `5s` | AP | Tiempo que un nodo queda SUSPECT antes de pasar a CONFIRM |
| `SWIM_RETRANSMIT_MULTIPLIER` | no | `4` | AP | λ de `λ·⌈log10(N+1)⌉`: reenvíos de cada novedad de membership |

`SWIM_PING_TIMEOUT + SWIM_PING_REQ_TIMEOUT` no debe superar `SWIM_PROBE_INTERVAL`.

### Load Balancer

| Variable | Requerida | Default | Descripción |
| --- | --- | --- | --- |
| `REGISTRY_ADDRESSES` | sí | — | Direcciones gRPC de los registry nodes, separadas por coma, usadas para Watch |
| `SERVICES` | sí | — | Servicios a balancear: `servicio:estrategia` separados por coma |
| `HTTP_BIND_ADDRESS` | no | `0.0.0.0:8000` | Escucha HTTP: `/{servicio}/...` se rutea a una instancia del servicio |
| `OBSERVABILITY_BIND_ADDRESS` | no | `0.0.0.0:10010` | Escucha de `/metrics`, `/healthz` y `/readyz` |
| `REGISTRY_RETRY_INTERVAL` | no | `1s` | Espera antes de reabrir un Watch que se cortó |
| `REGISTRY_KEEPALIVE_TIME` | no | `10s` | Inactividad tras la cual se envía un keepalive gRPC (mínimo `10s`) |
| `REGISTRY_KEEPALIVE_TIMEOUT` | no | `5s` | Espera máxima de la respuesta al keepalive |

Las estrategias disponibles son `round-robin`, `least-connections` y `weighted`.

### Mock Service

| Variable | Requerida | Default | Descripción |
| --- | --- | --- | --- |
| `NAME` | sí | — | Nombre del servicio con el que se registra la instancia |
| `INSTANCE_ID` | sí | — | Id de la instancia dentro del servicio |
| `ADVERTISE_ADDRESS` | sí | — | Dirección publicada en el registry, por la que el LB llega a la instancia |
| `REGISTRY_ADDRESSES` | sí | — | Registry nodes para Register y heartbeats, si uno falla se prueba el siguiente |
| `HTTP_BIND_ADDRESS` | no | `0.0.0.0:9000` | Escucha HTTP de las requests del load balancer |
| `OBSERVABILITY_BIND_ADDRESS` | no | `0.0.0.0:10020` | Escucha de `/metrics`, `/healthz` y `/readyz` |
| `WEIGHT` | no | `1` | Peso de la instancia para la estrategia weighted |
| `HEARTBEAT_INTERVAL` | no | `5s` | Intervalo entre heartbeats, el registry calcula el TTL a partir de él |

## Entorno local

El entorno de Docker Compose levanta 5 nodos de registry, el load balancer, tres servicios con tres
instancias cada uno (`service-a` con round-robin, `service-b` con least-connections y `service-c` con
weighted y pesos 1, 2 y 3), Prometheus y Grafana.

```sh
make docker-up              # backend CP
BACKEND=ap make docker-up   # backend AP
make docker-up-debug        # además levanta Dozzle para ver los logs
make docker-down            # detiene el entorno
make docker-clean           # además borra imágenes locales y volúmenes (estado de Raft)
```

| Servicio | URL |
| --- | --- |
| Load balancer | `http://localhost:8000/{servicio}/`, por ejemplo `curl localhost:8000/service-a/` |
| Métricas del load balancer | `http://localhost:10010/metrics` |
| Grafana | `http://localhost:8082` |
| Prometheus | `http://localhost:8083` |
| Dozzle (con `make docker-up-debug`) | `http://localhost:8081` |

## Desarrollo

```sh
make fmt-check    # verifica el formato
make vet          # análisis estático de go vet
make test-race    # tests con el detector de data races
make quality      # fmt-check, vet, staticcheck y test-race
make proto        # regenera gen/ desde proto/
```

`make help` lista todos los targets.
