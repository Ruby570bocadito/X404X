# X404X — Report Generator
# ========================
# Generates post-engagement reports with MITRE ATT&CK mapping.
# Output formats: JSON, Markdown, HTML, PDF (via headless Chrome)

import json
import os
import time
from datetime import datetime
from pathlib import Path
from typing import Any, Dict, List, Optional


class ReportGenerator:
    """Generates campaign reports for X404X."""

    def __init__(self, output_dir: str = "reports"):
        self.output_dir = Path(output_dir)
        self.output_dir.mkdir(parents=True, exist_ok=True)
        self._data = {
            "generated_at": datetime.now().isoformat(),
            "framework": "X404X v1.0",
            "author": "Rafael Gálvez — TFG Cybersecurity",
        }

    def add_campaign(self, campaign: Dict[str, Any]):
        self._data["campaign"] = campaign

    def add_kill_chain(self, entries: List[Dict[str, Any]]):
        self._data["kill_chain"] = entries

    def add_hosts(self, hosts: List[Dict[str, Any]]):
        self._data["hosts"] = hosts

    def add_vulnerabilities(self, vulns: List[Dict[str, Any]]):
        self._data["vulnerabilities"] = vulns

    def add_credentials(self, creds: List[Dict[str, Any]]):
        self._data["credentials"] = creds

    def add_metrics(self, metrics: Dict[str, Any]):
        self._data["metrics"] = metrics

    def add_ai_decisions(self, decisions: List[Dict[str, Any]]):
        self._data["ai_decisions"] = decisions

    def add_audit_log(self, entries: List[Dict[str, Any]]):
        self._data["audit_log"] = entries

    def to_json(self, path: Optional[str] = None) -> str:
        """Export report as JSON."""
        filename = path or str(self.output_dir / f"report_{timestamp()}.json")
        with open(filename, "w") as f:
            json.dump(self._data, f, indent=2, default=str)
        return filename

    def to_markdown(self, path: Optional[str] = None) -> str:
        """Export report as Markdown with MITRE ATT&CK mapping."""
        filename = path or str(self.output_dir / f"report_{timestamp()}.md")

        with open(filename, "w") as f:
            d = self._data

            f.write(f"# X404X — Post-Engagement Report\n\n")
            f.write(f"**Generated:** {d['generated_at']}\n")
            f.write(f"**Framework:** {d['framework']}\n")
            f.write(f"**Author:** {d['author']}\n\n")
            f.write("---\n\n")

            # Campaign Summary
            camp = d.get("campaign", {})
            if camp:
                f.write("## Campaign Summary\n\n")
                f.write(f"| Field | Value |\n|---|---|\n")
                for k, v in camp.items():
                    f.write(f"| {k} | {v} |\n")
                f.write("\n")

            # Kill Chain
            kc = d.get("kill_chain", [])
            if kc:
                f.write("## Kill Chain Execution\n\n")
                f.write("| Phase | Tactic | Technique | MITRE ID | Success | Detail |\n")
                f.write("|---|---|---|---|---|---|\n")
                for e in kc:
                    f.write(f"| {e.get('phase','')} | {e.get('tactic','')} | {e.get('technique','')} | {e.get('mitre_id','')} | {e.get('success','')} | {e.get('detail','')} |\n")
                f.write("\n")

            # MITRE ATT&CK Mapping
            f.write("## MITRE ATT&CK Mapping\n\n")
            f.write("| Tactic | Technique ID | Description |\n")
            f.write("|---|---|---|\n")
            mitre_map = {
                "Reconnaissance": "T1593", "Initial Access": "T1190",
                "Execution": "T1203", "Privilege Escalation": "T1068",
                "Persistence": "T1547.006", "Command and Control": "T1071.001",
                "Lateral Movement": "T1570", "Collection": "T1005",
                "Exfiltration": "T1041",
            }
            for tactic, tid in mitre_map.items():
                found = any(e.get("tactic") == tactic for e in kc)
                status = "✓ Executed" if found else "— Not executed"
                f.write(f"| {tactic} | [{tid}](https://attack.mitre.org/techniques/{tid.replace('.','/')}/) | {status} |\n")
            f.write("\n")

            # Hosts
            hosts = d.get("hosts", [])
            if hosts:
                f.write("## Discovered Hosts\n\n")
                f.write("| IP | Hostname | OS | Compromised |\n|---|---|---|---|\n")
                for h in hosts:
                    compromised = "●" if h.get("asset_value", 0) > 50 else "○"
                    f.write(f"| {h.get('ip','')} | {h.get('hostname','')} | {h.get('os','')} | {compromised} |\n")
                f.write("\n")

            # Vulnerabilities
            vulns = d.get("vulnerabilities", [])
            if vulns:
                f.write("## Discovered Vulnerabilities\n\n")
                f.write("| CVE | Severity | Service | Target |\n|---|---|---|---|\n")
                for v in vulns:
                    f.write(f"| {v.get('cve','')} | {v.get('severity','')} | {v.get('service','')}:{v.get('port','')} | {v.get('ip','')} |\n")
                f.write("\n")

            # Metrics
            metrics = d.get("metrics", {})
            if metrics:
                f.write("## BlueForge Detection Metrics\n\n")
                f.write(f"| Metric | Value |\n|---|---|\n")
                for k, v in metrics.items():
                    f.write(f"| {k} | {v} |\n")
                f.write("\n")

            # Ethical Notice
            f.write("---\n\n")
            f.write("## ⚠️ Ethical Notice\n\n")
            f.write("This report was generated exclusively for authorized security assessments. ")
            f.write("All actions were performed in controlled laboratory environments with explicit permission. ")
            f.write("The X404X Framework is an academic project (TFG Cybersecurity, Cisco NetAcad, Málaga).\n")

        return filename

    def to_html(self, path: Optional[str] = None) -> str:
        """Export report as HTML."""
        md_file = self.to_markdown()
        html_file = path or md_file.replace(".md", ".html")
        # In production: use markdown-to-html converter
        with open(html_file, "w") as f:
            f.write(f"<!DOCTYPE html><html><head><title>X404X Report</title>")
            f.write(f"<style>body{{font-family:monospace;background:#0a0a0f;color:#e0e0e0;padding:2em}}h1{{color:#6c63ff}}h2{{color:#00ff41}}table{{border-collapse:collapse}}td,th{{border:1px solid #333;padding:4px 8px}}</style>")
            f.write(f"</head><body><pre>")
            with open(md_file) as mdf:
                f.write(mdf.read())
            f.write(f"</pre></body></html>")
        return html_file

    def to_pdf(self, path: Optional[str] = None) -> str:
        """Export report as PDF (placeholder — requires headless Chrome)."""
        html_file = self.to_html()
        pdf_file = path or html_file.replace(".html", ".pdf")
        # In production: subprocess.run(["google-chrome", "--headless", "--print-to-pdf=" + pdf_file, html_file])
        return pdf_file


