# Changelog

All notable changes to X404X will be documented in this file.

## [v2.0] — 2026-06-06

### Added
- 14 new modules: credential_dump, bloodhound, responder, web_scan, cloud (AWS/Azure/GCP), cleanup, obfuscate, exfiltrate, phantom_xss, phantom_sw_persist, phantom_browser_mesh, phantom_socks5, payload_obfuscate
- PhantomWeb browser-native implant integration (XSS, Watering Hole, Service Worker, Browser Mesh, SOCKS5)
- Payload Builder CLI: `x404x payload generate` (multi-arch cross-compile)
- Listeners management: `x404x listeners` (HTTP/HTTPS/DNS/ICMP/SMB/TCP/WS/DoH)
- SQLite persistence (6 tables) via modernc.org/sqlite (pure-Go, no CGO)
- Real Module implementations (PostExploitModule, PrivescScanModule, ReconModule)
- Agent ↔ C2 gRPC connector wired in Agent
- PhantomWeb Pinia store + Browser Mesh dashboard tab
- 4 evasion profiles: none, balanced, stealth, maximum
- Campaign report generator (JSON/MD/HTML/PDF) with MITRE ATT&CK mapping
- Exfiltration manager (chunked 64KB encrypted transfer)
- Rise-Privilege binary wrapper (auto-compile, Scan/Exploit/FullChain)
- Auto-mode AI (auto-approve decisions when confidence > 0.85)
- CTF scenario: Active Directory Lab (6 containers)
- Bridge handlers: 11 → 19 (+8 new)
- Console modules: 28 → 42 (+14 new)
- Dashboard tabs: 7 → 8 (+Browser Mesh)
- Pinia stores: 6 → 7 (+phantom)

### Fixed
- Agent connector wired to C2 (was unassigned)
- escalate() now actually runs Rise-Privilege binary
- installCron() and installSystemd() now write real persistence
- Phantom API endpoints registered
- BrowserMesh.vue: 0% hardcodeo (uses phantom store)
- Bridge sys.path for proper module imports
- PROJECT_ROOT auto-added to Python path

### Security
- LICENSE added (MIT)
- SECURITY.md added
- Kill switch, geofencing, auto-destruct, max infections, no persistence by default

## [v1.0] — 2026-06-05

### Added
- Initial release: CLI (Cobra + Bubble Tea + msfconsole shell), Dashboard (Vue 3), API (REST + WebSocket)
- Orchestrator: Decision Engine (Rules 25% + A* 35% + AI 40%), WorldGraph, EventBus
- Agent: Go implant, gRPC Connector, BridgeClient (Python↔Go IPC)
- Crypto: X25519 + XChaCha20-Poly1305
- gRPC Proto: 4 services (Agent, C2, Bridge, Common)
- Python Bridge: 11 modules (recon, AI, privesc, persist, worm, relay, blue, evasion, report, exfil, health)
- 11 submodules integrated
- Docker lab (5 containers)
- SQLAlchemy models (11 tables)
- CI/CD (GitHub Actions)
- Documentation: Architecture, Roadmap, CLI Reference, TFG Memory, Benchmarks
