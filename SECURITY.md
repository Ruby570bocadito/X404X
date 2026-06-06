# Security Policy

## Reporting a Vulnerability

X404X is an academic project for cybersecurity education. If you discover a security vulnerability, please:

1. **DO NOT** open a public GitHub issue
2. Email the project maintainer with details
3. Allow 72 hours for initial response

## Supported Versions

| Version | Supported |
|---------|-----------|
| v2.0    | ✅ Active |
| v1.0    | ❌ End of life |

## Security Model

X404X is designed for **authorized use only** in:
- Controlled laboratory environments
- CTF competitions with explicit permission
- Penetration testing with written authorization

### Safety Controls

| Control | Default | Purpose |
|---------|---------|---------|
| Kill Switch | Enabled | Emergency stop all agents |
| Geofencing | Enabled | RFC 1918 private networks only |
| Auto-Destruct | 2 hours | Agents self-terminate |
| Max Infections | 1000 | Hard limit on compromised hosts |
| No Persistence | Enabled | Persistence requires explicit activation |
| Offline AI | Default | Ollama runs locally, no data exfiltration |

### Dependency Security

Dependencies are managed via:
- Go: `go.mod` + `go.sum` (cryptographic verification)
- Python: `requirements.txt` (pinned versions)
- Node.js: `package.json` + `package-lock.json` (audited)

### Responsible Disclosure

We follow the [RFPolicy](https://en.wikipedia.org/wiki/RFPolicy) for vulnerability disclosure.
