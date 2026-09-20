#!/usr/bin/env bash
# Measure resource usage (CPU%, RSS) of kivo vs Redis under identical
# redis-benchmark load on the same machine. macOS-compatible bash (3.2).
#
# Usage: bash scripts/bench-resources.sh
# Output: detail files in /tmp/bench-res + printed summary.
set -uo pipefail

cd "$(dirname "$0")/.."

OUT=/tmp/bench-res
mkdir -p "$OUT"

ROUNDS=3
N=300000          # enough for a steady sampling window under load
C=50
OPS="set,get,lpush,rpush,incr"
R_PORT=16390
K_PORT=16391

kill_strays() {
  pkill -f "redis.*16390" 2>/dev/null
  pkill -f "kivo.*16391" 2>/dev/null
  sleep 1
}

# sampler <pid> <prefix> — records RSS(KB) + CPU(%) every 0.15s until killed
sampler() {
  local pid=$1 prefix=$2
  local rss_file="$OUT/${prefix}-rss.txt"
  local cpu_file="$OUT/${prefix}-cpu.txt"
  : > "$rss_file"; : > "$cpu_file"
  while true; do
    rss=$(ps -o rss= -p "$pid" 2>/dev/null | tr -d ' ')
    cpu=$(ps -o %cpu= -p "$pid" 2>/dev/null | tr -d ' ')
    if [ -n "$rss" ]; then echo "$rss" >> "$rss_file"; fi
    if [ -n "$cpu" ]; then echo "$cpu" >> "$cpu_file"; fi
    sleep 0.15
  done
}

avg()  { awk '{s+=$1;n++} END {if(n>0) printf "%.1f", s/n; else print "0"}'; }
maxv() { sort -n | tail -1; }

# best_throughput <benchfile> — total ops/sec across the sampled ops
best_tput() {
  awk -F: '/requests per second/ {s+=$2} END {printf "%.0f", s}'
}

echo "=== kivo vs Redis resource usage on same machine ($(date -u +%Y-%m-%dT%H:%M:%SZ)) ==="
kill_strays

go build -o "$OUT/kivo" ./cmd/kivo 2>&1 || { echo "BUILD FAILED"; exit 1; }

echo "round | server | avg_cpu% | peak_rss KB | tput ops/s"
echo "------|--------|----------|-------------|-----------"

for round in 1 2 3; do
  # ---------- Redis ----------
  redis-server --port $R_PORT --save '' --appendonly no --daemonize yes
  sleep 1
  RPID=$(pgrep -f "redis-server.*$R_PORT" | head -1)
  # idle baseline (2s)
  sampler "$RPID" "r$round-redis-idle" & S1=$!
  sleep 2; kill $S1 2>/dev/null; wait $S1 2>/dev/null
  # under load: sample WHILE benchmark runs
  sampler "$RPID" "r$round-redis" & S2=$!
  redis-benchmark -p $R_PORT -n $N -c $C -t "$OPS" -q > "$OUT/r$round-redis-bench.txt" 2>&1
  kill $S2 2>/dev/null; wait $S2 2>/dev/null
  redis-cli -p $R_PORT shutdown nosave 2>/dev/null; sleep 1

  # ---------- kivo ----------
  "$OUT/kivo" -addr ":$K_PORT" -aof=false -webui=false > "$OUT/kivo.log" 2>&1 < /dev/null &
  KPID=$!
  sleep 1
  sampler "$KPID" "r$round-kivo-idle" & S3=$!
  sleep 2; kill $S3 2>/dev/null; wait $S3 2>/dev/null
  sampler "$KPID" "r$round-kivo" & S4=$!
  redis-benchmark -p $K_PORT -n $N -c $C -t "$OPS" -q > "$OUT/r$round-kivo-bench.txt" 2>&1
  kill $S4 2>/dev/null; wait $S4 2>/dev/null
  kill $KPID 2>/dev/null; wait $KPID 2>/dev/null; sleep 1

  # ---------- summarize round ----------
  kr=$(maxv < "$OUT/r$round-kivo-rss.txt")
  kc=$(avg   < "$OUT/r$round-kivo-cpu.txt")
  kt=$(best_tput < "$OUT/r$round-kivo-bench.txt")
  rr=$(maxv < "$OUT/r$round-redis-rss.txt")
  rc=$(avg  < "$OUT/r$round-redis-cpu.txt")
  rt=$(best_tput < "$OUT/r$round-redis-bench.txt")

  echo "  $round | kivo  | $kc% | $kr | $kt"
  echo "  $round | redis | $rc% | $rr | $rt"
done

kill_strays
echo "=== DONE. raw files in $OUT ==="