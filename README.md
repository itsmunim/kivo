# kivo

![Coverage](https://img.shields.io/badge/coverage-82.5%25-brightgreen)

A fast, lightweight, Redis-compatible in-memory data store written in Go.

> **Status:** v1 feature-complete. Single-node, reliable, Redis-protocol compatible.
> See [plans/v1-design.md](plans/v1-design.md) for full architectural walkthrough.

## Platform Support

kivo runs on all platforms. The only platform-specific behavior is a **startup warning** when `maxmemory` is configured:

- **Linux** — warns if `maxmemory` exceeds available system memory (via `syscall.Sysinfo`)
- **macOS** — warns if `maxmemory` exceeds available system memory (via `syscall.Sysctl`)
- **Windows / BSD** — no warning; server starts normally. If `maxmemory` exceeds physical RAM, the OS will handle it (OOM kill or swap).

All platforms support full kivo functionality.

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

### Not in v1
- Pub/Sub, Transactions, Lua scripting, Replication, Clustering, RDB snapshots, AOF rewrite

## Architecture

```
cmd/kivo/          # Server entry point
internal/
  resp/            # RESP2 protocol parser and encoder
  store/           # In-memory data engine
  commands/        # Command handlers
  persistence/     # AOF append-only file
  server/          # TCP server and connection handling
  config/          # Configuration
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

## Web Console

kivo ships with an embedded browser console for browsing keys and running Redis commands live (how cool is that?!). It starts automatically on its own HTTP port:

```bash
go run ./cmd/kivo            # then open http://localhost:3001
```

- **Keys panel** — lists all keys with type badges and TTL; click a key to `GET` it
- **Terminal** — run any supported Redis command; results render RESP-style (`+OK`, `-ERR`, `:1`, `"value"`, arrays)
- **History** — press ↑/↓ to recall previous commands; auto-refreshes every 5s (or hit refresh on top-right)

Disable it with `-webui=false`, or change the port with `-webui-addr=:8080`. 

**Note:** In production, keep it bound to localhost or shield it behind auth — the console has no auth of its own.

## Benchmarks

See [BENCHMARKS.md](BENCHMARKS.md) for performance results.

## Docker

### Pull from GitHub Container Registry

```bash
docker pull ghcr.io/itsmunim/kivo:latest
docker run -p 6379:6379 ghcr.io/itsmunim/kivo:latest
```

### Build from source

```bash
docker build -t kivo .
docker run -p 6379:6379 kivo
```

## Kubernetes

A ready-to-apply manifest lives at [`k8s/kivo-deploy.yaml`](k8s/kivo-deploy.yaml). It runs kivo as a 1-replica `StatefulSet` with an attached PVC, so the AOF survives pod restarts.

```bash
kubectl apply -f https://raw.githubusercontent.com/itsmunim/kivo/main/k8s/kivo-deploy.yaml
```

Connect from inside the cluster:

```bash
kubectl run -it --rm redis-client --image=redis --restart=Never -- redis-cli -h kivo PING
```

> **Storage note:** the manifest uses your cluster's default StorageClass, which on many clusters is node-local (hostPath). Node-local storage is fine for testing, but data is lost if the node is deleted or crashes. For durable persistence, download `k8s/kivo-deploy.yaml`, set `storageClassName` under `volumeClaimTemplates` to a detachable/network class (AWS EBS `gp2`, GCE PD `standard`, Azure Disk `managed-csi`, Longhorn, Rook...), then apply your copy.

## GitHub Pages

Visit the project landing page at: **https://itsmunim.github.io/kivo/**

## License

MIT
