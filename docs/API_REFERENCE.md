# X404X — API Reference (v3.2)

## Endpoints

All endpoints are served by the Go API server at `http://localhost:8443` (configurable via `--api-port`).
The dashboard proxies `/api` → Go server and `/ws` → WebSocket hub.

### REST API

#### Agents

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/agents` | List all agents. Query: `?campaign_id=X` |
| `GET` | `/api/agents/:id` | Get agent details |
| `POST` | `/api/agents/:id/kill` | Kill an agent. Body: `{"reason":"..."}` |

#### Campaigns

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/campaigns` | List all campaigns |
| `POST` | `/api/campaigns` | Create campaign. Body: `{"name","target_scope","goal","profile","auto_approve"}` |
| `GET` | `/api/campaigns/:id` | Get campaign details |
| `POST` | `/api/campaigns/:id/pause` | Pause a running campaign |
| `POST` | `/api/campaigns/:id/resume` | Resume a paused campaign |

#### Recon

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/hosts` | List discovered hosts. Query: `?campaign_id=X` |
| `GET` | `/api/services` | List discovered services. Query: `?campaign_id=X` |
| `GET` | `/api/vulnerabilities` | List vulnerabilities. Query: `?campaign_id=X` |
| `POST` | `/api/recon/scan` | Trigger scan. Body: `{"target":"10.0.0.1","mode":"quick"}` |

#### AI / Decisions

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/ai/chat` | Chat with AI. Body: `{"prompt":"analyze target..."}` |
| `GET` | `/api/decisions` | List pending decisions. Query: `?campaign_id=X` |
| `POST` | `/api/decisions/:id/approve` | Approve a decision |
| `POST` | `/api/decisions/:id/reject` | Reject a decision |

#### Metrics

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/metrics` | Get campaign KPIs. Query: `?campaign_id=X` |
| `GET` | `/api/blue/metrics` | Get BlueForge detection metrics |

#### Payloads

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/payload/generate` | Generate implant payload. Body: `{"os","arch","format","lhost","lport","amsi","unhook","encoder"}` |

#### Phantom (Browser Mesh)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/phantom/status` | Get phantom mesh status |
| `GET` | `/api/phantom/nodes` | List phantom browser nodes |
| `POST` | `/api/phantom/:action` | Execute phantom action (inject, steal, etc.) |

#### Config

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/config/ai` | Get AI configuration |
| `PUT` | `/api/config/ai` | Update AI config. Body: `{"model","temperature"}` |

#### Admin

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/health` | Health check: `{"status":"ok","version":"3.2","uptime":3600}` |
| `GET` | `/api/modules` | List all registered modules |
| `POST` | `/api/modules/push` | Push module to agent |
| `GET` | `/api/sessions` | List active sessions |
| `GET` | `/api/creds` | List captured credentials |

### WebSocket API

Connect to `ws://localhost:8443/ws?campaign_id=X`

**Event types broadcast by the server:**

| Event | Payload |
|-------|---------|
| `agent.checkin` | `{agent_id, hostname, os, ip, timestamp}` |
| `agent.dead` | `{agent_id, reason, timestamp}` |
| `campaign.started` | `{campaign_id, name, phase}` |
| `campaign.paused` | `{campaign_id}` |
| `campaign.resumed` | `{campaign_id}` |
| `decision.made` | `{decision_id, tactic, technique, mitre_id, confidence}` |
| `host.discovered` | `{ip, hostname, os, ports}` |
| `vuln.found` | `{cve, severity, target_ip, service}` |
| `recon.scan_complete` | `{target, hosts_found, vulns_found}` |
| `recon.scan_error` | `{target, error}` |
| `phase.changed` | `{campaign_id, from, to, progress}` |
| `blue.alert` | `{tool, alert_type, timestamp}` |
| `exploit.success` | `{target, exploit, cve}` |
| `exploit.failure` | `{target, exploit, error}` |
| `credential.captured` | `{username, domain, source}` |

### gRPC Services

The C2 server exposes three gRPC services at the configured port:

**AgentService** (`x404x.v1.AgentService`)
```
rpc CheckIn(CheckInRequest) returns (CheckInResponse)
rpc CommandStream(stream AgentMessage) returns (stream ServerMessage)
rpc Heartbeat(HeartbeatRequest) returns (HeartbeatResponse)
rpc Exfiltrate(stream ExfilChunk) returns (ExfilAck)
```

