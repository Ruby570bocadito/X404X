# X404X — ROADMAP COMPLETO

> Suite modular de operaciones ofensivas
> Última actualización: 2026-06-11

---

## ARQUITECTURA DEL PROYECTO

```
X404X/
├── cmd/                          # CLI + implante
│   ├── x404x/console.go          # Shell interactiva (~400 líneas)
│   ├── x404x/commands.go         # Comandos y dispatch
│   ├── x404x/payload.go          # Generación de payloads
│   ├── x404x/payload_commands.go # Control de ejecución
│   └── implant/main.go           # Agente C2 Go
│
├── internal/                     # CORE ENGINE
│   ├── ransomware/               # Motor de ransomware (35+ archivos Go)
│   │   ├── engine.go             # Pipeline principal
│   │   ├── engine_extended.go    # 14 sub-engines
│   │   ├── bootkit.go            # UEFI bootkit PE
│   │   ├── blockchain_c2.go      # C2 via BTC/ETH addresses
│   │   ├── dna_mutation.go       # Motor polimórfico
│   │   ├── scada_attack.go       # Modbus/TCP
│   │   ├── cloud_exploit.go      # Metadata services
│   │   ├── hardware_kill.go      # hdparm, nvme format
│   │   ├── propagation.go        # SSH, SMB, USB, Bluetooth
│   │   ├── stealth_c2.go         # DNS/ICMP/dead drops
│   │   ├── self_healing.go       # 7 mecanismos de persistencia
│   │   ├── antiforensics.go      # MFT/USN/prefetch
│   │   ├── bootkit_uefi.go       # DXE UEFI
│   │   ├── hydra.go              # Propagación Hydra
│   │   ├── extortion.go          # Notas de rescate
│   │   ├── destruction.go        # Auto-destrucción
│   │   ├── v26/  → evasion_deep, phantom_evasion, ...
│   │   ├── v27/  → hypervisor, uefi, kernel, ...
│   │   ├── v28/  → inception, iot, zombie, ...
│   │   ├── v29/  → firmware, smm, radiation, ...
│   │   ├── v210/ → apocalipsis (worm + botnet), phantom
│   │   └── v30/  → armageddon (próxima gen)
│   ├── agent/
│   │   ├── modules.go            # Post-exploit + privesc
│   │   └── v26_v210_modules.go   # Módulos v27-v29
│   ├── appstate/state.go         # FSM de campañas
│   ├── dispatch/dispatcher.go    # Orquestador de módulos
│   └── engine/                   # Orquestador interno
│
├── modules/                      # BRIDGE + EVASION
│   ├── bridge/
│   │   ├── bridge.py             # Router RPC Go↔Python
│   │   └── handlers/             # 11 archivos, ~135 handlers
│   │       ├── ransomware.py           # 9 handlers
│   │       ├── ransomware_advanced.py  # 17 handlers
│   │       ├── ransomware_blockz.py    # 14 handlers
│   │       ├── ransomware_v26.py       # 13 handlers
│   │       ├── ransomware_v27.py       # 10 handlers
│   │       ├── ransomware_v28.py       # 23 handlers
│   │       ├── ransomware_v29.py       # 24 handlers
│   │       ├── ransomware_v210.py      # 10 handlers
│   │       ├── bloodhound.py           # 5 functions
│   │       ├── cred_dump.py            # 6 functions
│   │       └── attacks.py              # 9 functions
│   ├── evasion/
│   │   ├── unified_evasion.py    # AMSI/ETW bypass
│   │   ├── lolbin_delivery.py    # 10 técnicas LOLBin
│   │   └── evasion.py
│   ├── notifications/webhook.py
│   └── phantom/controller.py
│
├── plugins/                      # PLUGINS
│   ├── worm/                     # Worm con RL
│   │   ├── worm_core/            # Core (mixins, config, standalone)
│   │   ├── exploits/             # 35 módulos de exploit
│   │   ├── payloads/             # Generación de payloads
│   │   ├── training/             # Escenarios de entrenamiento
│   │   ├── c2/                   # Multi-protocol C2
│   │   ├── rl_engine/            # Reinforcement learning
│   │   ├── scanner/              # Escáner de red
│   │   ├── infection/            # Módulo de infección
│   │   ├── monitoring/           # Dashboard + monitor
│   │   ├── utils/                # Logger, crypto, mitre, etc.
│   │   └── scripts/              # Scripts .bat/.ps1
│   ├── operations/               # Argos (operaciones)
│   │   ├── agents/cell/          # Agente célula (Go)
│   │   ├── agents/stager/        # Stager implant
│   │   ├── core/                 # EventBus, Director, KnowledgeTree
│   │   ├── ui/cli.py             # CLI rica con Rich
│   │   └── tests/                # Tests + demo
│   ├── blue/bluesky/             # Ataques Bluetooth
│   ├── ai/specter/               # SPECTER (AI orchestration)
│   ├── ai/apex/                  # APEX (planning)
│   ├── pulse-c2/                 # C2 en Go con evasión
│   ├── privesc/                  # Escalada de privilegios
│   └── recon/ui/dashboard.py     # Dashboard de reconocimiento
│
├── pkg/
│   ├── proto/gen/                # gRPC (agent, bridge, c2)
│   └── shared/database/          # Modelos DB
│
├── deploy/                       # Despliegue
│   ├── deploy.sh                 # One-click deploy
│   ├── Dockerfile                # Docker con healthchecks
│   └── docker-compose.yml
│
├── .github/workflows/            # CI/CD pipeline
│
└── docs/
    ├── USAGE.md
    ├── COMMANDS.md
    ├── CONSOLE.md
    ├── MODULES.md
    ├── CREATIVITY.md
    └── README.md
```

