# X404X — Kill Chain Collaboration Matrix

Cada fase de la kill chain muestra qué módulos colaboran y cómo fluyen los datos.

---

## FASE 1: RECONNAISSANCE
| Módulo | Rol | Entrada | Salida |
|--------|-----|---------|--------|
| **Horizon-Intel** | Escaneo de superficie | Target IP/CIDR | Hosts, puertos, servicios, OS guess |
| **Specter** | Análisis IA | Datos de escaneo | "DC(10.0.0.10) high-value, SMB+Kerberos expuestos" |
| **Apex** | Selección de vector | Vulnerabilidades | "EternalBlue MS17-010 → 10.0.0.10 (conf=0.92)" |
| **WorldGraph** | Topología | Hosts + servicios | Grafo con nodos, edges, asset values |

## FASE 2: WEAPONIZATION
| Módulo | Rol | Entrada | Salida |
|--------|-----|---------|--------|
| **Decision Engine** | Fusión de motores | Topología + vulns | Decisión rankeada: Rules(25%)+A*(35%)+AI(40%) |
| **Pulse-C2 Builder** | Compilar payload | Vector + OS + arch | Binario compilado (GOOS/GOARCH) |
| **Apex** | Seleccionar payload | Decisión del engine | "eternalblue + agent embedded" |

## FASE 3: DELIVERY
| Módulo | Rol | Entrada | Salida |
|--------|-----|---------|--------|
| **Breach-Entry** | Exploit inicial | CVE + target | Shell / acceso |
| **Evasion** | Bypass AV/EDR | Payload | AMSI/ETW bypass aplicado |
| **Agent** | Despliegue | Shell | Check-in a C2 vía gRPC + X25519 |
| **Pulse-C2** | Recibir check-in | Agent ID + hostname | Session ID asignado |
| **Specter** | Análisis post-acceso | Contexto del agente | "OS: Win2019, user: NT\SYSTEM" |

## FASE 4: EXPLOITATION (PrivEsc)
| Módulo | Rol | Entrada | Salida |
|--------|-----|---------|--------|
| **Rise-Privilege** | Escalada | Usuario no-root | 12 vectores escaneados + exploit |
| **Specter** | Análisis hallazgos | SUID/sudo/cron | "python3 SUID → GTFOBins shell (SAFE, 0.95)" |
| **Apex** | Ordenar vectores | Findings | Prioridad: safe → low → medium → high |
| **Vault-Kernel** | Escalada kernel | PID | IOCTL_GIVE_ROOT → cred manipulation |
| **Dashboard** | Visualización | Evento root | Kill chain avanza en tiempo real |

## FASE 5: INSTALLATION (Persistencia)
| Módulo | Rol | Entrada | Salida |
|--------|-----|---------|--------|
| **Vault-Kernel** | Rootkit | Root | Ocultar PID, file, port, module + keylogger + backdoor |
| **Agent** | Persistencia userland | Root | cron @reboot + systemd + .bashrc hook |
| **Evasion** | Sigilo | Perfil stealth | Sleep jitter + polymorphic + sandbox detect |
| **Pulse-C2** | Confirmación | Persistencia OK | Notificación al dashboard |

## FASE 6: COMMAND & CONTROL
| Módulo | Rol | Entrada | Salida |
|--------|-----|---------|--------|
| **Pulse-C2** | Canal C2 | Session ID | Canal cifrado X25519+XChaCha20Poly1305 |
| **Link-Relay** | Cadena relays | C2 server | Agent → Relay1 → Relay2 → C2 |
| **Specter** | Monitoreo | Telemetría | "Usuario admin logueado, Outlook abierto..." |
| **Apex** | Sugerir acciones | Credenciales capturadas | "Lateral a 10.0.0.20 vía SSH" |
| **Dashboard** | Dashboard | Eventos C2 | Uptime, procesos ocultos, teclas capturadas, conexiones |

## FASE 7: ACTIONS ON OBJECTIVE
| Módulo | Rol | Entrada | Salida |
|--------|-----|---------|--------|
| **Wormy-ML** | Propagación | Red objetivo | 44 exploits, RL engine, scan → infect → repeat |
| **Rise-Privilege** | Escalada en nuevos hosts | Cada nuevo host | Auto-root en cada máquina |
| **Vault-Kernel** | Sigilo en nuevos hosts | Root en cada host | Ocultación + persistencia en cadena |
| **Titan-Operations** | Campaña global | Todos los agentes | director.go: goal tracking, progress % |
| **BlueForge-Suite** | Métricas | Eventos de red | "Suricata: 2 detecciones, Wormy-ML: bypass 87%" |
| **Data Exfiltration** | Exfiltración | Archivos objetivo | Chunked encrypted transfer vía C2 |
| **Specter** | Análisis final | Datos completos | "8/23 hosts. Objetivo: PENDIENTE (DC2)" |

---

## Flujo de datos entre módulos

```
Horizon-Intel ──→ WorldGraph ──→ Decision Engine ──→ Pulse-C2 Builder
                                                         │
Breach-Entry ←── Evasion ←──────────────────────────────┘
     │
     ▼
   Agent ──→ Pulse-C2 ──→ Orchestrator ──→ Dashboard
     │           │              │
     │           │              ├──→ KillChainOrchestrator (fase auto-advance)
     │           │              └──→ EventBus → WebSocket → Dashboard
     │           │
     ├──→ Rise-Privilege ──→ Vault-Kernel ──→ Wormy-ML
     │         │                  │               │
     │         └─ root ──────────┘               │
     │                                            │
     └──→ Link-Relay ──→ Titan-Ops ──→ BlueForge │
                                                 │
     ┌───────────────────────────────────────────┘
     ▼
  PostExploitPipeline.FullChain()
  ┌─────────────────────────────────────┐
  │ Stage 1: escalate()                 │
  │ Stage 2: stealth()                  │
  │ Stage 3: persistence()              │
  │ Stage 4: propagate()                │
  │ Stage 5: askAI() → report to C2     │
  └─────────────────────────────────────┘
```
