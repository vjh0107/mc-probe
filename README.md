# mc-probe

A small static binary to embed in Minecraft server container images,
intended as Kubernetes liveness, readiness, and startup probes.

## Probes

| Probe | Check                                                                     |
|-------|---------------------------------------------------------------------------|
| `startup` | Ping status returns `players.max > 0` (server has finished bootstrapping) |
| `liveness` | TCP dial only — listening socket is open (avoids false positives)         |
| `readiness` | `startup` + `players.online < players.max`                                |

## Args

| Arg         | Default       | Description    |
|-------------|---------------|----------------|
| `--host`    | `127.0.0.1`   | Server host.   |
| `--port`    | `25565`       | Server port.   |
| `--timeout` | `1s`          | Ping timeout.  |

## K8s Usage

```yaml
livenessProbe:
  exec:
    command: ["mcprobe", "liveness"]
  periodSeconds: 10
  failureThreshold: 3

readinessProbe:
  exec:
    command: ["mcprobe", "readiness"]
  periodSeconds: 2
  failureThreshold: 1

startupProbe:
  exec:
    command: ["mcprobe", "startup"]
  initialDelaySeconds: 30
  periodSeconds: 3
  failureThreshold: 100
```

#### `readinessProbe.failureThreshold: 1` is intentional. 
A failure only removes the
pod from service endpoints (no restart), so flapping is cheap and load balancers
can route around full instances immediately.