def timestamp() -> str:
    return datetime.now().strftime("%Y%m%d_%H%M%S")


# ============================================================
# DEMO
# ============================================================

def demo_report() -> ReportGenerator:
    rg = ReportGenerator()
    rg.add_campaign({
        "name": "TFG-Demo", "target": "10.0.0.0/24",
        "goal": "domain_admin", "profile": "balanced",
        "status": "completed", "phase": "actions_on_objective",
        "progress": 1.0, "duration": "47m 12s",
    })
    rg.add_kill_chain([
        {"phase": "recon", "tactic": "Reconnaissance", "technique": "Network Scan", "mitre_id": "T1046", "success": True, "detail": "23 hosts discovered"},
        {"phase": "delivery", "tactic": "Initial Access", "technique": "EternalBlue MS17-010", "mitre_id": "T1210", "success": True, "detail": "NT\\SYSTEM on 10.0.0.10"},
        {"phase": "exploitation", "tactic": "Privilege Escalation", "technique": "SUID GTFOBins", "mitre_id": "T1548.001", "success": True, "detail": "Root on 10.0.0.20 via python3 SUID"},
        {"phase": "installation", "tactic": "Persistence", "technique": "Kernel Rootkit", "mitre_id": "T1547.006", "success": True, "detail": "Vault-Kernel LKM loaded"},
        {"phase": "c2", "tactic": "Command and Control", "technique": "Encrypted C2", "mitre_id": "T1071.001", "success": True, "detail": "X25519+XChaCha20"},
        {"phase": "actions", "tactic": "Lateral Movement", "technique": "Wormy-ML Propagation", "mitre_id": "T1570", "success": True, "detail": "8/23 hosts infected"},
    ])
    rg.add_hosts([
        {"ip": "10.0.0.10", "hostname": "DC", "os": "Windows 2019", "asset_value": 100},
        {"ip": "10.0.0.20", "hostname": "DB", "os": "Ubuntu 24.04", "asset_value": 70},
    ])
    rg.add_vulnerabilities([
        {"cve": "MS17-010", "severity": "critical", "service": "smb", "port": 445, "ip": "10.0.0.10"},
    ])
    rg.add_metrics({"stealth_rating": 0.87, "detections": 2, "total_exploits": 8, "successful": 6})
    return rg


if __name__ == "__main__":
    rg = demo_report()
    md = rg.to_markdown()
    print(f"Report generated: {md}")
