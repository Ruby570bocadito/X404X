# X404X — Architecture Documentation

## System Overview

RBYHACK Framework is a semi-autonomous Red Team platform covering the complete cyber kill chain: reconnaissance → weaponization → delivery → exploitation → installation → C2 → actions on objective → exfiltration.

### Three-Layer Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│ LAYER 1: INTERFACES                                                      │
│ ┌──────────────┐ ┌──────────────────┐ ┌──────────────┐                  │
│ │ rbyhack CLI  │ │ Web Dashboard    │ │ gRPC/REST API│                  │
│ │ (Go + Cobra) │ │ (Vue 3 + WS)     │ │ (external)   │                  │
│ └──────┬───────┘ └────────┬─────────┘ └──────┬───────┘                  │
│        │                  │                   │                          │
├────────┼──────────────────┼───────────────────┼──────────────────────────┤
│ LAYER 2: ORCHESTRATOR (Go)                                                │
│ ┌──────┴──────────────────┴───────────────────┴──────────────────────┐  │
│ │  Campaign Manager  │  Decision Engine  │  Event Bus  │  Database   │  │
│ │  (Titan-Ops)       │  (A* + CBR + AI)  │  (pub/sub)  │  (PG/SQLite)│  │
│ └────────────────────────────────┬────────────────────────────────────┘  │
│                                  │ gRPC (X25519 + XChaCha20-Poly1305)    │
├──────────────────────────────────┼───────────────────────────────────────┤
│ LAYER 3: FIELD AGENTS (Go + Python)                                       │
│ ┌────────────────────────────────┴────────────────────────────────────┐  │
│ │                    UNIFIED AGENT (Go)                                 │  │
│ │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐  │  │
│ │  │Recon     │ │Privilege │ │Kernel    │ │Python    │ │Evasion   │  │  │
│ │  │Scanner   │ │Escalation│ │Rootkit   │ │Bridge    │ │Engine    │  │  │
│ │  └──────────┘ └──────────┘ └─────┬────┘ └────┬─────┘ └──────────┘  │  │
│ │                                  │           │                       │  │
│ │                          ┌───────┴───┐ ┌─────┴──────────┐          │  │
│ │                          │Vault-Kernel│ │Python Modules  │          │  │
│ │                          │(C LKM)    │ │Horizon-Intel   │          │  │
│ │                          └───────────┘ │Specter-Terminal│          │  │
│ │                                        │Apex-Automation │          │  │
│ │                                        │Wormy-ML        │          │  │
│ │                                        │Link-Relay      │          │  │
│ │                                        │BlueForge-Suite │          │  │
│ │                                        └────────────────┘          │  │
│ └────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

## Component Map

| Component | Language | Layer | Phase | GitHub Repo |
|-----------|----------|-------|-------|-------------|
| **CLI (rbyhack)** | Go | 1 | All | NEW |
| **Web Dashboard** | Vue 3 | 1 | All | Extended from Pulse-C2 |
| **Orchestrator** | Go | 2 | All | NEW |
| **Database** | SQLAlchemy | 2 | All | NEW |
| **Unified Agent** | Go | 3 | All | NEW |
| **core/crypto** | Go | Shared | All | NEW |
| **core/proto** | Protobuf | Shared | All | NEW |
| **Pulse-C2** | Go + Vue 3 | 2,3 | C2 | Ruby570bocadito/Pulse-C2 |
| **Rise-Privilege** | Go | 3 | Exploitation | Ruby570bocadito/Rise-Privilege |
| **Vault-Kernel** | C + Go | 3 | Installation | Ruby570bocadito/Vault-Kernel |
| **Breach-Entry** | C + Python | 3 | Delivery | Ruby570bocadito/Breach-Entry |
| **Horizon-Intel** | Python | 3 | Recon | Ruby570bocadito/Horizon-Intel |
| **Specter-Terminal** | Python | 3 | C2 (AI) | Ruby570bocadito/Specter-Terminal |
| **Apex-Automation** | Python | 3 | C2 (AI) | Ruby570bocadito/Apex-Automation |
| **Wormy-ML** | Python | 3 | Lateral Movement | Ruby570bocadito/Wormy-ML-Network-Worm |
| **Link-Relay** | Python | 3 | C2 Relay | Ruby570bocadito/Link-Relay |
| **Titan-Operations** | Python + Go | 2 | Campaign Mgmt | Ruby570bocadito/Titan-Operations |
| **BlueForge-Suite** | Python | 3 | Metrics | Ruby570bocadito/BlueForge-Suite |
| **Python Bridge** | Python | 3 | IPC | NEW |
| **Evasion Module** | Python + Go | 3 | Evasion | NEW |

