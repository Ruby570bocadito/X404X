# X404X

Autonomous red team operations platform. Go core with Python bridge and Vue 3 dashboard.

## Features

- **Kill chain coverage**: recon, initial access, execution, persistence, privilege escalation, lateral movement, exfiltration
- **Multi-language**: Go (engine, agent, C2), Python (bridge, handlers, plugins), Vue 3 (dashboard)
- **Full-spectrum modules**: 45+ modules across evasion, C2, propagation, AI, and cross-platform
- **AI orchestration**: Q-learning FSM, federated learning, offline LLM (Ollama) for decision support
- **Post-quantum crypto**: Kyber-1024 KEM, X25519 ECDH, Ed25519 signing, XChaCha20-Poly1305
- **Dashboard**: Vue 3 with real-time WebSocket, D3 force graph, xterm.js terminal

## Quick Start

```
# Install dependencies
make setup

# Build CLI
make build

# Run console
./dist/x404x

# Start Docker lab
make lab-up
```

See [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) for full setup instructions.

## Architecture

```
CLI (x404x)  |  TUI (Bubble Tea)  |  Dashboard (Vue 3)
           |
      Orchestrator (Go)
      Campaign mgmt, decisions, events, world graph
           |
      Agent (Go)
      Module manager, bridge client, post-exploit pipeline
           |
      Python Bridge (IPC)
      ~170 handlers across recon, privesc, evasion, AI, exfil
```

## Project Structure

```
cmd/x404x/          CLI entry point
cmd/implant/        C2 agent binary
internal/           Go engine (crypto, agent, orchestrator, api, ransomware)
modules/bridge/     Python IPC bridge with handler modules
plugins/            Plugin ecosystem (worm, pulse-c2, AI, SDR, kernel)
web/                Vue 3 dashboard frontend
lab/                Docker lab environment
docs/               Documentation
```

## Documentation

| Document | Description |
|----------|-------------|
| [ARCHITECTURE.md](docs/ARCHITECTURE.md) | System architecture and component design |
| [COMMANDS.md](docs/COMMANDS.md) | CLI command reference |
| [MODULES.md](docs/MODULES.md) | Module matrix and phase coverage |
| [API_REFERENCE.md](docs/API_REFERENCE.md) | REST API endpoints |
| [DEPLOYMENT.md](docs/DEPLOYMENT.md) | Deployment and lab setup |

## Requirements

- Go 1.22+
- Python 3.11+
- Docker (optional, for lab environment)
- Node.js 22+ (optional, for dashboard development)

## License

MIT &copy; 2026 Rafael Galvez. See [LICENSE](LICENSE) for details.

> **Disclaimer**: This project is for educational purposes, authorized security research, and sanctioned red team engagements only. Unauthorized use against systems you do not own or have permission to test is illegal.