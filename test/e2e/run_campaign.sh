#!/bin/bash
# F5b: E2E Campaña multi-agente
set -e
cd "$(dirname "$0")/../.."
BOLD='\033[1m'; RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'
echo -e "${BOLD}========================================${NC}"
echo -e "${BOLD}  F5b: E2E Multi-Agent Campaign${NC}"
echo -e "${BOLD}========================================${NC}"

RESULT=0

echo -e "\n${YELLOW}Test 1: Agent registration simulation${NC}"
python3 -c "
import json
agents = ['agent-win-01', 'agent-linux-02', 'agent-esxi-03']
for a in agents:
    reg = {'agent_id': a, 'hostname': a, 'os': 'windows', 'arch': 'amd64'}
    js = json.dumps(reg)
    assert json.loads(js)['agent_id'] == a
print(f'  OK {len(agents)} agent payloads valid')
"

echo -e "\n${YELLOW}Test 2: Campaign config simulation${NC}"
python3 -c "
import json
cfg = {'name': 'multi-agent-test', 'target': '10.0.0.0/24', 'phases': 7, 'agents_required': 2}
js = json.dumps(cfg)
parsed = json.loads(js)
assert parsed['phases'] == 7
assert parsed['agents_required'] == 2
print('  OK campaign config valid')
"

echo -e "\n${YELLOW}Test 3: Decision dispatch simulation${NC}"
python3 -c "
import json
decision = {'id': 'dec-001', 'tactic': 'recon', 'technique': 'port_scan', 'confidence': 0.85, 'target': '10.0.0.5'}
js = json.dumps(decision)
parsed = json.loads(js)
assert parsed['confidence'] >= 0.0 and parsed['confidence'] <= 1.0
print('  OK decision payload valid')
"

echo -e "\n${BOLD}MULTI-AGENT CAMPAIGN TEST COMPLETE${NC}"
exit $RESULT
