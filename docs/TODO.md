# X404X — Plan de Desarrollo

## ✅ v2.3 (completado)
- Arquitectura monorepo Go con 11 submódulos
- Motor de ransomware completo: scanner (19 patrones, 30+ extensiones), hydra (RSA-4096, Shamir 3/3, AES-GCM + ChaCha20)
- C2 server via gRPC: AgentService (CheckIn, CommandStream, Heartbeat, Exfiltrate) + C2Service (10 RPCs)
- Orchestrator con Decision Engine (18 reglas, A*, AI heuristics) + WorldGraph + KillChain
- 46 módulos de consola registrados
- PhantomWeb: controller.py, Pinia store, Browser Mesh dashboard, 4 módulos consola
- Bridge Python con 20+ handlers
- TUI conectado a AppState live (agentes, hosts, vulns, campanas, decisiones)
- Console exploit handler: Bridge-first, offline-fallback, category-based dispatch
- Persistencia SQLite con 6 tablas via modernc.org/sqlite
- 13 commits, 5 tags (v1.0 -> v2.3)
- Binario 26MB, 16 tests PASS

## ✅ v2.4 — 14 Módulos Avanzados de Ransomware (completado)

### Bloque 1 — Maldad Psicológica y Reputacional Pura
- [x] 1.1 Desencriptación parcial + trampa de esperanza falsa (hope_trap.go + fake_decryptor.go)
- [x] 1.2 Destrucción de identidad digital: cookies/sesiones, hijack de cuentas, publicar contenido humillante
- [x] 1.3 RaaS inverso: panel Tor para múltiples atacantes, notas de rescate múltiples

### Bloque 2 — Propagación nivel Pandemia
- [x] 2.1 Gusano multi-plataforma: Windows/Linux/macOS/IoT, SSH + SMB + exploits + DDoS
- [x] 2.2 Infeccion cadena suministro: actualizadores, NuGet/pip/npm/git hooks, parches falsos
- [x] 2.3 Explotacion cloud: AWS/Azure/GCP, EC2 instances, AMIs maliciosas, S3 buckets publicos
- [x] 2.4 Propagacion Bluetooth/Wi-Fi Direct: BlueBorne, BLE MITM, APK malicioso

### Bloque 3 — Sabotaje Fisico e Infraestructura
- [x] 3.1 Ataque SCADA/PLC: Modbus/S7 commands, overwrite ladder logic, stop PLCs
- [x] 3.2 Destruccion hardware: overvoltage, fan kill, CPU burn, BIOS corruption
- [x] 3.3 ARP spoofing + MITM proxy + captive portal + SSL strip + Root CA install

### Bloque 4 — Automutacion y Resiliencia Extremas
- [x] 4.1 ADN digital: hibridar con DLLs legitimas, ROP gadgets, junk code, mutacion JIT
- [x] 4.2 Bootkit MBR/GPT: persistencia post-formateo, SMART error fake, disk interceptor
- [x] 4.3 Blockchain C2: Bitcoin OP_RETURN, canal de comunicacion inmutable y descentralizado

### Bonus
- [x] Survival game: empleados compiten por clave de descifrado, ultimo en pie gana

---

## ✅ v2.5 — Block Z: El Umbral de la Perdición (14 módulos)
- [x] Z.1 Genetic Evolution: Darwinian malware breeding with DLL genes
- [x] Z.2 Deepfake Identity: ONNX face+voice CEO impersonation pipeline
- [x] Z.3 Covert SCADA: Gradual parameter drift over months
- [x] Z.4 Firmware Worm: Router/switch/firewall backdoor surviving updates
- [x] Z.5 Medical Attack: Pacemaker/insulin/neurostimulator exploits
- [x] Z.6 AI Model Poison: Backdoor classifiers + label flipping
- [x] Z.7 Disinformation: Email/Slack/calendar corporate sabotage
- [x] Z.8 Air-Gap Jump: Ultrasound + LED optical exfiltration
- [x] Z.9 Post-Quantum: Kyber-1024 + AES-256-GCM encryption
- [x] Z.10 Dead Man Switch: 48h autonomous apocalypse trigger
- [x] Z.11 False Flag: APT framing + Mandiant forensic report generation
- [x] Z.12 EDR Control: Kill 10 EDRs + self-deploy via consoles
- [x] Z.13 Financial Attack: Insider harvest + put options + stock crash
- [x] Z.14 IoT Chain: Hospital/factory/powergrid cascading damage

## 🔜 v2.6 — Prioridad Alta (propuesto)

