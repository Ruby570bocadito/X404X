#!/bin/bash
# F2e: Tests Go v27 Blue Pill + Phishing
set -e
cd "$(dirname "$0")/../.."
BOLD='\033[1m'; RED='\033[0;31m'; GREEN='\033[0;32m'; NC='\033[0m'
echo -e "${BOLD}F2e: v27 Tests${NC}"
go test -v -count=1 -timeout 60s ./internal/ransomware/v27/... 2>&1
STATUS=$?
[ $STATUS -eq 0 ] && echo -e "${GREEN}  ALL v27 TESTS PASSED${NC}" || echo -e "${RED}  v27 TESTS FAILED${NC}"
exit $STATUS
