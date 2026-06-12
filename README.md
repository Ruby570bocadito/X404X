<p align="center">
  <img src="https://readme-typing-svg.herokuapp.com?font=Fira+Code&weight=700&size=48&duration=2000&pause=500&color=ff3366&center=true&vCenter=true&width=800&lines=X+4+0+4+X" alt="X404X" />
</p>

<p align="center">
  <img src="https://readme-typing-svg.herokuapp.com?font=Fira+Code&weight=500&size=16&duration=2400&pause=800&color=00d4ff&center=true&vCenter=true&width=720&lines=%C2%BB+Autonomous+Red+Team+Platform+%C2%AB;%C2%BB+Full-spectrum+offensive+operations+%C2%AB;%C2%BB+From+recon+to+post-exploitation+%C2%AB;%C2%BB+Go+%E2%80%A2+Python+%E2%80%A2+45+M%C3%B3dulos+%C2%AB" />
</p>

<br/>

<p align="center">
  <a href="#-what-is-x404x"><img src="https://img.shields.io/badge/About-ff3366?style=for-the-badge&logo=readthedocs&logoColor=white&labelColor=1a1a2e" /></a>
  <a href="#-architecture"><img src="https://img.shields.io/badge/Architecture-00d4ff?style=for-the-badge&logo=diagramsdotnet&logoColor=white&labelColor=1a1a2e" /></a>
  <a href="#-modules"><img src="https://img.shields.io/badge/Modules-00ff88?style=for-the-badge&logo=target&logoColor=white&labelColor=1a1a2e" /></a>
  <a href="#-quick-start"><img src="https://img.shields.io/badge/QuickStart-ff6b35?style=for-the-badge&logo=rocket&logoColor=white&labelColor=1a1a2e" /></a>
</p>

<br/>

---

<br/>

<p align="center">
  <img src="https://img.shields.io/badge/Go-v1.22+-00ADD8?style=flat-square&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Python-v3.11+-3776AB?style=flat-square&logo=python&logoColor=white" />
  <img src="https://img.shields.io/badge/Modules-45-ff3366?style=flat-square&logo=target&logoColor=white" />
  <img src="https://img.shields.io/badge/Lines-12,685-00d4ff?style=flat-square&logo=codecrafters&logoColor=white" />
  <img src="https://img.shields.io/badge/Phases-4-00ff88?style=flat-square&logo=checkmarx&logoColor=white" />
</p>

<p align="center">
  <a href="https://github.com/Ruby570bocadito/X404X/stargazers"><img src="https://img.shields.io/github/stars/Ruby570bocadito/X404X?style=social&logo=github" /></a>
  <a href="https://github.com/Ruby570bocadito/X404X/network"><img src="https://img.shields.io/github/forks/Ruby570bocadito/X404X?style=social&logo=github" /></a>
</p>

<br/>

<table align="center">
<tr>
<td align="center" width="25%">
  <h3>🛡️ KERNEL<br/>EVASION</h3>
  <sub>BYOVD · DKOM · Blue Pill</sub>
</td>
<td align="center" width="25%">
  <h3>🔐 POST-QUANTUM<br/>C2</h3>
  <sub>Kyber-1024 · Ed25519 · SPIFFE</sub>
</td>
<td align="center" width="25%">
  <h3>🦠 EXOTIC<br/>PROPAGATION</h3>
  <sub>Ultrasound · PLC · QR · PJL</sub>
</td>
<td align="center" width="25%">
  <h3>🧠 AI<br/>ORCHESTRATOR</h3>
  <sub>Q-Learning · FedAvg · Deepfake</sub>
</td>
</tr>
</table>

<br/>

## What is X404X?

> **X404X** is an **autonomous red team operations suite** that covers the complete offensive kill chain. From kernel-level evasion to post-quantum command & control — every module is implemented in real code, not simulated stubs.

Built as a **Go** core with a **Python** bridge layer and **WASM** extension support, X404X integrates 45 offensive modules across 4 development phases — all compiled, tested, and ready for authorized red team engagements.

---

## Architecture

<br/>

