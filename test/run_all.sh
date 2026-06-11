#!/bin/bash
# X404X — Master Test Suite Runner
# Ejecuta todas las fases de testing en orden
set -e
cd "$(dirname "$0")"
BOLD='\033[1m'; RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'

echo -e "${CYAN}${BOLD}"
echo "  ╔═══════════════════════════════════════╗"
echo "  ║       X404X Complete Test Suite       ║"
echo "  ╚═══════════════════════════════════════╝"
echo -e "${NC}"

START=$(date +%s%N)
TOTAL=0; PASS=0; FAIL=0

run_phase() {
  local phase=$1; local name=$2; local script=$3
  TOTAL=$((TOTAL+1))
  echo -e "\n${CYAN}${BOLD}[${phase}]${NC} ${BOLD}${name}${NC}"
  echo -e "${CYAN}──────────────────────────────────────────${NC}"
  if bash "$script"; then
    echo -e "${GREEN}  [${phase}] PASS${NC}"
    PASS=$((PASS+1))
  else
    echo -e "${RED}  [${phase}] FAIL${NC}"
    FAIL=$((FAIL+1))
  fi
}

# ── F1: Go Core ──
run_phase "F1" "Go Core Tests" "go/run_core.sh"

# ── F2: Go Ransomware ──
run_phase "F2a" "Go Ransomware Base" "go/run_ransomware.sh"
run_phase "F2b" "Go BlockZ" "go/run_blockz.sh"
run_phase "F2c" "Go v210 Apocalipsis" "go/run_v210.sh"
run_phase "F2d" "Go v26 POMDP" "go/run_v26.sh"
run_phase "F2e" "Go v27 Blue Pill" "go/run_v27.sh"
run_phase "F2f" "Go v28 Malice" "go/run_v28.sh"
run_phase "F2g" "Go v29 Network" "go/run_v29.sh"
run_phase "F2h" "Go v30 AD+Payroll" "go/run_v30.sh"
run_phase "F2i" "Go Hydra Vectors" "go/run_hydra_vectors.sh"

# ── F3: Python Bridge ──
run_phase "F3" "Python Bridge" "python/run_bridge.sh"

# ── F4: IntegraciOn ──
run_phase "F4" "Integration" "integration/run_integration.sh"

# ── F5: E2E ──
run_phase "F5a" "E2E Kill Chain" "e2e/run_killchain.sh"
run_phase "F5b" "E2E Campaign" "e2e/run_campaign.sh"

# ── F6: Security ──
run_phase "F6" "Security & Evasion" "security/run_evasion.sh"

# ── F7: Benchmarks ──
run_phase "F7" "Benchmarks" "benchmark/run_benchmarks.sh"

# ── Summary ──
ELAPSED=$((($(date +%s%N) - START)/1000000))
echo -e "\n${BOLD}${CYAN}═══════════════════════════════════════════${NC}"
echo -e "${BOLD}  RESULTS: ${GREEN}${PASS} PASS${NC} / ${RED}${FAIL} FAIL${NC} / ${TOTAL} TOTAL"
echo -e "${BOLD}  ELAPSED: ${ELAPSED}ms${NC}"
echo -e "${BOLD}${CYAN}═══════════════════════════════════════════${NC}"

[ $FAIL -eq 0 ]
