# X404X — Complete Implementation Roadmap

## Phase 1: Skeleton ✓
- [x] Git init, .gitignore, .gitmodules
- [x] Directory structure
- [x] go.work (Go workspace)
- [x] core/proto/ (gRPC definitions)
- [x] core/crypto/ (X25519 + XChaCha20-Poly1305)
- [x] shared/config/ (YAML config)
- [x] shared/logger/ (zap structured logging)
- [x] shared/types/ (shared domain types)
- [x] Makefile
- [x] CI/CD (.github/workflows)
- [x] Docs

## Phase 2: Core + DB (Week 1-2)
- [ ] shared/database/migrations (Alembic)
- [ ] core/crypto/ → integrate with Pulse-C2
- [ ] core/proto/ → generate Go/Python stubs
- [ ] scripts/setup.sh testing
- [ ] Submodules cloned and verified

## Phase 3: Orchestrator (Week 3-4)
- [ ] Event bus messaging
- [ ] Campaign manager implementation
- [ ] Decision engine (Rules + A* + AI)
- [ ] HITL workflow
- [ ] Integration with Pulse-C2

## Phase 4: Agent + Bridge (Week 5-6)
- [ ] Agent gRPC connector to Pulse-C2
- [ ] Module manager (load/unload modules)
- [ ] Python Bridge IPC (socket protocol)
- [ ] Adapt existing modules to bridge interface
- [ ] Evasion engine integration

## Phase 5: CLI + Dashboard (Week 7-8)
- [ ] rbyhack CLI (cobra commands)
- [ ] Dashboard Vue 3 extensions
- [ ] WebSocket real-time updates
- [ ] REST API endpoints
- [ ] Campaign visualization

## Phase 6: AI Integration (Week 9-10)
- [ ] Specter-Terminal → Orchestrator bridge
- [ ] Apex-Automation → Decision engine
- [ ] Auto-approval mode
- [ ] AI confidence scoring
- [ ] Prompt templates for attack context

## Phase 7: Lab + Tests (Week 11-12)
- [ ] Docker lab with all services
- [ ] End-to-end integration tests
- [ ] CTF scenario definitions
- [ ] BlueForge detection validation
- [ ] Performance benchmarks

## Phase 8: TFG Documentation (Week 13-14)
- [ ] Technical memory (Spanish)
- [ ] Architecture diagrams (Mermaid/PlantUML)
- [ ] Ethical/legal analysis
- [ ] User manual
- [ ] Defense recommendations
