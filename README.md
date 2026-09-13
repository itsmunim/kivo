# kivo

![Coverage](coverage-badge.svg)

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

## Benchmarks

See [BENCHMARKS.md](BENCHMARKS.md) for performance results.

## Docker

```bash
docker build -t kivo .
docker run -p 6379:6379 kivo
```

## License

MIT
