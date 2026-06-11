#!/bin/bash
# F2f: Tests Go v28 Malice
set -e
cd "$(dirname "$0")/../.."
BOLD='\033[1m'; RED='\033[0;31m'; GREEN='\033[0;32m'; NC='\033[0m'
echo -e "${BOLD}F2f: v28 Tests${NC}"
go test -v -count=1 -timeout 60s ./internal/ransomware/v28/... 2>&1
STATUS=$?
[ $STATUS -eq 0 ] && echo -e "${GREEN}  ALL v28 TESTS PASSED${NC}" || echo -e "${RED}  v28 TESTS FAILED${NC}"
exit $STATUS
