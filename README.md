<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=waving&color=gradient&customColorList=0,2,3,6,8&height=180&section=header&text=RBYHACK%20Framework&fontSize=50&fontColor=ffffff&animation=fadeIn&desc=Red%20Team%20Autonomous%20Platform&descAlignY=68&descSize=18" />
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Python-3.11+-3776AB?style=for-the-badge&logo=python&logoColor=white" />
  <img src="https://img.shields.io/badge/Vue-3-4FC08D?style=for-the-badge&logo=vue.js&logoColor=white" />
  <img src="https://img.shields.io/badge/C-Kernel-CC0000?style=for-the-badge&logo=c&logoColor=white" />
  <img src="https://img.shields.io/badge/gRPC-Protocol-4285F4?style=for-the-badge&logo=google&logoColor=white" />
  <img src="https://img.shields.io/badge/Docker-Lab-2496ED?style=for-the-badge&logo=docker&logoColor=white" />
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Kill%20Chain-Complete-FF4500?style=flat-square" />
  <img src="https://img.shields.io/badge/AI-Ollama%20Offline-00FF00?style=flat-square" />
  <img src="https://img.shields.io/badge/Crypto-X25519%20%2B%20XChaCha20--Poly1305-6C63FF?style=flat-square" />
  <img src="https://img.shields.io/badge/MITRE-ATT%26CK%20Mapped-FF6B35?style=flat-square" />
  <img src="https://img.shields.io/badge/License-MIT-blue?style=flat-square" />
</p>

---

> ⚠️ **AUTHORIZED USE ONLY** — This framework is exclusively for authorized security assessments, CTF competitions, academic research, and controlled lab environments. Unauthorized use may violate local and international laws. The author is not responsible for misuse. You are responsible for complying with all applicable laws.

---

## What is RBYHACK?

**RBYHACK Framework** is a **semi-autonomous Red Team platform** that covers the complete cyber kill chain — from reconnaissance to exfiltration. It integrates **11 specialized offensive security tools** into a unified monorepo with a shared gRPC protocol, encrypted communication, AI-powered decision making (Ollama offline), kernel-level persistence, and a Vue 3 web dashboard.

Built as a **TFG (Trabajo de Fin de Grado)** project in Cybersecurity at Cisco NetAcad, Málaga.

---

## Architecture

```
                    ┌─────────────────────────────────┐
                    │       USER INTERFACES            │
                    │  CLI (rbyhack)  │  Vue 3 Dashboard│
                    └────────┬───────────────┬─────────┘
                             │               │
                    ┌────────▼───────────────▼─────────┐
                    │          ORCHESTRATOR            │
                    │   Campaign Mgr │ Decision Engine │
                    │   Rules(25%) + A*(35%) + AI(40%) │
                    └────────┬─────────────────────────┘
                             │ gRPC (X25519+XChaCha20)
                    ┌────────▼─────────────────────────┐
                    │           PULSE-C2               │
                    │   Server (Go) │ Dashboard (Vue 3)│
                    └────────┬─────────────────────────┘
                             │ gRPC encrypted
                    ┌────────▼─────────────────────────┐
                    │        UNIFIED AGENT (Go)        │
                    └──┬──────┬───────┬──────┬─────────┘
                       │      │       │      │
                ┌──────┐ ┌────┐ ┌────┐ ┌──────────┐
                │Rise  │ │Vault│ │Breach│ │Python    │
                │Priv  │ │Kernel│ │Entry│ │Bridge    │
                └──────┘ └────┘ └────┘ └─────┬─────┘
                                             │
                    ┌────────────────────────┼──────────┐
                    │                        │          │
               ┌────▼────┐ ┌────────────┐ ┌──▼─────┐ ┌──▼──────────┐
               │Horizon  │ │Wormy-ML    │ │Specter │ │BlueForge    │
               │Intel    │ │(Lateral)   │ │+ Apex  │ │Suite        │
               │(Recon)  │ │            │ │(AI)    │ │(Metrics)    │
               └─────────┘ └────────────┘ └────────┘ └─────────────┘
```

