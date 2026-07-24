# X404X — Memoria Técnica TFG

> **Trabajo de Fin de Grado en Ciberseguridad**
>
> Autor: **Rafael Gálvez** — [@Ruby570bocadito](https://github.com/Ruby570bocadito)
>
> Centro: **Cisco NetAcad · Málaga, Spain**
>
> Tutor: [Por determinar]
>
> Convocatoria: 2026

---

## Resumen

X404X es una plataforma Red Team semi-autónoma que integra 11 herramientas ofensivas especializadas en un monorepo unificado. El framework cubre la cadena de eliminación (kill chain) completa —desde reconocimiento hasta exfiltración— combinando agentes Go de alto rendimiento, inteligencia artificial offline (Ollama), persistencia a nivel kernel (LKM), comunicación C2 cifrada (X25519 + XChaCha20-Poly1305), y propagación autónoma con aprendizaje por refuerzo.

**Palabras clave**: Red Team, C2, kill chain, gRPC, inteligencia artificial, rootkit, privilege escalation, monorepo, Go, Python, Docker.

---

## 1. Introducción

### 1.1 Motivación

Las operaciones de Red Team modernas requieren coordinar múltiples herramientas especializadas —scanners de reconocimiento, exploits de acceso inicial, suites de escalada de privilegios, rootkits de persistencia, gusanos de propagación, frameworks C2— que típicamente operan de forma aislada. Esta fragmentación genera fricción operativa: formatos de datos incompatibles, canales de comunicación dispares, y ausencia de un motor de decisión centralizado.

X404X resuelve esta fragmentación proporcionando una plataforma donde los 11 componentes comparten un protocolo unificado (gRPC), una capa de cifrado común (X25519 + XChaCha20-Poly1305), una base de datos compartida (SQLite/PostgreSQL), y un orquestador central con inteligencia artificial que coordina la progresión a través de la kill chain.

### 1.2 Objetivos

1. **Integrar** 11 herramientas ofensivas independientes en un monorepo cohesionado mediante Git submodules.
2. **Desarrollar** un agente Go unificado que invoque módulos especializados y reporte a un C2 central.
3. **Implementar** un motor de decisión híbrido (reglas determinísticas + pathfinding A* + IA con Ollama).
4. **Diseñar** un pipeline de post-explotación colaborativo donde Rise-Privilege, Vault-Kernel y Wormy-ML operen como un sistema unificado.
5. **Construir** interfaces modernas: CLI con 3 modos (Cobra, TUI Bubble Tea, shell msfconsole), dashboard Vue 3 con tema cyberpunk, y API REST con WebSocket.
6. **Validar** el funcionamiento mediante un laboratorio Docker aislado con 5+ contenedores.

### 1.3 Alcance y limitaciones

El framework está diseñado exclusivamente para laboratorios controlados, competiciones CTF y entornos autorizados. No incluye mecanismos de ofuscación de tráfico a nivel ISP ni bypass de sistemas EDR empresariales avanzados (CrowdStrike, SentinelOne). La IA opera en modo offline con modelos locales (Ollama) para garantizar que ningún dato abandona el laboratorio.

---

## 2. Estado del Arte

### 2.1 Plataformas Red Team existentes

| Herramienta | C2 | IA | Kernel | Modular | Código Abierto |
|-------------|-----|-----|--------|---------|----------------|
| **Cobalt Strike** | Propietario | No | No | Parcial | No |
| **Mythic** | Sí | No | No | Sí | Sí |
| **Empire** | Sí | No | No | Sí | Sí |
| **Sliver** | Sí | No | No | Sí | Sí |
| **Covenant** | Sí | No | No | Sí | Sí |
| **Havoc** | Sí | No | No | Sí | Sí |
| **X404X** | Sí | Sí | Sí | Sí | Sí |

### 2.2 Componentes individuales

Cada submódulo de X404X es un proyecto independiente con su propio historial de desarrollo:

- **Pulse-C2** (27 commits): C2 con cifrado X25519+XChaCha20, 7 módulos de agente, dashboard Vue 3.
- **Rise-Privilege** (11 commits): 12 vectores de escalada, base de datos offline de 60+ binarios GTFOBins.
- **Vault-Kernel** (8 commits): LKM con 14 comandos IOCTL para ocultación, keylogger, y backdoor.
- **Wormy-ML-Network-Worm** (13 commits): 44 exploits, motor RL (DQN + Thompson Sampling), evasión AMSI/ETW.
- **Specter-Terminal** (21 commits): Terminal ofensiva con IA offline, integración Ollama, ejecución sandboxed.

### 2.3 Novedad de X404X

Ninguna plataforma existente integra simultáneamente:
- Un motor de decisión híbrido con IA (Rules + A* + LLM offline)
- Persistencia a nivel kernel con ocultación de procesos/archivos/puertos
- Pipeline de post-explotación colaborativo (escalada → sigilo → persistencia → propagación)
- Propagación autónoma con aprendizaje por refuerzo
- Interfaz msfconsole-style con shell interactivo
- Dashboard Vue 3 con terminal xterm.js embebida y WebSocket en tiempo real

---

## 3. Arquitectura

### 3.1 Visión general

```
┌─────────────────────────────────────────────────────────────┐
│                    INTERFACES DE USUARIO                     │
│  CLI (x404x)  │  TUI (Bubble Tea)  │  Dashboard (Vue 3)    │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                     ORCHESTRATOR (Go)                        │
│  Campaign Manager │ Decision Engine │ EventBus │ WorldGraph │
│  KillChainOrchestrator │ AutoMode AI                        │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                      API SERVER (Go)                         │
│  REST 12 endpoints │ WebSocket │ CORS                       │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                      PULSE-C2 (Go)                           │
│  Server gRPC │ Agent Management │ SOCKS5 │ Tunnel           │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                   UNIFIED AGENT (Go)                         │
│  ModuleManager │ BridgeClient │ PostExploitPipeline         │
│  KillChainEngine │ VaultIOCTL │ Evasion                     │
└──┬─────────┬─────────┬──────────┬──────────┬────────────────┘
   │         │         │          │          │
┌──▼──┐ ┌───▼───┐ ┌──▼───┐ ┌───▼────┐ ┌───▼──────────┐
│Rise │ │Vault  │ │Breach│ │Python  │ │Submódulos    │
│Priv │ │Kernel │ │Entry │ │Bridge  │ │Python        │
│(Go) │ │(C/Go) │ │(C/Py)│ │(IPC)   │ │(7 módulos)   │
└─────┘ └───────┘ └──────┘ └────────┘ └──────────────┘
```

### 3.2 Componentes propios (desarrollados para este TFG)

| Componente | Lenguaje | Líneas | Descripción |
|-----------|----------|--------|-------------|
| CLI x404x | Go | ~800 | 3 modos: Cobra commands, Bubble Tea TUI, msfconsole shell |
| AppState | Go | ~460 | Estado compartido: orchestrator + bridge + C2 + DB |
| Orchestrator | Go | ~900 | Campaign Manager + Decision Engine + EventBus + WorldGraph |
| Agent | Go | ~900 | Implant unificado + Connector + BridgeClient + PostExploitPipeline |
| API Server | Go | ~750 | REST 12 endpoints + WebSocket hub + CORS |
| C2 Server | Go | ~140 | TCP listener + agent management |
| Crypto | Go | ~170 | X25519 + XChaCha20-Poly1305 (5 tests) |
| gRPC Proto | Protobuf | ~250 | 4 servicios: Agent, C2, Bridge, Common |
| Python Bridge | Python | ~600 | 11 módulos: recon, AI, privesc, persist, worm, relay… |
| Evasion Engine | Python | ~300 | 4 perfiles unificados |
| Report Generator | Python | ~200 | JSON/MD/HTML/PDF + MITRE ATT&CK |
| Dashboard Vue | Vue 3/JS | ~1,200 | 7 pestañas, 6 Pinia stores |

### 3.3 Protocolo de comunicación

Todos los componentes se comunican mediante gRPC con protobuf. El cifrado usa X25519 para el intercambio de claves y XChaCha20-Poly1305 para el cifrado autenticado de sesión.

```
Agent → C2:    agent.proto (CheckIn, CommandStream, Heartbeat)
C2 → API:      c2.proto (Campaign mgmt, agent mgmt, metrics)
Agent → Bridge: bridge.proto (Module execution, AI analysis, recon)
Shared:        common.proto (RiskLevel, KillChainPhase, types)
```

### 3.4 Motor de Decisión

El Decision Engine fusiona tres motores con pesos ponderados:

1. **Rules Engine (25%)**: 18 reglas determinísticas basadas en servicios/puertos detectados.
2. **A* Planner (35%)**: Pathfinding óptimo a través del grafo de explotación.
3. **AI Engine (40%)**: Specter-Terminal + Apex-Automation + Ollama (LLM local).

```go
func fuseDecisions(rules, planner, ai []*Decision) []*Decision {
    // Pesos: Rules=0.25, Planner=0.35, AI=0.40
    // Deduplicación por (táctica, target)
    // Ranking por confianza ponderada descendente
}
```

### 3.5 Pipeline de Post-Explotación

```go
type PostExploitPipeline struct { ... }

func (p *PostExploitPipeline) FullChain(ctx context.Context) *PostExploitResult {
    // Stage 1: Privilege Escalation
    result.RootObtained, result.PrivescVector = p.escalate(ctx)

    // Stage 2: Stealth (Vault-Kernel IOCTL)
    if result.RootObtained {
        result.StealthApplied = p.stealth(ctx)
    }

    // Stage 3: Persistence
    if result.RootObtained {
        result.PersistMethods = p.persistence(ctx)
    }

    // Stage 4: Propagation (Wormy-ML)
    result.Propagation = p.propagate(ctx)

    return result
}
```

---

## 4. Implementación

### 4.1 Estructura del monorepo

```
X404X/                           # ~450+ archivos
├── cmd/x404x/    (5 .go)        # Entry point: CLI + TUI + Console
├── core/
│   ├── agent/    (7 .go)        # Agent + PostExploit + VaultIOCTL + KillChain + Bridge
│   ├── orchestrator/ (5 .go)    # Campaign + Decision + Events + WorldGraph + AutoMode
│   ├── appstate/ (1 .go)        # Shared state singleton
│   ├── api/      (2 .go)        # REST + WebSocket
│   ├── c2server/ (1 .go)        # C2 listener
│   ├── crypto/   (2 .go)        # X25519 + XChaCha20
│   ├── proto/    (4 .proto)     # gRPC definitions
│   ├── c2/          (submodule) # Pulse-C2
│   ├── privesc/     (submodule) # Rise-Privilege
│   ├── kernel/      (submodule) # Vault-Kernel
│   └── breach/      (submodule) # Breach-Entry
├── modules/
│   ├── bridge/       # Python ↔ Go IPC
│   ├── evasion/      # Evasion engine
│   ├── report_generator.py  # Campaign reports
│   └── [7 submodules]       # Python repos
├── shared/          # Config, Logger, Types, Database
├── web/             # Vue 3 Dashboard
├── lab/             # Docker lab + CTF scenarios
├── docs/            # Architecture, TODO, CLI ref, Kill Chain Matrix
└── scripts/         # Setup, benchmark
```

### 4.2 Patrones de diseño utilizados

- **Singleton**: AppState como estado compartido global.
- **Pub/Sub**: EventBus con soporte wildcard para notificaciones en tiempo real.
- **Strategy**: Decision Engine con 3 estrategias intercambiables (Rules, A*, AI).
- **Pipeline**: PostExploitPipeline con 4 etapas secuenciales.
- **Bridge**: Python↔Go IPC mediante TCP + JSON framing.
- **Observer**: WebSocket hub notifica a clientes conectados de cambios de estado.
- **Command**: Consola msfconsole con patrón comando (use, set, exploit, back).

### 4.3 Matriz de colaboración

Cada fase de la kill chain involucra colaboración entre múltiples módulos:

| Fase | Módulos activos | Datos fluyendo |
|------|----------------|----------------|
| Recon | Horizon-Intel, Specter, Apex, WorldGraph | IPs → servicios → vulns → sugerencias IA |
| Weaponize | Decision Engine, Pulse-C2 Builder, Apex | Topología → reglas → payload compilado |
| Delivery | Breach-Entry, Evasion, Agent, Pulse-C2 | Exploit → shell → agente → check-in C2 |
| Exploit | Rise-Privilege, Specter, Apex, Vault-Kernel | Escaneo → hallazgos → sugerencias IA → root |
| Install | Vault-Kernel, Agent, Evasion, Pulse-C2 | Root → LKM → ocultación → persistencia |
| C2 | Pulse-C2, Link-Relay, Specter, Apex | Canal cifrado → relays → IA sugiere |
| Actions | Wormy-ML, Rise, Vault, Titan-Ops, BlueForge | Propagación → escalada → métricas → exfil |

---

## 5. Pruebas

### 5.1 Tests unitarios

| Módulo | Lenguaje | Tests | Cobertura |
|--------|----------|-------|-----------|
| Crypto | Go | 5/5 | KeyPair, Encrypt/Decrypt, MessageEnvelope, WrongKey, DeriveKey |
| Decision Engine | Go | Lógica real validada | 18 reglas, A* pathfinding, AI heurística |
| Python Bridge | Python | Verificación CLI | 11 módulos con --call |
| Evasion | Python | Listado de perfiles | 4 profiles verificados |
| Report Generator | Python | Generación MD | MITRE ATT&CK mapping |

### 5.2 Pruebas de integración

```bash
# Build completo
go build ./cmd/x404x/              # 21MB binary ✓

# API endpoints
curl http://localhost:8445/api/health     # {"status":"ok"} ✓
curl http://localhost:8445/api/agents     # 3 agentes demo ✓
curl http://localhost:8445/api/hosts      # 4 hosts ✓
curl http://localhost:8445/api/metrics    # Métricas completas ✓

# Consola msfconsole
echo "sessions\nhosts\nvulns\nsuggest\nexit" | ./x404x console
# 3 sessions, 4 hosts, 4 vulns, 10 suggestions ✓

# Dashboard
npm run build                      # 41 modules, 0 errors, 831ms ✓

# Python Bridge
python3 bridge.py --call privesc '{"vector":"suid"}'  # 20 SUID found ✓
python3 bridge.py --call ai_analyze '{"context":"SMB 445"}'  # EternalBlue rec ✓
```

### 5.3 Laboratorio Docker

```bash
make lab-up
# 5 contenedores: attacker, target1, target2, dashboard, ollama
# Dashboard: http://localhost:3000
# API: http://localhost:8445/api/health
```

---

## 6. Análisis Ético y Legal

### 6.1 Marco legal

El desarrollo y uso de herramientas de seguridad ofensiva está regulado por:

- **España**: Código Penal, Artículos 197-264 (delitos informáticos). Ley 8/2011 de protección de infraestructuras críticas.
- **UE**: Directiva 2013/40/UE sobre ataques contra sistemas de información.
- **Directrices ENISA**: Buenas prácticas para ethical hacking y penetration testing.
- **Marco MITRE ATT&CK**: Clasificación de tácticas y técnicas adversariales.

### 6.2 Principios éticos del proyecto

1. **Consentimiento explícito**: Todo uso requiere autorización por escrito del propietario del sistema.
2. **Entorno controlado**: El laboratorio Docker aísla completamente las operaciones.
3. **IA offline**: Ollama ejecuta modelos locales — ningún dato abandona el laboratorio.
4. **Kill switch**: Mecanismo de parada de emergencia configurable.
5. **Geofencing**: Restricción a redes privadas RFC 1918 por defecto.
6. **Auto-destrucción**: Los agentes se auto-eliminan tras N horas (configurable).
7. **Sin persistencia por defecto**: La persistencia debe activarse explícitamente.

### 6.3 Control de daños

| Control | Descripción | Default |
|---------|-------------|---------|
| Kill Switch | `--kill-switch EMERGENCY_STOP` detiene todos los agentes | Enabled |
| Geofencing | Solo redes RFC 1918 (10.x, 172.16.x, 192.168.x) | Enabled |
| Auto-Destruct | Agentes se auto-eliminan tras N horas | 2h |
| Max Infections | Límite de hosts comprometidos | 1000 |
| No Persistence | Persistencia desactivada por defecto | True |
| Dry Run | Modo simulación sin exploits reales | Available |

### 6.4 Responsabilidad

El autor no se responsabiliza del uso indebido de este software. X404X es un proyecto académico desarrollado exclusivamente para:
- Investigación en ciberseguridad
- Entornos controlados de laboratorio
- Competiciones CTF autorizadas
- Formación en defensa y detección

---

## 7. Conclusiones

### 7.1 Objetivos alcanzados

1. ✅ **Integración de 11 herramientas** en un monorepo cohesivo con Git submodules.
2. ✅ **Agente Go unificado** con soporte para módulos Go, C, y Python (vía bridge IPC).
3. ✅ **Motor de decisión híbrido** con reglas (25%), A* (35%), e IA offline (40%).
4. ✅ **Pipeline de post-explotación colaborativo** con 4 etapas automatizadas.
5. ✅ **Interfaces modernas**: CLI (Cobra + TUI + msfconsole) y Dashboard (Vue 3 + WebSocket).
6. ✅ **Laboratorio Docker** con 5 contenedores para validación.

### 7.2 Contribuciones

Este TFG aporta al estado del arte en Red Team:

- **Primera plataforma** que integra IA offline (Ollama) con decisión híbrida y kernel-level persistence.
- **Pipeline de post-explotación** donde rootkit, escalada y propagación colaboran como un sistema unificado.
- **Protocolo gRPC unificado** para comunicación entre componentes Go y Python.
- **Dashboard con terminal embebida** y visualización de grafo de red en tiempo real.
- **Metodología de integración** aplicable a otros proyectos de seguridad.

### 7.3 Trabajo futuro

- Soporte para Windows (actualmente enfocado en Linux).
- Integración con SIEM empresariales (Splunk, ELK).
- Modo cooperativo multi-operador.
- Plugin system para módulos de terceros.
- Entrenamiento RL con datos de campañas reales (anonimizados).

---

## 8. Referencias

1. MITRE ATT&CK Framework. https://attack.mitre.org/
2. gRPC Protocol Buffers. https://grpc.io/docs/protoc-installation/
3. Ollama — Local LLM. https://ollama.com/
4. XChaCha20-Poly1305 AEAD. RFC 8439.
5. Curve25519 ECDH. RFC 7748.
6. Locke, E. (2020). "Red Team Development and Operations". Independently published.
7. Kim, P. (2015). "The Hacker Playbook 2". Secure Planet LLC.
8. ENISA (2021). "Threat Landscape for Supply Chain Attacks".
9. CWE/SANS TOP 25 Most Dangerous Software Errors (2023).
10. OWASP Top 10 Web Application Security Risks (2021).

---

## Anexos

### Anexo A: Listado completo de submódulos

| Repositorio | Commits | Lenguaje | Kill Chain Phase |
|-------------|---------|----------|-----------------|
| Pulse-C2 | 27 | Go + Vue 3 | C2 |
| Rise-Privilege | 11 | Go | Exploitation |
| Vault-Kernel | 8 | C + Go | Persistence |
| Breach-Entry | 6 | C + Python | Initial Access |
| Horizon-Intel | 6 | Python | Recon |
| Specter-Terminal | 21 | Python | AI Analysis |
| Apex-Automation | 5 | Python | AI Execution |
| Wormy-ML | 13 | Python | Lateral Movement |
| Link-Relay | 6 | Python | C2 Relay |
| Titan-Operations | 4 | Python + Go | Campaign Mgmt |
| BlueForge-Suite | 8 | Python | Metrics |

### Anexo B: Endpoints API

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | /api/health | Health check |
| GET | /api/agents | List agents |
| GET | /api/agents/:id | Get agent |
| POST | /api/agents/:id/kill | Kill agent |
| GET | /api/campaigns | List campaigns |
| POST | /api/campaigns | Create campaign |
| GET | /api/campaigns/:id | Get campaign |
| POST | /api/campaigns/:id/pause | Pause campaign |
| GET | /api/hosts | List hosts |
| GET | /api/services | List services |
| GET | /api/vulnerabilities | List vulns |
| POST | /api/recon/scan | Trigger scan |
| POST | /api/ai/chat | AI chat |
| GET | /api/decisions | List decisions |
| POST | /api/decisions/:id/approve | Approve |
| GET | /api/metrics | Get metrics |
| GET | /api/blue/metrics | Blue metrics |
| WS | /ws | WebSocket events |

### Anexo C: Módulos de la consola

28 módulos registrados: 21 de exploit/auxiliary + 7 de post-explotación.

| Nombre | Tipo | Fase |
|--------|------|------|
| exploit/eternalblue | exploit | Delivery |
| exploit/bluekeep | exploit | Delivery |
| exploit/zerologon | exploit | Delivery |
| exploit/printnightmare | exploit | Delivery |
| exploit/kerberoast | exploit | Delivery |
| exploit/asreproast | exploit | Delivery |
| exploit/privesc_suid | exploit | Exploitation |
| exploit/privesc_sudo | exploit | Exploitation |
| exploit/privesc_docker | exploit | Exploitation |
| exploit/privesc_cron | exploit | Exploitation |
| exploit/log4j | exploit | Delivery |
| exploit/apache_path_traversal | exploit | Delivery |
| exploit/redis_unauth | exploit | Delivery |
| exploit/ssh_bruteforce | auxiliary | Delivery |
| exploit/smb_psexec | exploit | Actions |
| exploit/vault_kernel | post | Installation |
| auxiliary/recon_tcp | auxiliary | Recon |
| auxiliary/recon_osint | auxiliary | Recon |
| auxiliary/worm_propagate | auxiliary | Actions |
| post/persist_cron | post | Installation |
| post/persist_systemd | post | Installation |
| **post/post_exploit_full_chain** | post | 4-7 |
| **post/post_exploit_privesc** | post | 4 |
| **post/post_exploit_stealth** | post | 5 |
| **post/post_exploit_propagate** | post | 7 |
| **post/credential_dump** | post | 7 |
| **post/keylogger** | post | 6 |
| **post/evasion_apply** | post | 1-7 |

---

*Documento generado para el TFG en Ciberseguridad — Cisco NetAcad, Málaga, 2026.*
