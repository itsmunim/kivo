# kivo

A fast, lightweight, Redis-compatible in-memory data store written in Go.

> **Status:** Early development. v1 targets single-node, single-threaded (async) operation with core Redis commands, AOF persistence, and basic expiration.

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

## Architecture

```
cmd/kivo/          # Server entry point
internal/
  resp/            # RESP2 protocol parser and encoder
  store/           # In-memory data engine
  commands/        # Command handlers (strings, lists, sets, hashes, sorted sets, keys, pubsub)
  persistence/     # AOF append-only file
  server/          # TCP server and connection handling
  config/          # Configuration
pkg/api/           # Public API types (if any)
plans/             # Implementation plans and feature specs
docs/              # Documentation
```

## License

MIT