### 🔧 Orquestador — El Cerebro v2
- [ ] Planificacion dinamica con POMDPs (Partially Observable Markov Decision Process)
- [ ] Modelos TensorFlow Lite que predicen probabilidad de deteccion por accion
- [ ] Modo "Dios del Caos": inyectar fallos falsos para enganar al SOC
- [ ] Plan B automatico cuando el blue team responde

### 🕵️ Horizon-Intel — Reconocimiento Total
- [ ] Integracion con APIs: Shodan, Censys, BinaryEdge
- [ ] Busqueda de credenciales en leaks: Have I Been Pwned, Dehashed
- [ ] Recon pasivo de AD: LDAP anonimo, NetBIOS, mDNS, SSDP
- [ ] Generacion automatica de spear-phishing con Ollama LLM local

### 🚪 Breach-Entry — Acceso Inicial 2.0
- [ ] Modulos por protocolo: RDP (BlueKeep/CVE-2019-0708), SMB (EternalBlue/SMBGhost), HTTP (ProxyShell/Log4Shell), SSH, WinRM, VNC
- [ ] USB Rubber Ducky: generador de payloads ofuscados
- [ ] Ataque a impresoras: IPP, LPD, PJL como puente a red interna

### 🧠 AI — La Mente Maestra
- [ ] Agente de negociacion automatica via LLM (Evil ChatGPT mode)
- [ ] Generacion de exploits en caliente con Ollama
- [ ] Imitacion de trafico real (Netflix, Teams, Windows Update)

### 🐛 Wormy-ML — Propagacion Mutante
- [ ] Motor de mutacion polimorfica: cada gusano es unico (hash cambia en cada salto)
- [ ] Propagacion Bluetooth/Wi-Fi Direct con APK malicioso + Office macro
- [ ] Infeccion de Docker/K8s: pods maliciosos, envenenar imagenes base, API server takeover

### 🛡️ Evasion — Invisibilidad Absoluta
- [ ] Evasion de EDR via kernel hooking: parchear ObRegisterCallbacks, PsSetCreateProcessNotifyRoutine
- [ ] Hardware breakpoints (DR0-DR7) para API unhooking sin tocar .text
- [ ] Syscalls indirectas + codificacion XOR dinamico
- [ ] Suplantacion de firma digital en vivo: firmar binarios en memoria

### 💀 Persistence — Zombie Mode
- [ ] Bootkit UEFI + SMM (System Management Mode) fuera del alcance del SO
- [ ] Persistencia en dispositivos PCIe: GPU/FPGA firmware
- [ ] WMI Event Subscription + DGA para nombres de evento rotativos

### 📡 C2 — Comunicaciones del Infierno
- [ ] C2 sobre Blockchain extendido (BTC/ETH) con smart contracts
- [ ] C2 sobre redes sociales: tweets cifrados, respuestas en Reddit/Pinterest
- [ ] DoH tunneling (IP sobre DNS sobre HTTPS via Cloudflare/Google)

### 🧬 BlueForge-Suite — Validacion Ofensiva
- [ ] Simulador de Blue Team agresivo: contramedidas automaticas
- [ ] Generador de informes de cobertura ATT&CK

### 📱 MOBILE-X — Agente Android/iOS
- [ ] Agente nativo Java/Kotlin (Android) y Swift (iOS) via mismo gRPC
- [ ] Capacidades: audio, camara, SMS, GPS
- [ ] Explotacion de MDM: robar certificado, desplegar politicas maliciosas

### 🌩️ CLOUD-NEMESIS — Dominacion Cloud
- [ ] Escalada de privilegios en AWS/Azure/GCP
- [ ] Serverless C2: funciones Lambda/Azure Functions efimeras

### 🧰 FORGE — Taller de Exploits
- [ ] Integracion AFL++ y LibFuzzer para fuzzing en tiempo real
- [ ] Base de datos de gadgets ROP por SO y version

### 💣 Ransomware — El Toque Final v2
- [ ] Modo bomba de tiempo con chantaje progresivo: filtraciones diarias
- [ ] Descifrador con backdoor perpetuo: la victima queda esclavizada

---

## Pendiente Infraestructura
- [ ] Push a GitHub con credenciales
- [ ] Tests de integracion en Docker lab
- [ ] Documentacion de API REST
- [ ] Dashboard web React
- [ ] Modo evasion: detectar sandboxes/VMs antes de ejecutar
- [ ] Modulo de persistencia WMI + scheduled tasks
- [ ] Cifrado de comunicaciones bridge con E2E
- [ ] .gitignore actualizado (__pycache__, .DS_Store, *.log, etc.)
