# X404X — Architecture Documentation

## System Overview

X404X is a semi-autonomous Red Team platform covering the complete cyber kill chain: reconnaissance → weaponization → delivery → exploitation → installation → C2 → actions on objective → exfiltration.

### Updated Architecture (v2.3)

```
                     ┌─────────────────────────────────┐
                     │        USER INTERFACES           │
                     │ CLI · Console · TUI · Dashboard  │
                     └────────┬───────────────┬─────────┘
                              │               │
                     ┌────────▼───────────────▼─────────┐
                     │          ORCHESTRATOR             │
                     │  Campaign Mgr │ Decision Engine   │
                     │  Rules(25%) + A*(35%) + AI(40%)  │
                     │  WorldGraph · EventBus · KillChain│
                     └────────┬─────────────────────────┘
                              │ gRPC (X25519+XChaCha20)
                     ┌────────▼─────────────────────────┐
                     │        C2 SERVER (Go/gRPC)       │
                     │  AgentService · C2Service        │
                     │  CheckIn · CommandStream · Heart  │
                     └────────┬─────────────────────────┘
                              │ gRPC encrypted
                     ┌────────▼─────────────────────────┐
                     │        UNIFIED AGENT (Go)        │
                     │  Connector · BridgeClient        │
                     └──┬──────┬───────┬──────┬─────────┘
                        │      │       │      │
                 ┌──────┐ ┌────┐ ┌────┐ ┌──────────┐
                 │Rise  │ │Vault│ │Breach│ │Python    │
                 │Priv  │ │Kernel│ │Entry│ │Bridge    │
                 └──────┘ └────┘ └────┘ └─────┬─────┘
                                               │
                     ┌─────────────────────────┼─────────┐
                     │                         │         │
                ┌────▼────┐ ┌────────────┐ ┌──▼─────┐ ┌──▼──────────┐
                │Horizon  │ │Wormy-ML    │ │Specter │ │BlueForge    │
                │Intel    │ │(Lateral)   │ │+ Apex  │ │Suite        │
                │(Recon)  │ │            │ │(AI)    │ │(Metrics)    │
                └─────────┘ └────────────┘ └────────┘ └─────────────┘
```

### gRPC Service Architecture (NEW in v2.3)

The C2 server now uses two separate gRPC service interfaces:

**AgentService** (agent ↔ C2):
- `CheckIn(AgentInfo) → (SessionID, Tasks)`
- `CommandStream(stream ClientMessage) → stream ServerMessage` — bidirectional streaming for tasking + results
- `Heartbeat(HeartbeatRequest) → (HeartbeatResponse)`
- `Exfiltrate(stream Chunk) → (ExfilStatus)`

**C2Service** (management console ↔ C2):
- `ListAgents`, `GetAgent`, `KillAgent` — agent lifecycle
- `CreateCampaign`, `GetCampaign`, `ListCampaigns`, `PauseCampaign`, `ResumeCampaign` — campaign management
- `DecisionFeed` — AI decisions stream
- `GetMetrics` — C2 health metrics

All communication uses X25519 ECDH key exchange + XChaCha20-Poly1305 AEAD encryption over gRPC.

## Component Map

| Component | Language | Layer | Phase | Status |
|-----------|----------|-------|-------|--------|
| **CLI (x404x)** | Go | 1 | All | v2.3 |
| **Web Dashboard** | Vue 3 | 1 | All | v2.2 |
| **Orchestrator** | Go | 2 | All | v2.3 |
| **C2 Server (gRPC)** | Go | 2 | C2 | v2.3 |
| **Unified Agent** | Go | 3 | All | v2.3 |
| **core/crypto** | Go | Shared | All | v2.3 |
| **core/proto** | Protobuf | Shared | All | v2.3 |
| **Python Bridge** | Python | 3 | IPC | v2.3 |
| **Evasion Module** | Python + Go | 3 | Evasion | v2.3 |
| **PhantomWeb** | JS + Python | 1,3 | Browser | v2.2 |
| **Breach-Entry** | C + Python | 3 | Delivery | [Repo](https://github.com/Ruby570bocadito/Breach-Entry) |
| **Horizon-Intel** | Python | 3 | Recon | [Repo](https://github.com/Ruby570bocadito/Horizon-Intel) |
| **Rise-Privilege** | Go | 3 | Exploitation | [Repo](https://github.com/Ruby570bocadito/Rise-Privilege) |
| **Vault-Kernel** | C + Go | 3 | Installation | [Repo](https://github.com/Ruby570bocadito/Vault-Kernel) |
| **Specter-Terminal** | Python | 3 | AI Analysis | [Repo](https://github.com/Ruby570bocadito/Specter-Terminal) |
| **Apex-Automation** | Python | 3 | AI Execution | [Repo](https://github.com/Ruby570bocadito/Apex-Automation) |
| **Wormy-ML** | Python | 3 | Lateral Movement | [Repo](https://github.com/Ruby570bocadito/Wormy-ML-Network-Worm) |
| **Link-Relay** | Python | 3 | C2 Relay | [Repo](https://github.com/Ruby570bocadito/Link-Relay) |
| **Titan-Operations** | Python + Go | 2 | Campaign Mgmt | [Repo](https://github.com/Ruby570bocadito/Titan-Operations) |
| **BlueForge-Suite** | Python | 3 | Metrics | [Repo](https://github.com/Ruby570bocadito/BlueForge-Suite) |

## Communication Flow

```
                    ┌──────────────────┐
                    │   Orchestrator   │
                    └────────┬─────────┘
                             │ gRPC
                    ┌────────▼─────────┐
                    │   C2 Server      │
                    │  (Go/gRPC)       │
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
        │ (Rise,   │ │ (Vault,  │  │ (TCP JSON)   │
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
SQLite schema (6 tables) in `core/appstate/state.go`:

- **campaigns**: Red team operations
- **agents**: Implant tracking
- **targets**: Discovered hosts
- **vulnerabilities**: CVEs and misconfigurations
- **credentials**: Captured passwords/hashes
- **audit_log**: Complete action trail

## Security Model

1. **End-to-End Encryption**: X25519 ECDH + XChaCha20-Poly1305 AEAD per session
2. **gRPC Transport**: Service-to-service communication over gRPC
3. **Offline AI**: Ollama runs locally, no data leaves the lab
4. **Safety Controls**: Kill switch, geofencing, auto-destruct, max infections
5. **Audit Trail**: Every action logged with timestamp, agent, campaign, and result
