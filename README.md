# DANIEL_TPIII — Service Registry & Load Balancer

## Procesos

- `registry`: nodo del cluster de Service Registry. Expone la API gRPC (Register, Update, Deregister,
  Heartbeat, Lookup y Watch) y replica el estado con uno de dos backends: CP (Raft) o AP (Gossip, SWIM y
  Anti-Entropy).
- `load-balancer`: hace Watch de los servicios configurados y rutea `/{servicio}/{path}` hacia
  `http://{instancia}/{path}` con la estrategia de cada servicio.
- `mock-service`: instancia de prueba que se registra, envía heartbeats y responde las requests que le
  rutea el load balancer.
- `load-generator`: genera la carga de los experimentos. Corre ciclos Register → Heartbeat → Lookup →
  Deregister contra cada registry, sondas de visibilidad con un Watch abierto en cada nodo y tráfico HTTP
  al load balancer, y guarda cada operación en CSV.

## Requisitos

- Go 1.26.x
- Docker con Docker Compose
- `staticcheck` para ejecutar `make quality`
- `protoc`, `protoc-gen-go` y `protoc-gen-go-grpc`, solo para regenerar código con `make proto`
- Python 3 para los experimentos: `experiments.sh` exporta las métricas de Prometheus con la biblioteca
  estándar y `analyze.py` necesita además pandas y matplotlib (ver [Análisis](#análisis))

```sh
go install honnef.co/go/tools/cmd/staticcheck@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Los scripts de experimentos están pensados para Linux.

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

### Load Generator

| Variable | Requerida | Default | Descripción |
| --- | --- | --- | --- |
| `REGISTRY_ADDRESSES` | sí | — | Registry nodes, separados por coma; se abre una conexión fija a cada uno, sin failover |
| `LOAD_BALANCER_URL` | sí | — | URL del load balancer |
| `LOAD_BALANCER_SERVICES` | sí | — | Servicios a los que se envía tráfico HTTP, separados por coma |
| `DURATION` | no | `1m` | Duración de la carga |
| `REQUEST_TIMEOUT` | no | `2s` | Timeout de cada operación gRPC y de cada request HTTP |
| `REGISTRY_CYCLE_INTERVAL` | no | `200ms` | Intervalo entre ciclos Register → Heartbeat → Lookup → Deregister en cada nodo |
| `VISIBILITY_PROBE_INTERVAL` | no | `1s` | Intervalo entre sondas de visibilidad |
| `LOAD_BALANCER_REQUEST_INTERVAL` | no | `50ms` | Intervalo entre requests HTTP a cada servicio |
| `OUTPUT_DIR` | no | `/results` | Directorio donde se escriben los CSV |

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
| Grafana | `http://localhost:8082` (usuario `admin`, contraseña `admin`) |
| Prometheus | `http://localhost:8083` |
| Dozzle (con `make docker-up-debug`) | `http://localhost:8081` |

Grafana carga el dashboard **TPIII Overview**, que se actualiza cada 5 s y tiene cuatro filas:

- **Registry operations**: requests y errores por nodo, latencias p50 y p99 por método.
- **Cluster state**: líder de Raft (CP), peers vistos como SUSPECT o CONFIRM por SWIM (AP) e instancias
  HEALTHY, EXPIRED y DELETED en cada nodo.
- **Registry clients (load balancer and mock services)**: pool del LB por servicio, requests por
  instancia, respuestas 5xx y llamadas de los mocks al registry.
- **Processes**: targets activos, goroutines, memoria y uptime de cada proceso.

## Experimentos

Los experimentos comparan los dos backends bajo la misma carga y con fallas inyectadas en momentos fijos,
contados desde que arranca el generador de carga:

| Experimento | Backend | Falla | Duración |
| --- | --- | --- | --- |
| `e1-baseline` | ambos | ninguna | 120 s |
| `e2-kill-leader` | CP | mata al líder de Raft a los 30 s y lo levanta a los 60 s | 120 s |
| `e2-kill-follower` | CP | mata a un seguidor a los 30 s y lo levanta a los 60 s | 120 s |
| `e2-kill-node` | AP | mata `registry-3` a los 30 s y lo levanta a los 60 s | 120 s |
| `e3-partition` | ambos | partición 2 contra 3 nodos de los 30 s a los 90 s; en CP el líder queda en la minoría | 150 s |
| `e4-delay-100ms` | ambos | 100 ms de demora en la red saliente de los registries de los 30 s a los 90 s | 120 s |
| `e4-delay-250ms` | ambos | 250 ms de demora en la red saliente de los registries de los 30 s a los 90 s | 120 s |
| `e4-loss-10pct` | ambos | 10% de pérdida de paquetes en la salida de los registries de los 30 s a los 90 s | 120 s |
| `e5-kill-instance` | ambos | mata `service-c-3` a los 30 s y la levanta a los 60 s | 90 s |

La partición solo corta el tráfico entre registries: el generador y los clientes siguen llegando a todos
los nodos, lo que permite ver qué responde cada lado.

### Ejecución

```sh
scripts/experiments.sh run e3-partition ap 1   # un experimento, repetición 1
scripts/experiments.sh all                     # todos, 3 repeticiones por backend
scripts/experiments.sh all 1                   # todos, una repetición
```

Cada corrida es independiente y tarda entre 3 y 5 minutos; el set completo con 3 repeticiones son 45
corridas, unas 3 horas. Cada corrida:

1. Ejecuta `make docker-clean`: borra el entorno que esté levantado, incluido el estado de Raft.
2. Levanta el entorno con el backend elegido, construye el generador y levanta el `fault-injector`.
3. Espera a que los registries estén listos, haya líder (CP) y el LB tenga las 3 instancias de cada
   servicio.
4. Corre el generador de carga e inyecta las fallas del experimento.
5. Exporta a CSV las métricas de Prometheus de la corrida, porque la limpieza siguiente las borra.

Repetir una corrida con el mismo número la sobrescribe. Las mediciones son sensibles a la carga de la
máquina: conviene no correr procesos pesados durante los experimentos. Grafana se puede seguir en vivo,
pero se recrea en cada corrida (hay que recargar la página y volver a iniciar sesión).

### Resultados

Cada corrida escribe en `results/<experimento>/<backend>/run-<repetición>/`:

| Archivo | Contenido |
| --- | --- |
| `operations.csv` | Cada operación de los ciclos: nodo, método, código gRPC y latencia |
| `probes.csv` | Cada sonda de visibilidad escrita: nodo que la recibió y código |
| `sightings.csv` | Primera vez que cada nodo mostró cada sonda en su Watch |
| `load_balancer.csv` | Cada request HTTP al LB: servicio, código, instancia que respondió y latencia |
| `faults.csv` | Momento de cada falla inyectada |
| `metrics.csv` | Series de Prometheus: líder de Raft, miembros de SWIM, instancias por estado y pool del LB |
| `setup.log`, `load-generator.log` | Logs de la preparación del entorno y del generador |

Los tiempos son Unix en milisegundos. `results/` se versiona, salvo `results/demo/` (ver [Demo](#demo)).

### Análisis

```sh
python3 -m venv .venv
.venv/bin/pip install -r scripts/requirements.txt

.venv/bin/python scripts/analyze.py results/e3-partition/ap/run-1   # una corrida
.venv/bin/python scripts/analyze.py results                         # todas, más el resumen agregado
```

Por cada corrida genera en su carpeta `summary.md` y `summary.json` con disponibilidad, latencias,
visibilidad, convergencia, escrituras perdidas, errores del LB y estado del cluster, separados en las
fases antes, durante y después de la falla. También genera los gráficos `availability.png`,
`latency.png`, `visibility.png`, `load_balancer.png` y `cluster_state.png`. Se descartan los primeros
10 s de cada corrida como calentamiento y los últimos 5 s.

Sobre `results`, además escribe `results/summary.md` y `results/summary.csv` con la media y el desvío de
cada métrica por experimento y backend. Analizar las 45 corridas tarda unos 15 minutos, casi todo en
generar los gráficos.

### Fallas a mano

`scripts/faults.sh` inyecta las mismas fallas sobre el entorno levantado. Corre dentro del contenedor
`fault-injector`, que hay que levantar primero:

```sh
docker compose --profile experiments up -d fault-injector

scripts/faults.sh leader                                        # líder actual de Raft (CP)
scripts/faults.sh kill registry-2                               # SIGKILL al proceso del contenedor
scripts/faults.sh start registry-2
scripts/faults.sh partition registry-1,registry-2 registry-3,registry-4,registry-5
scripts/faults.sh heal
scripts/faults.sh degrade 100ms 5%                              # demora y pérdida en cada registry
scripts/faults.sh restore
```

### Demo

Para mostrar el sistema en vivo conviene no usar `experiments.sh`, que limpia y reconstruye el entorno en cada
corrida. El entorno se levanta una sola vez, el generador corre con una duración larga y las fallas se inyectan a
mano mientras se sigue Grafana con el rango en "Last 5 minutes":

```sh
make demo-up              # entorno, imagen del generador y fault-injector; BACKEND=ap make demo-up para AP
make demo-load            # generador por 20 min en primer plano; make demo-load DEMO_DURATION=30m lo cambia
```

Con el generador corriendo, en otra terminal:

| Qué mostrar | Falla | Vuelta atrás |
| --- | --- | --- |
| Caída del líder (CP) | `scripts/faults.sh kill <líder>` | `scripts/faults.sh start <líder>` |
| Partición | `scripts/faults.sh partition registry-1,registry-2 registry-3,registry-4,registry-5` | `scripts/faults.sh heal` |
| Red degradada | `scripts/faults.sh degrade 250ms 0%` | `scripts/faults.sh restore` |
| Caída de una instancia | `scripts/faults.sh kill service-c-3` | `scripts/faults.sh start service-c-3` |

`<líder>` es el nodo que muestra `scripts/faults.sh leader`.

Para pasar al otro backend sin reconstruir las imágenes: cortar el generador, `make docker-down`,
`BACKEND=ap make demo-up` y `make demo-load` de nuevo.

El generador escribe en `results/demo/`, que no se versiona y se sobrescribe en cada `make demo-load`. Grafana
muestra qué pasa; los tiempos y latencias precisos están en los resultados de los experimentos.

## Desarrollo

```sh
make fmt-check    # verifica el formato
make vet          # análisis estático de go vet
make test-race    # tests con el detector de data races
make quality      # fmt-check, vet, staticcheck y test-race
make proto        # regenera gen/ desde proto/
```

`make help` lista todos los targets.
