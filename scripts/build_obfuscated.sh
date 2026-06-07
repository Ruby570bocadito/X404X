#!/bin/bash
echo "=== X404X Obfuscated Build ==="
SEED=$(date +%s)
echo "Seed: $SEED"
if command -v garble &>/dev/null; then
    garble -literals -tiny -seed=random build -ldflags="-s -w" -o x404x_obfuscated ./cmd/x404x/ 2>&1 && echo "OK: x404x_obfuscated"
else
    go build -ldflags="-s -w" -o x404x_obfuscated ./cmd/x404x/ && echo "OK (go build fallback): x404x_obfuscated"
fi
ls -la x404x_obfuscated 2>/dev/null
