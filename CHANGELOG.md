# Changelog

All notable changes to X404X will be documented in this file.

## [v2.4] — 2026-06-06

### Added — 14 Advanced Ransomware Modules (Blocks 1-4 + Bonus)

#### Block 1: Psychological & Reputational
- **Hope Trap** (`psychological_advanced.go`): Partial decryption of 5 recoverable-looking media/doc files, forensic tool monitor (FTK/Encase/Autopsy), triggers re-encryption with doubled ransom on detection. Deploy fake decryptor binaries that destroy remaining keys if executed
- **Identity Destruction** (`identity_destruction.go`): Harvest browser cookies/sessions (Chrome/Firefox/Edge), hijack accounts (email/Amazon/Facebook/LinkedIn/Twitter/GitHub/Slack/Outlook), post humiliating content, send phishing emails to contacts, purchase with saved cards, enable attacker 2FA
- **Inverse RaaS** (`raas_inverse.go`): Tor-based panel inviting other attackers, multi-ransom notes from different groups, Shamir key shard distribution, per-group encryption variants

#### Block 2: Pandemic Propagation
- **Multi-Platform Worm** (`multiplatform_worm.go`): Scans network via SSH/SMB/HTTP, deploys platform-specific payloads (Win PS1, Linux shell, macOS Automator, IoT shell). Docker container escape, Kubernetes pod spread, SMB worm propagation. IoT botnet scanning via Telnet/RTSP/Dahua ports with DDoS capability
- **Supply Chain Poison** (`supply_chain.go`): Finds software updaters (NP++, 7-Zip, VLC, Python, Node.js), replaces with trojan. Poisons pip.conf/npmrc/NuGet.Config for internal Artifactory/Nexus repos. Scans for .git repos, poisons hooks, scorches repo READMEs. Deploys fake emergency patches to S3/Azure/Contact emails
- **Cloud Exploit** (`cloud_exploit.go`): Harvests AWS/Azure/GCP credentials from disk. Launches EC2 instances (3 regions), creates malicious AMIs, Azure VMs with startup scripts, GCP compute instances. Creates public S3 bucket with fake patch website
- **Bluetooth/Wi-Fi Direct** (`bluetooth_prop.go`): BT device scanning via PowerShell/hcitool/system_profiler. BlueBorne exploit, BLE MITM via gatttool, Apple SIP bypass (CVE-2021-30892), malicious APK push. Activates rogue Wi-Fi hotspot, scans WiFi Direct peers for KRACK attacks

#### Block 3: Physical & Infrastructure Sabotage
- **SCADA/PLC Attack** (`scada_attack.go`): Detects SCADA software (Siemens/Rockwell/Schneider/CODESYS/Wonderware), scans for PLCs on Modbus/S7/CIP/EtherNet/IP ports. Modbus stop/write coil/overwrite logic commands, S7 stop/DB delete/flash firmware, CIP generic commands. Modbus unit ID brute force
- **Hardware Kill** (`hardware_kill.go`): Firmware access detection (BIOS/UEFI/SMBIOS/WMI/MSR). CPU overvoltage via MSR manipulation, GPU fan kill with nvidia-smi, fan controller disable, thermal throttling disable. Infinite CPU burn loop on all cores, BIOS flash corruption
- **Network Poison** (`network_poison.go`): ARP spoofing via arpspoof/New-NetNeighbor, MITM proxy with iptables/netsh redirect, content injection into HTTP/HTTPS pages. Root CA generation (RSA-4096) + installation into system trust store. Captive portal with DNS redirect, SSL strip attack

#### Block 4: Automutation & Resilience
- **DNA Mutation** (`dna_mutation.go`): Extracts SHA256 fingerprints from system DLLs (kernel32, ntdll, user32, libc, ssl), hybridizes malware genome with legit code via crossover recombination. ROP gadget generation, 15% junk code NOP sled insertion, per-machine XOR key derivation from volume serial
- **Bootkit** (`bootkit.go`): MBR payload generation (512-byte real-mode stub), UEFI PE stub generation. MBR infection via raw disk write, disk write interception via filter driver/dm-setup. Bootkit stage2 with C2 endpoint + reinfection interval. Fake SMART error display
- **Blockchain C2** (`blockchain_c2.go`): Bitcoin OP_RETURN command extraction from blockstream.info API. AES-GCM encrypted commands embedded in BTC transactions. Monitoring loop, command queue with execution tracking. Commands: destroy, encrypt_more, change_note, exfil, propagate, self_destruct

#### Bonus: Survivor Game
- **Survivor Game** (`survivor_game.go`): Discovers workstations on network, broadcasts game start message via wall/msg, eliminates random stations every 90 seconds with full-screen lock, announces winner with free decryption key. Double ransom for eliminated participants

### Added (Files)
- `core/ransomware/psychological_advanced.go`, `identity_destruction.go`, `raas_inverse.go` — Block 1
- `core/ransomware/multiplatform_worm.go`, `supply_chain.go`, `cloud_exploit.go`, `bluetooth_prop.go` — Block 2
- `core/ransomware/scada_attack.go`, `hardware_kill.go`, `network_poison.go` — Block 3
- `core/ransomware/dna_mutation.go`, `bootkit.go`, `blockchain_c2.go` — Block 4
- `core/ransomware/survivor_game.go` — Bonus
- `core/ransomware/engine_extended.go` — Master orchestrator for all 14 new phases
- `core/agent/ransomware_advanced_module.go` — 10 new agent modules (Advanced, HopeTrap, Identity, RaaS, Worm, SCADA, Hardware, Bootkit, Blockchain, Survivor)
- `modules/bridge/handlers/ransomware_advanced.py` — 17 new Python RPC handlers for all blocks
- `core/ransomware/types.go` — 16 new phases, 30+ new config fields, extended report struct
- `core/appstate/state.go` — 16 new ModuleDef entries registered

