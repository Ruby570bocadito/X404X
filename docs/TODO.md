# X404X — Status & TODO

> Última actualización: $(date +%Y-%m-%d) | Build: 21MB | Commits: 3

---

## ✅ COMPLETADO

### Núcleo del Framework
- [x] **CLI x404x** — Cobra commands + Bubble Tea TUI + msfconsole shell (3 modos)
- [x] **AppState** — Estado compartido (orchestrator + bridge + C2 + DB)
- [x] **Orchestrator** — Campaign Manager + Decision Engine (Rules 25% + A* 35% + AI 40%) + WorldGraph + EventBus + KillChainOrchestrator
- [x] **Decision Engine real** — 18 reglas determinísticas, A* pathfinding con WorldGraph, AI heurística con fallback offline
- [x] **Agent** — Implant Go unificado: ModuleManager, gRPC Connector, BridgeClient, PostExploitPipeline, KillChainEngine
- [x] **Crypto** — X25519 + XChaCha20-Poly1305 (5/5 tests PASS)
- [x] **gRPC Proto** — 4 servicios (Agent, C2, Bridge, Common) + 8 stubs generados

### API & Dashboard
- [x] **API Server** — REST 12 endpoints + CORS + WebSocket hub + wildcard EventBus
- [x] **Dashboard Vue 3** — 7 pestañas, 6 Pinia stores, tema cyberpunk/glass, xterm.js
- [x] **Dashboard stores** — fetch() reales a API, sin fallback hardcodeado
- [x] **C2 Server** — TCP listener integrado, agent management

### Bridge & Python
- [x] **Python Bridge** — 9 módulos (recon, ai_analyze, privesc, persist, worm, relay, blue, exfil, health)
- [x] **Bridge auto-start** — AppState.Start() arranca el bridge automáticamente
- [x] **Evasion Module** — AMSI/ETW, polymorphic, sleep jitter, sandbox detect

### Post-Explotación
- [x] **PostExploitPipeline** — FullChain(): escalate() → stealth() → persistence() → propagate()
- [x] **VaultIOCTL** — 14 ioctl commands wrapper para /dev/vault_kernel
- [x] **KillChainEngine (Agent)** — Phase transitions automáticas
- [x] **KillChainOrchestrator** — ReconComplete() → WeaponizeComplete() → ... → ObjectiveComplete()
- [x] **Módulos post-explotación registrados** — 7 módulos: full_chain, privesc, stealth, propagate, credential_dump, keylogger, evasion_apply

### Submódulos
- [x] **Pulse-C2** (27 commits) — C2 server + agent + Vue dashboard
- [x] **Rise-Privilege** (11 commits) — 12 vectores, 60+ GTFOBins
- [x] **Vault-Kernel** (8 commits) — LKM rootkit
- [x] **Breach-Entry** (6 commits) — CVE-2026-XXXX
- [x] **Horizon-Intel** (6 commits) — OSINT recon
- [x] **Specter-Terminal** (21 commits) — Terminal IA offline
- [x] **Apex-Automation** (5 commits) — IA pentesting
- [x] **Wormy-ML-Network-Worm** (13 commits) — 44 exploits, RL engine
- [x] **Link-Relay** (6 commits) — C2 relay
- [x] **Titan-Operations** (4 commits) — ARGOS v2.0
- [x] **BlueForge-Suite** (8 commits) — Blue team metrics

### Infraestructura
- [x] **Git repo** — Inicializado, 3 commits
- [x] **go.work** — Workspace con 14 módulos
- [x] **npm install** — 159 packages, 0 vulnerabilidades
- [x] **Docker lab** — 5 contenedores (attacker + 2 targets + dashboard + ollama)
- [x] **CI/CD** — GitHub Actions: lint, test, build multi-arch
- [x] **Makefile** — setup, build, test, lab-up/down, proto, lint, clean
- [x] **SQLite schema** — 6 tablas (campaigns, agents, targets, vulns, credentials, audit_log)

---

## 🔶 EN PROGRESO