```mermaid
flowchart LR
    subgraph OPERATOR[" Operator Console"]
        CLI[Go CLI Shell]
        DASH[Vue 3 Dashboard]
    end

    subgraph C2[" C2 Infrastructure"]
        MTLS[SPIFFE mTLS]
        SIGN[Ed25519 Signer]
        KYBER[Kyber-1024 KEM]
        MULTI[5-Channel Stack]
    end

    subgraph CORE[" Core Engine Go"]
        RANSOM[Ransomware Engine]
        EVASION[Evasion Suite]
        PROPAG[Propagation Vectors]
        AI[AI Orchestrator]
    end

    subgraph BRIDGE[" Python Bridge"]
        RPC[Go↔Python RPC Router]
        H170[~170 Handlers]
    end

    subgraph PLUGINS[" Plugin Ecosystem"]
        WORM[Worm + RL]
        ARGOS[Argos Operations]
        BLUE[Bluesky BT]
        PULSE[Pulse C2]
        H_MIND[Hivemind AI]
        RF[RF Contagion SDR]
    end

    CLI --> MTLS
    DASH --> MTLS
    MTLS --> RANSOM
    MTLS --> EVASION
    MTLS --> PROPAG
    MTLS --> AI
    RANSOM --> RPC
    EVASION --> RPC
    PROPAG --> RPC
    AI --> RPC
    RPC --> WORM
    RPC --> ARGOS
    RPC --> BLUE
    RPC --> PULSE
    RPC --> H_MIND
    RPC --> RF

    style OPERATOR fill:#ff336620,stroke:#ff3366
    style C2 fill:#00d4ff20,stroke:#00d4ff
    style CORE fill:#00ff8820,stroke:#00ff88
    style BRIDGE fill:#ff6b3520,stroke:#ff6b35
    style PLUGINS fill:#a855f720,stroke:#a855f7
```

<br/>

---

## Kill Chain Coverage

<br/>

```
  ╔══════════════════════════════════════════════════════════════════════╗
  ║                        ATTACK KILL CHAIN                           ║
  ╠═══════╦══════════╦══════════╦══════════╦══════════╦════════════════╣
  ║ RECON ║ INITIAL  ║  EXEC    ║ PERSIST  ║ PRIVESC  ║   LATERAL     ║
  ║       ║ ACCESS   ║          ║          ║          ║                ║
  ╠═══════╬══════════╬══════════╬══════════╬══════════╬════════════════╣
  ║ OSINT ║ Phish AI ║ LOLBin   ║ WER      ║ BYOVD    ║ Kerberos      ║
  ║ Recon ║ CI/CD    ║ Chainer  ║ Triple   ║ DKOM     ║ Delegation    ║
  ║ DNS   ║ QR Worm  ║ Reflect  ║ MFT      ║ Token    ║ IMDSv2 (AWS)  ║
  ║ APIs  ║ USB ADB  ║ DLL      ║ Slack    ║ Steal    ║ VLAN Jump     ║
  ║       ║ PJL      ║ WASM     ║ Schtasks ║          ║ Chronos NTP   ║
  ╚═══════╩══════════╩══════════╩══════════╩══════════╩════════════════╝
                  │                          │
  ╔═══════════════╩══════════╦═══════════════╩═════════════╗
  ║          EVASION         ║         C2 & EXFIL          ║
  ╠══════════════════════════╬═════════════════════════════╣
  ║ WFP DNS Poisoning        ║ SPIFFE mTLS + Ed25519      ║
  ║ Blue Pill Hypervisor     ║ Kyber-1024 Post-Quantum    ║
  ║ Anti-Reversing Suite     ║ 5-Channel Stack            ║
  ║ Anti-Forensics (DoD 7p)  ║ Blockchain C2 (BTC/ETH)    ║
  ║ MFT Slack Hide           ║ MFT Slack Storage          ║
  ╚══════════════════════════╩═════════════════════════════╝
```

<br/>

---

## Phase Completion

<br/>

```
FASE 0  ████████████████████████  100%  Critical Stubs Fixed
FASE 1  ████████████████████████  100%  Evasion + Anti-Forensics
FASE 2  ████████████████████████  100%  C2 Hardened
FASE 3  ████████████████████████  100%  Advanced Propagation
FASE 4  ████████████████████████  100%  AI + Cross-Platform
        └──────────────────────┘
        45 Modules · 12,685 Lines
```

