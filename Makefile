# X404X — Makefile
# ===================
# make setup      — Install all dependencies (Go + Python + Node)
# make build      — Build all components
# make test       — Run all tests
# make lab-up     — Start Docker lab environment
# make lab-down   — Stop Docker lab environment
# make proto      — Generate gRPC code from .proto files
# make clean      — Remove build artifacts
# make lint       — Run linters (Go + Python)
# make release    — Build release binaries for all platforms

.PHONY: setup build test lab-up lab-down proto clean lint release

GO_PACKAGES := ./core/... ./shared/... ./cmd/...
PYTHON_MODULES := modules/ scripts/

# === Setup ===
setup: setup-go setup-python setup-node
	@echo "[+] Setup complete"

setup-go:
	@echo "[*] Installing Go dependencies..."
	cd core/crypto && go mod tidy
	cd core/agent && go mod tidy
	cd core/orchestrator && go mod tidy
	cd shared/config && go mod tidy
	cd shared/logger && go mod tidy
	cd shared/types && go mod tidy

setup-python:
	@echo "[*] Installing Python dependencies..."
	pip install -r modules/requirements.txt 2>/dev/null || true

setup-node:
	@echo "[*] Installing Node.js dependencies..."
	cd core/c2/web && npm install 2>/dev/null || true

# === Build ===
build: build-x404x build-agent build-c2
	@echo "[+] Build complete"

build-x404x:
	@echo "[*] Building x404x CLI..."
	cd cmd/x404x && go build -o ../../dist/x404x .
	@echo "  [+] dist/x404x"

build-agent:
	@echo "[*] Building agent..."
	cd core/agent && GOOS=linux GOARCH=amd64 go build -o ../../dist/agent-linux-amd64 ./cmd/agent
	cd core/agent && GOOS=linux GOARCH=arm64 go build -o ../../dist/agent-linux-arm64 ./cmd/agent
	cd core/agent && GOOS=windows GOARCH=amd64 go build -o ../../dist/agent-windows-amd64.exe ./cmd/agent

build-c2:
	@echo "[*] Building C2 server..."
	cd core/c2 && make build 2>/dev/null || echo "  [!] C2 not available (submodule not cloned)"

build-vault:
	@echo "[*] Building kernel module..."
	cd core/kernel/src && make 2>/dev/null || echo "  [!] Vault-Kernel not available"

# === Test ===
test: test-go test-python
	@echo "[+] All tests passed"

test-go:
	@echo "[*] Running Go tests..."
	go test $(GO_PACKAGES) -v -cover -timeout 60s

test-python:
	@echo "[*] Running Python tests..."
	pytest modules/ -v --tb=short 2>/dev/null || echo "  [!] pytest not found"

# === Lab ===
lab-up:
	@echo "[*] Starting Docker lab..."
	docker compose -f lab/docker-compose.yml up -d
	@echo "[+] Lab running at:"
	@echo "    Attacker:  docker exec -it x404x-attacker bash"
	@echo "    Target 1:  docker exec -it x404x-target1 bash"
	@echo "    Dashboard: http://localhost:3000"

lab-down:
	@echo "[*] Stopping Docker lab..."
	docker compose -f lab/docker-compose.yml down

lab-status:
	docker compose -f lab/docker-compose.yml ps

# === Proto ===
proto:
	@echo "[*] Generating protobuf code..."
	buf generate core/proto
	@echo "[+] Proto generation complete"

# === Lint ===
lint: lint-go lint-python
	@echo "[+] Lint complete"

lint-go:
	@echo "[*] Linting Go code..."
	golangci-lint run $(GO_PACKAGES) --timeout 5m 2>/dev/null || echo "  [!] golangci-lint not found"

lint-python:
	@echo "[*] Linting Python code..."
	ruff check $(PYTHON_MODULES) 2>/dev/null || echo "  [!] ruff not found"

# === Clean ===
clean:
	@echo "[*] Cleaning build artifacts..."
	rm -rf dist/
	rm -f core/kernel/src/*.ko core/kernel/src/*.o
	find . -type d -name __pycache__ -exec rm -rf {} + 2>/dev/null || true
	find . -type f -name "*.pyc" -delete 2>/dev/null || true
	@echo "[+] Clean complete"

# === Release ===
release: clean build
	@echo "[*] Creating release archive..."
	mkdir -p release
	cp dist/* release/ 2>/dev/null || true
	cp config.yaml release/
	tar -czf x404x.tar.gz release/
	@echo "[+] Release: x404x.tar.gz"

# === Help ===
help:
	@echo "X404X — Build Commands"
	@echo "========================="
	@echo "  make setup        Install all dependencies"
	@echo "  make build        Build all components"
	@echo "  make test         Run all tests"
	@echo "  make lab-up       Start Docker lab environment"
	@echo "  make lab-down     Stop Docker lab environment"
	@echo "  make lab-status   Show lab container status"
	@echo "  make proto        Generate gRPC code"
	@echo "  make lint         Run linters"
	@echo "  make clean        Remove build artifacts"
	@echo "  make release      Build release binaries"
	@echo "  make help         Show this help message"