---

## ESTADO ACTUAL (COMPLETADO)

| Componente | Estado | Detalle |
|-----------|--------|---------|
| Core Ransomware Engine | ✅ 100% | 35+ archivos Go, AES-256-GCM, bootkit, SCADA, blockchain C2, polimorfismo |
| Python Bridge Handlers | ✅ 100% | 11 archivos, ~135 handlers reales |
| Agent (post-exploit) | ✅ 100% | modules.go + v26_v210_modules.go |
| FSM Campañas | ✅ 100% | appstate/state.go |
| Evasión | ✅ 100% | LOLBin + unified evasion + AMSI/ETW bypass + BYOVD + DKOM + Blue Pill + LOLBin Chain |
| C2 Pulse-C2 | ✅ 100% | mTLS+SPIFFE + Ed25519 + Kyber + multi-channel + dashboard |
| Worm Plugin | ✅ 100% | 35 exploits + RL + cross-platform loader + JIT polymorphism |
| Operations Plugin | ✅ 100% | Agent cell/stager + CLI + EventBus + Dashboard |
| Bluesky Plugin | ✅ 100% | Bluetooth attacks + RF Contagion 4G/5G |
| Infra + CI/CD | ✅ 100% | Docker, deploy.sh, CI pipeline, EDR test lab |
| AI Plugins | ✅ 100% | Hivemind FedAvg + Autofactory fuzzer + Deepfake Vishing |
| Documentación | ✅ 100% | 6 docs en docs/ y raíz |

---

## FASES DE IMPLEMENTACIÓN

### FASE 0 (STUBS CRÍTICOS — COMPLETADA ✅) — ~500 líneas

Cerrar los últimos stubs y código incompleto del proyecto.

