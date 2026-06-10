#!/usr/bin/env python3
"""Smoke tests for all X404X bridge handler modules.

Imports all register_routes functions from all 8 handler files,
registers them into a unified registry, then calls each handler
with empty params to verify no exceptions and valid return types.
"""
import sys
import os
import traceback

sys.path.insert(0, os.path.join(os.path.dirname(__file__), ".."))

from handlers.ransomware import register_routes as reg_ransomware
from handlers.ransomware_advanced import register_routes as reg_advanced
from handlers.ransomware_v26 import register_routes as reg_v26
from handlers.ransomware_v27 import register_routes as reg_v27
from handlers.ransomware_v28 import register_routes as reg_v28
from handlers.ransomware_v29 import register_routes as reg_v29
from handlers.ransomware_v210 import register_routes as reg_v210
from handlers.ransomware_blockz import register_routes as reg_blockz


def build_registry():
    registry = {}
    reg_ransomware(registry)
    reg_advanced(registry)
    reg_v26(registry)
    reg_v27(registry)
    reg_v28(registry)
    reg_v29(registry)
    reg_v210(registry)
    reg_blockz(registry)
    return registry


def run_smoke_tests():
    registry = build_registry()

    total = 0
    passed = 0
    failed = 0
    failures = []

    for module_name, handlers in sorted(registry.items()):
        for handler_name, handler_fn in sorted(handlers.items()):
            total += 1
            test_id = f"{module_name}.{handler_name}"
            try:
                result = handler_fn({})
                if not isinstance(result, dict):
                    raise ValueError(f"expected dict, got {type(result).__name__}")
                if len(result) < 1:
                    raise ValueError("returned empty dict")
                passed += 1
                print(f"  PASS  {test_id} ({len(result)} keys)")
            except Exception as e:
                failed += 1
                tb = traceback.format_exc().splitlines()[-1]
                failures.append((test_id, str(e)))
                print(f"  FAIL  {test_id}: {tb}")

    print()
    print(f"{'=' * 60}")
    print(f"RESULTS: {passed}/{total} passed, {failed} failed")
    print(f"{'=' * 60}")

    if failures:
        print("\nFailed handlers:")
        for test_id, err in failures:
            print(f"  - {test_id}: {err}")

    return 0 if failed == 0 else 1


if __name__ == "__main__":
    exit_code = run_smoke_tests()
    sys.exit(exit_code)
