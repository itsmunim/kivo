# `kivo`

![Coverage](https://img.shields.io/badge/coverage-82.5%25-brightgreen)

A fast, lightweight, Redis-compatible in-memory data store written in Go.

> **Status:** v1 feature-complete. Single-node, reliable, Redis-protocol compatible.
> See [plans/v1-design.md](plans/v1-design.md) for full architectural walkthrough.

## When to use `kivo`

`kivo` is a single-node, Redis-protocol-compatible store. If your use case looks like one of these, feel free to use it for MVPs or production workloads:

- **App cache** — cache expensive computations, API responses, or rendered fragments
- **Session store** — sessions with TTL-backed expiry
- **Rate limiting** — atomic counters (`INCR` + `EXPIRE`)
- **Task queues** — `LPUSH`/`RPOP` FIFO workers
- **Leaderboards & ranked data** — sorted sets
- **Prototypes and MVPs** — get something reliable running in minutes, scale later

See the [benchmarks](#benchmarks) below: measured on identical hardware, `kivo` runs on par with Redis standalone and ahead on several operations (including 28–43% faster pipelined throughput) — so a single node is not a performance compromise.

> ## The story behind `kivo`
>
> In this era of AI it's not about building new things whenever you get a chance to learn and explore — the part that is even harder is to keep concentrating on a project you started on one fine weekend evening. You can try so many things, to learn and to build, that the urge to keep switching from one to another is real — so is the burden of context switching in your head.
>
> I started a few projects (a simple timeseries db, an AI task planner, etc.) and they now all sit in incomplete states. It's not easy to manage time on weekends with three little ones, while also doing groceries and taking them out — you know, all the fam stuff. Hence, I am really glad I could finish this one up.
>
> I did the initial implementation as a couple of files, no `git init` or anything — just to try out whether it works. I added a bunch of things from time to time whenever I felt like it, but wasn't really doing it with the intention of setting up git and GitHub. Then one weekend I started polishing things up, to seriously try this single-node setup as a cache/storage, and to see if others can check it out too with a first fat commit.
>
> I just wanted to try to build a Redis-like database in Go from scratch. But I also wanted to make it a reliable one, and simple, with a single-node AOF-only approach.

## Platform Support

`kivo` runs on all platforms. The only platform-specific behavior is a **startup warning** when `maxmemory` is configured:

- **Linux** — warns if `maxmemory` exceeds available system memory (via `syscall.Sysinfo`)
- **macOS** — warns if `maxmemory` exceeds available system memory (via `syscall.Sysctl`)
- **Windows / BSD** — no warning; server starts normally. If `maxmemory` exceeds physical RAM, the OS will handle it (OOM kill or swap).

All platforms support full `kivo` functionality.

## Goals

- **Redis protocol compatible** — works with existing Redis clients (redis-cli, go-redis, ioredis, etc.)
- **Simple and hackable** — small codebase, easy to understand and extend
- **Reliable** — AOF persistence, graceful shutdown, predictable behavior
- **Foundation for distributed** — single-node now, clustering and sharding later

## Quick Start

```bash
# Build and run
go run ./cmd/kivo

# Connect with redis-cli
redis-cli -p 6379 PING
```

## Web Console (no redis-cli required)

`kivo` ships with an embedded browser console: browse every key with its type and TTL, and run any supported Redis command live in a terminal. No `redis-cli`, no client library, nothing to install — start `kivo` and it is there.

```bash
go run ./cmd/kivo        # then open http://localhost:3001
```

- **Keys panel** — lists all keys with type badges and TTL; click a key to `GET` it
- **Terminal** — run any supported Redis command; results render RESP-style (`+OK`, `-ERR`, `:1`, `"value"`, arrays)
- **History** — press ↑/↓ to recall previous commands; auto-refreshes every 5s

Try the guided command tour (a full test suite with expected outputs) at the **[kivo web console docs](https://itsmunim.github.io/kivo/web-console.html)**.

Disable it with `-webui=false`, or change the port with `-webui-addr=:8080`.

**Note:** In production, keep it bound to localhost or shield it behind auth — the console has no auth of its own. On Kubernetes, reach it privately with `kubectl port-forward` (see below).

## v1 Features

### Data Types
- **Strings** — `GET`, `SET`, `MGET`, `MSET`, `INCR`, `DECR`, `APPEND`, `STRLEN`
- **Lists** — `LPUSH`, `RPUSH`, `LPOP`, `RPOP`, `LRANGE`, `LLEN`, `LINDEX`, `LREM`, `LTRIM`
- **Sets** — `SADD`, `SREM`, `SMEMBERS`, `SISMEMBER`, `SCARD`, `SPOP`, `SRANDMEMBER`, `SUNION`, `SINTER`, `SDIFF`
- **Hashes** — `HSET`, `HGET`, `HGETALL`, `HDEL`, `HLEN`, `HEXISTS`, `HKEYS`, `HVALS`, `HMGET`
- **Sorted Sets** — `ZADD`, `ZREM`, `ZRANGE`, `ZREVRANGE`, `ZRANGEBYSCORE`, `ZREVRANGEBYSCORE`, `ZCARD`, `ZSCORE`, `ZINCRBY`, `ZCOUNT`, `ZREMRANGEBYSCORE`, `ZREMRANGEBYRANK`

### Keys & Expiration
- `KEYS`, `FLUSHDB`, `DBSIZE`, `TYPE`, `RENAME`, `RENAMENX`
- `EXPIRE`, `PEXPIRE`, `TTL`, `PTTL`, `PERSIST`
- Lazy + active expiration (probabilistic sampling)

### Server
- RESP2 protocol with pipelining
- TCP server (one goroutine per connection)
- AOF persistence (`always`, `everysec`, `no`)
- `maxmemory` enforcement (no-eviction policy, with platform-specific memory detection on Linux and macOS)
- Graceful shutdown on SIGTERM/SIGINT
- Embedded web console (`:3001`)

### Not in v1
- Pub/Sub, Transactions, Lua scripting, Replication, Clustering, RDB snapshots, AOF rewrite

## Benchmarks

Both `kivo` and a real Redis (redis:7-alpine) run as containers on the same Linux/arm64 Docker host, driven by identical `redis-benchmark` commands (300K requests, 50 parallel clients, ops `set,get,lpush,rpush,incr`); CPU and memory sampled via `docker stats` during the run. Average of 3 rounds. See [BENCHMARKS.md](BENCHMARKS.md) for the full progression story — naive baseline, the LPUSH/deque fix, and both macOS and Linux resource comparisons.

| op | kivo (ops/sec) | Redis (ops/sec) | gap |
|----|---------------|-----------------|-----|
| `GET` | 130,300 | 134,200 | −3% |
| `INCR` | 125,100 | 139,800 | −11% |
| `LPUSH` | 126,200 | 127,700 | −1% |
| `RPUSH` | 129,900 | 128,200 | **+1%** |
| `SET` | 116,700 | 129,200 | −10% |
| **Avg CPU** | **116.7%** | **54.2%** | 2.2x |
| **Peak memory** | **57.6 MiB** | **17.7 MiB** | 3.3x |

**TL;DR:** on the same Linux host, `kivo` runs within 1–11% of Redis throughput on all core ops (RPUSH actually ahead) — with the same caveat as any single-node store: it's also using ~2.2x the CPU and ~3.3x the memory today. `LPUSH` went from 4,196 → 126,200 ops/sec (the O(n²) slice-preprend → ring-buffer deque fix, [design notes](plans/perf-plan.md)).
## Architecture

```
cmd/kivo/          # Server entry point
internal/
  resp/            # RESP2 protocol parser and encoder
  store/           # In-memory data engine
  commands/        # Command handlers
  persistence/     # AOF append-only file
  server/          # TCP server and connection handling
  webui/           # Embedded web console (React SPA + JSON API)
  config/          # Configuration
webui/             # Web console source (Vite + React)
landing/           # Public site (Vite, static) — GitHub Pages
k8s/               # Kubernetes manifests
plans/             # Implementation plans and design docs
```

## Design Document

For detailed explanations of every architectural decision — why RESP2 over RESP3, why map+sorted-slice for sorted sets, how expiration works, AOF tradeoffs, and more — see **[plans/v1-design.md](plans/v1-design.md)**.

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
./scripts/coverage.sh

# Run benchmarks (requires redis-benchmark)
./scripts/benchmark.sh
```

## Docker

### Pull from GitHub Container Registry

```bash
docker pull ghcr.io/itsmunim/kivo:latest
docker run -p 6379:6379 -p 3001:3001 ghcr.io/itsmunim/kivo:latest
```

### Build from source

```bash
docker build -t kivo .
docker run -p 6379:6379 -p 3001:3001 kivo
```

## Kubernetes

A ready-to-apply manifest lives at [`k8s/kivo-deploy.yaml`](k8s/kivo-deploy.yaml). It runs `kivo` as a 1-replica `StatefulSet` with an attached PVC, so the AOF survives pod restarts.

```bash
kubectl apply -f https://raw.githubusercontent.com/itsmunim/kivo/main/k8s/kivo-deploy.yaml
```

Connect from inside the cluster:

```bash
kubectl run -it --rm redis-client --image=redis --restart=Never -- redis-cli -h kivo PING
```

**Reach the private web console locally** — the console has no auth and binds inside the pod, so port-forward it instead of exposing it:

```bash
kubectl port-forward pod/kivo-0 3001:3001
# then open http://localhost:3001
```

> **Storage note:** the manifest uses your cluster's default StorageClass, which on many clusters is node-local (hostPath). Node-local storage is fine for testing, but data is lost if the node is deleted or crashes. For durable persistence, download `k8s/kivo-deploy.yaml`, set `storageClassName` under `volumeClaimTemplates` to a detachable/network class (AWS EBS `gp2`, GCE PD `standard`, Azure Disk `managed-csi`, Longhorn, Rook...), then apply your copy.

## GitHub Pages

- Landing page: **https://itsmunim.github.io/kivo/**
- Web console guide: **https://itsmunim.github.io/kivo/web-console.html**

## License

MIT