| # | Archivo | Línea | Descripción | Estado |
|---|---------|-------|-------------|--------|
| 1 | `plugins/privesc/exploit.go` | 275 | `exploitKernelCVE()` — busca binario local y lo ejecuta con `os/exec` | ✅ |
| 2 | `plugins/ai/specter/core/report_generator.py` | 162-168 | Export PDF via `fpdf2` + DOCX via `python-docx` | ✅ |
| 3 | `plugins/ai/specter/agents/orchestrator.py` | 104 | `_execute_task()` real: ReconAgent usa `nmap` + fallback socket scan | ✅ |
| 4 | `plugins/operations/ui/cli.py` | 490-496 | Timeline (events.jsonl), MITRE (ATTACK_TECHNIQUES), flags (flags.json) | ✅ |
| 5 | `plugins/operations/agents/cell/main.go` | 336 | `executeExploit()` con net.Dial + CVE routing + DropAndExecute | ✅ |
| 6 | `plugins/worm/payloads/specialized_payloads.py` | 269 | AES-256-CBC ransomware real con key exfiltration al C2 | ✅ |
| + | `plugins/worm/configs/config_simulation.yaml` | — | Archivo de configuración para modo simulación (referenciado por .bat) | ✅ |
| + | `scripts/run_demo.sh` | — | Demo cross-platform Linux con dry-run de punta a punta | ✅ |

---

### FASE 1 (EVASIÓN + ANTI-FORENSE — COMPLETADA ✅) — 3,599 líneas

| # | Módulo | Archivo | Líneas | Descripción |
|---|--------|---------|--------|-------------|
| 1.1 | BYOVD Loader | `internal/ransomware/byovd_loader.go` | 316 | 5 drivers (WinRing0, Gdrv, RTCore64, kprocesshacker, CPUID) con IOCTLs reales. Read/Write physical memory, MSR, handle elevation, EDR evasion. |
| 1.2 | DKOM | `internal/ransomware/dkom.go` | 391 | Ocultación de procesos vía ActiveProcessLinks, steal SYSTEM token, protección de proceso, downgrade handles EDR. |
| 1.3 | Anti-Reversing | `internal/ransomware/anti_reversing.go` | 402 | HW BP detect (DR0-DR7), INT3 scan, CRC integrity check, timing check, sandbox detect, MAC OUI, self-destruct. |
| 1.4 | Anti-Forense Avanzado | `internal/ransomware/anti_forensics_advanced.go` | 376 | VAD hide, DoD 5220.22-M wipe (7-pass), MFT corruption, crash dump disable, event log clear, prefetch/USN/shellbag/shimcache wipe. |
| 1.5 | WER Persistence | `internal/ransomware/wer_persistence.go` | 253 | WER Hangs hijack (DLL load on crash), Silent Process Exit, startup persistence, scheduled task. |
| 1.6 | MFT Slack Storage | `internal/ransomware/mft_slack.go` | 356 | Lectura/escritura en slack space de MFT, AES-GCM encrypt de fragments, hide agent/ransom note. |
| 1.7 | WFP DNS Poisoning | `internal/ransomware/wfp_dns_poison.go` | 349 | WFP provider + netsh fallback, DNS redirect via hosts, fake DNS server (UDP 53), DNS cache flush. |
| 1.8 | Blue Pill Lite | `internal/ransomware/v27/bluepill.go` | 408 | Hypervisor VT-x/AMD-V, VMXON region, VMCS setup con guest/host state, PatchGuard bypass, CPUID trap. |
| 1.9 | LOLBin Chain | `internal/ransomware/lolbin_chainer.go` | 365 | 28 LOLBins (20 Windows + 8 Linux), cadena aleatoria por hora, multi-layer base64 encoding, ejecución. |
| 1.10 | Kernel DNS Driver | `internal/ransomware/v29/wfp_driver.go` | 383 | NDIS filter driver, bloquea Defender updates, intercepta update services, redirect DNS a C2. |

---

### FASE 2 (C2 HARDENED — COMPLETADA ✅) — 2,374 líneas