- [ ] **SQLite driver real** — Schema listo, falta `mattn/go-sqlite3` (CGO). Datos en memoria funcionan.
- [ ] **Agent ↔ C2 gRPC handshake** — Connector escrito, falta arrancar Pulse-C2 y hacer check-in real
- [ ] **Pulse-C2 events → EventBus** — El C2 submodule tiene eventos, falta suscribirlos al orquestador

---

## ⬜ PENDIENTE

### FASE C — Dashboard + Tests (corto plazo)
- [ ] **Dashboard sin fallback** — Eliminar datos hardcodeados en componentes Vue
- [ ] **Dashboard WebSocket 100% real** — Live feed desde EventBus del orchestrator
- [ ] **x404x dashboard arranca TODO** — Un solo comando: API + Bridge + C2 + DB
- [ ] **Tests de integración** — Pipeline completo en Docker lab
- [ ] **Tests unitarios** — Cobertura para agent, orchestrator, appstate, bridge
- [ ] **CI/CD verde** — GitHub Actions con build + test + lint pasando

### FASE D — Medio plazo
- [ ] **Modo Automático AI** — Si `cfg.AI.AutoApproval = true`, ejecuta sin HITL cuando confianza > 0.85
- [ ] **Evasión unificada** — Integrar AMSI/ETW de Pulse-C2 + polimorfismo de Wormy
- [ ] **Escenarios CTF** — `x404x lab scenario ctf_ad` levanta DC, Exchange, DB, Web, WS, Firewall
- [ ] **Reportes PDF** — `x404x campaign report --format pdf` con MITRE ATT&CK mapping
- [ ] **Módulo credenciales** — Unificar Mimikatz + mimipenguin + /etc/shadow crack

### FASE E — TFG
- [ ] **Memoria TFG** — LaTeX: arquitectura, integración, pruebas, ética
- [ ] **Benchmarks** — Latencia Decision Engine, throughput C2, tasa evasión
- [ ] **Publicación GitHub** — `github.com/Ruby570bocadito/X404X` con CI verde
- [ ] **Video demo** — Kill chain completa en Docker lab
- [ ] **Plugin system** — API para módulos de terceros

---

## 🎯 PRÓXIMO (orden de prioridad)

1. **Instalar SQLite driver** → persistencia real (`mattn/go-sqlite3`)
2. **Dashboard: eliminar hardcodeo** → 100% datos de API
3. **Tests de integración** → Pipeline end-to-end en Docker
4. **Modo automático AI** → DecisionEngine sin HITL
5. **Publicación GitHub** → Subir repo público con CI
6. **Memoria TFG** → Documentación académica

---

## 🔗 Matriz de colaboración

Ver [docs/KILL_CHAIN_MATRIX.md](docs/KILL_CHAIN_MATRIX.md) para la matriz completa de interacción entre los 11 módulos en las 7 fases de la kill chain.

## 📁 Estructura

```
X404X/                          # 104 archivos propios
├── cmd/x404x/    (5 .go)       # CLI + TUI + Console + Dashboard
├── core/agent/    (7 .go)       # Agent + Connector + Bridge + PostExploit + VaultIOCTL + KillChain
├── core/orch/     (5 .go)       # Orchestrator + Decision + Events + WorldGraph + KillChainOrchestrator
├── core/appstate/ (1 .go)       # Shared state
├── core/api/      (2 .go)       # REST API + WebSocket
├── core/c2server/ (1 .go)       # C2 listener
├── core/crypto/   (2 .go)       # Crypto
├── core/proto/    (4 .proto)    # gRPC definitions
│   └── gen/       (8 .go)       # Generated stubs
├── modules/       (3 .py)       # Bridge + Evasion + DB models
│   └── [7 submodules]           # Python repos
├── core/          [4 submodules]# C/Go repos
├── shared/        (4 .go)       # Config + Logger + Types
├── web/           (16 .vue/.js) # Vue 3 Dashboard
├── lab/           (3 files)     # Docker lab
├── docs/          (4 .md)       # Documentation
└── scripts/       (1 .sh)       # Setup
```

---

**Total: ~1,150 archivos** (104 propios + 1,046 submódulos) | **21MB binary** | **11 repos integrados**