---

## Component Map

| Component | Language | Kill Chain Phase | Status |
|-----------|----------|-----------------|--------|
| **CLI (rbyhack)** | Go | All | NEW |
| **Orchestrator** | Go | All | NEW |
| **Unified Agent** | Go | Execution | NEW |
| **core/crypto** | Go | Shared (X25519+XChaCha20) | NEW |
| **core/proto** | Protobuf | Shared (gRPC) | NEW |
| **Python Bridge** | Python | IPC | NEW |
| **Evasion Module** | Python | Evasion | NEW |
| **Pulse-C2** | Go + Vue 3 | C2 | [Repo](https://github.com/Ruby570bocadito/Pulse-C2) |
| **Rise-Privilege** | Go | PrivEsc | [Repo](https://github.com/Ruby570bocadito/Rise-Privilege) |
| **Vault-Kernel** | C + Go | Persistence | [Repo](https://github.com/Ruby570bocadito/Vault-Kernel) |
| **Breach-Entry** | C + Python | Initial Access | [Repo](https://github.com/Ruby570bocadito/Breach-Entry) |
| **Horizon-Intel** | Python | Recon | [Repo](https://github.com/Ruby570bocadito/Horizon-Intel) |
| **Specter-Terminal** | Python | AI Analysis | [Repo](https://github.com/Ruby570bocadito/Specter-Terminal) |
| **Apex-Automation** | Python | AI Execution | [Repo](https://github.com/Ruby570bocadito/Apex-Automation) |
| **Wormy-ML** | Python | Lateral Movement | [Repo](https://github.com/Ruby570bocadito/Wormy-ML-Network-Worm) |
| **Link-Relay** | Python | C2 Relay | [Repo](https://github.com/Ruby570bocadito/Link-Relay) |
| **Titan-Operations** | Python + Go | Campaign Mgmt | [Repo](https://github.com/Ruby570bocadito/Titan-Operations) |
| **BlueForge-Suite** | Python | Defense/Detection | [Repo](https://github.com/Ruby570bocadito/BlueForge-Suite) |

---

## Quick Start

### Prerequisites

- **Go** 1.22+
- **Python** 3.11+
- **Node.js** 18+ & npm
- **Docker** & Docker Compose

### Clone with Submodules

```bash
git clone --recurse-submodules https://github.com/Ruby570bocadito/RBYHACK-Framework.git
cd RBYHACK-Framework
```

### Setup

```bash
# Automated setup (Go + Python + Node deps)
make setup

# Or manual
bash scripts/setup.sh
```

### Docker Lab

```bash
# Start isolated lab environment
make lab-up
# → Attacker:  docker exec -it rbyhack-attacker bash
# → Target 1:  docker exec -it rbyhack-target1 bash
# → Dashboard: http://localhost:3000

# Stop lab
make lab-down
```

### Build & Test

```bash
make build        # Build all components
make test         # Run all tests
make lint         # Run linters
```

---

## Kill Chain Flow

```
1. HORIZON-INTEL maps the target attack surface
        │
2. BREACH-ENTRY obtains initial access (CVE-2026-XXXX)
        │
3. AGENT deploys on target, checks in to PULSE-C2
        │
4. SPECTER-TERMINAL analyzes context (OS, user, privileges)
        │
5. APEX-AUTOMATION + DECISION ENGINE pick next move
        │
6. RISE-PRIVILEGE finds escalation vector → auto-root
        │
7. VAULT-KERNEL loads as LKM → kernel-level persistence
        │
8. WORMY-ML propagates to other network hosts
        │
9. LINK-RELAY chains C2 communication for evasion
        │
10. BLUEFORGE-SUITE validates what was (and wasn't) detected
```

---

## Project Structure

```
X404X/
├── core/
│   ├── agent/          # Unified Go implant (NEW)
│   │   └── cmd/agent/  # Agent entrypoint
│   ├── crypto/         # Shared crypto: X25519 + XChaCha20-Poly1305 (NEW)
│   ├── proto/          # gRPC definitions: agent, c2, bridge, common (NEW)
│   ├── orchestrator/   # Central coordination engine (NEW)
│   ├── c2/             # Pulse-C2 (submodule)
│   ├── privesc/        # Rise-Privilege (submodule)
│   ├── kernel/         # Vault-Kernel (submodule)
│   └── breach/         # Breach-Entry (submodule)
├── modules/
│   ├── bridge/         # Python-Go IPC bridge (NEW)
│   ├── evasion/        # AV/EDR bypass unified (NEW)
│   ├── recon/          # Horizon-Intel (submodule)
│   ├── ai/
│   │   ├── specter/    # Specter-Terminal (submodule)
│   │   └── apex/       # Apex-Automation (submodule)
│   ├── worm/           # Wormy-ML (submodule)
│   ├── relay/          # Link-Relay (submodule)
│   ├── operations/     # Titan-Operations (submodule)
│   └── blue/           # BlueForge-Suite (submodule)
├── shared/
│   ├── config/         # Central YAML configuration (NEW)
│   ├── logger/         # Structured logging (NEW)
│   ├── types/          # Shared domain types (NEW)
│   └── database/       # SQLAlchemy models (NEW)
├── docs/               # Architecture, roadmap, CLI reference
├── scripts/            # Setup, deployment scripts
├── lab/                # Docker lab environment
│   ├── docker-compose.yml
│   ├── Dockerfile.attacker
│   └── Dockerfile.target-linux
├── .github/workflows/  # CI/CD pipelines
├── .gitmodules         # Git submodule definitions
├── go.work             # Go workspace
├── config.yaml         # Default configuration
├── Makefile            # Build orchestration
└── README.md           # This file
```

---

## CLI Reference

```bash
rbyhack campaign start   -t 10.0.0.0/24 -g domain_admin -p stealth
rbyhack recon scan       <target> --stealth
rbyhack agent list       --status online
rbyhack ai chat                       # Interactive AI assistant
rbyhack exploit run      --risk safe
rbyhack lateral propagate --method smb
rbyhack persistence kernel load
rbyhack dashboard start  --port 3000
rbyhack lab up           --scenario ctf_basic
```

Full CLI reference: [docs/CLI_REFERENCE.md](docs/CLI_REFERENCE.md)

---

## Cryptography

| Component | Algorithm | Purpose |
|-----------|-----------|---------|
| Key Exchange | X25519 ECDH | Ephemeral session keys |
| Encryption | XChaCha20-Poly1305 | Authenticated symmetric encryption |
| Transport | TLS 1.3 (mTLS) | Service-to-service auth |
| Nonce | 192-bit random | Per-message uniqueness |

---

## AI Integration

- **Ollama** — Local LLM, fully offline, no data exfiltration
- **Specter-Terminal** — Offensive security context analysis
- **Apex-Automation** — Autonomous decision making and module orchestration
- **Decision Engine** — Weighted fusion: Rules (25%) + A* Planner (35%) + AI (40%)
- **HITL** — Human-in-the-Loop mode for manual approval (default)

---

## Safety Controls

| Control | Description | Default |
|---------|-------------|---------|
| Kill Switch | Emergency stop all agents | Enabled |
| Geofencing | RFC 1918 private networks only | Enabled |
| Auto-Destruct | Self-terminate after N hours | 2h |
| Max Infections | Stop after N compromised hosts | 1000 |
| No Persistence | Survive reboot = false | Enabled |

---

## TFG — Trabajo de Fin de Grado

This framework is the core project of the TFG in Cybersecurity. The technical memory documents:
- Each component and its integration
- Laboratory testing methodology
- Ethical and legal analysis
- Defense metrics (BlueForge-Suite validation)

> **Author:** Rafael Gálvez — [@Ruby570bocadito](https://github.com/Ruby570bocadito)
> **Center:** Cisco NetAcad · Málaga, Spain

---

## License

MIT License — see [LICENSE](LICENSE) for details.

---

<p align="center">
  <b>Built with passion for offensive security research</b><br/>
  <sub>"The best way to understand defense is to build the attack."</sub>
</p>
