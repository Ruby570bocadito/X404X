#!/bin/bash
# F2h: Tests Go v30 AD + Payroll
set -e
cd "$(dirname "$0")/../.."
BOLD='\033[1m'; RED='\033[0;31m'; GREEN='\033[0;32m'; NC='\033[0m'
echo -e "${BOLD}F2h: v30 Tests${NC}"
go test -v -count=1 -timeout 60s ./internal/ransomware/v30/... 2>&1
STATUS=$?
[ $STATUS -eq 0 ] && echo -e "${GREEN}  ALL v30 TESTS PASSED${NC}" || echo -e "${RED}  v30 TESTS FAILED${NC}"
exit $STATUS