| # | Módulo | Archivo | Líneas | Descripción |
|---|--------|---------|--------|-------------|
| 2.1 | mTLS + SPIFFE | `plugins/pulse-c2/.../crypto/spiffe.go` | 341 | SVID generation, trust bundle, peer SPIFFE ID verification, server/client mTLS configs, certificate rotation. |
| 2.2 | Multi-Canal C2 | `internal/ransomware/multi_channel_c2.go` | 560 | 5 canales: gRPC → WebSocket → DoH → Twitter → Blockchain, health check, auto-failover, beacon loop, censorship bypass. |
| 2.3 | Ed25519 Signing | `plugins/pulse-c2/.../crypto/ed25519.go` | 288 | Key generation, command signing/verification, nonce replay protection, trusted key ring, batch sign/verify. |
| 2.4 | Dashboard Ops | `plugins/pulse-c2/.../c2/dashboard.go` | 406 | HTTP+WebSocket API, agent registration, propagation map, signed command issuance, embedded HTML dashboard. |
| 2.5 | Kyber-1024+X25519 | `plugins/pulse-c2/.../crypto/kyber.go` | 419 | Hybrid KEM (ML-KEM-1024 + X25519), HKDF key derivation, AES-256-GCM + HMAC-SHA256 session encryption. |
| 2.6 | Proto Obfuscation | `pkg/proto/loader/obfuscated.go` | 360 | XOR + AES-CTR + GZIP, integrity verification, memory-only loading, vaporize buffers, embedded loader generation. |

---

### FASE 3 (PROPAGACIÓN AVANZADA — COMPLETADA ✅) — 3,274 líneas

| # | Módulo | Archivo | Líneas | Descripción |
|---|--------|---------|--------|-------------|
| 3.1 | Ultrasonido QPSK | `hydra_vectors/ultrasound.go` | 302 | Modulación QPSK >18kHz, WAV generation, preamble sync, speaker/mic RX/TX, arecord receiver. |
| 3.2 | Ethernet Powerline | `hydra_vectors/powerline.go` | 234 | HomePlug device scan, UPnP SSDP discover, PLC command injection via SOAP, SSRF port mapping. |
| 3.3 | USB ADB | `hydra_vectors/usb_adb.go` | 249 | ADB device enumeration, APK install, remote shell exec, contact/SMS dump, sdcard worm persistence. |
| 3.4 | DNS Rebinding | `hydra_vectors/dns_rebinding.go` | 267 | TTL=0 rebind server, SOP bypass JS payload, SSRF lateral via Host headers, ARP network scan. |
| 3.5 | CI/CD Webhooks | `hydra_vectors/cicd_webhooks.go` | 236 | GitHub Actions workflow injection, Jenkins job trigger, GitLab CI pipeline, 10 CI environment scanner. |
| 3.6 | VLAN Jump | `hydra_vectors/vlan_jump.go` | 236 | VLAN interface creation, double tagging, DTP frame negotiation, ARP flood, DHCP discover per VLAN. |
| 3.7 | QR Dinámico | `hydra_vectors/qr_worm.go` | 305 | QR matrix generation, finder/timing/alignment patterns, data encoding, PNG render, rotation channel. |
| 3.8 | PJL Worm | `hydra_vectors/pjl_worm.go` | 235 | Printer job language exploits, NVRAM read/write, firmware infect, ransom note print, PCL injection. |
| 3.9 | NTP Manipulation | `chronos_ntp.go` | 280 | Fake NTP server, time forward/rewind, scheduled task shift, log timestamp corruption, w32tm hijack. |
| 3.10 | Reflective DLL NASM | `stager/reflective_asm.go` | 335 | NtCreateSection+NtMapViewOfSection, 100-byte NASM stager, CreateRemoteThread, remote process injection. |
| 3.11 | Kerberos Delegation | `kerberos_delegation.go` | 270 | Unconstrained delegation discovery, LDAP/SMB/printerbug coercion, TGT dump, Pass-the-Ticket, Silver Ticket. |
| 3.12 | IMDSv2 Bypass | `imdsv2_bypass.go` | 325 | IMDSv2 token acquisition, IAM credential extraction, SSRF exploit, neighbor instance scan, STS AssumeRole. |

---

### FASE 4 (IA + CROSS-PLATFORM + LAB — COMPLETADA ✅) — 2,938 líneas

