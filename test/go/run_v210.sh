#!/bin/bash
# F2c: Tests Go v210 Apocalipsis
set -e
cd "$(dirname "$0")/../.."
BOLD='\033[1m'; RED='\033[0;31m'; GREEN='\033[0;32m'; NC='\033[0m'
echo -e "${BOLD}F2c: v210 Tests${NC}"
go test -v -count=1 -timeout 60s ./internal/ransomware/v210/... 2>&1
STATUS=$?
[ $STATUS -eq 0 ] && echo -e "${GREEN}  ALL v210 TESTS PASSED${NC}" || echo -e "${RED}  v210 TESTS FAILED${NC}"
exit $STATUS