**C2Service** (`x404x.v1.C2Service`)
```
rpc ListAgents(ListAgentsRequest) returns (ListAgentsResponse)
rpc GetAgent(GetAgentRequest) returns (AgentInfo)
rpc KillAgent(KillAgentRequest) returns (KillAgentResponse)
rpc CreateCampaign(CreateCampaignRequest) returns (Campaign)
rpc GetCampaign(GetCampaignRequest) returns (Campaign)
rpc ListCampaigns(ListCampaignsRequest) returns (ListCampaignsResponse)
rpc PauseCampaign(PauseCampaignRequest) returns (Campaign)
rpc ResumeCampaign(ResumeCampaignRequest) returns (Campaign)
rpc DecisionFeed(stream DecisionUpdate) returns (stream DecisionAck)
rpc GetMetrics(MetricsRequest) returns (MetricsResponse)
```

**BridgeService** (`x404x.v1.BridgeService`)
```
rpc ExecuteModule(ModuleRequest) returns (ModuleResponse)
rpc AIAnalyze(AIAnalyzeRequest) returns (stream AIAnalyzeResponse)
rpc ReconStream(ReconRequest) returns (stream ReconResponse)
rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse)
```

## Python Bridge Handlers

The bridge exposes 107 handlers across 12 registration files under `modules/bridge/handlers/`.

### Registry Groups

| Group | Handlers | File |
|-------|----------|------|
| `ransomware` | execute, scan, encrypt, exfil, status, decrypt, generate_note, propagate, destruct | `ransomware.py` |
| `ransomware_advanced` | hope_trap, identity_destroy, raas_panel, fake_decryptor, worm_deploy, supply_chain, cloud_exploit, bluetooth_prop, iot_botnet, scada_attack, hardware_kill, network_poison, captive_portal, dna_mutate, bootkit, blockchain_c2, survivor_game | `ransomware_advanced.py` |
| `ransomware_v26` | pomdp_decide, ai_negotiate, evasion_deep, bootkit_smm, mobile_x, cloud_nemesis, social_c2, block_omega | `ransomware_v26.py` |
| `ransomware_v27` | uefi_bootkit, hypervisor_ring1, pcie_rootkit, kernel_instrument, secure_boot_bypass, phishing_infra, spear_phish_ai, anti_phish_evasion, smishing_sms, vishing_voice | `ransomware_v27.py` |
| `ransomware_v28` | iot_identity_theft, false_memory, thousand_cuts, patchguard_bypass, keyboard_led, zombie_army, legacy_poison, seo_sabotage, fake_vulns, inception_hv, isp_bgp, anti_attribution, power_grid_harmonics, time_lock, vr_spyware, global_ai_poison, cdn_injection, bio_cyber_dna, browser_parasite, fake_documents, sound_panic, emotional_encrypt, false_redemption | `ransomware_v28.py` |
| `ransomware_v29` | hdd_firmware_destroy, vrm_overvoltage, acoustic_resonance, psu_corrupt, usb_killer, robot_sabotage, centrifuge_resonance, ui_shell_fake, deepfake_hallucinate, network_ghosts, medical_tamper, intel_me_flash, smm_handler, microcode_corrupt, nic_persist, mft_bitmap, backup_prune, journal_poison, dns_poison, bgp_phantom, ldap_intermittent, digital_thermite, honey_token, access_log_wipe | `ransomware_v29.py` |
| `ransomware_v210` | apocalipsis, phantom_evasion | `ransomware_v210.py` |
| `ransomware_blockz` | genetic_evolve, deepfake_generate, scada_covert, firmware_worm, medical_attack, model_poison, disinformation, airgap_exfil, post_quantum, deadman_arm, falseflag_plant, edr_kill, financial_crash, iot_chain | `ransomware_blockz.py` |
| `phase_1_4` | byovd_loader, dkom, amsi_patch, etw_patch, syscall_proxy, hollowing, unhook_ntdll, evasion_misc, otp, sandbox_detect, network_covert, persist_scheduled, persist_wmi, persist_registry, stego_config, x25519_wireguard, quic_tunnel, webrtc_p2p, beacon_dns, beacon_https, beacon_smb, obfuscate_code, packer_upx, crypter_xor, embed_payload, rsrc_hide, connect_back, bind_shell, pivot_socks5, relaying, ai_target, ai_phishing, ai_deepfake, ai_vishing, c2_waterfall, c2_cloudfront | `phase_1_4.py` |

### Calling a Handler

From Go:
```go
resp := bridge.CallRaw(ctx, "ransomware", "scan", map[string]interface{}{
    "root": "/home",
    "max_files": 500,
})
```

From Python:
```python
from handlers.ransomware import handle_scan
result = handle_scan({"root": "/home", "max_files": 500})
```

All handlers accept `params: dict` and return `dict` with at minimum `{"success": bool}`. Handlers respect `{"simulation": true}` (default) to avoid destructive operations.
