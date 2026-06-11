#!/bin/bash
# F2a: Tests Go Ransomware — módulos base
set -e
cd "$(dirname "$0")/../.."
BOLD='\033[1m'; RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'
echo -e "${BOLD}========================================${NC}"
echo -e "${BOLD}  F2a: Go Ransomware Base Tests${NC}"
echo -e "${BOLD}========================================${NC}"

go test -v -count=1 -timeout 120s ./internal/ransomware/... 2>&1
STATUS=$?

echo -e "\n${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
[ $STATUS -eq 0 ] && echo -e "${GREEN}  ALL RANSOMWARE TESTS PASSED${NC}" || echo -e "${RED}  SOME RANSOMWARE TESTS FAILED${NC}"
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
exit $STATUS
