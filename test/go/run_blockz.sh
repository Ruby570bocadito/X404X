#!/bin/bash
# F2b: Tests Go BlockZ
set -e
cd "$(dirname "$0")/../.."
BOLD='\033[1m'; RED='\033[0;31m'; GREEN='\033[0;32m'; NC='\033[0m'
echo -e "${BOLD}F2b: BlockZ Tests${NC}"
go test -v -count=1 -timeout 60s ./internal/ransomware/blockz/... 2>&1
STATUS=$?
[ $STATUS -eq 0 ] && echo -e "${GREEN}  ALL BLOCKZ TESTS PASSED${NC}" || echo -e "${RED}  BLOCKZ TESTS FAILED${NC}"
exit $STATUS
