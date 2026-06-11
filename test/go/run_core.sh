#!/bin/bash
# F1: Tests Go Core — appstate, dispatch, bridge, c2server, api, registry, defense
set -e
cd "$(dirname "$0")/../.."
BOLD='\033[1m'; RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'
echo -e "${BOLD}========================================${NC}"
echo -e "${BOLD}  F1: Go Core Tests${NC}"
echo -e "${BOLD}========================================${NC}"

PACKAGES=(
  "./internal/appstate/..."
  "./internal/dispatch/..."
  "./internal/bridge/..."
  "./internal/c2server/..."
  "./internal/api/..."
  "./internal/registry/..."
  "./internal/defense/..."
  "./internal/crypto/..."
  "./internal/agent/..."
  "./internal/orchestrator/..."
)

TOTAL=0; PASS=0; FAIL=0
for pkg in "${PACKAGES[@]}"; do
  TOTAL=$((TOTAL+1))
  echo -e "\n${YELLOW}Testing${NC} ${BOLD}$pkg${NC}"
  if go test -v -count=1 -timeout 60s "$pkg" 2>&1; then
    echo -e "${GREEN}  ✔ PASS${NC}"
    PASS=$((PASS+1))
  else
    echo -e "${RED}  ✘ FAIL${NC}"
    FAIL=$((FAIL+1))
  fi
done

echo -e "\n${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BOLD}  Core: $PASS/$TOTAL passed, $FAIL failed${NC}"
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
[ $FAIL -eq 0 ]