### Changed
- `core/ransomware/types.go` — Phase constants expanded from 7 to 23, config fields from 20 to 55, report fields from 13 to 24

## [v2.3] — 2026-06-06

### Added (Ransomware Engine — 8 features, 13 files)
- **Double extortion + selective exfiltration**: Heuristic content scanner (19 regex patterns for DNI, passports, credit cards, contracts, PST/OST, MDF/SQL, API keys, AWS keys, private keys, health data). ZIP with ChaCha20-password encrypted packaging. Exfil channels: DNS TXT fragments, CDN stego, S3 with stolen credentials. Shaming post generator + .onion negotiation URL
- **System poisoning & irreversible destruction**: MFT overwrite via raw volume access (`\\.\C:`), UEFI NVRAM/bootmgfw sabotage, shadow copies deletion (vssadmin/wmic/bcdedit), cloud backup API destruction (Veeam/Acronis/CommVault agents kill + config wipe), free space wiping, Linux MBR corruption
- **Propagation engine**: 6 exploit modules (Zerologon CVE-2020-1472, ProxyNotShell CVE-2023-23397, PrintNightmare CVE-2021-34527, BlueKeep CVE-2019-0708, EternalBlue MS17-010, SMBGhost CVE-2020-0796). Outlook COM propagation via fake thread conversations. WSUS/SCCM update poisoning, NuGet/NPM registry poisoning, Git hook poisoning
- **Real-time psychological attack**: TOPMOST full-screen countdown window (PowerShell WinForms), webcam capture via WIA COM, printer spam via Get-Printer, TTS audio threats via System.Speech, live file deletion display, desktop notification countdown loop
- **Anti-analysis & anti-IR**: 15 tool/process kill list (Procmon, Wireshark, x64dbg, IDA, Ghidra, etc.), PE header corruption, kernel driver kill (kprocesshacker.sys), sandbox hostname detection, kernel debugger detection (WMI + PowerShell), 2-hour sleep mode, C2 steganography via CDN image LSB + EXIF metadata
- **Binary polymorphism**: JIT mutation loop (function reordering, constant changes, junk code insertion), ROP gadget generation with NOP sleds, pre-packaging with per-machine key derivation (volume serial → SHA256 → XOR key)
- **Trust exploitation**: Self-signed code signing cert generation (RSA 4096), PFX/P12 certificate search on filesystem, WSUS fake update creation, SCCM malicious application deployment, NuGet/NPM registry poisoning, Git hook poisoning (pre-commit/post-commit/pre-push/post-merge)
- **Hydra multi-layer encryption**: 3× RSA-4096 key pairs, Shamir's Secret Sharing (3-of-3 split of master key), AES-256-GCM file encryption with double layer (AES-GCM + ChaCha20-Poly1305) for critical files (.mdf, .vhd, .pst), separate key per file, 3 shards sent to 3 independent C2 endpoints

### Added (Files)
- `core/ransomware/` — 12 Go files: engine.go, types.go, scanner.go, hydra.go, extortion.go, destruction.go, propagation.go, psychological.go, antianalysis.go, polymorph.go, trust.go, go.mod
- `core/agent/ransomware_module.go` — Agent Module interface for ransomware execution
- `modules/bridge/handlers/ransomware.py` — Python bridge handler (8 RPCs: execute, scan, encrypt, exfil, status, decrypt, note, propagate, destruct)
- `core/appstate/state.go` — 11 new ransomware ModuleDef entries in module registry
- `go.work` — Added `./core/ransomware`

### Added
- gRPC C2 Server with split service interfaces: `AgentService` (CheckIn, CommandStream, Heartbeat, Exfiltrate) and `C2Service` (ListAgents, GetAgent, KillAgent, CreateCampaign, GetCampaign, ListCampaigns, PauseCampaign, ResumeCampaign, DecisionFeed, GetMetrics)
- Real gRPC implementations in `core/c2server/agent_service.go` and `core/c2server/c2_service.go`
- Decision Engine integration in console exploit handler (Bridge + offline fallback)
- AppState-connected TUI (live agents, hosts, vulns, campaigns, decisions)

### Fixed
- C2 `server.go` rewritten from raw TCP to full `grpc.Server` with service registration
- Agent `connector.go` Send/Recv now use bidirectional `CommandStream` (was unconnected)
- Console `cmdExploit` — removed fake eternalblue/redis_unauth switch-case; replaced with generic Bridge + Decision Engine orchestration
- TUI — 100% hardcoded demo data replaced with live AppState queries
- `main.go` — TUI mode now initializes AppState and passes it to StartTUI
- Agent and c2server `go.mod` files updated with proper replace directives for proto dependencies

### Changed
- Console exploit handler: category-based dispatch (privesc, recon, post, auxiliary) with Bridge-first, offline-fallback strategy
- TUI kill chain rendering: dynamically derived from campaign phase instead of fixed
- All dashboard tabs render live state instead of hardcoded data

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
