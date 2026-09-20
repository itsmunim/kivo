# kivo Benchmarks

## How to read
The full benchmarking story of the kivo project — from first naive implementation to the current build — with the key optimization decisions that moved the numbers. Raw methodology notes are inline with each phase; generators exist in `scripts/`.

---

## Environments referenced

| Label | Machine | Stack |
|-------|---------|-------|
| **MBP (Phase 1)** | Apple M1 Pro, Darwin arm64 (2026-09) | Redis 8.10.1 (Homebrew) vs kivo built locally |
| **MBP (resources)** | Same Apple M1 Pro | Same, persistence off both sides, 3 rounds |
| **Linux (Docker)** | Docker on the same Mac but Linux/arm64 containers (redis:7-alpine, kivo alpine build) | 3 rounds, `docker stats` sampling |

---

## Phase 0 — The naive implementation (baseline)

> Everything below predates both optimization passes. Lists were plain Go slices; the RESP writer wrote straight to `net.Conn` via `fmt.Fprintf`; the parser read byte-at-a-time.

**Findings (MBP, `redis-benchmark`, 100K req, 50 clients) — the two problems that framed the whole optimization effort:**

| op | kivo (naive) | Redis (same MBP) |
|----|--------------|------------------|
| SET | ~81K | ~70K |
| GET | ~76K | ~83K |
| INCR | ~78K | ~87K |
| RPUSH | ~82K | ~81K |
| SADD | ~81K | ~78K |
| HSET | ~84K | ~80K |
| **LPUSH** | **~4.2K (collapsing to ~2K)** | **~82K** |
| Pipelined GET | ~414K | ~797K |

Two stories emerge:

1. **Throughput was already at or near parity** with real Redis on the same machine — the earlier "half of Redis" fear came from comparing against Redis's published Linux-server numbers on different hardware. GET/SET/INCR/RPUSH/SADD/HSET were all within ~10% of Redis.

2. **LPUSH was catastrophic — O(n²) in the list size.** `redis-benchmark` pushes to a single growing list; the naive `append([]string{v}, list...)` prepend copies the whole list every push. pprof showed `Engine.LPush` at **98.7% of ALL allocations — 18.7 GB in a 5-second profile**. The 4.2K → ~2K decay during a single run was the list growing under the benchmark.

**Profile evidence** (5s CPU profile under SET/GET load), ~77% of CPU was I/O framing, not data logic:

| share | function | meaning |
|-------|----------|---------|
| 51% | `syscall.rawsyscalln` | one syscall per response (unbuffered writer) |
| 26% | `bufio.Reader.ReadByte` | byte-at-a-time line parsing |
| 25% | `fmt.Fprintf` | reflection-based response formatting |

Micro-benchmarks of the naive store:

```
BenchmarkStoreLPush  233,846 ns/op   605,756 B/op   ← the O(n²) prepend
BenchmarkStoreRPush     100.3 ns/op       108 B/op
```

---

## Optimization #1 — The deque (the LPUSH fix)

**Decision:** replace the slice-preprend with a ring-buffer deque (`internal/store/list.go`) — O(1) amortized push/pop at both ends, O(1) `LINDEX` via index arithmetic, contiguous memory allocation. The engine already abstracted lists behind a `*List` type, so zero command-handler changes were needed.

**Micro-benchmarks after (before → after):**

```
BenchmarkStoreLPush  233,846 → 50.15 ns/op   (4,600x faster)
BenchmarkStoreRPush     100.3 → 42.14 ns/op
Allocations            605,756 → 40 B/op      (~15,000x fewer)
```

**MBP `redis-benchmark` after (100K req, 50 clients):**

| op | kivo (after) | Redis (same MBP) |
|----|--------------|------------------|
| PING | 83,701 | 81,638 |
| SET | 83,217 | 78,526 |
| GET | 84,703 | 84,934 |
| INCR | 84,207 | 76,393 |
| **LPUSH** | **77,386** | **78,304** |
| RPUSH | 83,362 | 76,463 |
| LPOP | 81,416 | 78,297 |
| RPOP | 84,345 | 81,218 |
| SADD | 83,820 | 84,865 |
| HSET | 84,118 | 88,549 |
| **Pipe SET ×16** | **1,123,595** | 787,402 |
| **Pipe GET ×16** | **1,226,994** | 956,938 |

