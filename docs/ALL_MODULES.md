# X404X — Complete Module Catalog (v2.8)

> **125 modules across 9 versions.** Autonomous Red Team Platform.
> Build: 26.6MB ELF x86-64 · Go 1.25 · Python 3.14

---

## v2.0-v2.3: Foundation + Ransomware Base (20 modules)

| Module | Description |
|--------|-------------|
| exploit/eternalblue | MS17-010 EternalBlue SMB RCE |
| exploit/bluekeep | CVE-2019-0708 BlueKeep RDP RCE |
| exploit/zerologon | CVE-2020-1472 Netlogon EoP |
| exploit/printnightmare | CVE-2021-34527 PrintNightmare RCE/LPE |
| exploit/kerberoast | Kerberoasting TGS extraction |
| exploit/asreproast | AS-REP Roasting without creds |
| exploit/privesc_suid | SUID binary privesc via GTFOBins |
| exploit/privesc_sudo | Sudo misconfig via GTFOBins |
| exploit/privesc_docker | Docker breakout |
| exploit/privesc_cron | Cron job injection |
| exploit/log4j | CVE-2021-44228 Log4Shell |
| exploit/apache_path_traversal | CVE-2021-41773 Apache RCE |
| exploit/redis_unauth | Redis SSH key injection |
| exploit/ssh_bruteforce | SSH brute force |
| exploit/smb_psexec | SMB PSExec lateral |
| exploit/vault_kernel | LKM rootkit persistence |
| auxiliary/recon_tcp | TCP port scanner |
| auxiliary/recon_osint | OSINT gathering |
| auxiliary/worm_propagate | Network propagation |
| post/persist_cron, post/persist_systemd | Cron+Systemd persistence |

## v2.4: Advanced Ransomware (16 modules)

| Block | Module | Description |
|-------|--------|-------------|
| 1 | ransomware/hope_trap | Partial decrypt trap + forensic monitor |
| 1 | ransomware/identity_destroy | Browser session theft + account hijack |
| 1 | ransomware/raas_inverse | Inverse RaaS multi-attacker panel |
| 1 | ransomware/fake_decryptor | Fake decryptor destroys keys |
| 2 | ransomware/worm | Multi-platform worm Win/Linux/macOS/IoT |
| 2 | ransomware/supply_chain | Updater + NuGet/pip/npm/git poisoning |
| 2 | ransomware/cloud_exploit | AWS/Azure/GCP creds + instances |
| 2 | ransomware/bluetooth_prop | BlueBorne/BLE/WiFi Direct |
| 2 | ransomware/iot_botnet | IoT botnet with DDoS |
| 3 | ransomware/scada_attack | Modbus/S7/CIP PLC attacks |
| 3 | ransomware/hardware_kill | Overvoltage/fan kill/BIOS corrupt |
| 3 | ransomware/network_poison | ARP spoof/MITM/captive portal |
| 4 | ransomware/dna_mutation | DLL gene hybridization |
| 4 | ransomware/bootkit | MBR/GPT persistence |
| 4 | ransomware/blockchain_c2 | Bitcoin OP_RETURN C2 |
| Bonus | ransomware/survivor_game | Employee competition |

## v2.5: Block Z — El Umbral de la Perdición (14 modules)

| # | Module | Description |
|---|--------|-------------|
| Z.1 | blockz/genetic_evolve | Darwinian malware breeding |
| Z.2 | blockz/deepfake | CEO ONNX face+voice impersonation |
| Z.3 | blockz/scada_covert | Gradual parameter drift sabotage |
| Z.4 | blockz/firmware_worm | Router/switch/firewall firmware worm |
| Z.5 | blockz/medical_attack | Pacemaker/insulin/neurostimulator CVE |
| Z.6 | blockz/model_poison | AI model backdoor poisoning |
| Z.7 | blockz/disinformation | Email/Slack/calendar sabotage |
| Z.8 | blockz/airgap_jump | Ultrasound + LED optical exfil |
| Z.9 | blockz/post_quantum | Kyber-1024 + AES-256-GCM |
| Z.10 | blockz/deadman | 48h autonomous apocalypse |
| Z.11 | blockz/falseflag | APT Lazarus/APT29/APT41 framing |
| Z.12 | blockz/edr_kill | Silence 10 EDRs + self-deploy |
| Z.13 | blockz/financial | Insider harvest + put options |
| Z.14 | blockz/iot_chain | Hospital/factory/powergrid cascade |

## v2.6: POMDPs + AI + Evasion + Cloud + Omega (15 modules)

