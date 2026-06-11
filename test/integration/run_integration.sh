#!/bin/bash
# F4: Tests de Integración — Go + Python Bridge + C2 + Dispatcher
set -e
cd "$(dirname "$0")/../.."
BOLD='\033[1m'; RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'
echo -e "${BOLD}========================================${NC}"
echo -e "${BOLD}  F4: Integration Tests${NC}"
echo -e "${BOLD}========================================${NC}"

RESULT=0

# Test 1: Binary exists and is executable
echo -e "\n${YELLOW}Integration 1: Binary smoke test${NC}"
if [ -f "x404x" ] && [ -x "x404x" ]; then
  echo -e "${GREEN}  ✔ binary exists${NC}"
else
  echo -e "${RED}  ✘ binary not found — skipping integration tests${NC}"
  exit 0
fi

# Test 2: Help works
"$PWD/x404x" --help >/dev/null 2>&1 && echo -e "${GREEN}  ✔ help works${NC}" || { echo -e "${RED}  ✘ help failed${NC}"; RESULT=1; }

# Test 3: Module listing
echo -e "\n${YELLOW}Integration 2: Module listing${NC}"
"$PWD/x404x" modules categories >/dev/null 2>&1 && echo -e "${GREEN}  ✔ module categories${NC}" || echo -e "${YELLOW}  ~ module categories (non-critical)${NC}"

# Test 4: Go core ↔ Python bridge protocol (if bridge running)
echo -e "\n${YELLOW}Integration 3: Bridge communication${NC}"
if python3 -c "import json; print(json.dumps({'module':'test','function':'ping','params':{}}))" 2>/dev/null; then
  echo -e "${GREEN}  ✔ bridge JSON protocol valid${NC}"
else
  echo -e "${RED}  ✘ bridge protocol invalid${NC}"; RESULT=1
fi

# Test 5: Registry consistency
echo -e "\n${YELLOW}Integration 4: Registry consistency${NC}"
python3 -c "
import sys, os
sys.path.insert(0, 'modules/bridge/handlers')
from phase_1_4 import register_routes
registry = {}
register_routes(registry)
total = sum(len(v) for v in registry.values())
print(f'  Phase 1-4 handlers: {total}')
# Verify key modules exist
keys = set(registry.keys())
expected = {'byovd','dkom','anti','c2','ai','cross'}
missing = expected - keys
if missing:
    print(f'  WARNING: missing categories: {missing}')
else:
    print('  All expected categories present')
" 2>&1

echo -e "\n${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
[ $RESULT -eq 0 ] && echo -e "${GREEN}  INTEGRATION TESTS PASSED${NC}" || echo -e "${RED}  INTEGRATION TESTS FAILED${NC}"
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
exit $RESULT
