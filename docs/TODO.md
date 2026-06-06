# X404X — Plan de Desarrollo

## ✅ v2.3 (completado)
- Arquitectura monorepo Go con 11 submódulos
- Motor de ransomware completo: scanner (19 patrones, 30+ extensiones), hydra (RSA-4096, Shamir 3/3, AES-GCM + ChaCha20)
- C2 server vía gRPC: AgentService (CheckIn, CommandStream, Heartbeat, Exfiltrate) + C2Service (10 RPCs)
- Orchestrator con Decision Engine (18 reglas, A*, AI heuristics) + WorldGraph + KillChain
- 46 módulos de consola registrados
- PhantomWeb: controller.py, Pinia store, Browser Mesh dashboard, 4 módulos consola
- Bridge Python con 20+ handlers
- TUI conectado a AppState live (agentes, hosts, vulns, campañas, decisiones)
- Console exploit handler: Bridge-first, offline-fallback, category-based dispatch
- Persistencia SQLite con 6 tablas vía modernc.org/sqlite
- 13 commits, 5 tags (v1.0 → v2.3)
- Binario 26MB, 16 tests PASS

## ✅ Bloque 1 — Maldad Psicológica y Reputacional Pura (v2.4)
- [x] 1.1 Desencriptación parcial + trampa de esperanza falsa (hope_trap.go + fake_decryptor.go)
- [x] 1.2 Destrucción de identidad digital: cookies/sesiones, hijack de cuentas, publicar contenido humillante
- [x] 1.3 RaaS inverso: panel Tor para múltiples atacantes, notas de rescate múltiples

## ✅ Bloque 2 — Propagación nivel Pandemia (v2.4)
- [x] 2.1 Gusano multi-plataforma: Windows/Linux/macOS/IoT, SSH + SMB + exploits + DDoS
- [x] 2.2 Infección cadena suministro: actualizadores, NuGet/pip/npm/git hooks, parches falsos
- [x] 2.3 Explotación cloud: AWS/Azure/GCP, EC2 instances, AMIs maliciosas, S3 buckets públicos
- [x] 2.4 Propagación Bluetooth/Wi-Fi Direct: BlueBorne, BLE MITM, APK malicioso

## ✅ Bloque 3 — Sabotaje Físico e Infraestructura (v2.4)
- [x] 3.1 Ataque SCADA/PLC: Modbus/S7 commands, overwrite ladder logic, stop PLCs
- [x] 3.2 Destrucción hardware: overvoltage, fan kill, CPU burn, BIOS corruption
- [x] 3.3 ARP spoofing + MITM proxy + captive portal + SSL strip + Root CA install

## ✅ Bloque 4 — Automutación y Resiliencia Extremas (v2.4)
- [x] 4.1 ADN digital: hibridar con DLLs legítimas, ROP gadgets, junk code, mutación JIT
- [x] 4.2 Bootkit MBR/GPT: persistencia post-formateo, SMART error fake, disk interceptor
- [x] 4.3 Blockchain C2: Bitcoin OP_RETURN, canal de comunicación inmutable y descentralizado

## ✅ Bonus — El Juego del Superviviente (v2.4)
- [x] Survival game: empleados compiten por clave de descifrado, último en pie gana

## Pendiente v2.5
- [ ] Push a GitHub con credenciales
- [ ] Tests de integración en Docker lab
- [ ] Documentación de API REST
- [ ] Dashboard web React
- [ ] Modo evasión: detectar sandboxes/VMs antes de ejecutar
- [ ] Módulo de persistencia WMI + scheduled tasks
- [ ] Cifrado de comunicaciones bridge con E2E
