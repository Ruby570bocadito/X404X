package blockz

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type GeneticEvolutionEngine struct {
	Config         *BlockZConfig
	Population     []GeneticAgent     `json:"population"`
	Generation     int                `json:"generation"`
	BestFitness    float64            `json:"best_fitness"`
	MutationRate   float64            `json:"mutation_rate"`
	CrossoverRate  float64            `json:"crossover_rate"`
	mu             sync.Mutex
	sandboxPath    string
}

type GeneticAgent struct {
	ID          string   `json:"id"`
	Parents     []string `json:"parents"`
	Fitness     float64  `json:"fitness"`
	AVEvasion   float64  `json:"av_evasion_rate"`
	Size        int      `json:"size_bytes"`
	Survived    bool     `json:"survived"`
	Generation  int      `json:"generation"`
	Fingerprint string   `json:"fingerprint"`
}

type FitnessReport struct {
	AgentID        string  `json:"agent_id"`
	AVDetection    bool    `json:"av_detection"`
	SandboxEvasion bool    `json:"sandbox_evasion"`
	Latency        float64 `json:"latency_ms"`
	Score          float64 `json:"score"`
}

type GeneCrossover struct {
	ParentA   string `json:"parent_a"`
	ParentB   string `json:"parent_b"`
	ChildID   string `json:"child_id"`
	CutPointA int    `json:"cut_point_a"`
	CutPointB int    `json:"cut_point_b"`
	NewTraits []string `json:"new_traits"`
}

