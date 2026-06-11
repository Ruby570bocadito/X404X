<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=waving&color=0:ff3366,100:00d4ff&height=220&section=header&text=X404X&fontSize=72&fontColor=ffffff&animation=twinkling&desc=Autonomous+Red+Team+Platform+%7C+Full+Attack+Chain+%7C+Post-Quantum+C2&descSize=16&descAlignY=65&fontAlignY=35" alt="X404X Header" />
</p>

<div align="center">

[![Typing SVG](https://readme-typing-svg.herokuapp.com?font=Fira+Code&weight=600&size=20&pause=1500&color=ff3366&center=true&vCenter=true&width=750&lines=Autonomous+Red+Team+Platform;Kernel+Evasion+%7C+Post-Quantum+C2+%7C+AI+Orchestrator;45+Offensive+Modules+%7C+Go+%2B+Python+%2B+WASM;Building+the+Full+Attack+Chain+Since+2025)](https://git.io/typing-svg)

<br/>

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Python](https://img.shields.io/badge/Python-3.11+-3776AB?style=for-the-badge&logo=python&logoColor=white)](https://python.org)
[![Lines of Code](https://img.shields.io/badge/Lines-12%2C685-ff3366?style=for-the-badge&logo=codecrafters&logoColor=white)]()
[![Modules](https://img.shields.io/badge/Modules-45-00d4ff?style=for-the-badge&logo=target&logoColor=white)]()
[![License](https://img.shields.io/badge/License-Educational-00ff88?style=for-the-badge&logo=bookstack&logoColor=white)]()
[![Status](https://img.shields.io/badge/Status-Active-00ff88?style=for-the-badge&logo=statuspage&logoColor=white)]()

<br/>

[![GitHub stars](https://img.shields.io/github/stars/Ruby570bocadito/X404X?style=flat&color=ff3366&logo=github)](https://github.com/Ruby570bocadito/X404X/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/Ruby570bocadito/X404X?style=flat&color=00d4ff&logo=github)](https://github.com/Ruby570bocadito/X404X/network)
[![GitHub last commit](https://img.shields.io/github/last-commit/Ruby570bocadito/X404X?style=flat&color=00ff88&logo=git)](https://github.com/Ruby570bocadito/X404X/commits)

</div>

---

## What is X404X?

> **X404X** is an autonomous red team platform that executes the complete offensive kill chain — from initial reconnaissance to post-exploitation and command & control. Built in **Go** and **Python** with a modular architecture, it combines kernel-level evasion, post-quantum C2 cryptography, 8 exotic propagation vectors, and AI-driven decision making into a single unified framework.

**Purpose-built for:** Red team operations · Adversarial simulation · Security research · EDR/AV evasion testing

```yaml
name: X404X
type: Autonomous Red Team Platform
architecture: Modular (Go Core + Python Bridge + WASM Extensions)
version: 3.0.0
phases: 4 (0-4)
modules: 45
languages:
  core: Go (internal/ransomware/ — 35+ files)
  bridge: Python (modules/bridge/handlers/ — 12 files, ~170 handlers)
  plugins: Go, Python, Vue 3, WASM
capabilities:
  - Kernel-level evasion (BYOVD, DKOM, Blue Pill hypervisor)
  - Post-quantum C2 (Kyber-1024 + X25519 + Ed25519 + SPIFFE mTLS)
  - Exotic propagation (ultrasound audio, powerline PLC, QR codes, printer PJL)
  - AI decision engine (Q-learning FSM, federated learning, deepfake vishing)
  - Cross-platform payloads (ELF, Mach-O, APK)
featured_in:
  - Red team adversarial simulation
  - Penetration testing & security assessments
  - Research into advanced evasion techniques
  - EDR/AV detection bypass validation
```

---

## Architecture

```
                            ┌─────────────────────────────────┐
                            │        DASHBOARD OPS            │
                            │   HTTP/WS API · Vue 3 · D3.js   │
                            │   Agent Map · Propagation Graph │
                            └──────────────┬──────────────────┘
                                           │
┌──────────────────────────────────────────┼──────────────────────────────────────┐
│                              C2 HARDENED                                        │
│  ┌────────────┐ ┌──────────┐ ┌────────────┐ ┌──────────┐ ┌──────────────────┐  │
│  │ SPIFFE     │ │ Ed25519  │ │ Kyber-1024 │ │ 5-Channel│ │ Proto Obfuscator │  │
│  │ mTLS+SVID  │ │ Signing  │ │ + X25519   │ │ C2 Stack │ │ XOR+AES-CTR+GZIP │  │
│  └────────────┘ └──────────┘ └────────────┘ └──────────┘ └──────────────────┘  │
└──────────────────────────────────────────┬──────────────────────────────────────┘
                                           │
┌──────────────────────────────────────────┼──────────────────────────────────────┐
│                              AI ORCHESTRATOR                                    │
│  ┌────────────┐ ┌──────────┐ ┌────────────┐ ┌──────────┐ ┌──────────────────┐  │
│  │ Q-Learning │ │ Federated│ │ Autofactory│ │ Deepfake │ │ Wazero Bridge    │  │
│  │ FSM Engine │ │ Learning │ │ AFL++      │ │ Vishing  │ │ WASM ↔ Go        │  │
│  └────────────┘ └──────────┘ └────────────┘ └──────────┘ └──────────────────┘  │
└──────────────────────────────────────────┬──────────────────────────────────────┘
                                           │
┌──────────────────────────────────────────┼──────────────────────────────────────┐
│                              CORE ENGINE (Go)                                   │
│                                                                                 │
│  ┌──────────────────────────────┐  ┌──────────────────────────────────────┐    │
│  │     RANSOMWARE ENGINE        │  │        EVASION SUITE                  │    │
│  │  ├─ AES-256-GCM Encryption   │  │  ├─ BYOVD Loader (5 drivers)         │    │
│  │  ├─ Bootkit (MBR/GPT/UEFI)   │  │  ├─ DKOM (process/token hiding)      │    │
│  │  ├─ SCADA/Modbus Attack      │  │  ├─ Anti-Reversing (HW BP, INT3)     │    │
│  │  ├─ Blockchain C2 (BTC/ETH)  │  │  ├─ Anti-Forensics (DoD 7-pass)      │    │
│  │  ├─ Polimorfismo JIT         │  │  ├─ WER Persistence                  │    │
│  │  ├─ Cross-Platform Loader    │  │  ├─ MFT Slack Storage                │    │
│  │  ├─ Multi-Channel C2         │  │  ├─ WFP DNS Poisoning                │    │
│  │  └─ 34 files total           │  │  ├─ Blue Pill Hypervisor             │    │
│  └──────────────────────────────┘  │  ├─ LOLBin Dynamic Chain             │    │
│                                     │  └─ Kernel DNS Driver                │    │
│  ┌──────────────────────────────┐  └──────────────────────────────────────┘    │
│  │    ADVANCED PROPAGATION      │                                               │
│  │  ├─ Ultrasound QPSK (19kHz)  │  ┌──────────────────────────────────────┐    │
│  │  ├─ Powerline PLC            │  │     LATERAL MOVEMENT                  │    │
│  │  ├─ USB ADB (Android)        │  │  ├─ Kerberos Delegation Abuse        │    │
│  │  ├─ DNS Rebinding (TTL=0)    │  │  ├─ IMDSv2 AWS Bypass                │    │
│  │  ├─ CI/CD Webhooks Injection │  │  ├─ Pass-the-Ticket / Silver Ticket  │    │
│  │  ├─ VLAN Jump (Double Tag)   │  │  ├─ Reflective DLL (100-byte stager)  │    │
│  │  ├─ QR Dynamic Worm          │  │  ├─ Chronos NTP Manipulation         │    │
│  │  └─ PJL Printer Worm         │  │  └─ Chronos NTP Manipulation         │    │
│  └──────────────────────────────┘  └──────────────────────────────────────┘    │
│                                                                                 │
└──────────────────────────────────────────┬──────────────────────────────────────┘
                                           │
┌──────────────────────────────────────────┼──────────────────────────────────────┐
│                            PYTHON BRIDGE                                        │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │  12 handler files · ~170 RPC handlers · v26↔v210 · BlockZ · Attacks       │  │
│  │  Phase 1-4 handlers (36 modules exposed via Go→Python RPC bridge)         │  │
│  └──────────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────┬──────────────────────────────────────┘
                                           │
┌──────────────────────────────────────────┼──────────────────────────────────────┐
│                            PLUGIN ECOSYSTEM                                     │
│  ┌────────┐ ┌──────────┐ ┌────────┐ ┌────────┐ ┌──────────┐ ┌──────────────┐  │
│  │ Worm+RL│ │Operations│ │Bluesky │ │ Pulse  │ │    AI    │ │ RF Contagion │  │
│  │  35    │ │ Argos    │ │  BT    │ │  C2    │ │ Hivemind │ │  SDR 4G/5G   │  │
│  │exploits│ │ Cell+UI  │ │ Attack │ │  Go    │ │Fuzzer+TT│ │ Baseband     │  │
│  └────────┘ └──────────┘ └────────┘ └────────┘ └──────────┘ └──────────────┘  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Kill Chain Coverage

```
  RECON          INITIAL ACCESS      EXECUTION       PERSISTENCE      PRIVESC
┌─────────┐    ┌──────────────┐   ┌───────────┐   ┌────────────┐   ┌──────────┐
│ OSINT   │───▶│ Phishing AI  │──▶│LOLBin     │──▶│ WER Hangs  │──▶│ BYOVD    │
│ DNS Map │    │ CI/CD Inject │   │Chainer    │   │ Triple     │   │ DKOM     │
│ Recon   │    │ QR Payload   │   │Reflect DLL│   │ Startup    │   │ Token    │
│ APIs    │    │ USB ADB      │   │WASM Exec  │   │ Schtasks   │   │ Steal    │
└─────────┘    └──────────────┘   └───────────┘   └────────────┘   └──────────┘
                                                                         │
  C2              EXFIL            LATERAL          EVASION              │
┌──────────┐   ┌───────────┐   ┌────────────┐   ┌────────────┐          │
│ SPIFFE   │◀──│ MFT Slack │◀──│ Kerberos   │◀──│ WFP DNS    │◀─────────┘
│ Ed25519  │   │ DNS DoH   │   │ Delegation │   │ Blue Pill  │
│ Kyber    │   │ Blockchain│   │ IMDSv2 AWS │   │ Anti-Revers│
│ 5-Chan   │   │ QR Exfil  │   │ VLAN Jump  │   │ Anti-Foren │
└──────────┘   └───────────┘   └────────────┘   └────────────┘
```

---

## Project Status

| Phase | Theme | Modules | Lines | Status |
|:-----:|-------|:-------:|-------|:------:|
| **0** | Critical Stubs | 8 | 500 | ✅ Complete |
| **1** | Evasion + Anti-Forensics | 10 | 3,599 | ✅ Complete |
| **2** | C2 Hardened | 6 | 2,374 | ✅ Complete |
| **3** | Advanced Propagation | 12 | 3,274 | ✅ Complete |
| **4** | AI + Cross-Platform | 9 | 2,938 | ✅ Complete |
| | **TOTAL** | **45** | **12,685** | **100%** |

---

## Module Catalog

<details open>
<summary><b>🔒 FASE 1 — Evasion & Anti-Forensics</b></summary>

| Module | File | Technique |
|--------|------|-----------|
| BYOVD Loader | `byovd_loader.go` | 5 vulnerable drivers (WinRing0, Gdrv, RTCore64, kprocesshacker, CPUID) — DeviceIoControl, R/W physical memory, MSR, handle elevation |
| DKOM | `dkom.go` | Process hiding via ActiveProcessLinks, SYSTEM token stealing, EPROCESS offsets per build |
| Anti-Reversing | `anti_reversing.go` | HW BP detect (DR0-DR7), INT3 scan, CRC32 integrity, RDTSC timing, virtual MAC sandbox check |
| Anti-Forensics | `anti_forensics_advanced.go` | DoD 5220.22-M 7-pass wipe, MFT $BITMAP corruption, VAD hide, crash dump/event log/prefetch/USN/Shellbag wipe |
| WER Persistence | `wer_persistence.go` | Windows Error Reporting Hangs hijack, Silent Process Exit, startup + Run key + schtasks triple persistence |
| MFT Slack Storage | `mft_slack.go` | NTFS slack space R/W via PowerShell, AES-GCM encrypted fragments, hidden agent/ransom note storage |
| WFP DNS Poisoning | `wfp_dns_poison.go` | WFP provider + netsh fallback, fake DNS server (UDP 53), hosts file injection, DNS cache flush |
| Blue Pill Lite | `v27/bluepill.go` | VMXON/VMCS VT-x hypervisor, PatchGuard bypass, CPUID trap, memory hiding |
| LOLBin Chainer | `lolbin_chainer.go` | 28 LOLBins (20 Windows + 8 Linux), randomized chain per hour, multi-layer base64 encoding |
| Kernel DNS Driver | `v29/wfp_driver.go` | NDIS filter driver, blocks Defender/Security updates, DNS redirect to C2 |

</details>

<details>
<summary><b>🔐 FASE 2 — C2 Hardened</b></summary>

| Module | File | Technique |
|--------|------|-----------|
| SPIFFE mTLS | `crypto/spiffe.go` | SVID generation, trust bundle, peer SPIFFE ID verification, certificate rotation, server/client mTLS |
| Multi-Channel C2 | `multi_channel_c2.go` | 5 channels: gRPC→WebSocket→DoH→Twitter→Blockchain, health check, auto-failover, beacon loop |
| Ed25519 Signing | `crypto/ed25519.go` | Command sign/verify, nonce replay protection, trusted key ring, batch sign/verify operations |
| Kyber-1024+X25519 | `crypto/kyber.go` | Hybrid KEM (ML-KEM-1024 + X25519), HKDF derivation, AES-256-GCM + HMAC-SHA256 sessions |
| Dashboard Ops | `c2/dashboard.go` | HTTP+WebSocket API, agent nodes, propagation map, signed command issuance, embedded HTML dashboard |
| Proto Obfuscation | `proto/loader/obfuscated.go` | XOR+AES-CTR+GZIP, integrity verification, vaporize buffers, memory-only loading |

</details>

<details>
<summary><b>🌊 FASE 3 — Advanced Propagation</b></summary>

| Module | File | Technique |
|--------|------|-----------|
| Ultrasound QPSK | `hydra_vectors/ultrasound.go` | >18kHz QPSK modulation, WAV generation, speaker/mic RX/TX, preamble synchronization |
| Powerline PLC | `hydra_vectors/powerline.go` | HomePlug device scan, UPnP SSDP discovery, SOAP injection over electrical grid |
| USB ADB Worm | `hydra_vectors/usb_adb.go` | ADB enumeration, APK install, remote shell exec, SMS/contacts dump, sdcard persistence |
| DNS Rebinding | `hydra_vectors/dns_rebinding.go` | TTL=0 rebind server, SOP bypass JavaScript payload, SSRF lateral via Host headers |
| CI/CD Webhooks | `hydra_vectors/cicd_webhooks.go` | GitHub Actions workflow injection, Jenkins/GitLab CI trigger, 10 CI environment scanner |
| VLAN Jump | `hydra_vectors/vlan_jump.go` | Double tagging, DTP negotiation, ARP flood, DHCP discover per VLAN interface |
| QR Dynamic Worm | `hydra_vectors/qr_worm.go` | QR matrix generation (finder/timing/alignment patterns), PNG rendering, rotation channel |
| PJL Printer Worm | `hydra_vectors/pjl_worm.go` | Printer Job Language exploits, NVRAM R/W, firmware infection, PCL ransom note |
| Chronos NTP | `chronos_ntp.go` | Fake NTP server, time forward/rewind, scheduled task shift, w32tm hijack |
| Reflective DLL | `stager/reflective_asm.go` | NtCreateSection+NtMapViewOfSection, 100-byte NASM stager, CreateRemoteThread |
| Kerberos Delegation | `kerberos_delegation.go` | Unconstrained delegation discovery, coercion (printerbug/PetitPotam), TGT dump, Silver Ticket |
| IMDSv2 Bypass | `imdsv2_bypass.go` | AWS IMDSv2 token acquisition, IAM credential extraction, SSRF exploit, AssumeRole |

</details>

<details>
<summary><b>🧠 FASE 4 — AI + Cross-Platform + Lab</b></summary>

| Module | File | Technique |
|--------|------|-----------|
| Cross-Platform Loader | `loader/cross.go` | ELF/Mach-O/APK generation, pack+encrypt, syscall hooks, anti-sandbox detection |
| JIT Polymorphism | `jit_polymorphism.go` | NOP-sleds, constant obfuscation, code crossover, register reordering, runtime mutation loop |
| AI FSM Orchestrator | `appstate/ai_orchestrator.go` | Q-learning state machine, exploration vs exploitation, risk prediction, Q-table import/export |
| Federated Learning | `ai/hivemind/federated.go` | FedAvg aggregation, victim profiling, optimal phishing time prediction, model export |
| Autofactory Fuzzer | `ai/autofactory/fuzzer.go` | AFL++ integration, 9 mutation strategies (bit/byte flip, arithmetic, splice, etc), crash detection |
| Wazero Bridge | `bridge/wazero_loader.go` | WASM module parsing, TinyGo compilation, Python handler migration to WASM, stub generation |
| RF Contagion | `rf_contagion/baseband.go` | SDR detection (RTL-SDR/HackRF), ModemManager, baseband injection, IMSI capture, SS7 |
| EDR Test Lab | `lab/docker-compose.edr.yml` | Windows Server 2022 + ELK Stack + Sysmon + automated evasion test suite |
| Deepfake Vishing | `ai/vishing/deepfake.go` | Coqui TTS voice cloning (tacotron2-DDC), VoIP call placement, SMS phishing, SE profiling |

</details>

---

## Quick Start

```bash
# Clone the repository
git clone https://github.com/Ruby570bocadito/X404X.git
cd X404X

# Install Python dependencies
pip install -r requirements.txt

# Install Go modules
cd internal/ransomware && go mod tidy && cd ../..

# Full demo (dry-run — no real exploits)
bash scripts/run_demo.sh

# Worm in simulation mode (Windows)
cd plugins/worm && python worm_core.py --config configs/config_simulation.yaml

# Worm in simulation mode (Linux)
cd plugins/worm && python3 worm_core.py --config configs/config_simulation.yaml

# Dashboard (port 9090)
cd plugins/pulse-c2/src/go && go run ./cmd/dashboard -port 9090

# EDR test lab
docker-compose -f lab/docker-compose.edr.yml up -d
```

---

## Directory Structure

```
X404X/
├── cmd/                           # CLI + implant agent
│   ├── x404x/                     # Interactive shell (~400 lines)
│   └── implant/                   # Go C2 agent
│
├── internal/                      # CORE ENGINE (Go)
│   ├── ransomware/                # 37 files — ransomware + evasion + propagation
│   │   ├── hydra_vectors/         # 8 exotic worm vectors
│   │   ├── loader/                # Cross-platform ELF/Mach-O/APK
│   │   ├── stager/                # Reflective DLL injection
│   │   ├── v27/, v29/, v210/      # Advanced version modules
│   │   └── *.go                   # Core engine files
│   ├── agent/                     # Post-exploitation + priv esc
│   ├── appstate/                  # FSM + AI orchestrator
│   ├── bridge/                    # Wazero WASM bridge loader
│   └── dispatch/                  # Module dispatcher + MITRE mapping
│
├── modules/bridge/                # PYTHON BRIDGE
│   ├── bridge.py                  # RPC router (Go ↔ Python)
│   └── handlers/                  # 12 files, ~170 RPC handlers
│       ├── ransomware*.py         # v26-v210, BlockZ
│       ├── attacks.py, bloodhound.py, cred_dump.py
│       └── phase_1_4.py           # 36 Phase 1-4 handlers
│
├── plugins/                       # MODULAR PLUGINS
│   ├── worm/                      # Worm with RL (35 exploits)
│   ├── operations/                # Argos cell/stager agents
│   ├── blue/bluesky/              # Bluetooth attacks
│   ├── pulse-c2/                  # C2 with crypto + dashboard
│   ├── ai/                        # AI plugins
│   │   ├── hivemind/              # Federated learning
│   │   ├── autofactory/           # AFL++ fuzzer
│   │   └── vishing/               # Deepfake voice phishing
│   └── rf_contagion/              # RF SDR 4G/5G
│
├── lab/                           # EDR test environment
├── pkg/proto/                     # gRPC protobuf definitions
├── scripts/run_demo.sh            # Cross-platform demo
├── deploy/deploy.sh               # One-click deployment
├── ROADMAP.md                     # Complete development roadmap
└── requirements.txt               # Python dependencies
```

---

## Tech Stack

<p align="center">
  <img src="https://skillicons.dev/icons?i=go,python,c,bash,powershell,docker,linux,git,github,vue,nodejs&theme=dark" alt="Core Stack"/>
</p>

<p align="center">
  <img src="https://skillicons.dev/icons?i=nmap,metasploit&theme=dark" alt="Security Tools"/>
</p>

<p align="center">
  <sub>
    <b>Go:</b> SPIFFE · Ed25519 · Kyber-1024 · gRPC · Wazero · TinyGo · AFL++ ·
    <b>Python:</b> cryptography · Rich · fpdf2 · python-docx · Coqui TTS · pyVoIP ·
    <b>Security:</b> Nmap · Metasploit · BloodHound · Mimikatz · Impacket · CrackMapExec
  </sub>
</p>

---

## Development Roadmap

See [`ROADMAP.md`](ROADMAP.md) for the complete development plan including:

- **Future Phases:** Enhanced AI capabilities, hardware-level attacks, satellite/SCADA integration
- **Module Specifications:** Detailed descriptions of all 45 modules
- **Integration Plan:** 12 proposed new modules (Kernel DNS WFP, Hypervisor stealth, etc.)
- **Code Conventions:** Go/Python style guides and naming conventions

---

## Legal Disclaimer

This project is developed for **educational purposes and authorized security testing only**. Unauthorized use against systems you do not own or have explicit permission to test is illegal. The authors assume no liability for misuse.

**Use only in:**
- Authorized penetration tests with written consent
- Academic research environments
- Your own isolated lab infrastructure
- CTF competitions

---

<p align="center">
  <a href="https://github.com/Ruby570bocadito">
    <img src="https://img.shields.io/badge/GitHub-Ruby570bocadito-181717?style=for-the-badge&logo=github&logoColor=white" alt="GitHub"/>
  </a>
  <a href="https://www.linkedin.com/in/rafael-g%C3%A1lvez-silipo-07445a409/">
    <img src="https://img.shields.io/badge/LinkedIn-Rafael_Gálvez-0A66C2?style=for-the-badge&logo=linkedin&logoColor=white" alt="LinkedIn"/>
  </a>
</p>

<p align="center">
  <sub>Built with ❤️ by <a href="https://github.com/Ruby570bocadito">RBYHACK</a> · Málaga, Spain · 2025-2026</sub>
</p>

<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=waving&color=0:ff3366,100:00d4ff&height=120&section=footer&text=X404X+%7C+Autonomous+Red+Team+Platform&fontSize=14&fontColor=ffffff&animation=fadeIn" alt="X404X Footer"/>
</p>
