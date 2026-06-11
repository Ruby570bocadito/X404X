```
 ██╗  ██╗██╗  ██╗ ██████╗ ██╗  ██╗
 ╚██╗██╔╝██║  ██║██╔═████╗╚██╗██╔╝
  ╚███╔╝ ███████║██║██╔██║ ╚███╔╝
  ██╔██╗ ╚════██║████╔╝██║ ██╔██╗
 ██╔╝ ██╗     ██║╚██████╔╝██╔╝ ██╗
 ╚═╝  ╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═╝
 Autonomous Red Team Platform
```

# X404X — Autonomous Red Team Platform

*Suite modular de operaciones ofensivas con evasión multi-capa, propagación avanzada y C2 post-cuántico.*

---

## ARQUITECTURA DEL SISTEMA

```
                           ┌──────────────────────────┐
                           │     DASHBOARD OPS        │
                           │  HTTP/WS + Vue 3 + D3    │
                           └───────────┬──────────────┘
                                       │
┌──────────────────────────────────────┼──────────────────────────────────────┐
│                         C2 HARDENED (Fase 2)                                │
│  ┌──────────┐ ┌───────────┐ ┌───────────┐ ┌──────────┐ ┌───────────────┐  │
│  │ SPIFFE   │ │ Ed25519   │ │ Kyber-1024│ │ Multi-   │ │ Proto         │  │
│  │ mTLS     │ │ Signing   │ │ + X25519  │ │ Canal    │ │ Obfuscation   │  │
│  └──────────┘ └───────────┘ └───────────┘ └──────────┘ └───────────────┘  │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │
┌──────────────────────────────────────┼──────────────────────────────────────┐
│                         AI ORCHESTRATOR (Fase 4)                            │
│  ┌──────────┐ ┌───────────┐ ┌───────────┐ ┌──────────┐ ┌───────────────┐  │
│  │ Q-Learn  │ │ Federated │ │ Autofact. │ │ Deepfake │ │ Wazero Bridge │  │
│  │ FSM      │ │ Learning  │ │ AFL++     │ │ Vishing  │ │ WASM→Go       │  │
│  └──────────┘ └───────────┘ └───────────┘ └──────────┘ └───────────────┘  │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │
┌──────────────────────────────────────┼──────────────────────────────────────┐
│                         CORE ENGINE (Go)                                    │
│  ┌──────────────────────┐  ┌──────────────────────┐  ┌──────────────────┐  │
│  │     RANSOMWARE       │  │       EVASION        │  │    PROPAGATION   │  │
│  │  AES-256-GCM         │  │  BYOVD (5 drivers)   │  │  Hydra (8 vect)  │  │
│  │  Bootkit UEFI        │  │  DKOM EPROCESS       │  │  Kerberos Deleg  │  │
│  │  SCADA Modbus        │  │  Anti-Reversing      │  │  IMDSv2 Bypass   │  │
│  │  Blockchain C2       │  │  Anti-Forensics      │  │  Reflective DLL  │  │
│  │  Polimorfismo JIT    │  │  Blue Pill HV        │  │  Chronos NTP      │  │
│  │  Cross-Platform      │  │  WFP DNS Poison      │  │  Ultrasound QPSK  │  │
│  └──────────────────────┘  └──────────────────────┘  └──────────────────┘  │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │
┌──────────────────────────────────────┼──────────────────────────────────────┐
│                         PYTHON BRIDGE (Fase 0)                              │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │  11 handler files · ~170 handlers · v26-v210 · BlockZ · Attacks      │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │
┌──────────────────────────────────────┼──────────────────────────────────────┐
│                         PLUGINS                                             │
│  ┌──────────┐ ┌───────────┐ ┌───────────┐ ┌──────────┐ ┌───────────────┐  │
│  │ Worm+RL  │ │ Operations│ │ Bluesky   │ │  Pulse   │ │ RF Contagion  │  │
│  │ 35 expl. │ │ Argos     │ │ BT Attk   │ │  C2 Go   │ │ SDR 4G/5G     │  │
│  └──────────┘ └───────────┘ └───────────┘ └──────────┘ └───────────────┘  │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │
┌──────────────────────────────────────┼──────────────────────────────────────┐
│                         INFRASTRUCTURE                                      │
│  ┌──────────┐ ┌───────────┐ ┌───────────┐ ┌──────────┐                    │
│  │ Docker   │ │ CI/CD     │ │ EDR Lab   │ │ Deploy   │                    │
│  │ Compose  │ │ GitHub    │ │ Win+ATP   │ │ One-click│                    │
│  └──────────┘ └───────────┘ └───────────┘ └──────────┘                    │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## ESTADO DEL PROYECTO

| Fase | Temática | Módulos | Líneas | Estado |
|------|----------|---------|--------|--------|
| **Fase 0** | Stubs críticos | 8 | 500 | ✅ COMPLETADA |
| **Fase 1** | Evasión + Anti-Forense | 10 | 3,599 | ✅ COMPLETADA |
| **Fase 2** | C2 Hardened | 6 | 2,374 | ✅ COMPLETADA |
| **Fase 3** | Propagación Avanzada | 12 | 3,274 | ✅ COMPLETADA |
| **Fase 4** | IA + Cross-Platform + Lab | 9 | 2,938 | ✅ COMPLETADA |
| **TOTAL** | | **45** | **12,685** | **100%** |

---

## CAPACIDADES POR DOMINIO

### Evasión (10 módulos)
| Módulo | Descripción |
|--------|-------------|
| BYOVD Loader | 5 drivers vulnerables (WinRing0, Gdrv, RTCore64, kprocesshacker, CPUID) |
| DKOM | Hide processes, steal SYSTEM token, EPROCESS manipulation |
| Anti-Reversing | HW BP detect (DR0-DR7), INT3 scan, CRC integrity, sandbox detect |
| Anti-Forensics | DoD 7-pass wipe, MFT corruption, crash dump disable, event log clear |
| WER Persistence | Hangs hijack + Silent Exit + triple persistence |
| MFT Slack Storage | Read/write in NTFS slack space, AES-GCM encrypted fragments |
| WFP DNS Poisoning | WFP provider + fake DNS server UDP 53 + hosts injection |
| Blue Pill Lite | VMXON/VMCS hypervisor, PatchGuard bypass, CPUID trap |
| LOLBin Chainer | 28 LOLBins, random chain per hour, multi-layer base64 |
| Kernel DNS Driver | NDIS filter, block Defender updates, DNS→C2 redirect |

### C2 Hardened (6 módulos)
| Módulo | Descripción |
|--------|-------------|
| SPIFFE mTLS | SVID generation, trust bundle, peer verification, cert rotation |
| Multi-Channel | gRPC→WebSocket→DoH→Twitter→Blockchain, auto-failover |
| Ed25519 | Command sign/verify, nonce replay protection, batch ops |
| Kyber-1024+X25519 | Hybrid KEM, AES-256-GCM + HMAC-SHA256 sessions |
| Dashboard Ops | HTTP+WS API, agent nodes, propagation map, command issuance |
| Proto Obfuscation | XOR+AES-CTR+GZIP, vaporize buffers, embedded loader |

### Propagación Avanzada (12 módulos)
| Módulo | Descripción |
|--------|-------------|
| Ultrasound QPSK | >18kHz modulation, speaker/mic RX/TX, WAV generation |
| Powerline PLC | HomePlug scan, UPnP SSDP, SOAP injection over electrical grid |
| USB ADB | ADB devices, APK install, shell exec, SMS/contacts dump |
| DNS Rebinding | TTL=0 rebind, SOP bypass JS, SSRF lateral movement |
| CI/CD Webhooks | GitHub Actions/Jenkins/GitLab CI injection, 10 env scanner |
| VLAN Jump | Double tagging, DTP negotiation, ARP flood per VLAN |
| QR Worm | QR matrix generation, PNG rendering, rotation channel |
| PJL Worm | Printer NVRAM, firmware infect, PCL ransom note |
| Chronos NTP | Fake NTP server, time shift, schedule hijack, w32tm |
| Reflective DLL | NtCreateSection+NtMapViewOfSection, 100-byte NASM stager |
| Kerberos Delegation | Unconstrained delegation, coercion, TGT dump, Silver Ticket |
| IMDSv2 Bypass | Token acquisition, IAM extraction, SSRF, neighbor scan |

### IA + Cross-Platform (9 módulos)
| Módulo | Descripción |
|--------|-------------|
| Cross-Platform Loader | ELF/Mach-O/APK generation, pack+encrypt, syscall hooks |
| JIT Polymorphism | NOP-sleds, code crossover, register reordering, runtime mutation |
| AI FSM Orchestrator | Q-learning FSM, exploration/exploitation, risk prediction |
| Federated Learning | FedAvg aggregation, victim profiling, phishing time prediction |
| Autofactory Fuzzer | AFL++ integration, 9 mutation strategies, crash detection |
| Wazero Bridge | WASM module parsing, TinyGo compilation, Python→WASM migration |
| RF Contagion | SDR detection, ModemManager, baseband injection, IMSI capture |
| EDR Test Lab | Windows Server 2022 + ELK + Sysmon + automated tests |
| Deepfake Vishing | Coqui TTS voice cloning, VoIP calls, SMS phishing |

---

## ESTRUCTURA DE DIRECTORIOS

```
X404X/
├── cmd/                        # CLI + agente implant
│   ├── x404x/console.go        # Shell interactiva
│   └── implant/main.go         # Agente C2 Go
│
├── internal/                   # CORE ENGINE
│   ├── ransomware/             # 35+ archivos Go
│   │   ├── engine.go           # Pipeline principal
│   │   ├── byovd_loader.go     # BYOVD 5 drivers
│   │   ├── dkom.go             # DKOM process hiding
│   │   ├── anti_reversing.go   # Anti-debugging
│   │   ├── anti_forensics_advanced.go  # DoD wipe + MFT
│   │   ├── wer_persistence.go  # WER hijack
│   │   ├── mft_slack.go        # MFT storage
│   │   ├── wfp_dns_poison.go   # DNS poisoning
│   │   ├── lolbin_chainer.go   # LOLBin chains
│   │   ├── chronos_ntp.go      # NTP manipulation
│   │   ├── kerberos_delegation.go  # Delegation abuse
│   │   ├── imdsv2_bypass.go    # AWS metadata bypass
│   │   ├── multi_channel_c2.go # Multi-channel C2
│   │   ├── jit_polymorphism.go # JIT mutations
│   │   ├── loader/cross.go     # Cross-platform loader
│   │   ├── stager/reflective_asm.go  # NASM stager
│   │   ├── hydra_vectors/      # 8 vectores Hydra
│   │   │   ├── ultrasound.go   # QPSK audio worm
│   │   │   ├── powerline.go    # PLC worm
│   │   │   ├── usb_adb.go      # ADB worm
│   │   │   ├── dns_rebinding.go # DNS rebind
│   │   │   ├── cicd_webhooks.go # CI/CD injection
│   │   │   ├── vlan_jump.go    # VLAN jump
│   │   │   ├── qr_worm.go      # QR worm
│   │   │   └── pjl_worm.go     # Printer worm
│   │   └── v27/, v29/, v210/   # Versiones avanzadas
│   ├── agent/                  # Post-exploit + privesc
│   ├── appstate/               # FSM + AI orchestradora
│   │   ├── state.go            # Registro de módulos
│   │   └── ai_orchestrator.go  # Q-learning FSM
│   ├── dispatch/dispatcher.go  # MITRE ATT&CK mapeo
│   └── bridge/wazero_loader.go # WASM bridge loader
│
├── modules/bridge/
│   ├── bridge.py               # Router RPC Go↔Python
│   └── handlers/               # 12 archivos, ~170 handlers
│       ├── ransomware.py       # 9 handlers
│       ├── ransomware_advanced.py  # 17 handlers
│       ├── ransomware_v26.py    # 13 handlers
│       ├── ransomware_v27.py    # 10 handlers
│       ├── ransomware_v28.py    # 23 handlers
│       ├── ransomware_v29.py    # 24 handlers
│       ├── ransomware_v210.py   # 10 handlers
│       ├── ransomware_blockz.py # 14 handlers
│       ├── attacks.py           # 9 handlers
│       ├── bloodhound.py        # 5 handlers
│       ├── cred_dump.py         # 6 handlers
│       └── phase_1_4.py         # 36 handlers (Fases 1-4)
│
├── plugins/                    # Plugins
│   ├── worm/                   # Worm con RL
│   ├── operations/             # Argos (cell/stager agents)
│   ├── blue/bluesky/           # Bluetooth attacks
│   ├── pulse-c2/               # C2 hardened
│   │   └── src/go/internal/
│   │       ├── crypto/         # SPIFFE, Ed25519, Kyber
│   │       └── c2/             # Dashboard
│   ├── ai/                     # AI plugins
│   │   ├── hivemind/federated.go    # Federated learning
│   │   ├── autofactory/fuzzer.go    # AFL++ fuzzer
│   │   └── vishing/deepfake.go      # Deepfake vishing
│   └── rf_contagion/baseband.go     # RF 4G/5G
│
├── lab/docker-compose.edr.yml  # EDR test lab
├── pkg/proto/                  # gRPC protos
├── deploy/deploy.sh            # One-click deploy
└── scripts/run_demo.sh         # Demo cross-platform
```

---

## USO

```bash
# Demo completa (dry-run)
bash scripts/run_demo.sh

# Worm simulation (Windows)
cd plugins/worm && python worm_core.py --config configs/config_simulation.yaml

# Worm simulation (Linux)
cd plugins/worm && python3 worm_core.py --config configs/config_simulation.yaml

# Dashboard operacional
cd plugins/pulse-c2/src/go && go run ./cmd/dashboard -port 9090

# EDR test lab
docker-compose -f lab/docker-compose.edr.yml up -d
```

---

## REQUISITOS

```bash
# Core
go >= 1.22.0
python >= 3.11

# Python packages
pip install -r requirements.txt

# Go modules
go mod tidy
```

---

## LICENCIA

Educational and research purposes only. Use only in authorized environments with explicit written permission.

---

*X404X v3.0 — 12,685 líneas — 45 módulos — 4 fases completas — 2026*
