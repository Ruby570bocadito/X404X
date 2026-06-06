package v26

type V26Config struct {
	Enabled       bool   `json:"enabled"`
	Simulation    bool   `json:"simulation"`
	POMDPEnable   bool   `json:"pomdp_enable"`
	AIGotchaEnable bool  `json:"ai_gotcha_enable"`
	EvasionDEEPEnable bool `json:"evasion_deep_enable"`
	BootkitSMMEnable bool `json:"bootkit_smm_enable"`
	MobileXEnable    bool `json:"mobile_x_enable"`
	CloudNemesisEnable bool `json:"cloud_nemesis_enable"`
	SocialC2Enable   bool `json:"social_c2_enable"`
	BlockOmegaEnable bool `json:"block_omega_enable"`
	C2Endpoint    string `json:"c2_endpoint"`
}

func DefaultV26Config() *V26Config {
	return &V26Config{
		Enabled: true,
		Simulation: true,
		C2Endpoint: "x404x-c2.online:8443",
	}
}
