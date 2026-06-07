#!/bin/bash
# X404X Binary Digital Signing
# Generates self-signed cert and signs the binary

BIN="${1:-x404x}"
OUT="${BIN}.signed"
KEY="/tmp/x404x_signing_key.pem"
CERT="/tmp/x404x_signing_cert.pem"

echo "=== X404X Digital Signature ==="

# Generate RSA-4096 signing key
openssl genrsa -out "$KEY" 4096 2>/dev/null && echo "Key generated: $KEY"

# Generate self-signed cert
openssl req -new -x509 -key "$KEY" -out "$CERT" -days 3650 \
    -subj "/CN=X404X Code Signing/O=X404X/C=ES" 2>/dev/null && echo "Cert generated: $CERT"

# Generate signature
if [ -f "$BIN" ]; then
    openssl dgst -sha256 -sign "$KEY" -out "$BIN.sig" "$BIN" 2>/dev/null && echo "Binary signed: $BIN.sig"
    openssl dgst -sha256 -verify <(openssl x509 -in "$CERT" -pubkey -noout) \
        -signature "$BIN.sig" "$BIN" 2>/dev/null && echo "Signature VERIFIED ✓" || echo "Signature FAILED ✗"
    HASH=$(sha256sum "$BIN" | cut -d' ' -f1)
    echo "SHA256: $HASH"
    echo "Size: $(du -h "$BIN" | cut -f1)"
else
    echo "Binary not found: $BIN"
fi
echo "=== Done ==="
