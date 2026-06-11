#!/bin/bash
# F2g: Tests Go v29 Network Immolation
set -e
cd "$(dirname "$0")/../.."
BOLD='\033[1m'; RED='\033[0;31m'; GREEN='\033[0;32m'; NC='\033[0m'
echo -e "${BOLD}F2g: v29 Tests${NC}"
go test -v -count=1 -timeout 60s ./internal/ransomware/v29/... 2>&1
STATUS=$?
[ $STATUS -eq 0 ] && echo -e "${GREEN}  ALL v29 TESTS PASSED${NC}" || echo -e "${RED}  v29 TESTS FAILED${NC}"
exit $STATUS