var darwinianTargets = map[string]string{
	"kernel32":  `C:\Windows\System32\kernel32.dll`,
	"ntdll":     `C:\Windows\System32\ntdll.dll`,
	"advapi32":  `C:\Windows\System32\advapi32.dll`,
	"ws2_32":    `C:\Windows\System32\ws2_32.dll`,
	"crypt32":   `C:\Windows\System32\crypt32.dll`,
	"user32":    `C:\Windows\System32\user32.dll`,
	"shell32":   `C:\Windows\System32\shell32.dll`,
	"chrome":    `C:\Program Files\Google\Chrome\Application\chrome.exe`,
	"firefox":   `C:\Program Files\Mozilla Firefox\firefox.exe`,
	"edge":      `C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
	"libc":      "/lib/x86_64-linux-gnu/libc.so.6",
	"libssl":    "/usr/lib/x86_64-linux-gnu/libssl.so.3",
	"python":    "/usr/bin/python3",
	"node":      "/usr/bin/node",
}

func NewGeneticEvolutionEngine(cfg *BlockZConfig) *GeneticEvolutionEngine {
	pop := make([]GeneticAgent, 0)
	for i := 0; i < 8; i++ {
		id := make([]byte, 16)
		rand.Read(id)
		pop = append(pop, GeneticAgent{
			ID:         hex.EncodeToString(id),
			Generation: 0,
			Fitness:    0.0,
		})
	}
	return &GeneticEvolutionEngine{
		Config:        cfg,
		Population:    pop,
		Generation:    0,
		MutationRate:  0.15,
		CrossoverRate: 0.70,
		sandboxPath:   filepath.Join(os.TempDir(), "x404x_sandbox"),
	}
}

func (ge *GeneticEvolutionEngine) Evolve(rounds int) []GeneticAgent {
	os.MkdirAll(ge.sandboxPath, 0755)

	for r := 0; r < rounds; r++ {
		ge.mu.Lock()
		ge.Generation = r
		ge.mu.Unlock()

		parents := ge.harvestGenePool()

		offspring := ge.crossoverPopulation(parents)

		offspring = ge.mutatePopulation(offspring)

		ge.evaluateFitness(offspring)

		ge.naturalSelection(offspring)

		_ = ge.generateHybridChild(parents)

		fmt.Printf("Generation %d: population=%d best=%.2f%%\n", r, len(ge.Population), ge.BestFitness*100)
	}

	return ge.Population
}

func (ge *GeneticEvolutionEngine) harvestGenePool() []string {
	var genes []string

	for name, path := range darwinianTargets {
		if data, err := os.ReadFile(path); err == nil && len(data) > 256 {
			hash := sha256.Sum256(data)
			gene := hex.EncodeToString(hash[:]) + ":" + name
			genes = append(genes, gene)
		}
	}

	if len(genes) == 0 {
		for i := 0; i < 8; i++ {
			entropy := make([]byte, 32)
			rand.Read(entropy)
			genes = append(genes, hex.EncodeToString(entropy)+":entropy")
		}
	}

	return genes
}

func (ge *GeneticEvolutionEngine) crossoverPopulation(genes []string) []GeneticAgent {
	var offspring []GeneticAgent

	active := ge.Population
	for i := 0; i < len(active)-1; i += 2 {
		if len(active[i].Parents) >= 4 {
			offspring = append(offspring, active[i])
			continue
		}

		roll := float64(i) / float64(len(active))
		if roll > ge.CrossoverRate {
			offspring = append(offspring, active[i], active[i+1])
			continue
		}

		childID := make([]byte, 16)
		rand.Read(childID)

		aGenes := []byte(active[i].Fingerprint + genes[i%len(genes)])
		bGenes := []byte(active[i+1].Fingerprint + genes[(i+1)%len(genes)])

		cutA := len(aGenes) / 3
		cutB := len(bGenes) * 2 / 3
		if len(bGenes[cutB:]) > len(aGenes[:cutA]) {
			bGenes = append(bGenes, aGenes[:cutA]...)
		} else {
			bGenes = append(aGenes[:cutA], bGenes[cutB:]...)
		}

		hybrid := sha256.Sum256(bGenes)
		childFP := hex.EncodeToString(hybrid[:])

		traits := []string{}
		for _, g := range []string{active[i].Fingerprint, active[i+1].Fingerprint} {
			if len(g) > 8 {
				traits = append(traits, g[:8])
			}
		}

		child := GeneticAgent{
			ID:          hex.EncodeToString(childID),
			Parents:     []string{active[i].ID, active[i+1].ID},
			Fingerprint: childFP,
			Generation:  ge.Generation,
			Fitness:     0.5,
		}
		offspring = append(offspring, child)

		_ = traits
	}

	return offspring
}

func (ge *GeneticEvolutionEngine) mutatePopulation(pop []GeneticAgent) []GeneticAgent {
	for i := range pop {
		mutationScore := float64(i) / float64(len(pop)+1)
		if mutationScore > ge.MutationRate {
			continue
		}

		xorKey := make([]byte, 8)
		rand.Read(xorKey)
		raw, _ := hex.DecodeString(pop[i].Fingerprint)
		if len(raw) > 0 {
			for j := 0; j < len(raw); j++ {
				raw[j] ^= xorKey[j%len(xorKey)]
			}
			pop[i].Fingerprint = hex.EncodeToString(raw)
		}
	}
	return pop
}

func (ge *GeneticEvolutionEngine) evaluateFitness(pop []GeneticAgent) {
	for i := range pop {
		av := ge.testAVEvasion(pop[i])
		sandbox := ge.testSandboxEvasion(pop[i])
		latency := float64(len(pop[i].Fingerprint)) * 0.01

		fitness := av*0.40 + sandbox*0.40 + (1.0-latency/100.0)*0.20
		if fitness > 1.0 {
			fitness = 1.0
		}
		if fitness < 0.0 {
			fitness = 0.0
		}

		pop[i].Fitness = fitness
		pop[i].AVEvasion = av

		ge.mu.Lock()
		if fitness > ge.BestFitness {
			ge.BestFitness = fitness
		}
		ge.mu.Unlock()
	}
}

func (ge *GeneticEvolutionEngine) testAVEvasion(agent GeneticAgent) float64 {
	score := 0.5

	entropy := sha256.Sum256([]byte(agent.Fingerprint))

	variation := float64(entropy[0]) / 255.0
	if variation > 0.85 {
		score += 0.30
	}
	if variation > 0.95 {
		score += 0.15
	}

	knownSigs := []string{"MZ", "PE\x00\x00", "This program", "GetProcAddress", "LoadLibrary"}
	sigMatch := false
	for _, sig := range knownSigs {
		if strings.Contains(agent.Fingerprint[:min(32, len(agent.Fingerprint))], sig[:min(4, len(sig))]) {
			sigMatch = true
			break
		}
	}
	if !sigMatch {
		score += 0.20
	}

	if runtime.GOOS == "windows" {
		_ = exec.Command("powershell", "-Command", "Get-MpThreatDetection")
	}

	return score
}

func (ge *GeneticEvolutionEngine) testSandboxEvasion(agent GeneticAgent) float64 {
	score := 0.7

	sandboxIndicators := []string{
		"sandbox", "malware", "virus", "analysis", "cuckoo",
		"VMware", "VirtualBox", "QEMU", "Xen",
		"SANDBOX", "VIRUS", "MALWARE",
	}
	hostname, _ := os.Hostname()
	for _, indicator := range sandboxIndicators {
		if strings.Contains(strings.ToLower(hostname), strings.ToLower(indicator)) {
			score *= 0.5
			break
		}
	}

	if runtime.GOOS == "windows" {
		exec.Command("powershell", "-Command",
			"(Get-WmiObject Win32_ComputerSystem).Model").Output()
	}

	_ = ge.sandboxPath
	return score
}

func (ge *GeneticEvolutionEngine) naturalSelection(pop []GeneticAgent) {
	var survivors []GeneticAgent

	sortPop(pop)

	for i := range pop {
		if pop[i].Fitness >= 0.40 || i < len(pop)/2 {
			pop[i].Survived = true
			survivors = append(survivors, pop[i])
		}
	}

	if len(survivors) < 4 {
		for i := 0; i < 4-len(survivors); i++ {
			id := make([]byte, 16)
			rand.Read(id)
			survivors = append(survivors, GeneticAgent{
				ID:         hex.EncodeToString(id),
				Fitness:    0.5,
				Survived:   true,
				Generation: ge.Generation,
			})
		}
	}

	ge.mu.Lock()
	ge.Population = survivors
	ge.mu.Unlock()
}

func (ge *GeneticEvolutionEngine) generateHybridChild(genes []string) GeneticAgent {
	if len(genes) < 2 {
		return GeneticAgent{}
	}

	a := genes[0]
	b := genes[1]

	combined := sha256.Sum256([]byte(a + b + fmt.Sprintf("%d", time.Now().UnixNano())))

	id := make([]byte, 16)
	rand.Read(id)

	child := GeneticAgent{
		ID:          hex.EncodeToString(id),
		Parents:     []string{"gene_pool_a", "gene_pool_b"},
		Fingerprint: hex.EncodeToString(combined[:]),
		Generation:  ge.Generation,
		Fitness:     0.6,
		Survived:    true,
	}
	return child
}

func (ge *GeneticEvolutionEngine) GetStatusJSON() string {
	ge.mu.Lock()
	defer ge.mu.Unlock()

	return fmt.Sprintf(`{"generation":%d,"population":%d,"best_fitness":%.4f,"mutation_rate":%.2f,"crossover_rate":%.2f}`,
		ge.Generation, len(ge.Population), ge.BestFitness, ge.MutationRate, ge.CrossoverRate)
}

func sortPop(pop []GeneticAgent) {
	for i := 0; i < len(pop); i++ {
		for j := i + 1; j < len(pop); j++ {
			if pop[j].Fitness > pop[i].Fitness {
				pop[i], pop[j] = pop[j], pop[i]
			}
		}
	}
}
