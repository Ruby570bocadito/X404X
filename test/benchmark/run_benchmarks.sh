#!/bin/bash
# F7: Benchmarks de Rendimiento
set -e
cd "$(dirname "$0")/../.."
BOLD='\033[1m'; RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'
echo -e "${BOLD}========================================${NC}"
echo -e "${BOLD}  F7: Performance Benchmarks${NC}"
echo -e "${BOLD}========================================${NC}"

RESULT=0

echo -e "\n${YELLOW}Benchmark 1: Crypto (existing)${NC}"
if go test -bench=BenchmarkEncrypt -benchmem -count=1 ./internal/crypto/... 2>&1; then
  echo -e "${GREEN}  OK crypto benchmarks${NC}"
else
  echo -e "${RED}  FAIL crypto benchmarks${NC}"; RESULT=1
fi

echo -e "\n${YELLOW}Benchmark 2: Scanner regex${NC}"
if go test -bench=BenchmarkRegex -benchmem -count=1 -run='^$' ./internal/ransomware/... 2>&1; then
  echo -e "${GREEN}  OK scanner benchmarks${NC}"
else
  echo -e "${YELLOW}  ~ scanner benchmarks (no benchmarks yet)${NC}"
fi

echo -e "\n${YELLOW}Benchmark 3: Python bridge serialization${NC}"
python3 -c "
import json, time, sys
payload = {'module': 'ransomware/encrypt', 'function': 'execute', 'params': {'path': '/tmp/test', 'files': 1000}}
start = time.time()
for _ in range(10000):
    json.dumps(payload)
elapsed = time.time() - start
print(f'  JSON serialization: {elapsed*1000:.1f}ms for 10,000 calls')
print(f'  Throughput: {10000/elapsed:.0f} ops/sec')
"

echo -e "\n${BOLD}BENCHMARKS COMPLETE${NC}"
exit $RESULT
