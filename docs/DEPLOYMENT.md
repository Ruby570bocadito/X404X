# X404X — Deployment Guide (v3.2)

## Quick Start

```bash
# Clone
git clone https://github.com/Ruby570bocadito/X404X.git && cd X404X

# Build
go build -o x404x ./cmd/x404x/

# Install Python deps
pip install -r requirements.txt

# Start
./x404x --interactive
```

## Prerequisites

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| Go | 1.25+ | 1.25+ |
| Python | 3.11+ | 3.13+ |
| Node.js | 22.x (for dashboard) | 22.12+ |
| Ollama | latest (for AI) | latest |
| Docker | 27+ (optional) | 27.3+ |

## Build

### Binary
```bash
go build -ldflags="-s -w" -o x404x ./cmd/x404x/
```

### Obfuscated Build (Garble)
```bash
bash scripts/build_obfuscated.sh
```

### Dashboard
```bash
cd web && npm install && npm run build
```
The built dashboard is in `web/dist/` and served by the Go server.

## Run Modes

### Interactive Console
```bash
./x404x --interactive
```

### Headless (CLI mode)
```bash
./x404x campaign start --name "Demo" --target-scope "10.0.0.0/24"
```

### With Dashboard
```bash
./x404x --api-port 8443 --dashboard
```

### With Python Bridge
```bash
./x404x --bridge --interactive
```

## Configuration

Configuration is in `config.yaml` or via CLI flags:

```yaml
agent:
  bridge_port: 9100       # Python bridge gRPC port
  c2_addr: "127.0.0.1:8443"
  checkin_interval: 30s

api:
  port: 8443
  enable_ws: true
  cors: true

ai:
  enabled: true
  ollama_host: "127.0.0.1"
  ollama_port: 11434
  model: "llama3.1:8b"
  temperature: 0.8
  auto_approval: false

campaign:
  auto_approve: false
  max_infections: 1000
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `X404X_API_PORT` | `8443` | API server port |
| `X404X_BRIDGE_PORT` | `9100` | Python bridge port |
| `X404X_C2_ADDR` | `127.0.0.1:8443` | C2 server address |
| `OLLAMA_HOST` | `127.0.0.1:11434` | Ollama API host |
| `X404X_AI_ENABLED` | `true` | Enable AI engine |

## Docker Deployment

```bash
# Build image
docker build -t x404x .

# Run
docker run -p 8443:8443 -p 9100:9100 \
  -v ./config.yaml:/app/config.yaml \
  x404x --bridge --dashboard
```

## Directory Layout After Deploy

```
.
├── x404x                    # Binary
├── config.yaml              # Configuration
├── modules/bridge/          # Python bridge handlers
├── plugins/                 # Specialized modules
├── web/dist/                # Dashboard static files
├── reports/                 # Campaign output
└── /tmp/x404x_attack_report.json  # ATT&CK coverage report
```

## Verifying Deployment

```bash
# Health check
curl http://localhost:8443/api/health

# List modules
curl http://localhost:8443/api/modules

# Dashboard
open http://localhost:8443
```

## Supported Platforms

| OS | Go Binary | Python Bridge | Dashboard |
|----|-----------|---------------|-----------|
| Linux (amd64) | ✅ | ✅ | ✅ |
| Linux (arm64) | ✅ | ✅ | ✅ |
| macOS (amd64) | ✅ | ✅ | ✅ |
| macOS (arm64) | ✅ | ✅ | ✅ |
| Windows (amd64) | ✅ | ⚠️ via WSL | ✅ |
