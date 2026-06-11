#!/bin/bash
# X404X Framework Setup Script
# Installs dependencies: Go, Node.js, Python, and Ollama (AI)

set -e

GREEN='\033[38;5;46m'
RED='\033[38;5;196m'
BLUE='\033[38;5;39m'
VIOLET='\033[38;5;99m'
NC='\033[0m'

print_info() { echo -e "  ${BLUE}[*]${NC} $1"; }
print_ok() { echo -e "  ${GREEN}[+]${NC} $1"; }
print_err() { echo -e "  ${RED}[!]${NC} $1"; }

DEFAULT_MODEL="llama3.2"
MODEL=$DEFAULT_MODEL

for arg in "$@"; do
    if [[ "$arg" == "--model="* ]]; then
        MODEL="${arg#*=}"
    fi
done

echo -e "\n${VIOLET}  ██╗  ██╗██╗  ██╗ ██████╗ ██╗  ██╗██╗  ██╗${NC}"
echo -e "${VIOLET}  ╚██╗██╔╝██║  ██║██╔═══██╗██║  ██║╚██╗██╔╝${NC}"
echo -e "${VIOLET}   ╚███╔╝ ███████║██║   ██║███████║ ╚███╔╝ ${NC}"
echo -e "${VIOLET}   ██╔██╗ ╚════██║██║   ██║╚════██║ ██╔██╗ ${NC}"
echo -e "${VIOLET}  ██╔╝ ██╗      ██║╚██████╔╝      ██║██╔╝ ██╗${NC}"
echo -e "${VIOLET}  ╚═╝  ╚═╝      ╚═╝ ╚═════╝       ╚═╝╚═╝  ╚═╝${NC}"
echo -e "\n  X404X Environment Setup\n"

# 1. Check/Install Python and venv
if ! command -v python3 &> /dev/null; then
    print_info "Installing Python3..."
    if command -v apt-get &> /dev/null; then
        sudo apt-get update && sudo apt-get install -y python3 python3-pip python3-venv
    else
        print_err "Please install Python3 manually."
    fi
else
    print_ok "Python3 is installed."
    # Check if venv is available
    if ! python3 -c "import ensurepip" &> /dev/null; then
        print_info "python3-venv not found. Installing..."
        if command -v apt-get &> /dev/null; then
            # Get specific python version like python3.14-venv or default python3-venv
            PY_VER=$(python3 -c 'import sys; print(f"python{sys.version_info.major}.{sys.version_info.minor}-venv")')
            sudo apt-get update && sudo apt-get install -y python3-venv $PY_VER || sudo apt-get install -y python3-venv
        else
            print_err "Please install python3-venv manually."
        fi
    fi
fi

# 2. Python Dependencies
print_info "Setting up Python Bridge dependencies..."
python3 -m venv venv
source venv/bin/activate
if [ -f "requirements.txt" ]; then
    pip install -r requirements.txt > /dev/null 2>&1
    print_ok "Python dependencies installed."
else
    print_info "No requirements.txt found, skipping."
fi
deactivate

# 3. Check/Install Go
if ! command -v go &> /dev/null; then
    print_info "Go not found. Please install Go 1.24+ manually: https://go.dev/dl/"
else
    print_ok "Go is installed: $(go version | awk '{print $3}')"
fi

# 4. Check/Install Node.js
if ! command -v npm &> /dev/null; then
    print_info "Node.js not found. Installing via apt..."
    if command -v apt-get &> /dev/null; then
        sudo apt-get install -y nodejs npm
    else
        print_err "Please install Node.js manually."
    fi
else
    print_ok "Node.js (npm) is installed."
fi

# 5. Build Frontend
print_info "Building Vue frontend..."
if [ -d "web" ]; then
    cd web
    npm install > /dev/null 2>&1
    npm run build > /dev/null 2>&1
    cd ..
    print_ok "Frontend built successfully."
else
    print_err "Directory 'web' not found, skipping frontend build."
fi

# 6. Build Backend
print_info "Building Go backend..."
GOOS=linux GOARCH=amd64 go build -o x404x ./cmd/x404x/
print_ok "Backend built (x404x binary created)."

# 7. Check/Install Ollama
if ! command -v ollama &> /dev/null; then
    print_info "Ollama not found. Installing..."
    curl -fsSL https://ollama.com/install.sh | sh
    print_ok "Ollama installed."
else
    print_ok "Ollama is installed."
fi

# 8. Pull AI Model
print_info "Ensuring Ollama is running..."
if ! pgrep -x "ollama" > /dev/null; then
    ollama serve &>/dev/null &
    sleep 3
fi

print_info "Pulling LLM model: $MODEL (this may take a while)..."
ollama pull $MODEL
print_ok "Model $MODEL is ready."

echo -e "\n${GREEN}  [✓] Setup Complete!${NC}"
echo -e "  To start the dashboard, run: ${VIOLET}./x404x dashboard${NC}\n"
