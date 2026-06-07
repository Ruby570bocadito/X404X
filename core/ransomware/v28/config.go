package v28

type V28Config struct {
	Enabled            bool   `json:"enabled"`
	Simulation         bool   `json:"simulation"`
	IoTIdentityTheft    bool `json:"iot_identity_theft"`
	FalseMemoryInjection bool `json:"false_memory_injection"`
	DeathByThousandCuts bool `json:"death_by_thousand_cuts"`
	PatchGuardBypass   bool   `json:"patch_guard_bypass"`
	KeyboardLEDExfil   bool   `json:"keyboard_led_exfil"`
	ZombieArmyPolitical bool  `json:"zombie_army_political"`
	LegacyPoison       bool   `json:"legacy_poison"`
	SEOSabotage        bool   `json:"seo_sabotage"`
	FakeVulnInjection  bool   `json:"fake_vuln_injection"`
	InceptionHypervisor bool  `json:"inception_hypervisor"`
	ISPBGPSubversion   bool   `json:"isp_bgp_subversion"`
	AntiAttributionClone bool `json:"anti_attribution_clone"`
	PowerGridHarmonics  bool  `json:"power_grid_harmonics"`
	TimeLockExtortion   bool  `json:"time_lock_extortion"`
	VRSpyware          bool   `json:"vr_spyware"`
	GlobalAIPoison     bool   `json:"global_ai_poison"`
	CDNMalwareInjection bool  `json:"cdn_malware_injection"`
	BioCyberDNA        bool   `json:"bio_cyber_dna"`
	BrowserParasite    bool   `json:"browser_parasite"`
	FakeDocumentsGen   bool   `json:"fake_documents_gen"`
	SoundPanicAttack   bool   `json:"sound_panic_attack"`
	EmotionalEncryption bool  `json:"emotional_encryption"`
	FalseRedemption    bool   `json:"false_redemption"`
	C2Endpoint         string `json:"c2_endpoint"`
}

func DefaultV28Config() *V28Config {
	return &V28Config{Enabled: true, Simulation: true, C2Endpoint: "x404x-c2.online:8443"}
}
