#!/bin/bash
# F2d: Tests Go v26 POMDP
set -e
cd "$(dirname "$0")/../.."
BOLD='\033[1m'; RED='\033[0;31m'; GREEN='\033[0;32m'; NC='\033[0m'
echo -e "${BOLD}F2d: v26 Tests${NC}"
go test -v -count=1 -timeout 60s ./internal/ransomware/v26/... 2>&1
STATUS=$?
[ $STATUS -eq 0 ] && echo -e "${GREEN}  ALL v26 TESTS PASSED${NC}" || echo -e "${RED}  v26 TESTS FAILED${NC}"
exit $STATUS
