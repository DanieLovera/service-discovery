# DANIEL_TPIII — Service Registry & Load Balancer

## Procesos

El diseño define tres procesos independientes:

- `registry`: nodo del cluster de Service Registry.
- `load-balancer`: proceso que consumirá el Registry y ruteará tráfico.
- `mock-service`: instancia de prueba que posteriormente se registrará, enviará
  heartbeats y recibirá carga sintética.

## Requisitos

- Go 1.23.x
- Docker (para imágenes)
- `staticcheck` para ejecutar `make quality`

Instalación de staticcheck:

```sh
go install honnef.co/go/tools/cmd/staticcheck@latest
```

## Configuración

La configuración se recibe por variables de entorno. `.env.example` documenta
un entorno local de referencia, los binarios no leen archivos `.env` automáticamente.

### Compartidas

| Variable | Default | Valores |
| --- | --- | --- |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warn`, `error` |
| `LOG_FORMAT` | `json` | `json`, `text` |

### Registry

| Variable | Requerida | Default |
| --- | --- | --- |
| `REGISTRY_INSTANCE_ID` | sí | — |
| `REGISTRY_BACKEND` | no | `ap` |
| `REGISTRY_GRPC_ADDRESS` | no | `0.0.0.0:7000` |
| `REGISTRY_CLUSTER_MEMBERS` | sí | — |
| `HEARTBEAT_TTL_MULTIPLIER` | no | `3` |

`HEARTBEAT_TTL_MULTIPLIER` define el TTL de cada instancia como
`heartbeat_interval × multiplier`, tolerando hasta `multiplier - 1` heartbeats perdidos.

`REGISTRY_CLUSTER_MEMBERS` es una lista separada por comas. Refleja el membership
estático definido en el diseño. En etapas posteriores se ampliará la configuración
con las direcciones específicas necesarias para Raft, SWIM, Gossip y Anti-Entropy.

### Load Balancer

| Variable | Requerida | Default |
| --- | --- | --- |
| `LB_HTTP_ADDRESS` | no | `0.0.0.0:8080` |
| `LB_REGISTRY_ADDRESSES` | sí | — |
| `SERVICES` | sí | — |
| `REGISTRY_RETRY_INTERVAL` | no | `1s` |
| `REGISTRY_KEEPALIVE_TIME` | no | `10s` |
| `REGISTRY_KEEPALIVE_TIMEOUT` | no | `5s` |

`SERVICES` es una lista separada por comas de `servicio:estrategia`, con estrategia
`round-robin`, `least-connections` o `weighted`. El Load Balancer rutea
`/{servicio}/{path}` hacia `http://{instancia}/{path}`. Las variables `REGISTRY_*`
configuran la conexión de Watch con el Registry. `REGISTRY_KEEPALIVE_TIME` debe ser
de al menos `10s` (mínimo que admite gRPC).

### Mock Service

| Variable | Requerida | Default |
| --- | --- | --- |
| `MS_NAME` | sí | — |
| `MS_INSTANCE_ID` | sí | — |
| `MS_HTTP_ADDRESS` | no | `0.0.0.0:9000` |
| `MS_REGISTRY_ADDRESSES` | sí | — |
| `MS_WEIGHT` | no | `1` |
| `MS_HEARTBEAT_INTERVAL` | no | `5s` |

Las duraciones usan el formato de `time.ParseDuration` de Go (`500ms`, `2s`,
`1m`, etc.).

## Desarrollo local

```sh
make test
make test-race
make vet
make build-all
```

Para ejecutar `staticcheck` junto con el resto de controles:

```sh
make quality
```

Los procesos validan configuración, inicializan logging estructurado y esperan
una señal de terminación. Ejemplo:

```sh
REGISTRY_INSTANCE_ID=registry-1 \
REGISTRY_CLUSTER_MEMBERS=registry-1:7000,registry-2:7000 \
make run SERVICE=registry
```