| Module | Description |
|--------|-------------|
| v26/pomdp | POMDPs orchestrator + God of Chaos |
| v26/ai_negotiation | AI negotiation with LLM templates |
| v26/evasion_deep | Indirect syscalls + AMSI/ETW + HW BP |
| v26/bootkit_smm | UEFI + SMM bootkit |
| v26/mobile_x | Android/iOS agents + MDM hijack |
| v26/cloud_nemesis | AWS privesc + serverless Lambda C2 |
| v26/social_c2 | Twitter/Reddit C2 + DoH tunneling |
| omega/backup_parasite | Infect ZIP/VHD/VMDK backups |
| omega/integrity_attack | Corrupt Tripwire/AIDE/FCIV checksums |
| omega/av_whitelist | AV exclusion injection via process |
| omega/multi_generational | 3-year anniversary ransom trap |
| omega/hvac_attack | Overheat server rooms via Modbus |
| omega/amt_implant | Intel AMT/AMD PSP firmware backdoor |
| omega/satcom_hijack | SATCOM firmware flash + redirect |

## v2.7: Total System Control + Phishing Arsenal (10 modules)

| Category | Module | Description |
|----------|--------|-------------|
| System | v27/uefi_bootkit | SPI flash DXE driver + NVRAM |
| System | v27/hypervisor_ring1 | Ring -1 Blue Pill/Vitriol |
| System | v27/pcie_rootkit | GPU VRAM + NIC firmware DMA |
| System | v27/kernel_instrument | eBPF syscall + ETW + BYOVD |
| System | v27/secure_boot_bypass | Shim + MOK + GRUB compromise |
| Phishing | v27/phishing_infra | DGA + Caddy + CF Workers + SOCKS5 |
| Phishing | v27/spear_phish_ai | Ollama LLM + fake M365/Google |
| Phishing | v27/anti_phish_evasion | Tokens + Safe Links bypass |
| Phishing | v27/smishing_sms | Twilio/Vonage + SS7 2FA |
| Phishing | v27/vishing_voice | Voice clone + Twilio calls |

## v2.8: Ultimate Malice Arsenal (24 modules)

| # | Module | Description |
|---|--------|-------------|
| 10 | v28/iot_identity_theft | Steal X.509 IoT device certs |
| 11 | v28/false_memory | Forge Teams/Slack/email evidence |
| 12 | v28/thousand_cuts | 90-day service degradation |
| 13 | v28/patchguard_bypass | KeBugCheckEx hook + DKOM |
| 14 | v28/keyboard_led | Morse LED exfil for air-gap |
| 15 | v28/zombie_army | Social media smear campaign |
| 16 | v28/legacy_poison | Frame victim for illegal activities |
| 1* | v28/seo_sabotage | Black-hat SEO fake sites |
| 2* | v28/fake_vulns | Trap 0-day vulnerabilities in repos |
| 3* | v28/inception_hv | Nested hypervisor matrioska |
| 4* | v28/isp_bgp | BGP prefix hijack Internet-scale |
| 5* | v28/anti_attribution | Clone identity to frame victim |
| 6* | v28/power_grid_harmonics | Resonant transformer destruction |
| 7* | v28/time_lock | 30-min window or key destroyed |
| 8* | v28/vr_spyware | VR passthrough + subliminal |
| 9* | v28/global_ai_poison | HuggingFace/Kaggle/OpenML backdoor |
| 10* | v28/cdn_injection | Cloudflare/Akamai/Fastly hijack |
| 11* | v28/bio_cyber_dna | Alter synthetic DNA orders |
| 12* | v28/browser_parasite | Hidden browser extension parasite |
| 13* | v28/fake_documents | Forge purchase orders/resolutions |
| 14* | v28/sound_panic | IP speaker fire/evacuation alarm |
| 15* | v28/emotional_encrypt | Sentimental file encryption |
| 16* | v28/false_redemption | Fake decryptor + permanent backdoor |

## Deployment via Dashboard

All 125 modules are available via the Vue3 Dashboard at `/api/modules`. Each module can be pushed to any connected victim agent via the dashboard UI or CLI:

```
x404x deploy victim01 ransomware/execute,blockz/genetic_evolve,v27/uefi_bootkit
x404x c2 listen  # listen-only mode
x404x modules list  # list all 125 modules
x404x victims list  # list registered victims
```

## Architecture

```
CLI / TUI / Dashboard (Vue3)
        │
   Orchestrator (POMDPs + Decision Engine)
        │
   C2 Server (gRPC AgentService + C2Service)
        │
   Agent (Go binary 26.6MB)
   ├── Module Registry (125 modules)
   ├── Deployment Manager (per-victim)
   ├── Python Bridge (20+ handlers)
   └── Kernel Modules (eBPF, BYOVD, SMM)
```