## Communication Flow

```
                    ┌──────────────────┐
                    │   Orchestrator   │
                    └────────┬─────────┘
                             │ gRPC
                    ┌────────▼─────────┐
                    │    Pulse-C2      │
                    │    Server        │
                    └────────┬─────────┘
                             │ gRPC (X25519 + XChaCha20-Poly1305)
                    ┌────────▼─────────┐
                    │  Unified Agent   │
                    └───┬──┬──┬──┬────┘
                        │  │  │  │
              ┌─────────┘  │  │  └────────┐
              ▼            ▼  ▼           ▼
        ┌──────────┐ ┌──────────┐  ┌──────────────┐
        │ Go Module│ │ C Module │  │Python Bridge │
        │ (Rise,   │ │ (Vault,  │  │(Unix Socket) │
        │  Crypto) │ │  Breach) │  └──────┬───────┘
        └──────────┘ └──────────┘         │
                               ┌──────────┼──────────┐
                               ▼          ▼          ▼
                          ┌────────┐ ┌────────┐ ┌────────┐
                          │Horizon │ │Specter │ │ Wormy  │
                          │Intel   │ │+ Apex  │ │ ML     │
                          └────────┘ └────────┘ └────────┘
```

## Decision Engine Flow

```
                    ┌─────────────────────────────┐
                    │       RECON DATA IN         │
                    │  (Horizon-Intel, Nmap, OSINT)│
                    └─────────────┬───────────────┘
                                  │
            ┌─────────────────────┼─────────────────────┐
            ▼                     ▼                     ▼
   ┌────────────────┐    ┌────────────────┐    ┌────────────────┐
   │ Rules Engine   │    │ A* Planner     │    │ AI Engine      │
   │ (Deterministic)│    │ (Pathfinding)  │    │ (Specter+Apex) │
   │ Weight: 25%    │    │ Weight: 35%    │    │ Weight: 40%    │
   └───────┬────────┘    └───────┬────────┘    └───────┬────────┘
           │                     │                     │
           └─────────────────────┼─────────────────────┘
                                 │
                      ┌──────────▼──────────┐
                      │  DECISION FUSION    │
                      │  (Weighted Rank)    │
                      └──────────┬──────────┘
                                 │
                      ┌──────────▼──────────┐
                      │  HUMAN-IN-THE-LOOP  │
                      │  (Approve / Reject) │
                      └──────────┬──────────┘
                                 │
                      ┌──────────▼──────────┐
                      │  EXECUTE ACTION     │
                      │  (via Agent → C2)   │
                      └─────────────────────┘
```

## Database Schema

See `shared/database/models.py` for the complete SQLAlchemy schema.

Key entities:
- **campaigns**: Red team operations
- **agents**: Implant tracking
- **targets**: Discovered hosts
- **vulnerabilities**: CVEs and misconfigurations
- **credentials**: Captured passwords/hashes
- **kill_chain**: Tactical log
- **decisions**: AI/Rules suggestions
- **audit_log**: Complete action trail
- **ai_analysis**: LLM interaction history
- **blue_metrics**: Detection validation

## Security Model

1. **End-to-End Encryption**: X25519 ECDH + XChaCha20-Poly1305 AEAD per session
2. **gRPC mTLS**: All service-to-service communication authenticated
3. **Offline AI**: Ollama runs locally, no data leaves the lab
4. **Safety Controls**: Kill switch, geofencing, auto-destruct, max infections
5. **Audit Trail**: Every action logged with timestamp, agent, campaign, and result
