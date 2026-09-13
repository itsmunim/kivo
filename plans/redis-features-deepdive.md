# Redis Features Deep Dive

This document explains every Redis feature that kivo supports, in plain English with real-world examples.

---

## Table of Contents

1. [Strings](#strings)
2. [Lists](#lists)
3. [Sets](#sets)
4. [Hashes](#hashes)
5. [Sorted Sets (ZSets)](#sorted-sets-zsets)
6. [Keys & Key Management](#keys--key-management)
7. [Expiration / TTL](#expiration--ttl)
8. [Connection Commands](#connection-commands)
9. [What kivo v1 Has vs. Full Redis](#what-kivo-v1-has-vs-full-redis)

---

## Strings

**What it is:** The simplest type. A key maps to a string value. That's it.

**What you use it for:**
- **Caching:** `SET user:123:profile '{"name":"Alice"}'` then `GET user:123:profile`
- **Counters:** `INCR page_views:homepage` — atomic increment, no race conditions even with 1000 clients hitting it simultaneously.
- **Rate limiting:** `INCR rate_limit:ip:1.2.3.4` + `EXPIRE rate_limit:ip:1.2.3.4 60` — count requests per minute per IP.
- **Session tokens:** `SET session:abc123 "user_id=456" EX 3600` — auto-expires in 1 hour.

**Kivo supports:** `GET`, `SET` (with `EX`/`PX` TTL), `MGET`, `MSET`, `INCR`, `DECR`, `INCRBY`, `DECRBY`, `APPEND`, `STRLEN`, `DEL`, `EXISTS`

**Why `MGET`/`MSET` matter:** Instead of 10 round trips for 10 keys, you do 1. When your app renders a page needing 20 cached fragments, `MGET` cuts latency dramatically.

---

## Lists

**What it is:** An ordered collection of strings. Think of it as an array that you can push to the front or back, and pop from either end.

**What you use it for:**
- **Task queues:** `LPUSH queue:emails "send_welcome_email"` then a worker does `RPOP queue:emails` to process oldest first (FIFO).
- **Activity feeds:** `LPUSH feed:user:123 "Alice liked your post"` — prepend new items, `LRANGE feed:user:123 0 9` to show the latest 10.
- **Bounded history:** `LPUSH search:user:123 "golang tutorials"` then `LTRIM search:user:123 0 99` — keeps only the last 100 searches.

**Kivo supports:** `LPUSH`, `RPUSH`, `LPOP`, `RPOP`, `LRANGE`, `LLEN`, `LINDEX`, `LREM`, `LTRIM`

**Important tradeoff in kivo:** `LPUSH` is O(n) because Go slices prepend by allocating a new array and copying everything. Redis uses a doubly-linked list, so it's O(1). For small lists (hundreds of items), you won't notice. For millions of items, it matters.

---

## Sets

**What it is:** An unordered collection of unique strings. No duplicates, no order.

**What you use it for:**
- **Tag systems:** `SADD post:123:tags "golang"` `SADD post:123:tags "redis"` — a post has multiple tags, no duplicates.
- **Follower lists:** `SADD followers:user:123 "user:456"` — check if user B follows user A with `SISMEMBER`.
- **IP blacklists:** `SADD blacklist:ips "1.2.3.4"` then `SISMEMBER blacklist:ips "1.2.3.4"` for O(1) lookup.
- **Set operations:**
  - `SINTER tags:golang tags:redis` — find posts tagged with BOTH.
  - `SUNION tags:golang tags:python` — find posts tagged with EITHER.
  - `SDIFF tags:golang tags:python` — find posts tagged golang but NOT python.

**Kivo supports:** `SADD`, `SREM`, `SMEMBERS`, `SISMEMBER`, `SCARD`, `SPOP`, `SRANDMEMBER`, `SUNION`, `SINTER`, `SDIFF`

---

## Hashes

**What it is:** A map of field-value pairs under a single key. Like a JSON object or a database row.

**What you use it for:**
- **User profiles:**
  ```
  HSET user:123 name "Alice" email "alice@example.com" age "30"
  HGET user:123 name        → "Alice"
  HGETALL user:123          → {"name":"Alice", "email":"alice@example.com", "age":"30"}
  ```
- **Shopping carts:**
  ```
  HSET cart:user:123 item:456 "2"   // 2 units of item 456
  HINCRBY cart:user:123 item:456 1  // add 1 more
  ```
- **Configuration:** `HSET config:app max_connections "100"` — fetch individual settings without loading the whole config object.

**Kivo supports:** `HSET`, `HGET`, `HGETALL`, `HDEL`, `HLEN`, `HEXISTS`, `HKEYS`, `HVALS`, `HMGET`

---

## Sorted Sets (ZSets)

**What it is:** A set where every member has a **score** (a float64). Members are kept sorted by score. If two members have the same score, they're sorted lexicographically by member name.

**Think of it as:** A leaderboard. Or a priority queue. Or a time-series index.

**What you use it for:**

**Leaderboards:**
```
ZADD leaderboard 100 "alice" 200 "bob" 150 "carol"
ZRANGE leaderboard 0 -1 WITHSCORES
→ alice (100), carol (150), bob (200)

ZREVRANGE leaderboard 0 2 WITHSCORES
→ bob (200), carol (150), alice (100)  // top 3
```

**Time-series data:**
```
ZADD events:clicks 1699900000 "user:123:clicked:buy"
ZADD events:clicks 1699900010 "user:456:clicked:view"
ZRANGEBYSCORE events:clicks 1699900000 1699900020
→ all clicks in that 20-second window
```

**Priority queues:**
```
ZADD queue:priority 1 "send_email" 5 "process_payment" 10 "alert_admin"
ZRANGE queue:priority 0 0  // lowest score = highest priority
→ "send_email"
```

**Geospatial (indirectly):** Redis proper uses sorted sets for `GEOADD`/`GEORADIUS` by encoding lat/lon into a geohash score. We don't support geospatial commands in v1, but the underlying structure is the same.

**Kivo supports:** `ZADD`, `ZREM`, `ZCARD`, `ZSCORE`, `ZINCRBY`, `ZCOUNT`, `ZRANGE`, `ZREVRANGE`, `ZRANGEBYSCORE`, `ZREVRANGEBYSCORE`, `ZREMRANGEBYSCORE`, `ZREMRANGEBYRANK`

**How kivo implements it:** `map[string]float64` for O(1) score lookup + a sorted slice for range queries. On every `ZRANGE`, we extract all members, sort by `(score, member)`, then return the slice. This is O(n log n) per range query. Redis uses a **skip list** for O(log n). Our approach is simpler and correct, but won't scale to millions of members in a single sorted set.

---

## Keys & Key Management

These aren't data types — they're operations on keys themselves.

| Command | What it does | Example |
|---------|-------------|---------|
| `DEL key [key ...]` | Delete keys | `DEL session:old123` |
| `EXISTS key [key ...]` | Count how many keys exist | `EXISTS user:123 user:456` |
| `KEYS pattern` | Find all keys matching a glob pattern | `KEYS user:*` → all user keys |
| `FLUSHDB` | Delete ALL keys in the database | ⚠️ Nuclear option |
| `DBSIZE` | Count total keys | |
| `TYPE key` | Tell you what type a key holds | `TYPE user:123` → "hash" |
| `RENAME old new` | Rename a key | |
| `RENAMENX old new` | Rename only if new key doesn't exist | |

**Warning about `KEYS`:** It scans the entire key space. In production with millions of keys, this blocks the server. Redis added `SCAN` for incremental iteration. We don't have `SCAN` in v1. Use `KEYS` sparingly.

---

## Expiration / TTL

**What it is:** Every key can have an expiration time. After that time, the key disappears.

**How it works in kivo:**
- **Lazy expiration:** When you `GET` a key, we check if it's expired. If yes, delete it and return "not found."
- **Active expiration:** A background goroutine runs every 100ms, samples 20 random keys, deletes expired ones. If >25% were expired, it repeats (up to 20 cycles). This is the exact same algorithm Redis uses.

**What you use it for:**
- **Session expiry:** `SET session:abc "user:123" EX 3600` — gone after 1 hour.
- **Cache invalidation:** `SET cache:product:456 "..." EX 300` — stale after 5 minutes.
- **Temporary locks:** `SET lock:resource:1 "owner:123" EX 10` — auto-release if the owner crashes.

**Kivo supports:** `EXPIRE` (seconds), `PEXPIRE` (milliseconds), `TTL` (seconds remaining), `PTTL` (milliseconds remaining), `PERSIST` (remove expiration)

---

## Connection Commands

| Command | What it does |
|---------|-------------|
| `PING` | Health check. Returns `PONG`. |
| `ECHO msg` | Echoes back the message. |
| `QUIT` | Closes the connection gracefully. |
| `SELECT index` | In Redis, switches databases (0-15). In kivo, it's a no-op — we only support one database. |
| `AUTH password` | In Redis, authenticates. In kivo v1, always returns `OK`. |

---

## What kivo v1 Has vs. Full Redis

| Feature | kivo v1 | Full Redis |
|---------|---------|-----------|
| Strings | ✅ Full | ✅ Full |
| Lists | ✅ Core | + `BLPOP`, `BRPOP`, `RPOPLPUSH`, etc. |
| Sets | ✅ Core | + `SMOVE`, `SSCAN`, etc. |
| Hashes | ✅ Core | + `HINCRBY`, `HSCAN`, etc. |
| Sorted Sets | ✅ Core | + `ZREVRANK`, `ZRANK`, `ZLEXRANK`, `ZSCAN`, etc. |
| Pub/Sub | ❌ | ✅ |
| Transactions | ❌ | `MULTI`/`EXEC`/`WATCH` |
| Lua scripting | ❌ | `EVAL`/`EVALSHA` |
| Streams | ❌ | Redis 5.0+ |
| Replication | ❌ | Master/slave |
| Clustering | ❌ | Redis Cluster |
| RDB snapshots | ❌ | Binary save/load |
| AOF rewrite | ❌ | Background compaction |

---

## See Also

- [v1-design.md](v1-design.md) — Architecture and implementation decisions
- [README.md](../README.md) — Quick start and feature list
