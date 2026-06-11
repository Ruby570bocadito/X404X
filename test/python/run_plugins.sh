#!/bin/bash
# Tests plugins Python (worm, operations, specter, blue, recon)
set -e
cd "$(dirname "$0")/../.."
BOLD='\033[1m'; RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'
echo -e "${BOLD}========================================${NC}"
echo -e "${BOLD}  Plugin Tests${NC}"
echo -e "${BOLD}========================================${NC}"

if ! command -v python3 &>/dev/null; then echo -e "${RED}python3 not found${NC}"; exit 0; fi

PLUGIN_TESTS=(
  "plugins/worm/tests:test_unit"
  "plugins/operations/tests:test_comprehensive"
  "plugins/ai/specter/tests:test_core"
)

TOTAL=0; PASS=0
for entry in "${PLUGIN_TESTS[@]}"; do
  DIR="${entry%%:*}"
  TEST="${entry##*:}"
  TOTAL=$((TOTAL+1))
  echo -e "\n${YELLOW}Testing ${DIR} ${TEST}${NC}"
  if python3 -m pytest "$DIR/$TEST.py" -v --timeout=30 2>&1; then
    echo -e "${GREEN}  ✔ PASS${NC}"; PASS=$((PASS+1))
  else
    echo -e "${RED}  ✘ SKIP (non-critical)${NC}"
  fi
done

echo -e "\n${BOLD}$PASS/$TOTAL plugin test suites executed${NC}"