<br/>

| Phase | Theme | 🧩 | 📝 Lines | Status |
|:-----:|-------|:--:|---------|:------:|
| **0** | Critical Stubs | 8 | 500 | `█████████░` 100% |
| **1** | Evasion + Anti-Forensics | 10 | 3,599 | `█████████░` 100% |
| **2** | C2 Hardened | 6 | 2,374 | `█████████░` 100% |
| **3** | Advanced Propagation | 12 | 3,274 | `█████████░` 100% |
| **4** | AI + Cross-Platform | 9 | 2,938 | `█████████░` 100% |

<br/>

---

## Module Matrix

<br/>

### 🔴 FASE 1 — Evasion & Anti-Forensics

| # | Module | Capability | OS |
|:--:|--------|------------|:--:|
| 1.1 | **BYOVD Loader** | 5 vulnerable drivers (WinRing0, Gdrv, RTCore64, kprocesshacker, CPUID) · IOCTLs · R/W physical memory · MSR · handle elevation | ![Win](https://img.shields.io/badge/Win-0078D6?style=flat-square&logo=windows&logoColor=white) |
| 1.2 | **DKOM** | Process hiding via ActiveProcessLinks unlink · SYSTEM token steal · EPROCESS offsets per build | ![Win](https://img.shields.io/badge/Win-0078D6?style=flat-square&logo=windows&logoColor=white) |
| 1.3 | **Anti-Reversing** | HW BP detect (DR0-DR7) · INT3 scan · CRC32 integrity · RDTSC timing · sandbox + virtual MAC check | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 1.4 | **Anti-Forensics** | DoD 5220.22-M 7-pass wipe · MFT $BITMAP corruption · VAD hide · crash dump/event log/prefetch/USN/Shellbag wipe | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 1.5 | **WER Persistence** | Windows Error Reporting Hangs hijack · Silent Process Exit · startup + Run key + schtasks | ![Win](https://img.shields.io/badge/Win-0078D6?style=flat-square&logo=windows&logoColor=white) |
| 1.6 | **MFT Slack** | NTFS slack space R/W via PowerShell · AES-GCM encrypted fragments · hidden agent/ransom note storage | ![Win](https://img.shields.io/badge/Win-0078D6?style=flat-square&logo=windows&logoColor=white) |
| 1.7 | **WFP DNS Poison** | WFP provider + netsh fallback · fake DNS server (UDP 53) · hosts file injection · cache flush | ![Win](https://img.shields.io/badge/Win-0078D6?style=flat-square&logo=windows&logoColor=white) |
| 1.8 | **Blue Pill HV** | VMXON/VMCS VT-x hypervisor · PatchGuard bypass · CPUID trap · memory hiding | ![Win](https://img.shields.io/badge/Win-0078D6?style=flat-square&logo=windows&logoColor=white) |
| 1.9 | **LOLBin Chainer** | 28 LOLBins (20 Win + 8 Linux) · randomized chain per hour · multi-layer base64 encoding | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 1.10 | **Kernel DNS Driver** | NDIS filter driver · blocks Defender/Security updates · DNS redirect to C2 | ![Win](https://img.shields.io/badge/Win-0078D6?style=flat-square&logo=windows&logoColor=white) |

### 🔵 FASE 2 — C2 Hardened

| # | Module | Capability | OS |
|:--:|--------|------------|:--:|
| 2.1 | **SPIFFE mTLS** | SVID generation · trust bundle · peer SPIFFE ID verification · certificate rotation · mTLS server/client | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 2.2 | **Multi-Channel C2** | 5 channels (gRPC→WebSocket→DoH→Twitter→Blockchain) · health check · auto-failover · beacon loop | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 2.3 | **Ed25519 Signing** | Command sign/verify · nonce replay protection · trusted key ring · batch operations | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 2.4 | **Dashboard Ops** | HTTP+WebSocket API · agent nodes · propagation map · signed command issuance · embedded HTML | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 2.5 | **Kyber-1024** | Hybrid KEM (ML-KEM-1024 + X25519) · HKDF derivation · AES-256-GCM + HMAC-SHA256 sessions | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 2.6 | **Proto Obfuscation** | XOR+AES-CTR+GZIP · integrity verification · vaporize buffers · memory-only loading | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |

### 🟢 FASE 3 — Advanced Propagation

| # | Module | Capability | OS |
|:--:|--------|------------|:--:|
| 3.1 | **Ultrasound QPSK** | >18kHz modulation · WAV generation · speaker/mic RX/TX · preamble sync | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 3.2 | **Powerline PLC** | HomePlug device scan · UPnP SSDP · SOAP injection over electrical grid | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 3.3 | **USB ADB** | ADB enumeration · APK install · remote shell exec · SMS/contacts dump | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 3.4 | **DNS Rebinding** | TTL=0 rebind server · SOP bypass JS payload · SSRF lateral via Host headers | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 3.5 | **CI/CD Webhooks** | GitHub Actions · Jenkins · GitLab CI injection · 10 CI scanner | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 3.6 | **VLAN Jump** | Double tagging · DTP negotiation · ARP flood · DHCP per VLAN | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 3.7 | **QR Dynamic Worm** | QR matrix generation · PNG rendering · rotation channel | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 3.8 | **PJL Printer Worm** | Printer job language exploits · NVRAM R/W · firmware infection · PCL ransom note | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 3.9 | **Chronos NTP** | Fake NTP server · time forward/rewind · schtask shift · w32tm hijack | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 3.10 | **Reflective DLL** | NtCreateSection+NtMapViewOfSection · 100-byte NASM stager · remote thread injection | ![Win](https://img.shields.io/badge/Win-0078D6?style=flat-square&logo=windows&logoColor=white) |
| 3.11 | **Kerberos Deleg** | Unconstrained delegation discovery · coercion · TGT dump · Silver Ticket | ![Win](https://img.shields.io/badge/Win-0078D6?style=flat-square&logo=windows&logoColor=white) |
| 3.12 | **IMDSv2 Bypass** | AWS token acquisition · IAM extraction · SSRF · AssumeRole · neighbor scan | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |

### 🟣 FASE 4 — AI + Cross-Platform

| # | Module | Capability | OS |
|:--:|--------|------------|:--:|
| 4.1 | **Cross-Platform Loader** | ELF · Mach-O · APK generation · pack+encrypt · syscall hooks | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 4.2 | **JIT Polymorphism** | NOP-sleds · code crossover · register reordering · runtime mutation loop | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 4.3 | **AI FSM Orchestrator** | Q-learning state machine · exploration vs exploitation · risk prediction | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 4.4 | **Federated Learning** | FedAvg aggregation · victim profiling · phishing time prediction · model export | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 4.5 | **Autofactory Fuzzer** | AFL++ integration · 9 mutation strategies · crash detection · exploit candidates | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 4.6 | **Wazero Bridge** | WASM module parsing · TinyGo compilation · Python→WASM migration | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |
| 4.7 | **RF Contagion** | SDR detection (RTL-SDR/HackRF) · ModemManager · baseband injection · IMSI · SS7 | ![Linux](https://img.shields.io/badge/Linux-FCC624?style=flat-square&logo=linux&logoColor=black) |
| 4.8 | **EDR Test Lab** | Docker Compose · Windows Server 2022 · ELK Stack · Sysmon · automated tests | ![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=docker&logoColor=white) |
| 4.9 | **Deepfake Vishing** | Coqui TTS (tacotron2-DDC) · VoIP · SMS phishing · SE profiling | ![All](https://img.shields.io/badge/All-333?style=flat-square&logo=linux&logoColor=white) |

<br/>

---

## 🚀 Quick Start

<br/>

### Prerequisites

```bash
# Required
go >= 1.22        # Core engine (Go)
python >= 3.11    # Bridge + handlers (Python)
git               # Version control

# Optional
docker >= 27.0    # EDR lab + containerized deployment
node >= 22.0      # Dashboard frontend (Vue 3)
```

### One-Click Deploy

```bash
# Full stack deployment (Linux/macOS)
bash deploy/deploy.sh

# This will:
#  1. Install Go & Python dependencies
#  2. Build the C2 server
#  3. Start the Python bridge
#  4. Launch the dashboard on port 9090
#  5. Start the worm engine in dry-run mode
```

### Manual Setup

```bash
# 1. Clone & enter
git clone https://github.com/Ruby570bocadito/X404X.git
cd X404X

# 2. Install dependencies
pip install -r requirements.txt
cd internal/ransomware && go mod tidy && cd ../..

# 3. Start the Python Bridge
cd modules/bridge && python3 bridge.py &
cd ../..

# 4. Build & launch C2 Dashboard
cd plugins/pulse-c2/src/go
go build -o x404x-dashboard ./cmd/dashboard
./x404x-dashboard -port 9090 &
cd ../../..
# → Open http://localhost:9090

# 5. Launch interactive console
cd cmd/x404x
go run . console
# → Type 'help' for available commands
```

### Run Modes

```bash
# ─── Dry-Run Demo (safe — no real exploits) ─────────────
bash scripts/run_demo.sh

# ─── Worm Simulation ────────────────────────────────────
cd plugins/worm
# Windows
python worm_core.py --config configs/config_simulation.yaml
# Linux/macOS
python3 worm_core.py --config configs/config_simulation.yaml

# ─── Live Campaign (⚠️ authorized targets only) ──────────
cd cmd/x404x
go run . campaign start --name "Operation Nightfall" --targets targets.json

# ─── EDR Test Lab ────────────────────────────────────────
docker-compose -f lab/docker-compose.edr.yml up -d
# → Windows Server 2022 + Defender ATP + ELK + Sysmon
# → Automated evasion test suite runs on boot
```

### Console Commands

```bash
# Inside the X404X interactive console:
help              # List all commands
modules           # Show available modules (45 total)
campaign start    # Start a new campaign (FSM-driven)
recon             # Run reconnaissance module
exploit           # Execute exploit chain
privesc           # Privilege escalation
persist           # Install persistence mechanisms
lateral           # Lateral movement
exfil             # Data exfiltration
dashboard         # Open dashboard URL
deploy            # One-click full deploy
status            # Show campaign status
killchain         # Display kill chain progress
ransomware        # Launch ransomware engine
propagate         # Start worm propagation
listeners         # List active C2 listeners
webhook           # Configure webhook notifications
```

<br/>

---

## Directory Map

<br/>

```
X404X/
│
├── cmd/                           ← CLI + Implant Agent
│   ├── x404x/console.go           Interactive shell
│   └── implant/main.go            Go C2 agent
│
├── internal/                      ← CORE ENGINE
│   ├── ransomware/                Engine principal (37 files)
│   │   ├── hydra_vectors/         8 vectores exóticos
│   │   ├── loader/cross.go        ELF · Mach-O · APK
│   │   ├── stager/reflective_asm  Reflective DLL NASM
│   │   ├── v27/ v29/ v210/        Versiones avanzadas
│   │   └── *.go                   Módulos core
│   ├── agent/                     Post-exploit + privesc
│   ├── appstate/                  FSM + AI orchestrator
│   ├── bridge/wazero_loader.go    WASM bridge
│   └── dispatch/dispatcher.go     MITRE ATT&CK mapper
│
├── modules/bridge/                ← PYTHON BRIDGE
│   ├── bridge.py                  Go↔Python RPC router
│   └── handlers/                  12 files · ~170 handlers
│
├── plugins/                       ← PLUGIN ECOSYSTEM
│   ├── worm/                      35 exploits + RL engine
│   ├── operations/                Argos agents + UI
│   ├── pulse-c2/                  Crypto + Dashboard
│   ├── ai/                        Hivemind · Fuzzer · Vishing
│   ├── blue/bluesky/              Bluetooth attacks
│   └── rf_contagion/              SDR 4G/5G baseband
│
├── lab/docker-compose.edr.yml     EDR test environment
├── deployments/deploy.sh          One-click deploy
├── scripts/run_demo.sh            Cross-platform demo
├── ROADMAP.md                     Development roadmap
└── requirements.txt               Dependencies
```

<br/>

---

## Tech Canvas

<p align="center">
  <table>
    <tr>
      <td align="center"><b>LANGUAGE</b></td>
      <td align="center"><b>FRAMEWORK</b></td>
      <td align="center"><b>CRYPTO</b></td>
      <td align="center"><b>TOOL</b></td>
    </tr>
    <tr>
      <td align="center"><img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white"/></td>
      <td align="center"><img src="https://img.shields.io/badge/gRPC-244c5a?style=for-the-badge&logo=google&logoColor=white"/></td>
      <td align="center"><img src="https://img.shields.io/badge/Kyber-PQ-ff3366?style=for-the-badge&logo=letsencrypt&logoColor=white"/></td>
      <td align="center"><img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white"/></td>
    </tr>
    <tr>
      <td align="center"><img src="https://img.shields.io/badge/Python-3776AB?style=for-the-badge&logo=python&logoColor=white"/></td>
      <td align="center"><img src="https://img.shields.io/badge/Vue.js-4FC08D?style=for-the-badge&logo=vuedotjs&logoColor=white"/></td>
      <td align="center"><img src="https://img.shields.io/badge/Ed25519-00d4ff?style=for-the-badge&logo=letsencrypt&logoColor=white"/></td>
      <td align="center"><img src="https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black"/></td>
    </tr>
    <tr>
      <td align="center"><img src="https://img.shields.io/badge/Bash-4EAA25?style=for-the-badge&logo=gnubash&logoColor=white"/></td>
      <td align="center"><img src="https://img.shields.io/badge/WASM-654FF0?style=for-the-badge&logo=webassembly&logoColor=white"/></td>
      <td align="center"><img src="https://img.shields.io/badge/X25519-00ff88?style=for-the-badge&logo=letsencrypt&logoColor=white"/></td>
      <td align="center"><img src="https://img.shields.io/badge/Git-F05032?style=for-the-badge&logo=git&logoColor=white"/></td>
    </tr>
    <tr>
      <td align="center"><img src="https://img.shields.io/badge/PowerShell-5391FE?style=for-the-badge&logo=powershell&logoColor=white"/></td>
      <td align="center"><img src="https://img.shields.io/badge/AFL++-F7931A?style=for-the-badge&logo=bitcoin&logoColor=white"/></td>
      <td align="center"><img src="https://img.shields.io/badge/AES_GCM-ff6b35?style=for-the-badge&logo=letsencrypt&logoColor=white"/></td>
      <td align="center"><img src="https://img.shields.io/badge/GitHub-181717?style=for-the-badge&logo=github&logoColor=white"/></td>
    </tr>
  </table>
</p>

<br/>

---

## Stats

<p align="center">
  <img src="https://img.shields.io/badge/Total_Files-450+-ff3366?style=for-the-badge&logo=files&logoColor=white" />
  <img src="https://img.shields.io/badge/Go_Files-319-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Python_Files-539-3776AB?style=for-the-badge&logo=python&logoColor=white" />
  <img src="https://img.shields.io/badge/Total_Lines-231K-00d4ff?style=for-the-badge&logo=codecrafters&logoColor=white" />
  <img src="https://img.shields.io/badge/Handlers-144-00ff88?style=for-the-badge&logo=target&logoColor=white" />
  <img src="https://img.shields.io/badge/Modules-76+-ff6b35?style=for-the-badge&logo=octopusdeploy&logoColor=white" />
</p>

<br/>

---

## Legal

> **⚠️ DISCLAIMER:** This project exists exclusively for **educational purposes, authorized security research, and red team operations with explicit written consent**. Unauthorized use against systems you do not own or have permission to test is **illegal** and **strictly prohibited**. The authors and contributors assume **zero liability** for any misuse, damage, or legal consequences resulting from improper use of this codebase.
>
> **By using this software, you acknowledge that you are solely responsible for complying with all applicable laws and regulations.**

<br/>

<p align="center">
  <sub>
    <a href="https://github.com/Ruby570bocadito">RBYHACK</a> © 2025-2026 ·
    Built in Málaga, Spain ·
    <a href="https://github.com/Ruby570bocadito/X404X/blob/main/ROADMAP.md">Roadmap</a> ·
    <a href="https://github.com/Ruby570bocadito/X404X/issues">Report Issue</a>
  </sub>
</p>

<br/>
