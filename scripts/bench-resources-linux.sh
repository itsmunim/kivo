#!/usr/bin/env bash
# Linux (Docker) resource benchmark: kivo vs Redis in containers,
# measured with docker stats alongside redis-benchmark. 3 rounds.
#
# Requires: docker (Linux kernel available), redis-benchmark (or the
# redis-tools image referenced below), kivo:bench image (build via
# scripts/bench-resources-linux.sh after `docker build -t kivo:bench .`).
set -uo pipefail
cd "$(dirname "$0")/.."

OUT=/tmp/bench-res-linux
mkdir -p "$OUT"

ROUNDS=3
N=300000
C=50
OPS="set,get,lpush,rpush,incr"
NET="kivo-bench-net"
R_NAME="bench-redis"
K_NAME="bench-kivo"

cleanup() {
  docker rm -f $R_NAME $K_NAME 2>/dev/null
  docker network rm $NET 2>/dev/null
}
trap cleanup EXIT

docker network create $NET > /dev/null 2>&1

echo "=== Linux (docker) resource usage: kivo vs Redis ($(date -u +%Y-%m-%dT%H:%M:%SZ)) ==="
echo "image: kivo:bench vs redis:7-alpine | ${N} req, ${C} clients, ops ${OPS}"

# The stats sampler: collects cpu% + mem (MiB) from `docker stats --no-stream`
sample_pid() { # $1=container  $2=prefix
  local c=$1 p=$2
  local cpu_f="$OUT/$p-cpu.txt" mem_f="$OUT/$p-mem.txt"
  : > "$cpu_f"; : > "$mem_f"
  while true; do
    line=$(docker stats --no-stream --format "{{.CPUPerc}} {{.MemUsage}}" "$c" 2>/dev/null)
    [ -n "$line" ] || continue
    cpu=$(echo "$line" | awk '{print $1}' | tr -d '%')
    mem=$(echo "$line" | awk '{print $2}' | grep -oE '^[0-9.]+' )
    echo "$cpu" >> "$cpu_f"
    echo "$mem" >> "$mem_f"
    sleep 0.2
  done
}

avg()  { awk '{s+=$1;n++}END{if(n)printf "%.1f",s/n;else print 0}'; }
maxv() { sort -n | tail -1; }
tput() { # $1=benchfile  → per-op table string
  local f=$1
  grep -aoE '[A-Z]+: [0-9]+\.[0-9]+ requests per second.*' "$f" | tail -5 |
    awk '{v[$1]=$2} END{for(o in v) printf "%s=%s ", o, v[o]}'
}

echo "round | server | avg_cpu% | peak_mem MiB | tput"

for round in 1 2 3; do
  # ---------- Redis ----------
  docker run -d --name $R_NAME --network $NET --rm \
    -e REDIS_PWD= -p 16390:6379 redis:7-alpine --save '' --appendonly no > /dev/null
  sleep 2
  sample_pid $R_NAME "r$round-redis" & S1=$!
  sleep 2; kill $S1 2>/dev/null; wait $S1 2>/dev/null    # idle baseline
  sample_pid $R_NAME "r$round-redis" & S2=$!
  perl -e 'alarm 65; exec @ARGV' docker run --rm --network $NET redis:7-alpine \
    redis-benchmark -h bench-redis -p 6379 -n $N -c $C -t "$OPS" -q \
    > "$OUT/r$round-redis-bench.txt" 2>&1
  kill $S2 2>/dev/null; wait $S2 2>/dev/null
  docker rm -f $R_NAME > /dev/null 2>&1; sleep 1

  # ---------- kivo ----------
  docker run -d --name $K_NAME --network $NET --rm -p 16391:6379 \
    --entrypoint ./kivo kivo:bench -addr :6379 -aof=false -webui=false > /dev/null
  sleep 2
  sample_pid $K_NAME "r$round-kivo" & S3=$!
  sleep 2; kill $S3 2>/dev/null; wait $S3 2>/dev/null
  sample_pid $K_NAME "r$round-kivo" & S4=$!
  perl -e 'alarm 65; exec @ARGV' docker run --rm --network $NET redis:7-alpine \
    redis-benchmark -h bench-kivo -p 6379 -n $N -c $C -t "$OPS" -q \
    > "$OUT/r$round-kivo-bench.txt" 2>&1
  kill $S4 2>/dev/null; wait $S4 2>/dev/null
  docker rm -f $K_NAME > /dev/null 2>&1; sleep 1

  # ---------- summarize ----------
  echo "  $round | kivo  | $(avg < "$OUT/r$round-kivo-cpu.txt")% | $(maxv < "$OUT/r$round-kivo-mem.txt") | $(tput "$OUT/r$round-kivo-bench.txt")"
  echo "  $round | redis | $(avg < "$OUT/r$round-redis-cpu.txt")% | $(maxv < "$OUT/r$round-redis-mem.txt") | $(tput "$OUT/r$round-redis-bench.txt")"
done

echo "=== DONE. raw files in $OUT ==="