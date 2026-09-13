# kivo v1 Design Document

> A fast, lightweight, Redis-compatible in-memory data store written in Go.

This document explains every architectural decision in kivo v1. Read this before hacking on the codebase.

---

## Table of Contents

1. [Project Goals](#project-goals)
2. [Why Go, Not Rust](#why-go-not-rust)
3. [Feature List](#feature-list)
4. [Architecture Overview](#architecture-overview)
5. [RESP2 Protocol](#resp2-protocol)
6. [Data Structures](#data-structures)
7. [Expiration](#expiration)
8. [AOF Persistence](#aof-persistence)
9. [Memory Management](#memory-management)
10. [Concurrency Model](#concurrency-model)
11. [Command Dispatch](#command-dispatch)
12. [Graceful Shutdown](#graceful-shutdown)
13. [Testing Strategy](#testing-strategy)
14. [What's Next](#whats-next)

---

## Project Goals

- **Redis protocol compatible** — works with existing Redis clients (redis-cli, go-redis, ioredis, etc.)
- **Simple and hackable** — small codebase, easy to understand and extend
- **Reliable** — AOF persistence, graceful shutdown, predictable behavior
- **Foundation for distributed** — single-node now, clustering and sharding later

---

## Why Go, Not Rust

We considered Rust, Go, Zig, C++, and Java. Here's why Go won for kivo v1.

### The core conflict: shared mutable state

A Redis server is essentially **one giant shared mutable state machine** — a single `map[string]Item` that every client connection reads and writes concurrently.

**In Go:**
```go
type Engine struct {
    mu   sync.RWMutex
    data map[string]Item
}
```
Simple, obvious, correct. Lock, mutate, unlock.

**In Rust:**
```rust
Arc<RwLock<HashMap<String, Item>>>  // coarse, contention
// or
DashMap<String, Item>               // external crate, less control
// or
tokio::sync::mpsc::channel + task   // architecturally different from Redis
```
Rust's ownership model fights the architecture. You either accept coarse locking, pull in external crates, or redesign the whole concurrency model. That's valuable learning, but it's **friction** when the goal is understanding distributed systems, not memory ownership patterns.

### Goroutines vs async/await

Redis uses a single event loop per core. Go uses lightweight threads per connection.

- **Go:** One goroutine per TCP connection. Parse command, send to store goroutine (or acquire mutex), return result. No lifetime annotations across await points. No `Pin<Box<dyn Future>>` when building pub/sub broadcast.
- **Rust:** Async/await with Tokio. Futures must be `Send` and `'static`. Shared state across tasks requires `Arc<Mutex<T>>` or message passing. The compiler becomes an active participant in design decisions.

For v1, Go's model is more forgiving and maps naturally to the problem.

### Distributed systems ecosystem

Go dominates production distributed systems. Look at what's written in Go:
- **etcd** — distributed KV with Raft consensus
- **Consul** — service discovery
- **NATS** — messaging
- **CockroachDB** — distributed SQL
- **TiKV** — distributed transactional KV

The Go ecosystem has battle-tested Raft libraries, gossip implementations, and well-documented patterns. Rust has them too, but they're harder to use and less battle-tested at scale.

### Prototyping speed

A working `GET`/`SET` server with RESP parsing is an afternoon in Go. In Rust, you might spend a week getting the parser to compile with proper lifetime management across TCP read buffers.

v1's goal is **shipping and learning**, not proving language proficiency.

### GC is a non-issue for this use case

Modern Go's GC pauses are sub-millisecond. For a cache/session store serving <1M ops/sec, you will never notice. If we later hyperscale, we optimize with `sync.Pool`. But "will GC pause hurt my cache?" is premature optimization at this stage.

### When Rust would be better

- **Embedded in a larger Rust system** — zero-overhead FFI
- **Maximum single-node throughput** — Rust squeezes more ops/sec, but Go is already "fast enough" (100K-500K ops/sec easily)
- **Learning memory safety discipline** — if the *primary* goal were learning Rust, not building a database

### Bottom line

Go gets us to the interesting problems (consensus, partitioning, replication) faster. Rust is a fine choice for a v2 rewrite once the architecture is understood and performance is the constraint, not clarity.

**Recommendation:** Build in Go. Understand the domain. Revisit Rust later if embedding or raw performance becomes the bottleneck.

---

---

## Feature List

### Protocol
- [x] RESP2 parser/encoder (all 5 types)
- [x] Array-based commands
- [x] Pipelining support
- [x] Null bulk strings and null arrays

### Strings
- [x] `GET`, `SET` (with `EX`/`PX` TTL)
- [x] `DEL`, `EXISTS`
- [x] `EXPIRE`, `PEXPIRE`, `TTL`, `PTTL`, `PERSIST`
- [x] `MGET`, `MSET`
- [x] `INCR`, `DECR`, `INCRBY`, `DECRBY`
- [x] `APPEND`, `STRLEN`

### Lists
- [x] `LPUSH`, `RPUSH`
- [x] `LPOP`, `RPOP`
- [x] `LRANGE`, `LLEN`, `LINDEX`
- [x] `LREM`, `LTRIM`

### Sets
- [x] `SADD`, `SREM`, `SMEMBERS`, `SISMEMBER`, `SCARD`
- [x] `SPOP`, `SRANDMEMBER`
- [x] `SUNION`, `SINTER`, `SDIFF`

### Hashes
- [x] `HSET`, `HGET`, `HGETALL`, `HDEL`
- [x] `HLEN`, `HEXISTS`, `HKEYS`, `HVALS`, `HMGET`

### Sorted Sets
- [x] `ZADD`, `ZREM`, `ZCARD`, `ZSCORE`, `ZINCRBY`, `ZCOUNT`
- [x] `ZRANGE`, `ZREVRANGE` (with `WITHSCORES`)
- [x] `ZRANGEBYSCORE`, `ZREVRANGEBYSCORE` (with `LIMIT`)
- [x] `ZREMRANGEBYSCORE`, `ZREMRANGEBYRANK`

### Keys
- [x] `KEYS` (with glob pattern `*` and `?`)
- [x] `FLUSHDB`, `DBSIZE`
- [x] `TYPE`, `RENAME`, `RENAMENX`

### Connection
- [x] `PING`, `ECHO`, `QUIT`
- [x] `AUTH` (hardcoded OK for v1)
- [x] `SELECT` (no-op, single DB)

### Server Features
- [x] TCP server with one goroutine per connection
- [x] AOF persistence (`always`, `everysec`, `no`)
- [x] Expiration (lazy + active)
- [x] Graceful shutdown on SIGTERM/SIGINT
- [x] `maxmemory` enforcement (no-eviction policy)

### Not in v1
- [ ] Pub/Sub (`SUBSCRIBE`, `PUBLISH`)
- [ ] Transactions (`MULTI`, `EXEC`)
- [ ] Lua scripting
- [ ] Replication / clustering
- [ ] RDB snapshots
- [ ] AOF rewrite
- [ ] LRU/LFU eviction policies

---

## Architecture Overview

```
cmd/kivo/main.go          → Entry point, wires everything together
internal/
  resp/resp.go            → RESP2 protocol parser/encoder
  store/store.go          → In-memory data engine (all types + expiration)
  commands/commands.go    → Command handlers (Redis command → store call)
  persistence/aof.go      → Append-only file persistence
  server/server.go        → TCP server, one goroutine per connection
  config/config.go        → Configuration struct
```

**Request flow:**

```
TCP connection → RESP reader → command lookup → acquire engine lock
                                                    ↓
RESP writer ← serialize result ← execute command ← store mutation
                                                    ↓
                                              AOF append (if write)
```

---

## RESP2 Protocol

Redis speaks RESP2 (REdis Serialization Protocol). It's text-based and simple:

| Type | Format | Example |
|------|--------|---------|
| Simple string | `+OK\r\n` | `+PONG\r\n` |
| Error | `-ERR ...\r\n` | `-ERR unknown command\r\n` |
| Integer | `:123\r\n` | `:42\r\n` |
| Bulk string | `$5\r\nhello\r\n` | `$0\r\n\r\n` (empty), `$-1\r\n` (null) |
| Array | `*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n` | Commands are arrays of bulk strings |

**Why RESP2 not RESP3?** RESP3 adds more types (doubles, booleans, maps, etc.) but breaks compatibility with older clients. RESP2 is the safe common denominator. Every Redis client library (go-redis, ioredis, redis-py) speaks it.

**Key design decisions:**
- **Array-based commands only.** No inline commands (space-separated for `telnet`). `redis-cli` uses arrays by default.
- **Pipelining supported natively.** The reader doesn't block waiting for a response. It parses values sequentially from the TCP stream.

**The parser** (`resp.ReadValue`) reads exactly N bytes + `\r\n` for bulk strings. This matters because TCP is a stream — a `SET` command might arrive fragmented across multiple packets. The parser must handle partial reads correctly.

---

## Data Structures

### 4.1 Strings

**Structure:** `string` (Go native string).

**Why:** Strings are immutable byte sequences. `Get`/`Set` are O(1). `Append` creates a new string (Go strings are immutable, so `item.Value.(string) + value` allocates a new string). `INCR`/`DECR` parse the string as int64, increment, format back to string.

**Redis uses:** SDS (Simple Dynamic String), a length-prefixed mutable string buffer. We use Go strings because they're simpler and the performance difference is negligible for v1.

**Code:** `store.go` lines ~170-240.

### 4.2 Lists

**Structure:** `[]string` (Go slice).

**Why:** Redis lists are doubly-linked lists (for O(1) push/pop at both ends and O(1) insert at arbitrary positions). Go slices give us O(1) append (`RPUSH`) but O(n) prepend (`LPUSH`) because we allocate a new slice and copy. For v1 workloads (queues, small lists), this is fine.

**Key methods:**
- `LPush`: Prepend via `append([]string{v}, list...)`
- `RPush`: `append(list, values...)`
- `LRange`: Slice operation with index normalization
- `LRem`: Iterate and filter

**Index normalization:** Redis supports negative indices (-1 = last element). `normalizeRange` and `normalizeIndex` handle this.

**Code:** `store.go` lines ~360-530.

### 4.3 Sets

**Structure:** `map[string]struct{}`.

**Why:** O(1) membership testing, O(1) insert/delete. The `struct{}` value takes zero bytes. This is the idiomatic Go set.

**Set operations:**
- `SUnion`/`SInter`/`SDiff`: Iterate over maps, build result maps, then convert to sorted slices for deterministic output.

**Code:** `store.go` lines ~540-750.

### 4.4 Hashes

**Structure:** `map[string]string`.

**Why:** Each hash key stores field→value pairs. O(1) field access.

**Code:** `store.go` lines ~760-890.

### 4.5 Sorted Sets (the most complex)

**Structure:** `map[string]float64` (member → score) + **sorted slice for range queries**.

**Why Redis uses skip lists:** O(log n) insert, delete, range query. A skip list is a probabilistic data structure with multiple levels of linked lists.

**Why we use map + sorted slice:**
- Simpler to implement correctly.
- For v1 workloads (thousands to low millions of members), sorting on query is acceptable.
- **Tradeoff:** `ZRange` is O(n log n) due to sort, vs O(log n) for skip list.

**How it works:**
1. `ZAdd`: Insert into map. No sorting yet.
2. `ZRange`/`ZRevRange`/`ZRangeByScore`: Call `sortedZSetMembers()` which extracts all members, sorts by `(score, member)`, then returns the requested range.
3. `ZRem`/`ZRemRangeByScore`/`ZRemRangeByRank`: Delete from map, then if needed, re-sort.

**Code:** `store.go` lines ~900-1150.

**Future:** Replace with a skip list or treap for O(log n) operations.

---

## Expiration

Two mechanisms, both required:

### 5.1 Lazy expiration
On every read (`Get`, `HGet`, `LIndex`, etc.), check if `ExpiresAt` is in the past. If so, delete the key and return "not found".

**Why:** Prevents returning stale data.

### 5.2 Active expiration
A background goroutine runs every 100ms. It implements Redis's probabilistic algorithm:
1. Sample 20 random keys.
2. Delete expired ones.
3. If >25% of sampled keys were expired, repeat (up to 20 cycles).

**Why both?** Lazy handles keys that are read. Active handles keys that are written once and never read again (prevents memory leaks).

**Code:** `activeExpiration()` and `expireSample()` in `store.go`.

---

## AOF Persistence

### How it works
- Every write command is appended to a file as a **RESP array**.
- Format: `*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n`
- This is **text-based** (human-readable-ish), not binary.

### Sync strategies
- `always`: `fsync` after every write. Safest, slowest.
- `everysec` (default): Buffer writes, flush + fsync every second. Good balance.
- `no`: Let OS decide when to flush. Fastest, riskiest.

### Is there a background processor?
No. Writes happen synchronously in the request handler, then the AOF writer appends. For `everysec`, a background goroutine does periodic `fsync`.

### Is it performant?
For v1, yes. The file is append-only, so writes are sequential disk I/O (fast). But:
- No compression. The AOF grows forever.
- No rewrite. If you set the same key 1000 times, the AOF has 1000 `SET` commands. On replay, it executes all 1000.

**Redis solves this with AOF rewrite:** A background process reads the current in-memory state and writes a minimal AOF (only the final values). We explicitly deferred this to v1.1.

**SSTable / LSM tree?** No, and we shouldn't use them for AOF. SSTables and LSM trees are for **random-read heavy, write-optimized key-value storage** (like RocksDB, LevelDB). AOF is a **write-ahead log** — fundamentally an append-only stream. They're different problem domains.

**Where AOF should evolve:**
1. **AOF rewrite** (v1.1) — compact the file by dumping current state.
2. **RDB snapshots** (v1.2) — periodic binary snapshots for fast restarts.
3. **Copy-on-write fork** (much later) — for non-blocking rewrite, like Redis does.

**Code:** `internal/persistence/aof.go`.

---

## Memory Management

### maxmemory

v1 implements a simple `maxmemory` cap with a **no-eviction policy**:
- If `maxmemory` is set and a write would exceed it, the write is rejected with an OOM error.
- The process still starts even if `maxmemory` exceeds available system memory (same as Redis).
- A warning is logged at startup if `maxmemory` > available memory.

**Why no-eviction?** LRU/LFU eviction requires tracking access recency, which adds significant complexity. For v1, we reject writes and let the operator decide.

**Memory estimation:** Approximate by summing key lengths + value sizes. This is not exact (Go runtime overhead, map buckets, etc.) but it's close enough for a safety cap.

**Redis behavior reference:**
- Redis logs a warning if `maxmemory` > available memory but still starts.
- Redis supports `maxmemory-policy` (allkeys-lru, volatile-lru, etc.). We defer this.

**Code:** `store.go` `checkMemory()` and `memory.go`.

### Platform-specific memory detection

At startup, kivo checks whether `maxmemory` exceeds available system memory and logs a warning if so. This detection is **platform-specific:**

**Linux:** Uses `syscall.Sysinfo` to read `Totalram` from the kernel. This is the same approach Redis uses internally.

**Other platforms (macOS, Windows, BSD):** Returns 0 with no error, effectively skipping the check. The process starts normally.

**Why a no-op on non-Linux?**
- macOS has no `sysinfo` syscall. Detecting total memory requires `sysctl` or cgo, adding complexity.
- Windows uses a completely different API.
- The check is a **best-effort warning**, not a hard requirement. Redis itself doesn't refuse to start when it can't determine available memory.

**Build tags used:**
```
internal/store/memory_linux.go  // go:build linux
internal/store/memory_other.go  // go:build !linux
```

This keeps the code clean and avoids conditional compilation complexity in the main store logic.

**Code:** `internal/store/memory_linux.go`, `internal/store/memory_other.go`.

---

## Concurrency Model

**One goroutine per TCP connection.**
- Each goroutine: read RESP → lookup command → acquire engine lock → execute → write RESP → release lock.
- The engine uses a **global `sync.RWMutex`**.

**Why not a single event loop like Redis?** Go's scheduler is excellent. One goroutine per connection is idiomatic Go, simpler, and performant enough for v1.

**Lock granularity:** A single global mutex means all commands are serialized. This is the simplest correct approach. For v1.1, we can consider sharding the lock by key hash (like a striped lock) to allow concurrent operations on different keys.

**Code:** `server.go` and `store.go`.

---

## Command Dispatch

Commands are registered in a map: `map[string]Handler`.
- Key: uppercase command name.
- Value: function `func(e *store.Engine, args []resp.Value) resp.Value`.

The server parses the RESP array, extracts the command name, looks it up in the registry, and calls the handler.

**Why a registry instead of a switch?** Easier to test, easier to extend, and makes command plugins possible in the future.

**Code:** `commands.go`.

---

## Graceful Shutdown

1. Stop accepting new connections (`listener.Close()`).
2. Wait for in-flight handlers to finish (`sync.WaitGroup`).
3. Flush AOF buffer.
4. Close all connections.

**Code:** `server.go` `Stop()` and `main.go` signal handling.

---

## Testing Strategy

- **Unit tests** for RESP parser (37 tests, 4 benchmarks).
- **Unit tests** for store engine (46 tests covering all types, expiration, concurrency).
- **Integration tests** for server (14 tests over real TCP).
- **AOF tests** for write/read/replay round-trips.

**Benchmarking:** Use `redis-benchmark` for end-to-end performance validation.

---

## What's Next

| Feature | Target |
|---------|--------|
| Pub/Sub | v1.1 |
| AOF rewrite | v1.1 |
| Password auth | v1.1 |
| Config file / CLI flags | v1.1 |
| RDB snapshots | v1.2 |
| LRU/LFU eviction | v1.2 |
| Replication | v2.0 |
| Clustering / sharding | v2.0 |