| # | Módulo | Archivo | Líneas | Descripción |
|---|--------|---------|--------|-------------|
| 4.1 | Cross-Platform Loader | `loader/cross.go` | 365 | ELF/Mach-O/APK generation, pack+encrypt, syscall hooks, anti-sandbox detection. |
| 4.2 | JIT Polymorphism | `jit_polymorphism.go` | 380 | NOP-sleds, constant obfuscation, code crossover, register reordering, instruction substitution, runtime mutation loop. |
| 4.3 | AI FSM Orchestrator | `appstate/ai_orchestrator.go` | 400 | Q-learning FSM, state transitions, exploration/exploitation, risk prediction, Q-table import/export. |
| 4.4 | Federated Learning | `ai/hivemind/federated.go` | 372 | FedAvg aggregation, victim profiling, optimal phishing time prediction, model export. |
| 4.5 | Autofactory Fuzzer | `ai/autofactory/fuzzer.go` | 476 | AFL++ integration, 9 mutation strategies, crash detection, exploit candidate generation. |
| 4.6 | Bridge→Go/Wazero | `bridge/wazero_loader.go` | 302 | WASM module parsing, TinyGo compilation, handler migration from Python, stub WASM generation. |
| 4.7 | RF Contagion | `rf_contagion/baseband.go` | 281 | SDR detection (RTL-SDR/HackRF), ModemManager, baseband payload injection, IMSI capture, SS7 attack. |
| 4.8 | EDR Test Lab | `lab/docker-compose.edr.yml` | 94 | Windows Server 2022 + ELK stack + Sysmon + automated evasion test suite, 5/5 pass target. |
| 4.9 | Deepfake Vishing | `ai/vishing/deepfake.go` | 268 | Coqui TTS voice cloning, VoIP call placement, SMS phishing, social engineering profiling. |

---

## MÓDULOS NUEVOS (12 PROPUESTOS) — MAPA DE INTEGRACIÓN

| # | Módulo | Fase | Prioridad |
|---|--------|------|-----------|
| ① | Kernel DNS WFP Driver | Fase 1 | Alta |
| ② | Blue Pill Lite Hypervisor | Fase 1 | Alta |
| ③ | Kerberos Unconstrained Delegation | Fase 3 | Media |
| ④ | Deepfake Voice Calls (Vishing) | Fase 4 | Baja |
| ⑤ | LOLBin Chain Dinámica | Fase 1 | Alta |
| ⑥ | Reflective DLL Injection NASM | Fase 3 | Media |
| ⑦ | NTP Time Manipulation | Fase 3 | Media |
| ⑧ | WER Persistence | Fase 1 | Alta |
| ⑨ | IMDSv2 Bypass | Fase 3 | Media |
| ⑩ | MFT Slack Space Storage | Fase 1 | Alta |
| ⑪ | Federated Learning Hivemind | Fase 4 | Baja |
| ⑫ | Proto Obfuscation XOR | Fase 2 | Media |

---

## RESUMEN DE ESFUERZO

| Fase | Items | Líneas | Estado |
|------|-------|--------|--------|
| Fase 0 | 8 fixes | ~500 | ✅ **COMPLETADA** |
| Fase 1 | 10 módulos | **3,599** | ✅ **COMPLETADA** |
| Fase 2 | 6 módulos | **2,374** | ✅ **COMPLETADA** |
| Fase 3 | 12 módulos | **3,274** | ✅ **COMPLETADA** |
| Fase 4 | 9 módulos | **2,938** | ✅ **COMPLETADA** |
| **TOTAL** | **45 items** | **12,685** | — |

---

## CONVENCIONES DE CÓDIGO

- **Go**: `internal/` para core, `plugins/` para módulos externos. Syscalls directas con `syscall`.
- **Python**: `modules/bridge/handlers/` para handlers RPC. Tipo `dict` para todos los params/returns.
- **Nombres**: snake_case en Python, PascalCase en Go. Archivos descriptivos en minúscula.
- **Errores**: `return error` en Go, `raise` con mensaje descriptivo en Python.
- **Simulación**: Flag `--dry-run` en CLI, `simulation=True` en bridge handlers.

---

*Este roadmap se actualizará conforme avancen las implementaciones.*