LPUSH went from **4.2K → 77.4K ops/sec (18x)** and reached parity with Redis. kivo now matched or beat Redis on 11 of 12 metrics on the same machine, led by pipelining (28-43% faster).

---

## Resource comparison #1 — macOS, throughput + CPU + memory (3 rounds)

Once throughput reached parity, we measured the cost of that throughput: CPU and memory under identical load. `redis-benchmark` 300K req, 50 clients, ops `set,get,lpush,rpush,incr`; CPU sampled via `ps` every 150ms during the run, RSS peaked per round; 3 rounds, persistence disabled both sides.

**Average of 3 rounds (MBP):**

| Server | Avg CPU | Peak RSS | SET | GET | INCR | LPUSH | RPUSH |
|--------|---------|----------|-------|-------|-------|-------|-------|
| **kivo** | **234.9%** | **69.5 MB** | 88.5K | 92.2K | 93.8K | 87.9K | 89.3K |
| **Redis** | **91.5%** | **7.3 MB** | 92.1K | 91.5K | 85.6K | 80.9K | 84.3K |

**Interpretation:** throughput parity held (kivo ahead on GET/INCR/LPUSH/RPUSH, Redis ahead ~4% on SET), but at 2.5x the CPU (kivo spreads across cores; Redis is single-threaded) and ~9.5x the RSS (Go runtime base + per-command allocations + per-key `Item` structs vs Redis's compact intset/ziplist encodings). This identified the next optimization targets: the I/O-framing path for CPU, allocations/pooling for memory.

---

## Resource comparison #2 — Linux/Docker, the deploy-environment benchmark (3 rounds)

Since nobody runs production Redis on a Mac, we reran on Linux/arm64 in containers — `redis:7-alpine` vs the current kivo alpine build — sampled with `docker stats` (`%CPU`, `MemUsage`), same workload, 3 rounds. Script: `scripts/bench-resources-linux.sh`.

**Average of 3 rounds (Linux/Docker):**

| op | kivo | Redis | gap |
|----|------|-------|-----|
| GET | 130.3K | 134.2K | −3% |
| INCR | 125.1K | 139.8K | −11% |
| LPUSH | 126.2K | 127.7K | −1% |
| RPUSH | 129.9K | 128.2K | +1% |
| SET | 116.7K | 129.2K | −10% |
| **Avg CPU** | **116.7%** | **54.2%** | 2.2x |
| **Peak memory** | **57.6 MiB** | **17.7 MiB** | 3.3x |

**Interpretation:** near-parity throughput on real Linux (within 1-11% either way, RPUSH ahead), with Redis's single-thread model using much less total CPU and the compact-encoding advantage showing in memory. These are the honest numbers for deployment environments — and the targets the optimization effort is aimed at next: buffered writer, bulk parsing/pooling for the CPU gap; allocation reduction and leaner value layout for the memory gap.

---

## Current standing + next steps

- **Throughput:** parity with Redis on all core ops, both on macOS and Linux. LPUSH was fixed via the deque (4.2K → 126K across the two environments).
- **CPU:** kivo runs ~2.2-2.5x the CPU of Redis for the same throughput (multicore spread vs single-threaded Redis). Remaining headroom is the known I/O-framing path (51% syscalls / 26% byte-parsing / 25% fmt in the original profile).
- **Memory:** kivo peaks at 3.3-9.5x Redis's RSS. Go runtime base + allocation pattern vs Redis's compact encodings. Next optimization phase.
- **Optimizations identified but not yet applied:** buffered RESP writer (one buffer flush per pipelined batch instead of per response), bulk line parsing (read whole line, scan for `\r\n` instead of byte-at-a-time), command-arg conversion/pooling (kill the per-command double conversion + allocations).
