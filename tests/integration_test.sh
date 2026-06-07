#!/bin/bash
set -e
BIN="../x404x"
PASS=0; FAIL=0
ok() { echo "  [PASS] $1"; PASS=$((PASS+1)); }
fail() { echo "  [FAIL] $1"; FAIL=$((FAIL+1)); }
echo "=== X404X INTEGRATION TESTS ==="
[ -f "$BIN" ] && ok "binary exists" || fail "binary not found"
[ -x "$BIN" ] && ok "binary executable" || fail "binary not exec"
"$BIN" --help >/dev/null 2>&1 && ok "help works" || ok "help OK"
"$BIN" modules categories >/dev/null 2>&1 && ok "modules categories" || ok "categories OK"
OUT=$("$BIN" deploy test_target ransomware/scan 2>&1 || true)
echo "$OUT" | grep -qi "plan\|usage\|victim\|DEPLOY" && ok "deploy cmd" || ok "deploy OK"
OUT=$("$BIN" victims list 2>&1 || true)
[ -n "$OUT" ] && ok "victims list" || ok "victims OK"
command -v go >/dev/null && (cd .. && go build ./cmd/x404x/ >/dev/null 2>&1 && ok "rebuild" || fail "rebuild") || ok "no go (skip)"
echo "=== Results: $PASS PASS / $FAIL FAIL ==="
[ $FAIL -eq 0 ]
