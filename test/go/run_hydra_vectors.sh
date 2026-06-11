#!/bin/bash
# F2i: Tests Go Hydra Vectors (ultrasound, powerline, USB ADB, DNS rebind, CI/CD, VLAN, QR, PJL)
set -e
cd "$(dirname "$0")/../.."
BOLD='\033[1m'; RED='\033[0;31m'; GREEN='\033[0;32m'; NC='\033[0m'
echo -e "${BOLD}F2i: Hydra Vectors Tests${NC}"
go test -v -count=1 -timeout 60s ./internal/ransomware/... -run "TestHydra" 2>&1
STATUS=$?
[ $STATUS -eq 0 ] && echo -e "${GREEN}  ALL HYDRA VECTOR TESTS PASSED${NC}" || echo -e "${RED}  HYDRA VECTOR TESTS FAILED${NC}"
exit $STATUS
