#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

REPORT_FILE="BENCHMARKS.md"
PORT=16379

echo "=== Building kivo ==="
go build -o kivo-bench ./cmd/kivo

echo "=== Starting kivo on port ${PORT} ==="
./kivo-bench -addr ":${PORT}" -aof=false > /tmp/kivo-bench.log 2>&1 &
KIVO_PID=$!

echo "=== Waiting for kivo to be ready ==="
for i in {1..10}; do
    if redis-cli -p ${PORT} PING > /dev/null 2>&1; then
        echo "kivo is ready"
        break
    fi
    if [ $i -eq 10 ]; then
        echo "ERROR: kivo failed to start"
        cat /tmp/kivo-bench.log
        kill $KIVO_PID 2>/dev/null || true
        rm -f kivo-bench
        exit 1
    fi
    sleep 0.5
done

echo "=== Running redis-benchmark ==="

cat > "${REPORT_FILE}" << 'EOF'
# Benchmark Results

Tests run against kivo using `redis-benchmark`.

## Environment
EOF

echo "- Date: $(date -u +"%Y-%m-%d %H:%M:%S UTC")" >> "${REPORT_FILE}"
echo "- Host: $(uname -a)" >> "${REPORT_FILE}"
echo "- Go Version: $(go version)" >> "${REPORT_FILE}"
echo "" >> "${REPORT_FILE}"

# Run benchmark tests
echo "## Basic Operations (100000 requests, 50 parallel clients)" >> "${REPORT_FILE}"
echo "" >> "${REPORT_FILE}"
echo '```' >> "${REPORT_FILE}"
redis-benchmark -p ${PORT} -n 100000 -c 50 -t ping,set,get,incr,lpush,rpush,lpop,rpop,sadd,hset 2>&1 | tee -a "${REPORT_FILE}"
echo '```' >> "${REPORT_FILE}"
echo "" >> "${REPORT_FILE}"

echo "## String Operations (100000 requests, 50 parallel clients)" >> "${REPORT_FILE}"
echo "" >> "${REPORT_FILE}"
echo '```' >> "${REPORT_FILE}"
redis-benchmark -p ${PORT} -n 100000 -c 50 -t set,get 2>&1 | tee -a "${REPORT_FILE}"
echo '```' >> "${REPORT_FILE}"
echo "" >> "${REPORT_FILE}"

echo "## List Operations (100000 requests, 50 parallel clients)" >> "${REPORT_FILE}"
echo "" >> "${REPORT_FILE}"
echo '```' >> "${REPORT_FILE}"
redis-benchmark -p ${PORT} -n 100000 -c 50 -t lpush,lpop,rpush,rpop 2>&1 | tee -a "${REPORT_FILE}"
echo '```' >> "${REPORT_FILE}"
echo "" >> "${REPORT_FILE}"

echo "## Set & Hash Operations (100000 requests, 50 parallel clients)" >> "${REPORT_FILE}"
echo "" >> "${REPORT_FILE}"
echo '```' >> "${REPORT_FILE}"
redis-benchmark -p ${PORT} -n 100000 -c 50 -t sadd,hset 2>&1 | tee -a "${REPORT_FILE}"
echo '```' >> "${REPORT_FILE}"
echo "" >> "${REPORT_FILE}"

echo "## Pipelined Operations (100000 requests, 50 parallel clients, 16 pipeline)" >> "${REPORT_FILE}"
echo "" >> "${REPORT_FILE}"
echo '```' >> "${REPORT_FILE}"
redis-benchmark -p ${PORT} -n 100000 -c 50 -P 16 -t set,get 2>&1 | tee -a "${REPORT_FILE}"
echo '```' >> "${REPORT_FILE}"
echo "" >> "${REPORT_FILE}"

echo "Killing kivo server (PID: ${KIVO_PID})" >> "${REPORT_FILE}"
kill $KIVO_PID 2>/dev/null || true
wait $KIVO_PID 2>/dev/null || true

# Cleanup
rm -f kivo-bench

echo ""
echo "=== Benchmark complete ==="
echo "Results saved to ${REPORT_FILE}"
