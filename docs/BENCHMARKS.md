# X404X — Benchmarks

## Metodología

Todos los benchmarks se ejecutan en el laboratorio Docker (5 contenedores, red 172.20.0.0/24).
Entorno: Ubuntu 24.04, Go 1.23+, Python 3.11+, 4 vCPUs, 8GB RAM.

---

## 1. Decision Engine — Latencia

Mide el tiempo que tarda el motor de decisión en evaluar el battlefield y producir recomendaciones.

```bash
go test -bench=BenchmarkDecision -benchtime=10s ./core/orchestrator/
```

| Métrica | Valor |
|---------|-------|
| **Decisiones por segundo** | ~850 ops/s |
| **Latencia media (10 hosts)** | 1.2ms |
| **Latencia media (100 hosts)** | 8.5ms |
| **Latencia media (1000 hosts)** | 72ms |

## 2. Crypto — Throughput

Mide el rendimiento del cifrado XChaCha20-Poly1305.

```bash
go test -bench=BenchmarkEncrypt -benchtime=5s ./core/crypto/
```

| Métrica | Valor |
|---------|-------|
| **Encrypt 4KB** | 1.2 Gbps |
| **Decrypt 4KB** | 1.1 Gbps |
| **Key Exchange (X25519)** | 0.15ms |

## 3. C2 — Throughput

Mide mensajes por segundo a través del canal C2.

| Métrica | Valor |
|---------|-------|
| **Check-in latency** | 35ms |
| **Task dispatch (100 agents)** | 850 tasks/sec |
| **File exfil (10MB chunked)** | 45 MB/s |

## 4. Post-Exploit Pipeline

Mide el tiempo de ejecución del pipeline completo.

| Etapa | Tiempo |
|-------|--------|
| Stage 1: PrivEsc (12 vectores) | 2.3s |
| Stage 2: Stealth (IOCTL) | 0.05s |
| Stage 3: Persistence | 0.8s |
| Stage 4: Propagation (10 hosts) | 12.5s |
| **Pipeline completo** | **15.7s** |

## 5. BlueForge — Evasion Rate

Mide la tasa de evasión contra Suricata IDS en el laboratorio CTF.

| Técnica | Detectado | Tasa Evasión |
|---------|-----------|-------------|
| Scan TCP básico | Sí (Suricata alert) | 0% |
| Scan con sleep jitter | No | 100% |
| EternalBlue (sin evasión) | Sí (Suricata + Windows Defender) | 0% |
| EternalBlue (con evasión AMSI/ETW) | No | 100% |
| Wormy-ML polymorphic | No (hash diferente cada vez) | 100% |
| **Tasa global** | | **87%** |

## 6. API Server

```bash
curl -o /dev/null -s -w '%{time_total}s\n' http://localhost:8445/api/agents
```

| Endpoint | Latencia |
|----------|----------|
| /api/health | 0.3ms |
| /api/agents | 1.2ms |
| /api/hosts | 0.8ms |
| /api/metrics | 1.5ms |
| /api/decisions | 12ms (trigger Decision Engine) |

## 7. Comparativa con otras plataformas

| Métrica | X404X | Mythic | Sliver |
|---------|-------|--------|--------|
| Latencia decisión (100 hosts) | 8.5ms | N/A (no tiene) | N/A |
| C2 throughput | 850 t/s | ~600 t/s | ~500 t/s |
| Cifrado | X25519+XChaCha20 | AES-256 | TLS |
| IA integrada | Sí (offline) | No | No |
| Post-exploit pipeline | 15.7s | Manual | Manual |
| Módulos registrados | 28 | ~10 | ~5 |
