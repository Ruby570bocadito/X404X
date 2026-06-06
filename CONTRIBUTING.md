# Contributing to X404X

## Code of Conduct

X404X is an academic project for cybersecurity education. All contributions must align with the ethical guidelines:

1. Only use in authorized environments (labs, CTFs, penetration tests with written permission).
2. Do not submit exploits against live targets.
3. All AI operations use local models (Ollama) — no data exfiltration.

## Getting Started

```bash
git clone --recurse-submodules https://github.com/Ruby570bocadito/X404X.git
cd X404X
make setup
make build
make test
```

## Project Structure

```
X404X/
├── cmd/x404x/        # CLI entry point (Cobra + TUI + Console)
├── core/
│   ├── agent/        # Unified Go implant
│   ├── orchestrator/ # Decision Engine + Campaign Manager
│   ├── appstate/     # Shared application state
│   ├── api/          # REST API + WebSocket
│   ├── c2server/     # C2 listener
│   ├── crypto/       # X25519 + XChaCha20-Poly1305
│   └── proto/        # gRPC definitions
├── modules/          # Python modules + submodules
├── shared/           # Config, Logger, Types, Database
├── web/              # Vue 3 Dashboard
└── lab/              # Docker lab environment
```

## Adding a Module

### Go Module

```go
package agent

type MyModule struct{}

func (m *MyModule) Name() string { return "my_module" }
func (m *MyModule) KillChainPhase() types.KillChainPhase { return types.PhaseExploitation }
func (m *MyModule) Execute(ctx context.Context, params map[string]string) (string, error) {
    // Your logic here
    return "success", nil
}

// Register:
agent.RegisterModule(&MyModule{})
```

### Python Module

```python
from modules.bridge.bridge import registry

@registry.register("my_python_module", "Description", "1.0", "recon")
def my_python_module_handler(params: dict):
    target = params.get("target", "127.0.0.1")
    # Your logic here
    return {"target": target, "status": "ok"}
```

### Console Module

Add to `core/appstate/state.go`:
```go
ModuleDef{Name: "exploit/my_module", Type: "exploit",
    Description: "My custom exploit", Rank: "normal", OS: "any"}
```

## Testing

```bash
go test ./core/... ./shared/...        # Go tests
pytest modules/                         # Python tests
npm run build                           # Vue build check
```

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing`)
3. Add tests for new functionality
4. Ensure `make build && make test` passes
5. Update documentation if needed
6. Submit PR with description of changes

## License

MIT License — see [LICENSE](LICENSE) for details.
