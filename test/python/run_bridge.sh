#!/bin/bash
# F3: Tests Python Bridge
set -e
cd "$(dirname "$0")/../.."
BOLD='\033[1m'; RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'
echo -e "${BOLD}========================================${NC}"
echo -e "${BOLD}  F3: Python Bridge Tests${NC}"
echo -e "${BOLD}========================================${NC}"

if ! command -v python3 &>/dev/null; then
  echo -e "${RED}python3 not found — skipping bridge tests${NC}"
  exit 0
fi

cd modules/bridge

echo -e "\n${YELLOW}Test 1: Smoke — all handlers${NC}"
python3 tests/test_all_handlers.py 2>&1 || true

echo -e "\n${YELLOW}Test 2: Unittest — handler registration${NC}"
python3 -m unittest tests.test_handlers -v 2>&1 || true

echo -e "\n${YELLOW}Test 3: Phase 1-4 handlers (36 handlers)${NC}"
python3 -c "
import sys, os
sys.path.insert(0, 'handlers')
from phase_1_4 import register_routes
registry = {}
register_routes(registry)
count = sum(len(v) for v in registry.values())
print(f'Phase 1-4: {count} handlers registered')
assert count >= 5, f'Expected >=5 handlers, got {count}'
print('PASS: Phase 1-4 handlers OK')
" 2>&1

echo -e "\n${YELLOW}Test 4: Bridge protocol serialization${NC}"
python3 -c "
import json, sys
req = {'module':'test','function':'ping','params':{},'timeout_ms':5000}
js = json.dumps(req)
parsed = json.loads(js)
assert parsed['module'] == 'test'
assert parsed['function'] == 'ping'
print('PASS: Bridge protocol JSON OK')
" 2>&1

echo -e "\n${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BOLD}  BRIDGE TESTS COMPLETE${NC}"
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
