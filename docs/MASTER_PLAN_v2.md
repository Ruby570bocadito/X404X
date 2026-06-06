# X404X — Master Improvements Plan v2.0

> Documento maestro de todas las mejoras a implementar.
> Cada tarea tiene: archivo, código, tiempo estimado.

---

## BLOQUE 1: Conectar end-to-end (~11h)

### 1.1 SQLite real (30min)
- **Archivo:** `core/appstate/state.go`
- **Driver:** `modernc.org/sqlite` (pure-Go, no CGO)
- **Cambio:** `sql.Open("sqlite", dbPath)` en initDB()
- **Verificación:** `x404x db status` muestra "SQLite | Connected | 6 tables"

### 1.2 Agent ↔ C2 gRPC real (2h)
- **Archivos:** `core/agent/connector.go`, `core/c2server/server.go`
- **Cambio:** Implementar `agent.proto` AgentService en C2 Server. El agente llama `CheckIn()` con gRPC real.
- **Verificación:** `x404x dashboard` → agente hace check-in → logs confirman session_id

### 1.3 Rise-Privilege binario (1.5h)
- **Archivo:** `core/agent/post_exploit.go`
- **Cambio:** `escalate()` llama `os/exec` al binario Rise-Privilege compilado con flag `--json`
- **Verificación:** `use post/post_exploit_full_chain; exploit` → muestra vector real

### 1.4 Vault-Kernel .ko (1h)
- **Archivo:** `lab/Dockerfile.target-linux`
- **Cambio:** Añadir compilación del LKM en el Dockerfile del target
- **Verificación:** `insmod vault_kernel.ko; ls /dev/vault_kernel`

### 1.5 Wormy-ML runtime (1h)
- **Archivo:** `modules/bridge/bridge.py`, `modules/requirements.txt`
- **Cambio:** Añadir dependencias scapy, torch. Verificar `import worm_core` funciona
- **Verificación:** `bridge.py --call worm '{"action":"scan","target":"127.0.0.1"}'`

### 1.6 Ollama real (1h)
- **Archivo:** `modules/bridge/bridge.py`
- **Cambio:** Instalar Ollama en el lab Docker. Probar `ollama.chat()` desde el handler `ai_analyze`
- **Verificación:** `bridge.py --call ai_analyze '{"context":"target SMB 445"}'` → respuesta LLM real

### 1.7 Exfiltración real (2h)
- **Archivos:** `core/agent/exfil.go` (nuevo), `modules/bridge/bridge.py`
- **Cambio:** Chunked file upload sobre C2. Handler `exfil` divide archivo en chunks de 64KB, encripta, envía.
- **Verificación:** `use post/exfiltrate; set PATH /etc/shadow; exploit` → chunks received at C2

### 1.8 Consola exploit real (2h)
- **Archivo:** `cmd/x404x/console.go`
- **Cambio:** `exploit/eternalblue` y otros ya no simulan. Llaman al Decision Engine + Pulse-C2 builder.
- **Verificación:** `use exploit/eternalblue; set RHOSTS 10.0.0.10; exploit` → compila payload + entrega

---

## BLOQUE 2: Herramientas Nivel 1 — Esenciales (~14h)

### 2.1 Credential Dumper (3h)
- **Archivos:** `modules/bridge/handlers/cred_dump.py` (nuevo), `core/appstate/state.go`
- **Cambio:** Handler Python que ejecuta mimipenguin (Linux) + parsea /etc/shadow + extrae cookies/historial. Go module registrado como `post/credential_dump`.
- **Verificación:** `use post/credential_dump; exploit` → creds aparecen en AppState + DB

### 2.2 BloodHound Collector (3h)
- **Archivos:** `modules/bridge/handlers/bloodhound.py` (nuevo), `core/appstate/state.go`
- **Cambio:** Subprocess SharpHound.exe (Windows) o Python BloodHound collector. Parse JSON output. Popula WorldGraph con relaciones AD (MemberOf, AdminTo, HasSession).
- **Verificación:** `use auxiliary/bloodhound; set DOMAIN corp.local; exploit` → grafo AD en dashboard

### 2.3 Payload Builder CLI (2.5h)
- **Archivos:** `cmd/x404x/payload.go` (nuevo)
- **Cambio:** Comando Cobra `x404x payload generate` que compila agentes Go multi-arch con opciones: os, arch, evasion level, C2 address, output format.
- **Verificación:** `x404x payload generate --os linux --arch amd64 --stealth` → binario en dist/

### 2.4 Listeners Multi-Proto (2h)
- **Archivos:** `cmd/x404x/listeners.go` (nuevo), `core/c2server/server.go`
- **Cambio:** `x404x listeners` manage HTTP/HTTPS/DNS/ICMP/SMB. Expone configuración de Pulse-C2 transports.
- **Verificación:** `x404x listeners add --type https --port 443 --cert cert.pem`

### 2.5 SOCKS5 Proxy (1.5h)
- **Archivos:** `cmd/x404x/socks5.go` (nuevo), `core/c2server/socks5.go` (nuevo)
- **Cambio:** `x404x socks5 start --port 1080`. Túnel SOCKS5 a través del C2 para pivoting.
- **Verificación:** `curl --socks5 localhost:1080 http://10.0.0.20`

### 2.6 Responder/NTLM Relay (2h)
- **Archivos:** `modules/bridge/handlers/responder.py` (nuevo), `core/appstate/state.go`
- **Cambio:** Handler Python que ejecuta Responder en la red local. Captura hashes NTLM. Auto-registra en AppState creds.
- **Verificación:** `use auxiliary/responder; exploit` → hashes capturados → `creds` muestra resultados

---

## BLOQUE 3: Diferenciadores (~16h)

### 3.1 Obfuscator/Packer Pipeline (3h)
- **Archivos:** `cmd/x404x/obfuscate.go` (nuevo), `modules/bridge/handlers/packer.py` (nuevo)
- **Cambio:** `x404x payload obfuscate --method polymorphic --packer upx --encrypt aes`. Pipeline Go+Python que aplica: mutación polimórfica (Wormy engine) + UPX packing + AES encryption.
- **Verificación:** `x404x payload obfuscate dist/agent` → hash diferente del original

### 3.2 Cleanup / Anti-Forensics (2h)
- **Archivos:** `core/agent/cleanup.go` (nuevo), `modules/bridge/handlers/cleanup.py` (nuevo)
- **Cambio:** `x404x cleanup --wipe-logs --clear-timestamps --remove-persistence --secure-delete`. Borra logs (syslog, auth.log, bash_history), timestamps de archivos, elimina persistencia instalada.
- **Verificación:** `use post/cleanup; exploit` → logs vacíos, persistencia eliminada

### 3.3 Web App Scanner (3h)
- **Archivos:** `modules/bridge/handlers/webscan.py` (nuevo), `core/appstate/state.go`
- **Cambio:** Scanner de SQLi, XSS, LFI/RFI. Basado en payloads de Wormy-ML web modules. Reporta vulnerabilidades al WorldGraph.
- **Verificación:** `use auxiliary/web_scan; set URL http://10.0.0.40; exploit` → vulns detectadas

### 3.4 Cloud Attack Modules (3h)
- **Archivos:** `modules/bridge/handlers/cloud.py` (nuevo), `core/appstate/state.go`
- **Cambio:** Módulos: `exploit/aws_imds` (metadata exfil), `exploit/azure_managed_identity` (token steal), `exploit/gcp_service_account` (key leak).
- **Verificación:** `use exploit/aws_imds; set TARGET 169.254.169.254; exploit` → creds AWS

### 3.5 Multi-Operator Support (3h)
- **Archivos:** `core/api/server.go` (mod), `web/src/views/Operators.vue` (nuevo)
- **Cambio:** WebSocket rooms por campaña. Token JWT para auth. Varios operadores en misma campaña ven mismos agentes/hosts en tiempo real.
- **Verificación:** Dos navegadores → mismo dashboard → cambios en uno visibles en otro

### 3.6 CI/CD Payload Pipeline (1h)
- **Archivo:** `.github/workflows/build-payloads.yml` (nuevo)
- **Cambio:** GitHub Actions compila payloads multi-arch (linux/windows/darwin, amd64/arm64) en cada push. Artefactos disponibles para descarga.
- **Verificación:** Push → Actions → artifacts con payloads compilados

### 3.7 Mobile/OT Stubs (1h)
- **Archivos:** `modules/bridge/handlers/mobile.py` (nuevo), `core/appstate/state.go`
- **Cambio:** Stubs iniciales para `exploit/mqtt_default`, `exploit/modbus_read`. Preparados para expansión futura.
- **Verificación:** `search mqtt` → módulo encontrado en consola

---

## BLOQUE 4: PhantomWeb (~6h)

### 4.1 Submódulo PhantomWeb (5min)
- **Archivo:** `.gitmodules`
- **Cambio:** `git submodule add <url> modules/phantom`
- **Verificación:** `ls modules/phantom/` muestra archivos

### 4.2 Bridge Handler PhantomWeb (1h)
- **Archivo:** `modules/bridge/bridge.py`
- **Cambio:** Handler `phantom` para controlar PhantomWeb C2 desde X404X bridge.
- **Verificación:** `bridge.py --call phantom '{"action":"status"}'`

### 4.3 Módulos Consola PhantomWeb (30min)
- **Archivo:** `core/appstate/state.go`
- **Cambio:** Añadir `exploit/phantom_xss`, `post/phantom_sw_persist`, `auxiliary/phantom_browser_mesh` al module registry.

### 4.4 Dashboard Tab "Browser" (2h)
- **Archivos:** `web/src/views/BrowserMesh.vue` (nuevo), `web/src/App.vue` (mod)
- **Cambio:** Nueva pestaña en dashboard: Browser Mesh Map (nodos = navegadores infectados), SW status, cookies/sessions robadas.

### 4.5 PhantomWeb C2 ↔ X404X Orchestrator (1.5h)
- **Archivo:** `core/orchestrator/events.go`
- **Cambio:** Eventos de PhantomWeb (browser_infected, cookie_stolen, sw_persisted) se publican en EventBus X404X.

### 4.6 Decision Engine Dual OS/Browser (1h)
- **Archivo:** `core/orchestrator/decision.go`
- **Cambio:** El motor evalúa si atacar por OS (EternalBlue, SSH) o por Browser (XSS, watering hole) según fingerprint del target.

---

## 📊 Resumen total

| Bloque | Tareas | Tiempo | Archivos nuevos | Archivos modificados |
|--------|--------|--------|----------------|---------------------|
| B1: End-to-End | 8 | 11h | 2 | 8 |
| B2: Nivel 1 | 6 | 14h | 6 | 4 |
| B3: Nivel 2-3 | 7 | 16h | 9 | 4 |
| B4: PhantomWeb | 6 | 6h | 3 | 5 |
| **TOTAL** | **27** | **47h** | **20** | **21** |

## 🎯 Prioridad de ejecución

```
1. BLOQUE 1 → Conectar lo que YA existe (máximo impacto, mínimo código nuevo)
2. BLOQUE 2 → Herramientas esenciales (sin esto no es un framework real)
3. BLOQUE 3 → Diferenciadores (lo hacen único frente a otras plataformas)
4. BLOQUE 4 → PhantomWeb (nuevo vector de ataque: navegador)
```
